package rsp

import (
	"RSP/lib/tcp"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// Server inlcudes necessary worker and statistics of current state
type Server struct {
	port   int
	worker *Worker // worker to serve all requests

	// the following numbers are the statistics of current server status,
	// including # of connections, load, served requests, ...
	numConn       uint64 // number of connections
	numJobWaiting uint64 // number of jobs waiting
	numJobDone    uint64 // number of jobs that has been done
	running       bool   // server is running
}

// Job defines rsp job
// which actually is a wrapper of coming reqeust and response writter
type Job struct {
	res *tcp.ResponseWriter
	req *tcp.Request
}

// Worker is used to serve incoming jobs
// by default there will be only one worker that serves all jobs
type Worker struct {
	server *Server       // point back to server that owns this worker
	rate   time.Duration // request per second limit
	jobQ   chan *Job     // job queue implemented via channel
}

// NewServer construct server instance and relative resources
func NewServer(port int, ratePerSec int, jobsCap int) *Server {
	s := &Server{
		port: port,
		worker: &Worker{
			rate: time.Second / time.Duration(ratePerSec),
			jobQ: make(chan *Job, jobsCap),
		},
	}
	s.worker.server = s
	return s
}

// IsRunning return the state of current server
// return true if server is running; false if server is not runnig.
func (s *Server) IsRunning() bool {
	return s.running
}

// NumConn return the number of current connections
func (s *Server) NumConn() uint64 {
	return s.numConn
}

// NumJobWaiting reurn the number of job that is waiting
func (s *Server) NumJobWaiting() uint64 {
	return s.numJobWaiting
}

// NumJobDone returns the job that has been completed
func (s *Server) NumJobDone() uint64 {
	return s.numJobDone
}

// Run will start another routine for async job processing
// then the server blocks in the main rountine waiting for new connection
// once there is new connection coming in, a separate routine will started to serve that connection.
// Parameter requestTimeout decides the amount of time (in seconds) before a request is
// regarded as timed out and discared.
func (s *Server) Run(requestTimeout int) {
	// process jobs asyncly
	go s.worker.processJobs()

	// start tcp server
	s.running = true
	tcp.ListenAndServe(fmt.Sprintf(":%d", s.port), s.worker, requestTimeout)
	s.running = false
}

// dispatch job to different worker according to different commands
func dispatchJob(command string) *tcp.Response {
	argv := strings.Split(command, " ")
	res := &tcp.Response{}

	switch argv[0] {
	case "shorten":
		res.Result, res.Status = CreateShortURL(argv[1])
	case "wait":
		seconds, _ := strconv.Atoi(argv[1])
		time.Sleep(time.Second * time.Duration(seconds))
		res.Result, res.Status = "OK", 200
	default:
		res.Result, res.Status = "Unknown command", 404
	}
	return res
}

func (h *Worker) processJobs() {
	throttle := time.Tick(h.rate)
	for req := range h.jobQ {
		<-throttle
		go h.launchJob(req)

		atomic.StoreUint64(&h.server.numJobWaiting, uint64(len(h.jobQ)))
		atomic.AddUint64(&h.server.numJobDone, 1)
	}
}

func (h *Worker) launchJob(job *Job) {
	// before dispatching the job, check if the job had been canceled
	select {
	case <-job.req.Done():
		log.Printf("canceled by caller %s", job.req.Command)
		return
	default:
		res := dispatchJob(job.req.Command)
		job.res.Write(res)
	}
}

// ProcessRequest implement interface tcp.TcpHandler.ProcessRequest
func (h *Worker) ProcessRequest(res *tcp.ResponseWriter, req *tcp.Request) {
	// push new rsp job to job queue
	select {
	case h.jobQ <- &Job{res, req}:
		atomic.StoreUint64(&h.server.numJobWaiting, uint64(len(h.jobQ)))
	default:
		log.Println("job queue is full")
	}
}

// ConnArrived acts as a receiver to get notified that a new connection arrived
func (h *Worker) ConnArrived() {
	atomic.AddUint64(&h.server.numConn, 1)
}

// ConnClosed acts as a receiver to get notified that a connection had been closed
func (h *Worker) ConnClosed() {
	atomic.AddUint64(&h.server.numConn, ^uint64(0))
}
