package tcp

import (
	"fmt"
	"testing"
	"time"
)

func TestEchoClient(t *testing.T) {
	port, err := GetRandomPort()
	if err != nil {
		t.Fatal("could not get random port")
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	handler := &EchoServer{}
	go ListenAndServe(addr, handler, 5)

	// allow retrying because server may be not ready yet!
	retry := 5
	client, err := NewClient(addr)
	for err != nil && retry > 0 {
		fmt.Printf("%s", err)
		time.Sleep(time.Second)
		retry--
		client, err = NewClient(addr)
	}
	if err != nil {
		t.Fatalf("failed to connect to server")
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
