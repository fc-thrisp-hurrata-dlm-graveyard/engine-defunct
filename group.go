package engine

import (
	"net/http"
	"path/filepath"

	"github.com/thrisp/engine/router"
)

type (
	Group struct {
		prefix string
		engine *Engine
		parent *Group
		HttpStatuses
	}
)

func (group *Group) pathFor(path string) string {
	joined := filepath.ToSlash(filepath.Join(group.prefix, path))
	// Append a '/' if the last component had one, but only if it's not there already
	if len(path) > 0 && path[len(path)-1] == '/' && joined[len(joined)-1] != '/' {
		return joined + "/"
	}
	return joined
}

// NewGroup creates a group with no parent and the provided prefix, usually as
// a primary engine Group
func NewGroup(prefix string, engine *Engine) *Group {
	return &Group{prefix: prefix,
		engine:       engine,
		HttpStatuses: defaultHttpStatuses()}
}

// New creates a group from an existing group using the component string as a
// prefix for all subsequent route attachments to the group. The existing group
// will be the parent of the new group, and both will share the same engine
func (group *Group) New(component string) *Group {
	prefix := group.pathFor(component)
	newgroup := NewGroup(prefix, group.engine)
	newgroup.parent = group
	return newgroup
}

// Handle provides a route, method, and HandlerFunc to the router, and creates
// a function using the handler when the router matches the route and method.
func (group *Group) Handle(route string, method string, handler HandlerFunc) {
	route = group.pathFor(route)
	group.engine.router.Handle(method, route, func(w http.ResponseWriter, req *http.Request, params router.Params) {
		c := group.engine.getContext(w, req, params)
		c.group = group
		handler(c)
		group.engine.putContext(c)
	})
}
