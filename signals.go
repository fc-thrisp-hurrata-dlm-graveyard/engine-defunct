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

func (e *Engine) SendSignal(head string, content string) {
	msg := &Msg{head, content}
	e.signals <- msg
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
