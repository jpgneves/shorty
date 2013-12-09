package main

import (
	"github.com/jpgneves/microbe/routers"
	"log"
	"net"
	"net/http"
	"strconv"
)

func main() {
	router := routers.NewMatchingRouter()
	config := ReadConfig("./shorty.config")
	router.AddRoute("/", &SiteResource{config})
	shorty := NewShortyResource(config)
	router.AddRoute("/:id", shorty)
	router.AddRoute("/create", shorty)
	rh := routers.MakeRoutingHandler(router)
	addr := net.JoinHostPort(*config.ListenAddr.Address, strconv.Itoa(config.ListenAddr.Port))
	log.Printf("Starting server on %s", addr)
	http.ListenAndServe(addr, rh)
}
