package routers

import (
	"fmt"
	"github.com/jpgneves/shorty/resources"
	"net/http"
)

type MatchingRouter struct {
	trie* Trie
}

func NewMatchingRouter() Router {
	return MatchingRouter{trie: CreateTrie()}
}

func (r MatchingRouter) AddRoute(route string, resource resources.Resource) {
	r.trie.Insert(route, resource)
}

func (r MatchingRouter) RemoveRoute(route string) {

}

func (r MatchingRouter) Route(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	path := request.URL.Path
	match := r.trie.Find(path)
	if resource, ok := match.value.(resources.Resource); ok {
		switch method {
		case "GET":
			fmt.Fprintf(writer, resource.Get(path))
		case "POST":
			fmt.Fprintf(writer, resource.Post(path, request.Body))
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(writer, "Error 405 not allowed")
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "Error 404 not found")
	}
}