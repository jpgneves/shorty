package routers

import (
	"fmt"
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

func (r StaticRouter) Route(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.Method)
	method := request.Method
	url := request.URL.String()
	if resource, ok := r.routes[url]; ok {
		switch method {
		case "GET":
			fmt.Fprintf(writer, resource.Get(url))
		case "POST":
			fmt.Fprintf(writer, resource.Post(url, request.Body))
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(writer, "Error 405 not allowed")
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "Error 404 not found")
	}
}