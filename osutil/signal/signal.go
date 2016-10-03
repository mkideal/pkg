package signal

import (
	"os"
	"os/signal"
	"sync"
)

type SignalHandler func(os.Signal) bool

var (
	mu       sync.Mutex
	handlers = map[os.Signal][]SignalHandler{}
	sigChan  = make(chan os.Signal)
)

func Register(sig os.Signal, handler SignalHandler) {
	mu.Lock()
	defer mu.Unlock()
	if g, ok := handlers[sig]; ok {
		handlers[sig] = append(g, handler)
	} else {
		handlers[sig] = []SignalHandler{handler}
	}
	signal.Notify(sigChan, sig)
}

func Listen() {
	for {
		select {
		case sig := <-sigChan:
			mu.Lock()
			g, ok := handlers[sig]
			mu.Unlock()
			if ok {
				for _, h := range g {
					if h(sig) {
						return
					}
				}
			}
		}
	}
}
