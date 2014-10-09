package engine

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
)

const (
	statusHtml = `<!DOCTYPE HTML>
<title>%d %s</title>
<h1>%s</h1>
<p>%s</p>
`
	panicBlock = `<h1>%s</h1>
<pre style="font-weight: bold;">%s</pre>
`
	panicHtml = `<html>
<head><title>Internal Server Error</title>
<style type="text/css">
html, body {
font-family: "Roboto", sans-serif;
color: #333333;
margin: 0px;
}
h1 {
color: #2b3848;
background-color: #ffffff;
padding: 20px;
border-bottom: 1px dashed #2b3848;
}
pre {
font-size: 1.1em;
margin: 20px;
padding: 20px;
border: 2px solid #2b3848;
background-color: #ffffff;
}
pre p:nth-child(odd){margin:0;}
pre p:nth-child(even){background-color: rgba(216,216,216,0.25); margin: 0;}
</style>
</head>
<body>
%s
</body>
</html>
`
)

type (
	// Status code, message, and handlers for a http status.
	HttpStatus struct {
		code     int
		message  string
		handlers []Manage
	}

	// A map of HttpStatus instances, keyed by status code
	HttpStatuses map[int]*HttpStatus
)

func NewHttpStatus(code int, message string) *HttpStatus {
	n := &HttpStatus{code: code, message: message}
	n.handlers = append(n.handlers, n.before(), n.after())
	return n
}

func (h *HttpStatus) name() string {
	return http.StatusText(h.code)
}

func (h *HttpStatus) before() Manage {
	return func(c *Ctx) {
		c.RW.WriteHeader(h.code)
	}
}

func (h *HttpStatus) after() Manage {
	return func(c *Ctx) {
		if !c.RW.Written() {
			if c.engine.HTMLStatus {
				c.RW.Header().Set("Content-Type", "text/html")
				c.RW.Write(h.format())
			} else {
				c.RW.WriteHeaderNow()
			}
		}
	}
}

func (h *HttpStatus) format() []byte {
	return []byte(fmt.Sprintf(statusHtml, h.code, h.name(), h.name(), h.message))
}

// Adds any number of custom Manage to the HttpStatus, between the
// default status before & after handlers.
func (h *HttpStatus) Update(handlers ...Manage) {
	s := len(h.handlers) + len(handlers)
	newh := make([]Manage, 0, s)
	newh = append(newh, h.handlers[0])
	if len(h.handlers) > 2 {
		newh = append(newh, h.handlers[1:(len(h.handlers)-2)]...)
	}
	newh = append(newh, handlers...)
	newh = append(newh, h.handlers[len(h.handlers)-1])
	h.handlers = newh
}

func defaultHttpStatuses() HttpStatuses {
	hss := make(HttpStatuses)
	hss.New(NewHttpStatus(400, "The browser (or proxy) sent a request that this server could not understand."))
	hss.New(NewHttpStatus(401, "The server could not verify that you are authorized to access the URL requested.\nYou either supplied the wrong credentials (e.g. a bad password), or your browser doesn't understand how to supply the credentials required."))
	hss.New(NewHttpStatus(403, "You do not have the permission to access the requested resource.\nIt is either read-protected or not readable by the server."))
	hss.New(NewHttpStatus(404, "The requested URL was not found on the server. If you entered the URL manually please check your spelling and try again."))
	hss.New(NewHttpStatus(405, "The method is not allowed for the requested URL."))
	hss.New(NewHttpStatus(418, "I'M A TEAPOT, NOT A COFFEE MACHINE."))
	hss.New(NewHttpStatus(500, "The server encountered an internal error and was unable to complete your request. Either the server is overloaded or there is an error in the application."))
	hss[500].Update(PanicHandle)
	hss.New(NewHttpStatus(502, "The proxy server received an invalid response from an upstream server."))
	hss.New(NewHttpStatus(503, "The server is temporarily unable to service your request due to maintenance downtime or capacity problems. Please try again later."))
	hss.New(NewHttpStatus(504, "The connection to an upstream server timed out."))
	hss.New(NewHttpStatus(505, "The server does not support the HTTP protocol version used in the request."))
	return hss
}

// PanicHandler is called with Http status 500 that gets all ErrorTypePanic from
// *Ctx.Errors, logs to logger if LoggingOn is true(general logging otherwise,
// you need to be informed of panics), and and serves a basic html page if engine
// ServePanic is true.
func PanicHandle(c *Ctx) {
	panics := c.Errors.ByType(ErrorTypePanic)
	var auffer bytes.Buffer
	for _, p := range panics {
		sig := fmt.Sprintf("encountered an internal error: %s\n-----\n%s\n-----\n", p.Err, p.Meta)
		go c.engine.SendSignal(sig) // this will hang in another package (e.g. Flotilla) without making it go
		if !c.engine.LoggingOn {
			log.Printf("[ENGINE]\n %s", sig)
		}
		reader := bufio.NewReader(bytes.NewReader([]byte(fmt.Sprintf("%s", p.Meta))))
		var err error
		lineno := 0
		var buffer bytes.Buffer
		for err == nil {
			lineno++
			l, _, err := reader.ReadLine()
			if lineno%2 == 0 {
				buffer.WriteString(fmt.Sprintf("\n%s</p>\n", l))
			} else {
				buffer.WriteString(fmt.Sprintf("<p>%s\n", l))
			}
			if err != nil {
				break
			}
		}
		pb := fmt.Sprintf(panicBlock, p.Err, buffer.String())
		auffer.WriteString(pb)
	}
	if c.engine.ServePanic {
		servePanic := fmt.Sprintf(panicHtml, auffer.String())
		c.RW.Header().Set("Content-Type", "text/html")
		c.RW.Write([]byte(servePanic))
	}
}

// New creates a HttpStatus in the HttpStatuss map.
func (hs HttpStatuses) New(h *HttpStatus) {
	hs[h.code] = h
}
