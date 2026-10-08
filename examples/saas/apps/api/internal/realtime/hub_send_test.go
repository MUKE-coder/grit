package realtime

import (
	"io"
	"log"
	"os"
	"sync"
	"testing"
)

// A send racing a closing connection must not panic.
//
// deliverLocal and broadcastLocal used to copy their targets under the read
// lock and send after releasing it. An Unregister in that gap closed Send
// first, the send panicked with "send on closed channel", and a panic on a
// goroutine takes the whole API process down, not just one socket.
func TestSendsRacingUnregisterDoNotPanic(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	hub := NewHub()
	const rounds, clientsPerRound = 50, 200
	sends := 0
	for round := 0; round < rounds; round++ {
		clients := make([]*Client, clientsPerRound)
		for i := range clients {
			clients[i] = &Client{UserID: "u1", Send: make(chan []byte, 1)}
			hub.Register(clients[i])
		}
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			defer wg.Done()
			for _, c := range clients {
				hub.Unregister(c)
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				hub.Broadcast(Event{Type: "ping"})
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				hub.SendToUsers([]string{"u1"}, Event{Type: "ping"})
			}
		}()
		wg.Wait()
		sends += 40
	}
	t.Logf("%d sends raced %d closing connections without a panic", sends, rounds*clientsPerRound)
}
