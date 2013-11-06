package main

import (
	"github.com/jpgneves/shorty/routers"
	"net/http"
)

type ShortyResource struct{}

func (r ShortyResource) Get(url string) string {
	return "shorty"
}

func (r ShortyResource) Post(url string, data interface{}) string {
	return "shorty post"
}

func main() {
	router := routers.NewMatchingRouter()
	router.AddRoute("/", new(ShortyResource))
	router.AddRoute("/{id}", new(ShortyResource))
	rh := routers.MakeRoutingHandler(router)
	http.ListenAndServe(":55384", rh)
}
