package engine

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/thrisp/engine/router"
)

type (
	// Engine is the the core struct containing Groups, sync.Pool cache, router,
	// and a signal channel, in addition to configuration options.
	Engine struct {
		*Group
		cache   sync.Pool
		router  *router.Router
		signals signal
		logger  *log.Logger
		*conf
	}

	conf struct {
		MaxFormMemory int64
	}
)

func newconf() *conf {
	return &conf{
		MaxFormMemory: 1000000,
	}
}

// Empty returns an empty Engine with no Router, for you to build up from.
func Empty() *Engine {
	return &Engine{
		conf:    newconf(),
		signals: make(signal, 1),
		logger:  log.New(os.Stdout, "[Engine]", 0),
	}
}

// Returns a new engine, with the least configuration.
func New() *Engine {
	engine := Empty()
	engine.router = router.New()
	engine.Group = NewGroup("/", engine)
	engine.cache.New = engine.newContext
	go engine.ReadSignal()
	engine.SendSignal("new engine")
	return engine
}

// ServeHTTP makes the engine implement the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	engine.router.ServeHTTP(w, req)
}

func (engine *Engine) Run(addr string) {
	if err := http.ListenAndServe(addr, engine); err != nil {
		panic(err)
	}
}
