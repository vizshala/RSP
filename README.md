# RSP
Really Simple Protocol on top of TCP/IP. For Demonstration purpose only.

# Introduction

# Architecture



# Build

Change directory to source code root, assuming located at $(GOPATH)/src/RSP.
From shell, run the following commands

```
go build RSP/service/server

go build RSP/service/client

go build RSP/service/cache

go build RSP/service/httpserver
```

The above commands will build server, client, cache, httpserver, where:

server is rsp server,
client is rsp client,
cache is memory cache server,
httpserver is http server


# 3rd-party libraries
* [C3.js][c3]
* [D3.js][d3]

[c3]: https://c3js.org/ 
[d3]: https://d3js.org/ 
