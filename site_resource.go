package main

import (
	"bytes"
	"github.com/jpgneves/microbe/config"
	"github.com/jpgneves/microbe/resources"
	"github.com/jpgneves/microbe/requests"
	"html/template"
	"log"
	"net/http"
)

type SiteResource struct{
	config *config.Configuration
}

func (r *SiteResource) Init(config *config.Configuration) resources.Resource {
	r.config = config
	return r
}

func (r *SiteResource) Get(request *requests.Request) *requests.Response {
	filepath := *(r.config.WwwRoot) + request.RawRequest.URL.Path
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