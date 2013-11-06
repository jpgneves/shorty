package routers

import (
	"github.com/jpgneves/shorty/resources"
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

func (r MatchingRouter) Route(path string) *RouteMatch {
	return r.trie.Find(path)
}
