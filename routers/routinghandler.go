package routers

import (
	"net/http"
)

type RoutingHandler struct {
	router Router
}

func MakeRoutingHandler() *RoutingHandler {
	return &RoutingHandler{}
}

func (rh *RoutingHandler) SetRouter(r Router) {
	rh.router = r
}

func (rh *RoutingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rh.router != nil {
		rh.router.Route(w, r)
	} else {
		return
	}
}