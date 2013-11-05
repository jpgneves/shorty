package main

import (
	"fmt"
	"net/http"
)

type Trie struct {
	value interface{}
	children map[rune]*Trie
}

func CreateTrie() *Trie {
	return &Trie{value: nil, children: make(map[rune]*Trie)}
}

func (t Trie) Lookup(r rune) *Trie {
	if child, ok := t.children[r]; ok {
		return child
	}
	return nil
}

func (t Trie) Find(s string) interface{} {
	node := &t
	for _, r := range s {
		n := node.Lookup(r)
		if n != nil {
			node = n
		} else {
			return nil
		}
	}
	return node.value
}

func (t Trie) Insert(key string, value interface{}) {
	node := &t
	for _, r := range key {
		n := node.Lookup(r)
		if n != nil {
			node = n
		} else {
			new_node := CreateTrie()
			node.children[r] = new_node
			node = new_node
		}
	}
	node.value = value
}

func (t Trie) Print() {
	for key, _ := range t.children {
		fmt.Printf("%#U - %T %s\n", key, t.children[key].value, t.children[key].children)
		t.children[key].Print()
	}
}

type Resource interface {
	Get(url string) string
	Post(url string, data interface{}) string
}

type Router interface {
	AddRoute(route string, resource Resource)
	RemoveRoute(route string)
	Route(writer http.ResponseWriter, req *http.Request)
}

type SimpleRouter struct {
	routes map[string]Resource
}

func (r SimpleRouter) AddRoute(route string, resource Resource) {
	r.routes[route] = resource
}

func (r SimpleRouter) RemoveRoute(route string) {
	delete(r.routes, route)
}

func (r SimpleRouter) Route(writer http.ResponseWriter, request *http.Request) {
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

type MatchingRouter struct {
	trie* Trie
}

func (r MatchingRouter) AddRoute(route string, resource Resource) {
	r.trie.Insert(route, resource)
}

func (r MatchingRouter) RemoveRoute(route string) {

}

func (r MatchingRouter) Route(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	url := request.URL.String()
	resource := r.trie.Find(url).(Resource)
	if resource != nil {
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

type ShortyResource struct {}

func (r ShortyResource) Get(url string) string {
	return "shorty"
}

func (r ShortyResource) Post(url string, data interface{}) string {
	fmt.Println(data)
	return "shorty post"
}

func main() {
	//router := SimpleRouter{routes: make(map[string]Resource)}
	router := MatchingRouter{trie: CreateTrie()}
	router.AddRoute("/", new(ShortyResource))
	router.trie.Print()
	rh := MakeRoutingHandler()
	rh.SetRouter(router)
	fmt.Println(rh)
	http.ListenAndServe(":55384", rh)
}