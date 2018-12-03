package main

import (
	"flag"
	"github.com/vizshala/RSP/lib/memcache"
	"log"
)

var (
	serverAddr string
)

func main() {
	flag.StringVar(&serverAddr, "addr", "localhost:1234", "designate the address to connect to")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// start cache server
	log.Println("start cache server")
	memcache.LaunchServer(serverAddr)
}
