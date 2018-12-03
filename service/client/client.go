package main

import (
	"flag"
	"github.com/vizshala/RSP/lib/rsp"
	"log"
)

var (
	serverAddr string
)

func main() {
	flag.StringVar(&serverAddr, "addr", "localhost:1314", "designate the address to connect to")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// create rsp client
	rspClient, err := rsp.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("error[%s]", err)
		return
	}
	defer rspClient.Close()

	rspClient.Run()
}
