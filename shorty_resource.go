package main

import (
	"fmt"
	"github.com/jpgneves/microbe/requests"
	"github.com/jpgneves/shorty/storage"
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
	lock      *sync.RWMutex
	config    *Configuration
	db        *storage.DB
}

func NewShortyResource(config *Configuration) *ShortyResource {
	db, err := storage.OpenDB(*config.StorageConf.Backend, *config.StorageConf.Location)
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
	return &ShortyResource{cache, rev_cache, counter, new(sync.RWMutex), config, &db}
}

func (r *ShortyResource) Get(request *requests.Request) *requests.Response {
	id := request.Params["id"]
	log.Println(id)
	r.lock.RLock()
	defer r.lock.RUnlock()
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
		r.lock.RLock()
		if cached, ok := r.rev_cache[url]; ok {
			shorturl = cached
		} else {
			r.lock.RUnlock()
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
		resp := fmt.Sprintf("<a href=%v>%v</a>", shorturl, shorturl)
		return &requests.Response{http.StatusOK, &resp}
	}
	return &requests.Response{http.StatusBadRequest, nil}
}

func (r *ShortyResource) shorten(url string) string {
	r.counter += 1
	db := *(r.db)
	db.Insert("counter", strconv.FormatUint(r.counter, 10))
	return strconv.FormatUint(r.counter, 36)
}