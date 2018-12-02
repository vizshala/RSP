RSP
=============
Really Simple Protocol on top of TCP/IP. For Demonstration purpose only.

Introduction 
=============

RSP stands for Really Simple Protocol. This protocol is built on top of TCP. I coded this project in order to further understand how the go's asynchronous features work.

Architecture
=============
![RSP Architecture](/doc/rsp_architecture.svg)

The grayed rectangle components are implemented in this project.
The grayed database componet is not part of this project. Just to show the possiblity of plugin persistent storage.

Also, the memory cache can be replaced with other in-memory database such redis and sync data with persistent database if necessary.
It is also possible to add subsriber pattern message broker service such as kafka, RabbitMQ, and so on, to scale this architecture to build more robust communication among services.

Finally, there can be an global configuration and service discovery server such as ectd help to automate service start/stop gracefully.

Build
=============
Change directory to source code root, assuming located at $(GOPATH)/src/RSP. From shell, run the following commands:

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

Start Server
=============

The launching order of server, httpserver and cache doesn't matter. Just make sure client must be luanched after server is running.

httpserver is a plugin to help to monitor the rsp server state. If httpserver is running, you can open http://127.0.0.1:8080 with web browser such as Chrome to check the server status. Currenly only # of connection, # of waiting jobs, # of completed jobs are shown. As the following image shows. The page shows the past 30 records of server state with per record updated at a 5-second interval.

![server monitor](/doc/server_monitor.png)


## Usage

Start client and you will see something like this
```
Enter command: 
```

which is the prompt text asking you to type something. You can try some commands with this. The Format is:

```
command [payload]
```

Each line is taken as a complete command plus its payload.

Command|Payload  (optional)| Description |
-------|-------------------|-------------|
shorten|a URL              | invoke an external API for shortening long URL and returning a short URL |
quit   |n/a                | exit applicatoin |

### Example
```
Enter command:
shorten https://www.google.com
[200 http://bit.ly/2SmO2Qr]
```

Where 

`Enter command` is prompt text.

`shorten https://www.google.com` is the command I entered

`[200 http://bit.ly/2SmO2Qr]` is the result returned from rsp server

```
Enter command:
quit
Exit now...
```

Where 

`Enter command` is prompt text.

`quit` is the command I entered, meaning I want to exit the rsp client.

`Exit now...` is the execution result telling me that the connection with server is closed and the application is going to exit.

## Server

3rd-party libraries
====================
* [C3.js][c3]
* [D3.js][d3]

[c3]: https://c3js.org/ 
[d3]: https://d3js.org/ 
