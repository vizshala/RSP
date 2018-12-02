package memcache

// RspServer rpc type
type RspServer int

// RspServerID is an global uniquo number to identify separate rsp server
type RspServerID int

// RspServerStat is the data representing rsp server state
type RspServerStat struct {
	ID            RspServerID
	NumConn       uint64 // number of connections
	NumJobWaiting uint64 // number of jobs waiting
	NumJobDone    uint64 // number of jobs that has been done
}
