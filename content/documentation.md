+++
title = "Documentation"
+++
# engine
    import "github.com/thrisp/engine"




## Constants
``` go
const (
    ErrorTypeInternal = 1 << iota
    ErrorTypeExternal = 1 << iota
    ErrorTypePanic    = 1 << iota
    ErrorTypeAll      = 0xffffffff
)
```
``` go
const (
    NotWritten = -1
)
```


## func CleanPath
``` go
func CleanPath(p string) string
```
CleanPath is the URL version of path.Clean, it returns a canonical URL path
for p, eliminating . and .. elements.

The following rules are applied iteratively until no further processing can
be done:


	1. Replace multiple slashes with a single slash.
	2. Eliminate each . path name element (the current directory).
	3. Eliminate each inner .. path name element (the parent directory)
	   along with the non-.. element that precedes it.
	4. Eliminate .. elements that begin a rooted path:
	   that is, replace "/.." by "/" at the beginning of a path.

If the result of this process is an empty string, "/" is returned


## func PanicHandle
``` go
func PanicHandle(c *Ctx)
```
PanicHandle is the default Manage for 500 & internal panics. Retrieves all
ErrorTypePanic from \*Ctx.Errors, sends signal, logs to stdout or logger, and
serves a basic html page if engine.ServePanic is true.



## type Ctx
``` go
type Ctx struct {
    RW      ResponseWriter
    Request *http.Request
    Params  Params
    Form    url.Values
    // Files tbd
    Errors errorMsgs
    // contains filtered or unexported fields
}
```
Ctx is the core request-response context passed between any Manage
handlers, useful for storing & persisting data within a request & response.











### func (\*Ctx) Abort
``` go
func (c *Ctx) Abort(code int)
```
Immediately abort the context, writing out the code to the response



### func (\*Ctx) Error
``` go
func (c *Ctx) Error(err error, meta interface{})
```
Attaches an error to a list of errors. Call Error for each error that occurred
during the resolution of a request.



### func (\*Ctx) Fail
``` go
func (c *Ctx) Fail(code int, err error)
```
Fail is the same as Abort plus an error message.
Calling `c.Fail(500, err)` is equivalent to:
```
c.Error(err, "Failed.")
c.Abort(500)
```



### func (\*Ctx) LastError
``` go
func (c *Ctx) LastError() error
```
Returns the last error for the Ctx.



### func (\*Ctx) Status
``` go
func (c *Ctx) Status(code int)
```
Calls an HttpStatus in the current group by integer code from the Context,
if the status exists.



## type Engine
``` go
type Engine struct {
    *Group
    // contains filtered or unexported fields
}
```
Engine is the the core struct containing Groups, sync.Pool cache, and
signaling, in addition to configuration options.









### func Basic
``` go
func Basic() *Engine
```
Basic produces a new engine with LoggingOn set to true and visible logging.


### func Empty
``` go
func Empty() *Engine
```
Empty returns an empty Engine with zero configuration.


### func New
``` go
func New() *Engine
```
New produces a new engine, with default configuration, a base group, method
for retrieving a new Ctx, and signalling.




### func (\*Engine) Handler
``` go
func (e *Engine) Handler(method, path string, handler http.Handler)
```
Handler allows the usage of a http.Handler as request manage.



### func (\*Engine) HandlerFunc
``` go
func (e *Engine) HandlerFunc(method, path string, handler http.HandlerFunc)
```
HandlerFunc allows the use of a http.HandlerFunc as request manage.



### func (\*Engine) LogSignal
``` go
func (e *Engine) LogSignal()
```


### func (\*Engine) Lookup
``` go
func (e *Engine) Lookup(method, path string) (Manage, Params, bool)
```
Lookup allows the manual lookup of a method + path combo.



### func (\*Engine) Manage
``` go
func (e *Engine) Manage(method string, path string, m Manage)
```
Registers a new request Manage function with the given path and method.



### func (\*Engine) NewSignaller
``` go
func (e *Engine) NewSignaller() signal
```


### func (\*Engine) Run
``` go
func (engine *Engine) Run(addr string)
```


### func (\*Engine) SendSignal
``` go
func (e *Engine) SendSignal(msg string)
```


### func (\*Engine) ServeFiles
``` go
func (e *Engine) ServeFiles(path string, root http.FileSystem)
```
ServeFiles serves files from the given file system root. The path must end
with "/\*filepath", files are then served from the local path
/defined/root/dir/\*filepath.

e.g., if root is "/etc" and \*filepath is "passwd", the local file
"/etc/passwd" would be served.

Internally a http.FileServer is used, therefore http.NotFound is used instead
of the Router's NotFound handler.

To use the operating system's file system implementation,
use http.Dir:


	router.ServeFiles("/src/\*filepath", http.Dir("/var/www"))



### func (\*Engine) ServeHTTP
``` go
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)
```
ServeHTTP makes the engine implement the http.Handler interface.



## type EngineError
``` go
type EngineError struct {
    // contains filtered or unexported fields
}
```










### func (\*EngineError) Error
``` go
func (e *EngineError) Error() string
```


## type Group
``` go
type Group struct {
    HttpStatuses
    // contains filtered or unexported fields
}
```








### func NewGroup
``` go
func NewGroup(prefix string, engine *Engine) *Group
```
NewGroup creates a group with no parent and the provided prefix, usually as
a primary engine Group




### func (\*Group) Handle
``` go
func (group *Group) Handle(route string, method string, handler Manage)
```
Handle provides a route, method, and Manage to the router, and creates
a function using the handler when the router matches the route and method.



### func (\*Group) New
``` go
func (group *Group) New(component string) *Group
```
New creates a group from an existing group using the component string as a
prefix for all subsequent route attachments to the group. The existing group
will be the parent of the new group, and both will share the same engine



## type HttpStatus
``` go
type HttpStatus struct {
    Code     int
    Message  string
    Handlers []Manage
}
```
Status code, message, and Manage handlers for a http status.









### func NewHttpStatus
``` go
func NewHttpStatus(code int, message string) *HttpStatus
```
Create new HttpStatus with the code, message, and default Manage handlers.




### func (\*HttpStatus) Update
``` go
func (h *HttpStatus) Update(handlers ...Manage)
```
Adds any number of custom Manage to the HttpStatus, between the
default status before & after manage.



## type HttpStatuses
``` go
type HttpStatuses map[int]*HttpStatus
```
A map of HttpStatus instances, keyed by status code











### func (HttpStatuses) New
``` go
func (hs HttpStatuses) New(h *HttpStatus)
```
New adds a new HttpStatus to HttpStatuses keyed by status code.



## type Manage
``` go
type Manage func(*Ctx)
```
Manage is a function that can be registered to a route to handle HTTP
requests. Like http.HandlerFunc, but takes a \*Ctx











## type Param
``` go
type Param struct {
    Key   string
    Value string
}
```
Param is a single URL parameter, consisting of a key and a value.











## type Params
``` go
type Params []Param
```
Params is a Param-slice, as returned by the router. The slice is ordered,
the first URL parameter is also the first slice value. It is safe to read
values by the index.











### func (Params) ByName
``` go
func (ps Params) ByName(name string) string
```
ByName returns the value of the first Param which key matches the given name.
If no matching Param is found, an empty string is returned.



## type ResponseWriter
``` go
type ResponseWriter interface {
    http.ResponseWriter
    http.Hijacker
    http.Flusher
    http.CloseNotifier

    Status() int
    Size() int
    Written() bool
    WriteHeaderNow()
}
```
