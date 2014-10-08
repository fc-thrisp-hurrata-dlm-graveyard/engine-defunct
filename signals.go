package engine

type (
	signal chan string
)

func (e *Engine) NewSignaller() signal {
	s := make(signal, 1)
	return s
}

func (e *Engine) SendSignal(msg string) {
	e.signals <- msg
}

func (e *Engine) LogSignal() {
	for msg := range e.signals {
		e.logger.Printf(" %s", msg)
	}
}
