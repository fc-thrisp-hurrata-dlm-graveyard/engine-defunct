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
		HTMLStatus    bool
		SignalsOn     bool
		MaxFormMemory int64
	}
)

func newconf() *conf {
	return &conf{
		HTMLStatus:    false,
		SignalsOn:     false,
		MaxFormMemory: 1000000,
	}
}

// Empty returns an empty Engine with default configuration.
func Empty() *Engine {
	return &Engine{conf: newconf()}
}

// Returns a new engine, with a method for retrieving a new Ctx, signals & logging.
func New() *Engine {
	engine := Empty()
	engine.router = router.New()
	engine.Group = NewGroup("/", engine)
	engine.cache.New = engine.newContext
	engine.SignalsOn = true
	engine.signals = engine.NewSignaller()
	engine.logger = log.New(os.Stdout, "[Engine]", 0)
	go engine.ReadSignal()
	return engine
}

// NotFound sets a http.HandlerFunc as default NotFound with the router.
func (engine *Engine) NotFound(h http.HandlerFunc) {
	engine.router.NotFound = h
}

// Panic sets a PHandle function as default panic handler with the router.
func (engine *Engine) Panic(h router.PHandle) {
	engine.router.PanicHandler = h
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
