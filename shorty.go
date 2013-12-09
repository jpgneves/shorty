package main

import (
	"github.com/jpgneves/microbe"
)

func main() {
	microbe := microbe.Init("./shorty.config")
	site_res := &SiteResource{}
	microbe.InitResource(site_res)
	microbe.AddRoute("/", &SiteResource{})
	shorty_res := &ShortyResource{}
	microbe.InitResource(shorty_res)
	microbe.AddRoute("/:id", shorty_res)
	microbe.AddRoute("/create", shorty_res)
	microbe.Start(false)
}
