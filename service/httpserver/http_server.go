package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vizshala/RSP/lib/memcache"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	numConn       uint64 // number of connections
	numJobWaiting uint64 // number of jobs waiting
	numJobDone    uint64 // number of jobs that has been done
)

// state
func state(res http.ResponseWriter, req *http.Request) {
	data := map[string]uint64{
		"current_conn": numConn,
		"job_waiting":  numJobWaiting,
		"job_done":     numJobDone,
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(data)
}

// index page
func index(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	t, err := template.ParseFiles("./template/index.html")
	if err != nil {
		log.Println(err)
		http.NotFound(res, req)
		return
	}

	data := struct {
		Title string
	}{
		Title: "RSP Server Statistics",
	}
	err = t.Execute(res, data)
	if err != nil {
		log.Println(err)
		http.NotFound(res, req)
		return
	}
}

// FileSystem custom file system handler
type FileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

// updateCache periodically
func updateCache(cacheAddr string) {
	// cache client
	cacheClient := memcache.NewClient(10)
	var wait = 0
	for {
		if cacheClient.IsConnected() {
			rspStat, err := cacheClient.GetRspState(memcache.RspServerID(0))
			if err != nil {
				log.Println("Get rsp state:", err)
			} else {
				numConn = rspStat.NumConn
				numJobDone = rspStat.NumJobDone
				numJobWaiting = rspStat.NumJobWaiting
			}
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
		serverPort int    // server port to listen to
		cacheAddr  string // cache server address
	)

	flag.IntVar(&serverPort, "port", 8080, "port to listen to")
	flag.StringVar(&cacheAddr, "cache", "localhost:1234", "cache server address")
	flag.Parse()

	http.HandleFunc("/", index)
	http.HandleFunc("/api/v1/state/", state)

	fileServer := http.FileServer(FileSystem{http.Dir("statics")})
	http.Handle("/statics/", http.StripPrefix("/statics", fileServer))

	go updateCache(cacheAddr)

	log.Println("http server started")
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil)
}
