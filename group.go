package engine

import (
	"net/http"
	"path/filepath"
)

type (
	Group struct {
		engine   *Engine
		parent   *Group
		children []*RouterGroup
		HttpExceptions
	}
)

func (group *Group) pathFor(path string) string {
	joined := filepath.Join(group.prefix, path)
	// Append a '/' if the last component had one, but only if it's not there already
	if len(path) > 0 && path[len(path)-1] == '/' && joined[len(joined)-1] != '/' {
		return joined + "/"
	}
	return joined
}

func NewGroup(prefix string, engine *Engine) *Group {
	//engine.HttpExceptions = defaulthttpexceptions()
	return &Group{}
}

func (g *Group) New(component string) *Group {
	prefix := group.pathFor(component)
	newgroup := NewGroup(prefix, group.engine)
	newrgroup.parent = group
}

func (group *Group) Handle(route string, method string) {
	routepath = group.pathFor(route)
	group.engine.router.Handle(method, routepath, func(w http.ResponseWriter, req *http.Request, params router.Params) {
		c := group.engine.getCtx(w, req, params)
		c.currentgroup = group
		c.handler()
		group.engine.cache.Put(c)
	})
}
