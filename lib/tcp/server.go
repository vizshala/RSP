package tcp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// A Server defines parameters for running a TCP server.
type Server struct {
	Addr    string  // TCP address to listen on
	Handler Handler // handler to invoke

	requestTimeout int          // request timeout in seconds
	listener       net.Listener // listener for this server
	running        bool
}

// A Handler responds to an TCP request.
type Handler interface {
	ProcessRequest(res *ResponseWriter, req *Request)
	ConnArrived()
	ConnClosed()
}

// Response defines the data to be returned to client
type Response struct {
	Result string
	Status int
}

// Request defines the data to be sent to custom handler
type Request struct {
	Command string
	ctx     context.Context
}

// Done is used to test if the request had been canceled or timed out
//
func (req *Request) Done() <-chan struct{} {
	return req.ctx.Done()
}

// ResponseWriter is provided for application on top of TCP to write back response
type ResponseWriter struct {
	responseChannel chan *Response
}

// Write response
func (rw *ResponseWriter) Write(res *Response) {
	rw.responseChannel <- res
}

// helper function to get remote address as identification string for remote peer.
func getAddr(conn net.Conn) string {
	return conn.RemoteAddr().String()
}

// get a random port that is not currently used
func GetRandomPort() (int, error) {
	// bind port 0 will cause a random port to be picked by system
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

// serve is called as a separate goroutine to handle different connection
func serve(conn net.Conn, srv *Server) {
	handler := srv.Handler
	defer conn.Close()
	defer handler.ConnClosed()

	reader := bufio.NewReader(conn)

	chRes := make(chan *Response)
	writter := &ResponseWriter{
		responseChannel: chRes,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			log.Println(getAddr(conn), err)
			break
		}

		log.Println("command", command)

		// Provide a writer for delegator writing back response
		// The writer has a method Write that will direct the response to responseChannel
		// Therefore, we need to go another routine to avoid deadlock
		go handler.ProcessRequest(writter,
			&Request{
				Command: command[0 : len(command)-1],
				ctx:     ctx,
			})

		// The request is an external api call
		// If the job wrapping the request is not get handled within a specific time, cancel the request.
		select {
		case response := <-chRes:
			fmt.Fprintf(conn, "%d %s\n", response.Status, response.Result)
		// default timeout for request is defined by requestTimeout
		case <-time.After(time.Duration(srv.requestTimeout) * time.Second):
			cancel()
			response := &Response{Result: "timeout waiting", Status: 408}
			fmt.Fprintf(conn, "%d %s\n", response.Status, response.Result)
		}

	}
}

// start the loop for accepting new connectoin
func (srv *Server) start() error {
	if srv.listener == nil {
		return errors.New("listener not ready")
	}

	for {
		conn, err := srv.listener.Accept()

		if !srv.running {
			log.Printf("closing")
			break
		}

		if err != nil {
			log.Println("Error accepting connection")
			continue
		}

		// notify handler there is new connection
		srv.Handler.ConnArrived()

		go serve(conn, srv)
	}
	return nil
}

// Stop may be invoked by externally by server application
// or internally due to interruption or other reason
func (srv *Server) Stop() error {
	if !srv.running {
		return errors.New("closing is still in progress")
	}

	srv.running = false
	// trigger returning from listener.Accept()
	srv.listener.Close()
	return nil
}

// ListenAndServe waits for new connection and go another routine to serve this connection
func (srv *Server) ListenAndServe() error {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sig)

	// if Addr is empty, try random port
	if len(srv.Addr) == 0 {
		port, err := GetRandomPort()
		if err != nil {
			log.Printf("Failed to get a random port")
			return err
		}
		srv.Addr = fmt.Sprintf(":%d", port)
	}

	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Printf("Failed to open port on %s", srv.Addr)
		return err
	}

	srv.running = true
	srv.listener = listener
	log.Println("Listen on address", srv.Addr)

	go func() {
		s := <-sig
		log.Printf("signal %d, exit now...", s)
		srv.Stop()
	}()

	return srv.start()
}

// ListenAndServe is used to listen on specific port and serve requests
// it will block on current routine until server exits
func ListenAndServe(addr string, handler Handler, requestTimeout int) error {
	server := &Server{Addr: addr, Handler: handler, requestTimeout: requestTimeout}
	return server.ListenAndServe()
}
