package routers

import (
	"net/http"
	"github.com/jpgneves/shorty/resources"
)

type Router interface {
	AddRoute(route string, resource resources.Resource)
	RemoveRoute(route string)
	Route(writer http.ResponseWriter, req *http.Request)
}