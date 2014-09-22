package engine

import "net/http"

type (
	HandlerFunc func(*Ctx)

	Ctx struct {
		engine       *Engine
		currentgroup *Group
		handler      HandlerFunc
		rwmem        responseWriter
		rw           ResponseWriter
		Request      *http.Request
		Errors       errorMsgs
	}
)

func (engine *Engine) newContext() interface{} {
	c := &Ctx{engine: engine, handler: engine.Handler}
	c.rw = &c.rwmem
	return c
}

func (engine *Engine) getContext(w http.ResponseWriter, req *http.Request, params httprouter.Params) *Ctx {
	c := engine.cache.Get().(*Ctx)
	c.rwmem.reset(w)
	c.Request = req
	// set params
	return c
}
