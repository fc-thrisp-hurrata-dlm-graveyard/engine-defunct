package engine

import (
	"bufio"
	"bytes"
	"fmt"
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
	// Status code, message, and Manage handlers for a http status.
	HttpStatus struct {
		Code     int
		Message  string
		Handlers []Manage
	}

	// A map of HttpStatus instances, keyed by status code
	HttpStatuses map[int]*HttpStatus
)

// Create new HttpStatus with the code, message, and default Manage handlers.
func NewHttpStatus(code int, message string) *HttpStatus {
	n := &HttpStatus{Code: code, Message: message}
	n.Handlers = append(n.Handlers, n.before(), n.after())
	return n
}

func (h *HttpStatus) name() string {
	return http.StatusText(h.Code)
}

func (h *HttpStatus) before() Manage {
	return func(c *Ctx) {
		c.RW.WriteHeader(h.Code)
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
	return []byte(fmt.Sprintf(statusHtml, h.Code, h.name(), h.name(), h.Message))
}

// Adds any number of custom Manage to the HttpStatus, between the
// default status before & after manage.
func (h *HttpStatus) Update(handlers ...Manage) {
	s := len(h.Handlers) + len(handlers)
	newh := make([]Manage, 0, s)
	newh = append(newh, h.Handlers[0])
	if len(h.Handlers) > 2 {
		newh = append(newh, h.Handlers[1:(len(h.Handlers)-2)]...)
	}
	newh = append(newh, handlers...)
	newh = append(newh, h.Handlers[len(h.Handlers)-1])
	h.Handlers = newh
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

// PanicHandle is the default Manage for 500 & internal panics. Retrieves all
// ErrorTypePanic from *Ctx.Errors, sends signal, logs to stdout or logger, and
// serves a basic html page if engine.ServePanic is true.
func PanicHandle(c *Ctx) {
	panics := c.Errors.ByType(ErrorTypePanic)
	var auffer bytes.Buffer
	for _, p := range panics {
		sig := fmt.Sprintf("encountered an internal error: %s\n-----\n%s\n-----\n", p.Err, p.Meta)
		c.engine.Send("panics-now", sig)
		if c.engine.ServePanic {
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
	}
	if c.engine.ServePanic {
		servePanic := fmt.Sprintf(panicHtml, auffer.String())
		c.RW.Header().Set("Content-Type", "text/html")
		c.RW.Write([]byte(servePanic))
	}
}

// New adds a new HttpStatus to HttpStatuses keyed by status code.
func (hs HttpStatuses) New(h *HttpStatus) {
	hs[h.Code] = h
}
