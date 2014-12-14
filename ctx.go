package engine

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context"
)

type (
	// Ctx is the core request-response context passed between any Manage
	// handlers, useful for storing & persisting data within a request & response.
	Ctx struct {
		engine  *Engine
		group   *Group
		rwmem   responseWriter
		RW      ResponseWriter
		request *http.Request
		Params  Params
		form    url.Values
		files   map[string][]*multipart.FileHeader
		Errors  errorMsgs
		*recorder
	}

	recorder struct {
		start     time.Time
		stop      time.Time
		latency   time.Duration
		status    int
		method    string
		path      string
		requester string
	}
)

func (engine *Engine) newCtx() interface{} {
	c := &Ctx{engine: engine}
	c.RW = &c.rwmem
	return c
}

func (engine *Engine) getCtx(w http.ResponseWriter, req *http.Request) *Ctx {
	c := engine.cache.Get().(*Ctx)
	c.group = engine.groups["/"]
	c.rwmem.reset(w)
	c.recorder = &recorder{}
	c.Start()
	c.request = req
	c.parseform()
	return c
}

func (engine *Engine) putCtx(c *Ctx) {
	c.PostProcess(c.request, c.RW)
	if engine.LoggingOn {
		engine.Send("message", c.LogFmt())
	}
	engine.Send("recorder", c.Fmt())
	c.group = nil
	c.request = nil
	c.Params = nil
	c.form = nil
	c.recorder = nil
	c.Errors = nil
	engine.cache.Put(c)
}

func (c *Ctx) parseform() {
	c.request.ParseMultipartForm(c.engine.MaxFormMemory)
	c.form = c.request.Form
	if c.request.MultipartForm != nil {
		c.files = c.request.MultipartForm.File
	}
}

func (c *Ctx) Request() *http.Request {
	return c.request
}

func (c *Ctx) Data() map[string]interface{} {
	ret := make(map[string]interface{})
	for _, p := range c.Params {
		ret[p.Key] = p.Value
	}
	return ret
}

func (c *Ctx) Form() url.Values {
	return c.form
}

func (c *Ctx) Files() map[string][]*multipart.FileHeader {
	return c.files
}

func (c *Ctx) Writer() ResponseWriter {
	return c.RW
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

func (c *Ctx) StatusFunc() (func(int), bool) {
	return c.Status, true
}

// Calls an HttpStatus in the current group by integer code from the Context,
// if the status exists.
func (c *Ctx) Status(code int) {
	if status, ok := c.group.HttpStatuses[code]; ok {
		s := len(status.Handlers)
		for i := 0; i < s; i++ {
			status.Handlers[i](context.WithValue(CurrentContext, "Current", c))
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
	return fmt.Sprintf("recorder	%s	%s	%s	%3d	%s	%s	%s", r.start, r.stop, r.latency, r.status, r.method, r.path, r.requester)
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
