package routers

import (
	"fmt"
	"github.com/jpgneves/shorty/resources"
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
		match := rh.router.Route(r)
		if resource, ok := match.value.(resources.Resource); ok {
			path := r.URL.Path
			switch r.Method {
			case "GET":
			fmt.Fprintf(w, resource.Get(path))
			case "POST":
				fmt.Fprintf(w, resource.Post(path, r.Body))
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	return
}
