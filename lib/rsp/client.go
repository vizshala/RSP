package rsp

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/vizshala/RSP/lib/tcp"
	"os"
	"strings"
)

// Client is a warpper of tcp.Client
type Client struct {
	tcp *tcp.Client
}

// NewClient create a new rsp client
func NewClient(serverAddr string) (*Client, error) {

	tcpClient, err := tcp.NewClient(serverAddr)
	if err != nil {
		return nil, err
	}
	return &Client{tcpClient}, nil
}

// read user input subroutine
func (c *Client) readInput(chUsrInput chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				panic(err)
			}
		}
		chUsrInput <- scanner.Text()
	}
}

// read connection subroutine
func (c *Client) readConn(chRes chan string, chErr chan error) {
	for {
		data, err := c.tcp.Read()
		if err != nil {
			chErr <- err
			return
		}

		if len(data) == 0 {
			chErr <- errors.New("connection closed")
			return
		}

		chRes <- data
	}
}

// Close the connection
func (c *Client) Close() {
	c.tcp.Close()
}

// Run to serve the client connection
func (c *Client) Run() {
	// because the connection may be closed when blocked on stdin
	// use two more goroutines to allow detection of connection close
	chUsrInput := make(chan string)
	chData := make(chan string)
	chErr := make(chan error)

	go c.readConn(chData, chErr)
	go c.readInput(chUsrInput)

	for {
		fmt.Println("Enter command:")

		select {
		// wait for user input
		case command := <-chUsrInput:
			argv := strings.Split(command, " ")
			switch argv[0] {
			case "quit":
				fmt.Println("Exit now...")
				return
			default:
				c.tcp.Write(command)
				select {
				// timeout (let it be 2 seconds, just for testing)
				//case <-time.After(2 * time.Second):
				//	fmt.Println("timeout!")
				// got data from remote connection
				case data := <-chData:
					response := strings.Split(data[:len(data)-1], " ")
					fmt.Println(response)
				}
			}
		// detect connection close during stdin blocking
		case data := <-chData:
			response := strings.Split(data[:len(data)-1], " ")
			fmt.Println(response)
		// something was wrong with the connection
		case err := <-chErr:
			fmt.Println("connection closed", err)
			return
		}
	}
}
