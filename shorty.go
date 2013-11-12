package main

import (
	"bytes"
	"fmt"
	"github.com/jpgneves/shorty/requests"
	"github.com/jpgneves/shorty/routers"
	"github.com/jpgneves/shorty/storage"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type ShortyResource struct {
	cache     map[string]string
	rev_cache map[string]string
	counter   uint64
	config    *Configuration
}

func (r *ShortyResource) Get(request *requests.Request) *requests.Response {
	log.Println(request.Params["id"])
	if redirect, ok := r.cache[request.Params["id"]]; ok {
		return &requests.Response{http.StatusTemporaryRedirect, &redirect}
	}
	return &requests.Response{http.StatusNotFound, nil}
}

func (r *ShortyResource) Post(request *requests.Request) *requests.Response {
	url := request.RawRequest.FormValue("url")
	if url != "" {
		if cached, ok := r.rev_cache[url]; ok {
			return &requests.Response{http.StatusOK, &cached}
		}
		short := r.shorten(url)
		log.Printf("Caching %v as %v\n", url, short)
		r.rev_cache[url] = short
		r.cache[short] = url
		shorturl := fmt.Sprintf("http://%v:%v/%v", *(r.config.Hostname), r.config.Port, short)
		return &requests.Response{http.StatusOK, &shorturl}
	}
	return &requests.Response{http.StatusBadRequest, nil}
}

func (r *ShortyResource) shorten(url string) string {
	r.counter += 1
	return strconv.FormatUint(r.counter, 36)
}

type SiteResource struct{
	config *Configuration
}

func (r *SiteResource) Get(request *requests.Request) *requests.Response {
	filepath := *(r.config.SiteRoot) + request.RawRequest.URL.Path
	if filepath[len(filepath) - 1] == '/' {
		filepath += "index.html"
	}
	log.Println(filepath)
	t, _ := template.ParseFiles(filepath)
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
	config := ReadConfig("./shorty.config")
	db, err := storage.OpenDB(*config.StorageConf.Backend, *config.StorageConf.Hostname)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Flush()
	router.AddRoute("/", &SiteResource{config})
	shorty := &ShortyResource{make(map[string]string), make(map[string]string), 13370, config}
	router.AddRoute("/:id", shorty)
	router.AddRoute("/create", shorty)
	rh := routers.MakeRoutingHandler(router)
	addr := fmt.Sprintf("%v:%v", *config.Hostname, config.Port)
	log.Printf("Starting server on %s", addr)
	http.ListenAndServe(addr, rh)
}
