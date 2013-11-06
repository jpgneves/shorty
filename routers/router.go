package routers

import (
	"github.com/jpgneves/shorty/resources"
)

type Router interface {
	AddRoute(route string, resource resources.Resource)
	RemoveRoute(route string)
	Route(path string) *RouteMatch
}
