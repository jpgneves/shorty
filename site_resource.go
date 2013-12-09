package main

import (
	"bytes"
	"github.com/jpgneves/microbe/requests"
	"github.com/jpgneves/microbe/resources"
	"html/template"
	"log"
	"net/http"
)

type SiteResource struct{
	config *Configuration
}

func NewSiteResource(config *Configuration) resources.Resource {
	return &SiteResource{config}
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
