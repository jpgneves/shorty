package routers

import (
	"net/http"
)

type RoutingHandler struct {
	router Router
}

func MakeRoutingHandler(router Router) *RoutingHandler {
	return &RoutingHandler{router}
}

func (rh *RoutingHandler) SetRouter(r Router) {
	rh.router = r
}

func (rh *RoutingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rh.router != nil {
		rh.router.Route(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}