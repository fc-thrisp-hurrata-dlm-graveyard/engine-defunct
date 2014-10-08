package engine

import (
	"errors"
	"net/http"
	"net/url"
)

type (
	Ctx struct {
		engine  *Engine
		group   *Group
		rwmem   responseWriter
		RW      ResponseWriter
		Request *http.Request
		Params  Params
		Form    url.Values
		// Files
		Errors errorMsgs
	}
)

func (engine *Engine) newContext() interface{} {
	c := &Ctx{engine: engine}
	c.RW = &c.rwmem
	return c
}

func (engine *Engine) getContext(w http.ResponseWriter, req *http.Request) *Ctx {
	c := engine.cache.Get().(*Ctx)
	c.rwmem.reset(w)
	c.Request = req
	req.ParseMultipartForm(engine.MaxFormMemory)
	c.Form = req.Form
	return c
}

func (engine *Engine) putContext(c *Ctx) {
	//if engine.SignalsOn {
	//	engine.SendSignal(fmt.Sprintf("%d %s %s", c.RW.Status(), c.Request.Method, c.Request.URL.Path))
	//}
	c.group = nil
	c.Request = nil
	c.Params = nil
	c.Form = nil
	engine.cache.Put(c)
}

func (c *Ctx) errorTyped(err error, typ uint32, meta interface{}) {
	c.Errors = append(c.Errors, errorMsg{
		Err:  err.Error(),
		Type: typ,
		Meta: meta,
	})
}

// Attaches an error that is pushed to a list of errors. It's a good idea
// to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database
// together, print a log, or append it in the HTTP response.
func (c *Ctx) Error(err error, meta interface{}) {
	c.errorTyped(err, ErrorTypeExternal, meta)
}

// Returns the last error for the Ctx.
func (c *Ctx) LastError() error {
	s := len(c.Errors)
	if s > 0 {
		return errors.New(c.Errors[s-1].Err)
	} else {
		return nil
	}
}

// Immediately abort the context writing out the code to the response
func (c *Ctx) Abort(code int) {
	if code >= 0 {
		c.RW.WriteHeader(code)
	}
}

// Fail is the same as Abort plus an error message.
// Calling `c.Fail(500, err)` is equivalent to:
// ```
// c.Error(err, "Failed.")
// c.Abort(500)
// ```
func (c *Ctx) Fail(code int, err error) {
	c.Error(err, err.Error())
	c.Abort(code)
}

// Calls an HttpException in the current group by integer code from the Context,
// if the status exists.
func (c *Ctx) Status(code int) {
	if status, ok := c.group.HttpStatuses[code]; ok {
		s := len(status.handlers)
		for i := 0; i < s; i++ {
			status.handlers[i](c)
		}
	}
}
