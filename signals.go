package engine

import "log"

type (
	signal chan *Msg

	Msg struct {
		Head    string
		Content string
	}
)

func (e *Engine) NewSignaller() signal {
	s := make(signal, 1)
	return s
}

func (e *Engine) sendsignal(m *Msg) {
	e.signals <- m
}

func (e *Engine) Signal(m *Msg) {
	go e.sendsignal(m)
}

func DoSignal(head string, content string) *Msg {
	return &Msg{head, content}
}

func DoLog(content string) *Msg {
	return &Msg{"do-log", content}
}

func (e *Engine) LogSignal(watch string) {
	for msg := range e.signals {
		if msg.Head == watch {
			if e.Logger != nil {
				e.Logger.Printf(" %s", msg.Content)
			} else {
				log.Printf("[ENGINE] %s", msg.Content)
			}
		}
	}
}
