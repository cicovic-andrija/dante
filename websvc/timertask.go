package websvc

import (
	"fmt"
	"sync"
	"time"
)

type taskFn func(args ...interface{}) (string, bool)

type timerTask struct {
	name    string
	execute taskFn
	stopped bool
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
				t.log.info("[timer task %q] stopped", t.name)
				wg.Done()
				t.stopped = true
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

	tasks []*timerTask
	wg    *sync.WaitGroup
}

// This is just a help method to be used during server boot.
// It should not be used later.
func (t *timerTaskManager) addTask(name string, fn taskFn, period time.Duration, log *logstruct) {
	t.Lock()
	t.tasks = append(
		t.tasks,
		&timerTask{name: name, execute: fn, period: period, log: log},
	)
	t.Unlock()
}

// This is just a help method to be used during server boot.
// It should not be used later.
func (t *timerTaskManager) runAll() {
	t.Lock()
	for _, task := range t.tasks {
		task.run(t.wg)
		t.wg.Add(1)
	}
	t.Unlock()
}

func (t *timerTaskManager) scheduleTask(task *timerTask, args ...interface{}) {
	t.Lock()
	t.tasks = append(t.tasks, task)
	task.run(t.wg, args...)
	t.wg.Add(1)
	t.Unlock()
}

func (t *timerTaskManager) stopAll() {
	t.Lock()
	for _, task := range t.tasks {
		task.stop()
	}
	t.Unlock()
}
