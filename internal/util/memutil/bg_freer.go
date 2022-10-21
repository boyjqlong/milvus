package memutil

import (
	"runtime/debug"
	"time"
)

type bgFreer struct {
	ch chan struct{}
}

func newBgFreer() *bgFreer {
	return &bgFreer{ch: make(chan struct{})}
}

func (bg *bgFreer) freeOSMemoryPeriodically() {
	const interval = time.Second * 5
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			debug.FreeOSMemory()
		case <-bg.ch:
			return
		}
	}
}

func (bg *bgFreer) stop() {
	close(bg.ch)
}
