package memcache

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync/atomic"
	"time"
)

// statMap is a global map to store all rsp server state
// map from RspServerID to RspServerStat
var statMap map[RspServerID]RspServerStat

// store (sum) all state number
var allStat RspServerStat

var uniqueID int32

// Register will register with this memory cache service and get an unique id
func (t *RspServer) Register(args *int, id *RspServerID) error {
	// TODO: should have a mechanism for checking whether the rsp server is already registered

	atomic.AddInt32(&uniqueID, 1)
	*id = RspServerID(uniqueID)
	return nil
}

// UpdateState will update rsp server state, typically called by rsp server
func (t *RspServer) UpdateState(args *RspServerStat, reply *int) error {
	old := statMap[args.ID]
	diff := RspServerStat{
		NumConn:       args.NumConn - old.NumConn,
		NumJobWaiting: args.NumJobWaiting - old.NumJobWaiting,
		NumJobDone:    args.NumJobDone - old.NumJobDone,
	}
	statMap[args.ID] = *args

	allStat.NumConn += diff.NumConn
	allStat.NumJobDone += diff.NumJobDone
	allStat.NumJobWaiting += diff.NumJobWaiting

	*reply = 1

	return nil
}

// GetState return latest rsp server state, typically called by http server
func (t *RspServer) GetState(args *RspServerID, reply *RspServerStat) error {
	// if RspServerID is 0, return aggregating statistics
	// otherwise, return specific state for this id
	if *args != RspServerID(0) {
		if _, ok := statMap[*args]; !ok {
			return errors.New("invalid rsp server id")
		}
		*reply = statMap[*args]
	} else {
		*reply = allStat
	}

	return nil
}

func display(exit chan struct{}) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Printf("conn: %d | waiting job: %d | job done: %d",
				allStat.NumConn, allStat.NumJobWaiting, allStat.NumJobDone)
		case <-exit:
			ticker.Stop()
			return
		}
	}
}

// LaunchServer starts a rpc server acts as a memory cache server
func LaunchServer(serverAddr string) error {
	// create map
	statMap = make(map[RspServerID]RspServerStat)

	// display state
	exit := make(chan struct{})
	defer close(exit)
	go display(exit)

	// start rpc server
	rspserver := new(RspServer)
	rpc.Register(rspserver)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", serverAddr)
	if e != nil {
		return e
	}
	return http.Serve(l, nil)
}
