package websvc

import (
	"time"
)

type taskFn func()

type timerTask struct {
	ticker *time.Ticker
	quit   chan struct{}
}

func (t *timerTask) run(task taskFn, period time.Duration) {
	t.ticker = time.NewTicker(period)
	t.quit = make(chan struct{})
	go func() {
		for {
			select {
			case <-t.ticker.C:
				task()
			case <-t.quit:
				t.ticker.Stop()
				return
			}
		}
	}()
}

func (t *timerTask) stop() {
	close(t.quit)
}
