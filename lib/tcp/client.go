package tcp

import (
	"bufio"
	"fmt"
	"net"
)

// Client is the simple wrapper of net.Conn
type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

// Read from the connection
func (c *Client) Read() (string, error) {
	return c.reader.ReadString('\n')
}

// Write data to connection
func (c *Client) Write(data string) (int, error) {
	return fmt.Fprintf(c.conn, "%s\n", data)
	//s	return c.writer.WriteString(fmt.Sprintf(c.conn, "%s\n", data))
}

// Close the connection
func (c *Client) Close() {
	c.conn.Close()
}

// NewClient creates a new tcp connection, and return Client
func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}, nil
}
