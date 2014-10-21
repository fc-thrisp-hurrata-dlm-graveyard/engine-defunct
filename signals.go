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
				e.Message(fmt.Sprintf("%s", sig))
			}
		}
	}()
}

func (e *Engine) Message(message string) {
	if e.LoggingOn {
		e.Logger.Printf(" %s", message)
	}
}

func (e *Engine) PanicsNow(message string) {
	log.Println(fmt.Errorf("[ENGINE] %s", message))
	e.Signals <- []byte("Engine-Panic")
}

// Send sends a message to the specified queue.
func (e *Engine) Send(queue string, message string) {
	go func() {
		if q, ok := e.Queues[queue]; ok {
			q(message)
		} else {
			e.SendSignal(message)
		}
	}()
}

// SendSignal sends a signal as []byte directly to engine.Signals
func (e *Engine) SendSignal(message string) {
	e.Signals <- []byte(message)
}

func (e *Engine) defaultqueues() queues {
	q := make(queues)
	q["messages"] = e.Message
	q["panics-now"] = e.PanicsNow
	q["signal"] = e.SendSignal
	return q
}
