package main

import (
	"github.com/jpgneves/microbe"
)

func main() {
	microbe := microbe.Init("./microbe.config")
	shorty_conf := ReadConfig("./shorty.config")
	site_res := NewSiteResource(shorty_conf)
	microbe.AddRoute("/", site_res)
	shorty_res := NewShortyResource(shorty_conf)
	microbe.AddRoute("/:id", shorty_res)
	microbe.AddRoute("/create", shorty_res)
	microbe.Start(false)
}
