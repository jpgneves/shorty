package main

import (
	"bytes"
	"github.com/jpgneves/shorty/requests"
	"github.com/jpgneves/shorty/routers"
	"html/template"
	"net/http"
)

type ShortyResource struct{}

func (r ShortyResource) Get(request *requests.Request) *requests.Response {
	return &requests.Response{http.StatusTemporaryRedirect, "http://www.google.com"}
}

func (r ShortyResource) Post(request *requests.Request) *requests.Response {
	return &requests.Response{http.StatusOK, "shorty post"}
}

type SiteResource struct{}

func (r SiteResource) Get(request *requests.Request) *requests.Response {
	t, _ := template.ParseFiles("templates/index.tmpl")
	buf := new(bytes.Buffer)
	t.Execute(buf, request)
	return &requests.Response{http.StatusOK, buf.String()}
}

func (r SiteResource) Post(request *requests.Request) *requests.Response {
	return &requests.Response{http.StatusMethodNotAllowed, ""}
}

func main() {
	router := routers.NewMatchingRouter()
	router.AddRoute("/", new(SiteResource))
	router.AddRoute("/{id}", new(ShortyResource))
	rh := routers.MakeRoutingHandler(router)
	http.ListenAndServe(":55384", rh)
}
