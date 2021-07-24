package websvc

import (
	"time"
)

// FIXME: Confirm timer task stopped.

type taskFn func()

type timerTask struct {
	name   string
	task   taskFn
	period time.Duration
	log    *logstruct
	iter   uint64
	quit   chan struct{}
}

func (t *timerTask) run() {
	t.quit = make(chan struct{})
	go func() {
		ticker := time.NewTicker(t.period)
		for {
			select {
			case <-ticker.C:
				t.iter += 1
				t.log.info("timer task %q iteration %d", t.name, t.iter)
				t.task()
				t.log.info("timer task %q iteration %d done", t.name, t.iter)
			case <-t.quit:
				t.log.info("timer task %q stopping", t.name)
				ticker.Stop()
				return
			}
		}
	}()
}

func (t *timerTask) stop() {
	close(t.quit)
}
