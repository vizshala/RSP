package tcp

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

// reserver string helper
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// An Simple echo client for testing tcp server

// An simple echo server for testing tcp server
// implement tcp.Handler interface
type EchoServer struct {
	out string
	in  string
}

func (srv *EchoServer) ProcessRequest(res *ResponseWriter, req *Request) {
	res.Write(&Response{
		Reverse(req.Command),
		200,
	})
}

func (srv *EchoServer) ConnArrived() {

}
func (srv *EchoServer) ConnClosed() {

}
func TestEchoServer(t *testing.T) {
	port, err := GetRandomPort()
	if err != nil {
		t.Fatal("could not get random port")
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	handler := &EchoServer{}
	go ListenAndServe(addr, handler, 5)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("could not connect to %s, %s", addr, err)
	}

	expect := "200 olleH\n"
	fmt.Fprintf(conn, "Hello\n")

	reader := bufio.NewReader(conn)
	res, err := reader.ReadString('\n')
	if err != nil {
		t.Fatal("failed to read from conn")
	}

	if res != expect {
		t.Error(
			"expected", expect,
			"got", res,
		)
	}
}
