package engine

import (
	"net/http"
	"sync"

	"github.com/thrisp/engine/router"
)

type (
	Engine struct {
		*Group
		cache  sync.Pool
		router *router.Router
	}
)

// Returns a new engine, with the least configuration.
func New() *Engine {
	engine := &Engine{router: router.New()}
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
