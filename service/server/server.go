package main

import (
	"flag"
	"github.com/vizshala/RSP/lib/memcache"
	"github.com/vizshala/RSP/lib/rsp"
	"log"
	"time"
)

func updateCache(cacheAddr string, rspServer *rsp.Server) {
	// cache client
	cacheClient := memcache.NewClient(10)
	var wait = 0
	for {
		if cacheClient.IsConnected() {
			// FIXME: the cache server can be down after connected
			// should be able to detect this situation and try to reconnect
			cacheClient.UpdateRspState(&memcache.RspServerStat{
				NumConn:       rspServer.NumConn(),
				NumJobDone:    rspServer.NumJobDone(),
				NumJobWaiting: rspServer.NumJobWaiting(),
			})
		} else {
			wait--
			if wait <= 0 {
				err := cacheClient.Connect(cacheAddr)
				if err != nil {
					// try to connect to cache server (the server may be down for some reason)
					wait = 10
					log.Printf("cache server is not reachable, retry after %d seconds", wait)
				} else {
					log.Println("connected to cache server")
					cacheClient.Register()
				}
			}
		}

		time.Sleep(time.Second)
	}
}

func main() {
	var (
		serverPort     int    // server port to listen to
		ratePerSec     int    // external api rate per second
		jobCapacity    int    // job queue capacity
		requestTimeout int    // job queue capacity
		cacheAddr      string // cache server address
	)

	flag.IntVar(&serverPort, "port", 1314, "designate port to listen to")
	flag.IntVar(&ratePerSec, "rps", 30, "external api rate per second")
	flag.IntVar(&jobCapacity, "job", 100, "job queue capacity")
	flag.IntVar(&requestTimeout, "req_timeout", 5, "reqeust timeout")
	flag.StringVar(&cacheAddr, "cache", "localhost:1234", "cache server address")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// create server
	rspServer := rsp.NewServer(serverPort, ratePerSec, jobCapacity)

	// start cache client
	go updateCache(cacheAddr, rspServer)

	rspServer.Run(requestTimeout)
}
