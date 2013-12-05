package main

import (
	"bytes"
	"fmt"
	"github.com/jpgneves/microbe/requests"
	"github.com/jpgneves/microbe/routers"
	"github.com/jpgneves/shorty/storage"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type ShortyResource struct {
	cache     map[string]string
	rev_cache map[string]string
	counter   uint64
	lock      *sync.Mutex
	config    *Configuration
	db        *storage.DB
}

func NewShortyResource(config *Configuration) *ShortyResource {
	db, err := storage.OpenDB(*config.StorageConf.Backend, *config.StorageConf.Hostname)
	if err != nil {
		log.Fatal(err)
	}
	var counter uint64
	if c := db.Find("counter"); c == nil {
		counter = 13370
	} else {
		counter, err = strconv.ParseUint(c.(string), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	}
	cache := make(map[string]string)
	rev_cache := make(map[string]string)
	for kv := range db.Iterator() {
		str_v := kv.Value.(string)
		cache[kv.Key] = str_v
		rev_cache[str_v] = kv.Key
	}
	return &ShortyResource{cache, rev_cache, counter, new(sync.Mutex), config, &db}
}

func (r *ShortyResource) Get(request *requests.Request) *requests.Response {
	id := request.Params["id"]
	log.Println(id)
	if redirect, ok := r.cache[id]; ok {
		return &requests.Response{http.StatusTemporaryRedirect, &redirect}
	}
	db := *(r.db)
	url := db.Find(id)
	if url != nil {
		str := url.(string)
		return &requests.Response{http.StatusTemporaryRedirect, &str}
	}
	return &requests.Response{http.StatusNotFound, nil}
}

func (r *ShortyResource) Post(request *requests.Request) *requests.Response {
	url := request.RawRequest.FormValue("url")
	if url != "" {
		db := *(r.db)
		var shorturl string
		if cached, ok := r.rev_cache[url]; ok {
			shorturl = cached
		} else {
			r.lock.Lock()
			defer r.lock.Unlock()
			short := r.shorten(url)
			log.Printf("Caching %v as %v\n", url, short)
			r.rev_cache[url] = short
			r.cache[short] = url
			shorturl = short
			db.Insert(short, url)
			defer db.Flush()
		}
		if host, err := os.Hostname(); err == nil {
			hostport := net.JoinHostPort(host, strconv.Itoa(r.config.Port))
			shorturl = fmt.Sprintf("http://%v/%v", hostport, shorturl)
		} else {
			log.Fatal(err)
		}
		return &requests.Response{http.StatusOK, &shorturl}
	}
	return &requests.Response{http.StatusBadRequest, nil}
}

func (r *ShortyResource) shorten(url string) string {
	r.counter += 1
	db := *(r.db)
	db.Insert("counter", strconv.FormatUint(r.counter, 10))
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
	router.AddRoute("/", &SiteResource{config})
	shorty := NewShortyResource(config)
	router.AddRoute("/:id", shorty)
	router.AddRoute("/create", shorty)
	rh := routers.MakeRoutingHandler(router)
	addr := net.JoinHostPort(*config.ListenAddr, strconv.Itoa(config.Port))
	log.Printf("Starting server on %s", addr)
	http.ListenAndServe(addr, rh)
}
