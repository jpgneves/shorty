package routers

import (
	"github.com/jpgneves/shorty/resources"
	"net/http"
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

func (r StaticRouter) Route(request *http.Request) *RouteMatch {
	url := request.URL.Path
	if resource, ok := r.routes[url]; ok {
		return &RouteMatch{resource, nil}
	}
	return nil
}

func NewStaticRouter() Router {
	return StaticRouter{make(map[string]resources.Resource)}
}
