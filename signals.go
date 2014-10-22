package engine

import (
	"fmt"
	"log"
)

var (
	EnginePanic = []byte("engine-panic")
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
		for sig := range e.Signals {
			e.Message(fmt.Sprintf("%s", sig))
		}
	}()
}

// Message goes directly to a logger, if enbaled.
func (e *Engine) Message(message string) {
	if e.LoggingOn {
		e.Logger.Printf(" %s", message)
	}
}

// PanicMessage goes to a standard and unavaoidable log, then emits a signal.
func (e *Engine) PanicMessage(message string) {
	log.Println(fmt.Errorf("[ENGINE] %s", message))
	e.Signals <- EnginePanic
	//e.Signals <- []byte(message)
}

// Emit goes as []byte directly to engine.Signals
func (e *Engine) Emit(message string) {
	e.Signals <- []byte(message)
}

// Send sends a message to the specified queue.
func (e *Engine) Send(queue string, message string) {
	go func() {
		if q, ok := e.Queues[queue]; ok {
			q(message)
		} else {
			e.Emit(message)
		}
	}()
}

func (e *Engine) defaultqueues() queues {
	q := make(queues)
	q["message"] = e.Message
	q["panic"] = e.PanicMessage
	q["emit"] = e.Emit
	return q
}
