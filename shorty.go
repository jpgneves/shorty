package main

import (
	"github.com/jpgneves/shorty/requests"
	"github.com/jpgneves/shorty/routers"
	"net/http"
)

type ShortyResource struct{}

func (r ShortyResource) Get(request *requests.Request) *requests.Response {
	return &requests.Response{http.StatusOK, "shorty"}
}

func (r ShortyResource) Post(request *requests.Request) *requests.Response {
	return &requests.Response{http.StatusOK, "shorty post"}
}

func main() {
	router := routers.NewMatchingRouter()
	router.AddRoute("/", new(ShortyResource))
	router.AddRoute("/{id}", new(ShortyResource))
	rh := routers.MakeRoutingHandler(router)
	http.ListenAndServe(":55384", rh)
}
