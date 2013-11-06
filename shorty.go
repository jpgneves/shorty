package main

import (
	"fmt"
	"net/http"
)

type TrieNode interface {
	IsPattern() bool
	PatternName() *string
	Value() interface{}
	Lookup(r rune) *TrieNode
	Find(s string) *Match
	Insert(key string, value interface{})
}

type Match struct {
	value interface{}
	matches map[string]string
}

type Trie struct {
	pattern_name *string
	value interface{}
	children map[rune]*Trie
}

func CreateTrie() *Trie {
	return &Trie{value: nil, children: make(map[rune]*Trie)}
}

func (t Trie) IsPattern() bool {
	return t.pattern_name != nil
}

func (t Trie) PatternName() *string {
	return t.pattern_name
}

func (t Trie) Lookup(r rune) *Trie {
	if child, ok := t.children[r]; ok {
		return child
	}
	return nil
}

func (t Trie) Find(s string) *Match {
	matched_patterns := make(map[string]string)
	var current_pattern *string
	var current_value []rune
	node := &t
	for _, r := range s {
		n := node.Lookup(r)
		if n != nil {
			node = n
			if current_pattern != nil {
				matched_patterns[*current_pattern] = string(current_value)
				current_pattern = nil
			}
		} else if current_pattern != nil {
			current_value = append(current_value, r)
		} else {
			for key := range node.children {
				child := node.children[key]
				if child.IsPattern() {
					current_pattern = child.PatternName()
					current_value = make([]rune, 10, 255)
					current_value = append(current_value, r)
					node = child
					break
				}
			}
			if current_pattern != nil {
				continue
			}
			return &Match{nil, nil}
		}
	}
	if current_pattern != nil {
		matched_patterns[*current_pattern] = string(current_value)
	}
	return &Match{node.value, matched_patterns}
}

func (t Trie) Insert(key string, value interface{}) {
	in_pattern := false
	var last_rune rune
	var pattern_name []rune
	node := &t
	for _, r := range key {
		if r == '{' && !in_pattern {
			in_pattern = true
			pattern_name = make([]rune, 15)
		} else if r != '}' && in_pattern {
			pattern_name = append(pattern_name, r)
		} else if r == '}' && in_pattern {
			in_pattern = false
			new_node := CreateTrie()
			pattern_name_str := string(pattern_name)
			new_node.pattern_name = &pattern_name_str
			node.children[last_rune] = new_node
			node = new_node
		} else {
			n := node.Lookup(r)
			if n != nil {
				node = n
			} else {
				new_node := CreateTrie()
				node.children[r] = new_node
				node = new_node
			}
			last_rune = r
		}
	}
	node.value = value
}

func (t Trie) Print() {
	for key, _ := range t.children {
		if t.children[key].IsPattern() {
			fmt.Printf("{%+v} - %T %s\n", *(t.children[key].PatternName()), t.children[key].value, t.children[key].children)
		} else {
			fmt.Printf("%#U - %T %s\n", key, t.children[key].value, t.children[key].children)
		}
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

type StaticRouter struct {
	routes map[string]Resource
}

func (r StaticRouter) AddRoute(route string, resource Resource) {
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
	path := request.URL.Path
	match := r.trie.Find(path)
	if resource, ok := match.value.(Resource); ok {
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

type FooResource struct {}

func (r FooResource) Get(url string) string {
	return "foo"
}

func (r FooResource) Post(url string, data interface{}) string {
	return "foo post"
}

func main() {
	//router := StaticRouter{routes: make(map[string]Resource)}
	router := MatchingRouter{trie: CreateTrie()}
	router.AddRoute("/", new(ShortyResource))
	router.AddRoute("/{id}", new(ShortyResource))
	router.AddRoute("/foo/{foo}/{bar}/baz", new(FooResource))
	router.trie.Print()
	rh := MakeRoutingHandler()
	rh.SetRouter(router)
	fmt.Println(rh)
	http.ListenAndServe(":55384", rh)
}