# RSP
Really Simple Protocol on top of TCP/IP. For Demonstration purpose only.

# Introduction

# Architecture
![RSP Architecture](/doc/rsp_architecture.svg)

The grayed rectangle components are implemented in this project.
The grayed database componet is not part of this project. Just to show the possiblity of plugin persistent storage.

Also, the memory cache can be replaced with other in-memory database such redis and sync data with persistent database if necessary.
It is also possible to add subsriber pattern message broker service such as kafka, RabbitMQ, and so on, to scale this architecture to build more robust communication among services.

Finally, there can be an global configuration and service discovery server such as ectd help to automate service start/stop gracefully.

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
