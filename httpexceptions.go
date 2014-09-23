package engine

import (
	"fmt"
	"net/http"
)

const (
	exceptionHtml = `<!DOCTYPE HTML>
<title>%d %s</title>
<h1>%s</h1>
<p>%s</p>
`
)

type (
	// Status code, message, and handlers for a http exception.
	HttpException struct {
		code     int
		message  string
		handlers []HandlerFunc
	}

	// A map of HttpException instances, keyed by status code
	HttpExceptions map[int]*HttpException
)

func NewHttpException(code int, message string) *HttpException {
	n := &HttpException{code: code, message: message}
	n.handlers = append(n.handlers, n.before(), n.after())
	return n
}

func (h *HttpException) name() string {
	return http.StatusText(h.code)
}

func (h *HttpException) before() HandlerFunc {
	return func(c *Ctx) {
		c.RW.WriteHeader(h.code)
	}
}

func (h *HttpException) after() HandlerFunc {
	return func(c *Ctx) {
		if !c.RW.Written() {
			if c.RW.Status() == h.code {
				c.RW.Header().Set("Content-Type", "text/html")
				c.RW.Write(h.format())
			} else {
				c.RW.WriteHeaderNow()
			}
		}
	}
}

func (h *HttpException) format() []byte {
	return []byte(fmt.Sprintf(exceptionHtml, h.code, h.name(), h.name(), h.message))
}

// Adds any number of custom HandlerFunc to the HttpException, between the
// default exception before & after handlers.
func (h *HttpException) Update(handlers ...HandlerFunc) {
	s := len(h.handlers) + len(handlers)
	newh := make([]HandlerFunc, 0, s)
	newh = append(newh, h.handlers[0])
	if len(h.handlers) > 2 {
		newh = append(newh, h.handlers[1:(len(h.handlers)-2)]...)
	}
	newh = append(newh, handlers...)
	newh = append(newh, h.handlers[len(h.handlers)-1])
	h.handlers = newh
}

func defaultHttpExceptions() HttpExceptions {
	httpexceptions := make(HttpExceptions)
	httpexceptions.New(NewHttpException(400, "The browser (or proxy) sent a request that this server could not understand."))
	httpexceptions.New(NewHttpException(401, "The server could not verify that you are authorized to access the URL requested.\nYou either supplied the wrong credentials (e.g. a bad password), or your browser doesn't understand how to supply the credentials required."))
	httpexceptions.New(NewHttpException(403, "You do not have the permission to access the requested resource.\nIt is either read-protected or not readable by the server."))
	httpexceptions.New(NewHttpException(404, "The requested URL was not found on the server. If you entered the URL manually please check your spelling and try again"))
	httpexceptions.New(NewHttpException(405, "The method is not allowed for the requested URL."))
	httpexceptions.New(NewHttpException(418, "This server is a teapot, not a coffee machine"))
	httpexceptions.New(NewHttpException(500, "The server encountered an internal error and was unable to complete your request. Either the server is overloaded or there is an error in the application."))
	httpexceptions.New(NewHttpException(502, "The proxy server received an invalid response from an upstream server."))
	httpexceptions.New(NewHttpException(503, "The server is temporarily unable to service your request due to maintenance downtime or capacity problems. Please try again later."))
	httpexceptions.New(NewHttpException(504, "The connection to an upstream server timed out."))
	httpexceptions.New(NewHttpException(505, "The server does not support the HTTP protocol version used in the request"))
	return httpexceptions
}

// New creates a HttpException in the HttpExceptions map.
func (hs HttpExceptions) New(h *HttpException) {
	hs[h.code] = h
}
