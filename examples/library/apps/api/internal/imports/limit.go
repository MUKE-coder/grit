// Package imports keeps CSV imports from all running at once.
package imports

import (
	"os"
	"strconv"
	"sync"
)

// Concurrency is how many imports run at the same time, from IMPORT_CONCURRENCY
// (default 2). The others wait their turn.
//
// Each import holds a database connection for its whole run and writes in
// batches, so a burst of them took the connections every request shares.
var Concurrency = concurrencyFromEnv()

var slots = make(chan struct{}, Concurrency)

func concurrencyFromEnv() int {
	if n, err := strconv.Atoi(os.Getenv("IMPORT_CONCURRENCY")); err == nil && n > 0 {
		return n
	}
	return 2
}

// Wait blocks until an import may run, and returns the function that ends its
// turn. onQueued runs first when the import has to wait, so its job can say so.
func Wait(onQueued func()) (release func()) {
	select {
	case slots <- struct{}{}:
	default:
		if onQueued != nil {
			onQueued()
		}
		slots <- struct{}{}
	}
	var once sync.Once
	return func() { once.Do(func() { <-slots }) }
}
