package engine

import (
	"path/filepath"

	"golang.org/x/net/context"
)

type (
	groups map[string]*Group

	Group struct {
		prefix string
		parent *Group
		engine *Engine
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

// NewGroup creates a group with no parent and the provided prefix.
func NewGroup(prefix string, engine *Engine) *Group {
	if group, exists := engine.groups[prefix]; exists {
		return group
	} else {
		newgroup := &Group{prefix: prefix,
			engine:       engine,
			HttpStatuses: defaultHttpStatuses()}
		engine.groups[prefix] = newgroup
		return newgroup
	}
}

// New creates a group from an existing group using the the groups prefix and
// the provided component string as a prefix. The existing group will be the
// parent of the new group.
func (group *Group) New(component string) *Group {
	prefix := group.pathFor(component)
	newgroup := NewGroup(prefix, group.engine)
	newgroup.parent = group
	return newgroup
}

// Handle provides a route, method, and Manage to the router, and creates
// a function using the handler when the router matches the route and method.
func (group *Group) Take(route string, method string, handler func(context.Context)) {
	group.engine.Manage(method, group.pathFor(route), func(c context.Context) {
		curr := currentCtx(c)
		curr.group = group
		handler(context.WithValue(c, "Current", curr))
	})
}

func (group *Group) TakeStatus(code int, statushandler func(context.Context)) {
	if ss, ok := group.HttpStatuses[code]; ok {
		ss.Update(statushandler)
	} else {
		ns := NewHttpStatus(code, string(code))
		ns.Update(statushandler)
		group.HttpStatuses.New(ns)
	}
}
