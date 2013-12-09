package main

import (
	"github.com/jpgneves/microbe"
)

func main() {
	microbe := microbe.Init("./shorty.config")
	microbe.AddRoute("/", &SiteResource{config})
	shorty := NewShortyResource(config)
	microbe.AddRoute("/:id", shorty)
	microbe.AddRoute("/create", shorty)
	microbe.Start(false)
}
