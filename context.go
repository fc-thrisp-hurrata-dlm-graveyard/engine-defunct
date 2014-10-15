package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type (
	// Ctx is the core request-response context passed between any Manage
	// handlers, useful for storing & persisting data within a request & response.
	Ctx struct {
		engine  *Engine
		group   *Group
		rwmem   responseWriter
		RW      ResponseWriter
		Request *http.Request
		Params  Params
		Form    url.Values
		// Files tbd
		Errors errorMsgs
		*recorder
	}

	recorder struct {
		start     time.Time     `json:"start,string"`
		stop      time.Time     `json:"stop,string"`
		latency   time.Duration `json:"latency,omitempty"`
		status    int           `json:"status,string"`
		method    string        `json:"method"`
		path      string        `json:"path,omitempty"`
		requester string        `json:"requester,omitempty"`
		Extra     string        `json:"extra,omitempty,string"`
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
	c.recorder = &recorder{}
	c.Start()
	c.Request = req
	req.ParseMultipartForm(engine.MaxFormMemory)
	c.Form = req.Form
	return c
}

func (engine *Engine) putContext(c *Ctx) {
	c.PostProcess(c.Request, c.RW)
	if engine.LoggingOn {
		go engine.DoLog(c.LogFmt())
	} else {
		engine.SendSignal("recorder", c.Fmt())
	}
	c.group = nil
	c.Request = nil
	c.Params = nil
	c.Form = nil
	c.recorder = nil
	engine.cache.Put(c)
}

func (c *Ctx) errorTyped(err error, typ uint32, meta interface{}) {
	c.Errors = append(c.Errors, errorMsg{
		Err:  err.Error(),
		Type: typ,
		Meta: meta,
	})
}

// Attaches an error to a list of errors. Call Error for each error that occurred
// during the resolution of a request.
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

// Immediately abort the context, writing out the code to the response
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

// Calls an HttpStatus in the current group by integer code from the Context,
// if the status exists.
func (c *Ctx) Status(code int) {
	if status, ok := c.group.HttpStatuses[code]; ok {
		s := len(status.Handlers)
		for i := 0; i < s; i++ {
			status.Handlers[i](c)
		}
	}
}

func (r *recorder) Start() {
	r.start = time.Now()
}

func (r *recorder) Stop() {
	r.stop = time.Now()
}

func (r *recorder) Requester(req *http.Request) {
	rqstr := req.Header.Get("X-Real-IP")

	if len(rqstr) == 0 {
		rqstr = req.Header.Get("X-Forwarded-For")
	}

	if len(rqstr) == 0 {
		rqstr = req.RemoteAddr
	}

	r.requester = rqstr
}

func (r *recorder) Latency() time.Duration {
	return r.stop.Sub(r.start)
}

func (r *recorder) PostProcess(req *http.Request, rw ResponseWriter) {
	r.Stop()
	r.latency = r.Latency()
	r.Requester(req)
	r.method = req.Method
	r.path = req.URL.Path
	r.status = rw.Status()
}

func (r *recorder) Fmt() string {
	b, err := json.Marshal(r)
	if err != nil {
		return newError("recorder formatting error: %s", err).Error()
	}
	return string(b)
}

func (r *recorder) LogFmt() string {
	return fmt.Sprintf("%v |%s %3d %s| %12v | %s |%s %s %-7s %s",
		r.stop.Format("2006/01/02 - 15:04:05"),
		StatusColor(r.status), r.status, reset,
		r.latency,
		r.requester,
		MethodColor(r.method), reset, r.method,
		r.path)
}
