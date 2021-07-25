package websvc

import (
	"sync"
	"time"
)

type taskFn func() string

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
		for {
			select {
			case <-ticker.C:
				iter += 1
				status := t.execute()
				t.log.info("[timer task %s] iteration %d: %s", t.name, iter, status)
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
