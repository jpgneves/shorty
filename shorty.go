package main

import (
	"fmt"
	"net/http"
	"github.com/jpgneves/shorty/routers"
)

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
	router := routers.NewMatchingRouter()
	router.AddRoute("/", new(ShortyResource))
	router.AddRoute("/{id}", new(ShortyResource))
	router.AddRoute("/foo/{foo}/{bar}/baz", new(FooResource))
	rh := routers.MakeRoutingHandler()
	rh.SetRouter(router)
	fmt.Println(rh)
	http.ListenAndServe(":55384", rh)
}