package routers

import (
	"github.com/jpgneves/shorty/resources"
)

type StaticRouter struct {
	routes map[string]resources.Resource
}

func (r StaticRouter) AddRoute(route string, resource resources.Resource) {
	r.routes[route] = resource
}

func (r StaticRouter) RemoveRoute(route string) {
	delete(r.routes, route)
}

func (r StaticRouter) Route(path string) *RouteMatch {
	if resource, ok := r.routes[path]; ok {
		return &RouteMatch{resource, nil}
	}
	return nil
}

func NewStaticRouter() Router {
	return StaticRouter{make(map[string]resources.Resource)}
}
