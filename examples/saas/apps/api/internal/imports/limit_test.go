package imports

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestImportsTakeTurns(t *testing.T) {
	var running, peak, queued atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < Concurrency+3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			release := Wait(func() { queued.Add(1) })
			defer release()
			n := running.Add(1)
			for {
				p := peak.Load()
				if n <= p || peak.CompareAndSwap(p, n) {
					break
				}
			}
			time.Sleep(20 * time.Millisecond)
			running.Add(-1)
		}()
	}
	wg.Wait()
	if int(peak.Load()) > Concurrency {
		t.Errorf("%d imports ran at once; the limit is %d", peak.Load(), Concurrency)
	}
	if queued.Load() == 0 {
		t.Error("no import was told it was waiting")
	}
}

func TestReleasingTwiceFreesOneTurn(t *testing.T) {
	release := Wait(nil)
	release()
	release()
	if got := len(slots); got != 0 {
		t.Errorf("%d turns held after one import released twice", got)
	}
}
