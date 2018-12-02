package memcache

import (
	"net/rpc"
)

// CacheClient is the cache client, including rsp server, http server
type CacheClient struct {
	id        RspServerID
	rpcClient *rpc.Client
	retryWait int // when connection with rpc server is closed, wait retryWait seconds before next retry
}

// NewClient creats new cache client. If error occurred, CacheClient will be nil
func NewClient(retryWait int) *CacheClient {
	return &CacheClient{
		retryWait: retryWait,
	}
}

// IsConnected checks if current rpc client connection is ok
func (c *CacheClient) IsConnected() bool {
	return c.rpcClient != nil
}

// Connect tries create new connection with remote rpc server
func (c *CacheClient) Connect(serverAddr string) error {
	client, err := rpc.DialHTTP("tcp", serverAddr)
	if err != nil {
		return err
	}
	c.rpcClient = client
	return nil
}

// Close rpt client connection
func (c *CacheClient) Close() {
	c.rpcClient.Close()
	c.rpcClient = nil
}

// Register tries to get an unique id for this server
func (c *CacheClient) Register() (RspServerID, error) {
	var (
		rspID RspServerID
		dummy int
	)
	call := c.rpcClient.Go("RspServer.Register", &dummy, &rspID, nil)
	reply := <-call.Done
	if reply.Error != nil {
		c.Close()
		return 0, reply.Error
	}

	// store id
	c.id = rspID

	return rspID, nil
}

// GetRspState tries to get rsp server state
func (c *CacheClient) GetRspState(rspID RspServerID) (*RspServerStat, error) {
	rspStat := new(RspServerStat)
	call := c.rpcClient.Go("RspServer.GetState", &rspID, rspStat, nil)
	reply := <-call.Done
	if reply.Error != nil {
		c.Close()
		return nil, reply.Error
	}
	return rspStat, nil
}

// UpdateRspState tries to update rsp server state
func (c *CacheClient) UpdateRspState(stat *RspServerStat) error {
	args := stat

	// set id
	args.ID = c.id

	var ret int
	call := c.rpcClient.Go("RspServer.UpdateState", args, &ret, nil)
	reply := <-call.Done
	if reply.Error != nil {
		c.Close()
		return reply.Error
	}
	return nil
}
