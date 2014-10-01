package engine

import (
	"net/http"
	"sync"

	"github.com/thrisp/engine/router"
)

type (
	// Engine is the the core struct containing Groups, sync.Pool cache, and
	// router, in addition to configuration options.
	Engine struct {
		*Group
		cache  sync.Pool
		router *router.Router
	}
)

// Empty returns an empty Engine with no Router, for you to build up from.
func Empty() *Engine {
	return &Engine{}
}

// Returns a new engine, with the least configuration.
func New() *Engine {
	engine := Empty()
	engine.router = router.New()
	engine.Group = NewGroup("/", engine)
	engine.cache.New = engine.newContext
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
