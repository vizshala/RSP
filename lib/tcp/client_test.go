package tcp

import (
	"fmt"
	"testing"
)

func TestEchoClient(t *testing.T) {
	port, err := GetRandomPort()
	if err != nil {
		t.Fatal("could not get random port")
	}

	addr := fmt.Sprintf(":%d", port)
	handler := &EchoServer{}
	go ListenAndServe(addr, handler, 5)

	client, err := NewClient(addr)
	if err != nil {
		t.Fatalf("failed to create new client")
	}

	expect := "200 olleH\n"
	client.Write("Hello")
	res, err := client.Read()
	if err != nil {
		t.Fatalf("failed to read from conn")
	}

	if res != expect {
		t.Error(
			"expected", expect,
			"got", res,
		)
	}
}
