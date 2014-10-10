package engine

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type (
	// Manage is a function that can be registered to a route to handle HTTP
	// requests. Like http.HandlerFunc, but takes a *Ctx
	Manage func(*Ctx)

	// Engine is the the core struct containing Groups, sync.Pool cache, and
	// signaling, in addition to configuration options.
	Engine struct {
		trees map[string]*node
		*Group
		cache   sync.Pool
		signals signal
		logger  *log.Logger
		*conf
	}

	conf struct {
		ServePanic            bool
		RedirectTrailingSlash bool
		RedirectFixedPath     bool
		HTMLStatus            bool
		LoggingOn             bool
		MaxFormMemory         int64
	}
)

func defaultconf() *conf {
	return &conf{
		ServePanic:            true,
		RedirectTrailingSlash: true,
		RedirectFixedPath:     true,
		HTMLStatus:            false,
		LoggingOn:             false,
		MaxFormMemory:         1000000,
	}
}

// Empty returns an empty Engine with zero configuration.
func Empty() *Engine {
	return &Engine{}
}

// New produces a new engine, with default configuration, a base group, method
// for retrieving a new Ctx, and signalling.
func New() *Engine {
	engine := Empty()
	engine.conf = defaultconf()
	engine.Group = NewGroup("/", engine)
	engine.cache.New = engine.newContext
	engine.signals = engine.NewSignaller()
	return engine
}

// Basic produces a new engine with LoggingOn set to true and visible logging.
func Basic() *Engine {
	engine := New()
	engine.LoggingOn = true
	engine.logger = log.New(os.Stdout, "[Engine]", 0)
	go engine.LogSignal()
	return engine
}

// Registers a new request Manage function with the given path and method.
func (e *Engine) Manage(method string, path string, m Manage) {
	if path[0] != '/' {
		panic("path must begin with '/'")
	}

	if e.trees == nil {
		e.trees = make(map[string]*node)
	}

	root := e.trees[method]
	if root == nil {
		root = new(node)
		e.trees[method] = root
	}

	root.addRoute(path, m)
}

// Handler allows the usage of a http.Handler as request manage.
func (e *Engine) Handler(method, path string, handler http.Handler) {
	e.Manage(method, path,
		func(c *Ctx) {
			handler.ServeHTTP(c.RW, c.Request)
		},
	)
}

// HandlerFunc allows the use of a http.HandlerFunc as request manage.
func (e *Engine) HandlerFunc(method, path string, handler http.HandlerFunc) {
	e.Manage(method, path,
		func(c *Ctx) {
			handler(c.RW, c.Request)
		},
	)
}

// Lookup allows the manual lookup of a method + path combo.
func (e *Engine) Lookup(method, path string) (Manage, Params, bool) {
	if root := e.trees[method]; root != nil {
		return root.getValue(path)
	}
	return nil, nil, false
}

// ServeFiles serves files from the given file system root. The path must end
// with "/*filepath", files are then served from the local path
// /defined/root/dir/*filepath.
//
// e.g., if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
//
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
//
// To use the operating system's file system implementation,
// use http.Dir:
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (e *Engine) ServeFiles(path string, root http.FileSystem) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath")
	}

	fileServer := http.FileServer(root)

	e.Manage("GET", path, func(c *Ctx) {
		c.Request.URL.Path = c.Params.ByName("filepath")
		fileServer.ServeHTTP(c.RW, c.Request)
	})
}

// internal "recover"
func (e *Engine) rcvr(c *Ctx) {
	if rcv := recover(); rcv != nil {
		p := newError(fmt.Sprintf("%s", rcv))
		c.errorTyped(p, ErrorTypePanic, stack(3))
		c.group = e.Group
		c.Status(500)
	}
}

// internal "not found"
func (e *Engine) ntfnd(c *Ctx) {
	c.group = e.Group
	c.Status(404)
}

// internal "servehttp"
func (engine *Engine) srvhttp(w http.ResponseWriter, req *http.Request, c *Ctx) {
	defer engine.rcvr(c)
	if root := engine.trees[req.Method]; root != nil {
		path := req.URL.Path
		if manage, ps, tsr := root.getValue(path); manage != nil {
			c.Params = ps
			manage(c)
			return
		} else if req.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && engine.RedirectTrailingSlash {
				if path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if engine.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					engine.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = string(fixedPath)
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	}

	engine.ntfnd(c)
	return
}

// ServeHTTP makes the engine implement the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.getContext(w, req)
	engine.srvhttp(w, req, c)
	engine.putContext(c)
}

func (engine *Engine) Run(addr string) {
	if err := http.ListenAndServe(addr, engine); err != nil {
		panic(err)
	}
}
