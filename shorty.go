package main

import (
	"bytes"
	"github.com/jpgneves/shorty/requests"
	"github.com/jpgneves/shorty/routers"
	"html/template"
	"net/http"
	"strconv"
)

type ShortyResource struct{
	cache map[string]string
	rev_cache map[string]string
	counter uint64
}

func (r *ShortyResource) Get(request *requests.Request) *requests.Response {
	if redirect, ok := r.cache[request.Params["id"]]; ok {
		return &requests.Response{http.StatusTemporaryRedirect, &redirect}
	}
	return &requests.Response{http.StatusNotFound, nil}
}

func (r *ShortyResource) Post(request *requests.Request) *requests.Response {
	if url, ok := request.Params["url"]; ok {
		if cached, ok := r.rev_cache[url]; ok {
			return &requests.Response{http.StatusOK, &cached}
		}
		short := r.shorten(url)
		r.rev_cache[url] = short
		r.cache[short] = url
		return &requests.Response{http.StatusOK, &short}
	}
	return &requests.Response{http.StatusBadRequest, nil}
}

func (r *ShortyResource) shorten(url string) string {
	r.counter += 1
	return strconv.FormatUint(r.counter, 36)
}

type SiteResource struct{}

func (r *SiteResource) Get(request *requests.Request) *requests.Response {
	t, _ := template.ParseFiles("templates/index.tmpl")
	buf := new(bytes.Buffer)
	t.Execute(buf, request)
	str := buf.String()
	return &requests.Response{http.StatusOK, &str}
}

func (r *SiteResource) Post(request *requests.Request) *requests.Response {
	return &requests.Response{http.StatusMethodNotAllowed, nil}
}

func main() {
	router := routers.NewMatchingRouter()
	router.AddRoute("/", new(SiteResource))
	shorty := &ShortyResource{make(map[string]string), make(map[string]string), 13370}
	router.AddRoute("/{id}", shorty)
	router.AddRoute("/create/{url}", shorty)
	rh := routers.MakeRoutingHandler(router)
	http.ListenAndServe(":55384", rh)
}
