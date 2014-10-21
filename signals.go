package engine

import (
	"fmt"
	"log"
)

type (
	// Signal denotes a byte signal
	Signal []byte

	// Signals is a channel for Signal
	Signals chan Signal

	queue func(string)

	queues map[string]queue
)

// A simple Signal queue outputting everything emitted to engine.Signals.
func SignalQueue(e *Engine) {
	go func() {
		for {
			select {
			case sig := <-e.Signals:
				e.Log(fmt.Sprintf("%s", sig))
			}
		}
	}()
}

func (e *Engine) Log(message string) {
	if e.LoggingOn {
		e.Logger.Printf(" %s", message)
	}
}

func (e *Engine) PanicsNow(message string) {
	log.Println(fmt.Errorf("[ENGINE-PANIC] %s", message))
	e.Signals <- []byte("Engine-Panic")
}

func (e *Engine) Send(q string, message string) {
	go func() {
		if queue, ok := e.Queues[q]; ok {
			queue(message)
		} else {
			e.Signals <- []byte(message)
		}
	}()
}

func (e *Engine) defaultqueues() queues {
	q := make(queues)
	q["messages"] = e.Log
	q["panics-now"] = e.PanicsNow
	return q
}
