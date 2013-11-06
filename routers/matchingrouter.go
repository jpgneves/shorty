package routers

import (
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
	r.trie.Remove(route)
}

func (r MatchingRouter) Route(request *http.Request) *RouteMatch {
	path := request.URL.Path
	return r.trie.Find(path)
}
