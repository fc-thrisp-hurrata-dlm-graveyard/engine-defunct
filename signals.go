package engine

type (
	signal chan string
)

func (e *Engine) NewSignaller() signal {
	s := make(signal, 1)
	//defer close(s)
	return s
}

func (e *Engine) SendSignal(msg string) {
	e.signals <- msg
}

func (e *Engine) ReadSignal() {
	for msg := range e.signals {
		e.logger.Printf(" %-25s", msg)
	}
}
