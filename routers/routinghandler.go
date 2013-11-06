package routers

import (
	"fmt"
	"github.com/jpgneves/shorty/requests"
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
	var response *requests.Response
	if rh.router != nil {
		path := r.URL.Path
		match := rh.router.Route(path)
		if resource, ok := match.value.(resources.Resource); ok {
			req := &requests.Request{r, match.matches}
			switch r.Method {
			case "GET":
				response = resource.Get(req)
			case "POST":
				response = resource.Post(req)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(response.StatusCode)
			fmt.Fprintf(w, response.Data)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	return
}
