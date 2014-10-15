package engine

import (
	"log"
	"os"
	"reflect"
)

type (
	// A configuration function that takes an engine pointer, configures the
	// engine within the function, and returns an error.
	Conf func(*Engine) error

	conf struct {
		ServePanic            bool
		RedirectTrailingSlash bool
		RedirectFixedPath     bool
		HTMLStatus            bool
		LoggingOn             bool
		MaxFormMemory         int64
	}
)

func defaultconf() *conf {
	return &conf{
		ServePanic:            true,
		RedirectTrailingSlash: true,
		RedirectFixedPath:     true,
		HTMLStatus:            false,
		LoggingOn:             false,
		MaxFormMemory:         1000000,
	}
}

func (e *Engine) SetConf(opts ...Conf) error {
	for _, opt := range opts {
		if err := opt(e); err != nil {
			return err
		}
	}
	return nil
}

func ServePanic(b bool) Conf {
	return func(e *Engine) error {
		return e.SetConfBool("ServePanic", b)
	}
}

func RedirectTrailingSlash(b bool) Conf {
	return func(e *Engine) error {
		return e.SetConfBool("RedirectTrailingSlash", b)
	}
}

func RedirectFixedPath(b bool) Conf {
	return func(e *Engine) error {
		return e.SetConfBool("RedirectFixedPath", b)
	}
}

func HTMLStatus(b bool) Conf {
	return func(e *Engine) error {
		return e.SetConfBool("HTMLStatus", b)
	}
}

func Logger(l *log.Logger) Conf {
	return func(e *Engine) error {
		e.Logger = l
		LoggingOn(true)
		return nil
	}
}

func LoggingOn(b bool) Conf {
	return func(e *Engine) error {
		if b == true {
			go e.LogSignal("do-log")
		}
		e.Logger = log.New(os.Stdout, "[Engine]", 0)
		return e.SetConfBool("LoggingOn", b)
	}
}

func MaxFormMemory(byts int64) Conf {
	return func(e *Engine) error {
		return e.SetConfInt64("MaxFormMemory", byts)
	}
}

func (e *Engine) elem() reflect.Value {
	v := reflect.ValueOf(e)
	return v.Elem()
}

func (e *Engine) getfield(fieldname string) reflect.Value {
	return e.elem().FieldByName(fieldname)
}

func (e *Engine) SetConfInt64(fieldname string, as int64) error {
	f := e.getfield(fieldname)
	if f.CanSet() {
		f.SetInt(as)
		return nil
	}
	return newError("Engine could not set field %s as %d", fieldname, as)
}

func (e *Engine) SetConfBool(fieldname string, as bool) error {
	f := e.getfield(fieldname)
	if f.CanSet() {
		f.SetBool(as)
		return nil
	}
	return newError("Engine could not set field %s as %t", fieldname, as)
}
