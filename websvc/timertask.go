package websvc

import (
	"fmt"
	"sync"
	"time"
)

type taskFn func() (string, bool)

type timerTask struct {
	name    string
	execute taskFn
	period  time.Duration
	log     *logstruct
	quit    chan struct{}
}

func (t *timerTask) run(wg *sync.WaitGroup) {
	t.quit = make(chan struct{})
	go func() {
		ticker := time.NewTicker(t.period)
		iter := 0
		t.log.info("[task %s] started", t.name)
		for {
			select {
			case <-ticker.C:
				iter += 1
				status, failed := t.execute()
				t.log.info("[task %s] iteration %d: %s", t.name, iter, status)
				if failed {
					t.log.err("[task %s] iteration %d: %s", t.name, iter, status)
				}
			case <-t.quit:
				ticker.Stop()
				t.log.info("[timer task %q] stopped", t.name)
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
