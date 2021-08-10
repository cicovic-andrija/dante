package websvc

import (
	"fmt"
	"sync"
	"time"
)

const (
	stopChBufferSz = 32
)

type taskFn func(args ...interface{}) (string, bool)

type timerTask struct {
	name    string
	execute taskFn
	period  time.Duration
	log     *logstruct
	quit    chan struct{}
}

func (t *timerTask) run(wg *sync.WaitGroup, args ...interface{}) {
	t.quit = make(chan struct{})
	go func() {
		ticker := time.NewTicker(t.period)
		iter := 0
		t.log.info("[task %s] started", t.name)
		for {
			select {
			case <-ticker.C:
				iter += 1
				status, failed := t.execute(args...)
				t.log.info("[task %s] iteration %d: %s", t.name, iter, status)
				if failed {
					t.log.err("[task %s] iteration %d: %s", t.name, iter, status)
				}
			case <-t.quit:
				ticker.Stop()
				t.log.info("[task %s] stopped", t.name)
				wg.Done()
				return
			}
		}
	}()
}

func (t *timerTask) stop() {
	close(t.quit)
}

func timerTaskSuccess(message string) (string, bool) {
	return fmt.Sprintf("successful: %s", message), false
}

func timerTaskFailure(err error) (string, bool) {
	return fmt.Sprintf("failed with error: %v", err), true
}

type timerTaskManager struct {
	sync.Mutex

	// self management
	quit chan struct{}

	// task management
	tasks  map[string]*timerTask
	stopCh chan string
	wg     *sync.WaitGroup
}

func newTimerTaskManager() *timerTaskManager {
	return &timerTaskManager{
		quit:   make(chan struct{}),
		tasks:  make(map[string]*timerTask),
		stopCh: make(chan string, stopChBufferSz),
		wg:     &sync.WaitGroup{},
	}
}

// Run task manager itself.
func (t *timerTaskManager) run(wg *sync.WaitGroup) {
	go func() {
		for {
			select {
			case taskName := <-t.stopCh:
				t.Lock()
				if task, ok := t.tasks[taskName]; ok {
					task.stop()
					delete(t.tasks, taskName) // is this ok?
				}
				t.Unlock()
			case <-t.quit:
				wg.Done()
				return
			}
		}
	}()
}

// Stop task manager itself.
func (t *timerTaskManager) stop() {
	close(t.quit)
	close(t.stopCh)
}

func (t *timerTaskManager) scheduleTask(task *timerTask, args ...interface{}) {
	t.Lock()
	t.tasks[task.name] = task
	t.wg.Add(1)
	task.run(t.wg, args...)
	t.Unlock()
}

func (t *timerTaskManager) stopTask(name string) {
	t.Lock()
	if task, ok := t.tasks[name]; ok {
		task.stop()
		delete(t.tasks, name)
	}
	t.Unlock()
}

func (t *timerTaskManager) stopAll() {
	t.Lock()
	for _, task := range t.tasks {
		task.stop()
	}
	t.Unlock()
}

// This is just a help method to be used during server boot.
// It should not be used later.
func (t *timerTaskManager) addTask(name string, fn taskFn, period time.Duration, log *logstruct) {
	t.Lock()
	t.tasks[name] = &timerTask{
		name:    name,
		execute: fn,
		period:  period,
		log:     log,
	}
	t.Unlock()
}

// This is just a help method to be used during server boot.
// It should not be used later.
func (t *timerTaskManager) runAllTasks() {
	t.Lock()
	for _, task := range t.tasks {
		task.run(t.wg)
		t.wg.Add(1)
	}
	t.Unlock()
}
