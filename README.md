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

Also, the memory cache can be replaced with other in-memory database such [redis] and sync data with persistent database if necessary.
It is also possible to add subsriber pattern message broker service such as [Kafka], [RabbitMQ], and so on, to scale this architecture to build more robust communication among services.

Finally, there can be an global configuration and service discovery server such as [etcd] help to automate service start/stop gracefully.

## Design Concept
The idea behind this desing is to implement an flexible architecture that services are low-coupled. In this graph, rsp server, cache server and http server can be run independently. The depencendies are as following:

* http server pulls information from the cache server.
* rsp server pushs information to the cache server.
* http server and rsp server has no dependency on each other and even don't know each.

This is the goal that I try to achieve: create a scalable system with low coupling.

By introducing an extra layer (in this case, the `cache server` act as the middleware), it is able to decouple the relation between rsp server and http server that would had been highly coupled if they were to communicate with each other.

## Components

### rsp server
The server implement the RSP protocol.

### cache server
This is a rpc protocol based implementation simple in-memory cache. 

### http server
An simple http server that can be pluged in anytime. This server will try to find cache server and pull information from cache server.
It servers an RESTful API by providing this URL http://127.0.0.1:8080/api/v1/state. 

### rsp client
The client implement the RSP protocol. Provide an simple shell for user to interact with rsp server.

RSP protocol
============
The protocol is simple

## client
```
command_string payload
```

## server
```
status_code result_string
```
status_code|description|
-------|----|
200|command succeeded|
400|command failed or no such command|
408|timed out|


Build
=============
Build Environment: 
go version go1.11.2 windows/amd64
go version go1.11.2 linux/amd64

Change directory to source code root, assuming located at $(GOPATH)/src/RSP. From shell, run the following commands:

```
go build ./service/server

go build ./service/client

go build ./service/cache

go build ./service/httpserver
```

The above commands will build server, client, cache, httpserver, where:

server is rsp server,

client is rsp client,

cache is memory cache server,

httpserver is http server

The above executable files will be located at $(GOPATH)/src/RSP

Start Server
=============

The launching order of server, httpserver and cache doesn't matter. Just make sure client must be luanched after server is running.

httpserver is a plugin to help to monitor the rsp server state. If httpserver is running, you can open http://127.0.0.1:8080 with web browser such as Chrome to check the server status. Currenly only # of connection, # of waiting jobs, # of completed jobs are shown. As the following image shows. The page shows the past 30 records of server state with per record updated at a 5-second interval. Unnder $(GOPATH)/src/RSP there is an folder called `template`. The displayed web page is generated based on template pages in this folder. If you see `404 page not found`, check httpserver and template folder is correctly structured.

The following figure shows the screenshot when httpserver and template folder configured correctly.

![server monitor](/doc/server_monitor.png)

## run erver

### start rsp server
```
$(GOPATH)/src/RSP/server
```
option|type|description|
-------|----|-------------------|
cache|string|cache server address (default "localhost:1234")|
job|int|job queue capacity (default 100)|
port|int|designate port to listen to (default 1314)|
req_timeout|int|reqeust timeout (default 5)|
rps|int|external api rate per second (default 30)|
        
### start cache server
```
$(GOPATH)/src/RSP/cache
```
option|type|description|
------|----|-------------------|
addr|string|designate the address to connect to (default "localhost:1234")|
        
### start http server
```
$(GOPATH)/src/RSP/httpserver
```

option|type|description|
------|----|-------------------|
cache |string|cache server address (default "localhost:1234")|
port |int|port to listen to (default 8080)|

## run client

### start client 
```
$(GOPATH)/src/RSP/client
```

option|type|description|
------|----|-------------------|
addr |string|designate the address to connect to (default "localhost:1314")|  

you will see something like this
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
shorten|a URL              | invoke an external API for shortening long URL and returning a short URL. Take usage of bit.ly API|
quit   |n/a                | exit applicatoin |

#### Example
```
Enter command:
shorten https://www.google.com
200 http://bit.ly/2SmO2Qr
```

Where 

`Enter command` is prompt text.

`shorten https://www.google.com` is the command I entered

`200 http://bit.ly/2SmO2Qr` is the result returned from rsp server

```
Enter command:
quit
Exit now...
```

Where 

`Enter command` is prompt text.

`quit` is the command I entered, meaning I want to exit the rsp client.

`Exit now...` is the execution result telling me that the connection with server is closed and the application is going to exit.

Userful tools 
===============
[mock api][mockapi]: A great tool to mock up RESTful API. Considering the scenario when frontend and backend developers need to conform to the same API spec define for their product and to avoid dependencies, API callers may utilize this kind of tool to mock the API before the backend is ready for testing. 

Further Reading
===============
[AMQP][AMQP]

Todo
====
* Integrate CI/CD to this project. Currently take usage of circleci for this purpose. Still have some problem with the procedure. The dashboard can be found at https://circleci.com/gh/vizshala/RSP. 

3rd-party libraries
====================
* [C3.js][c3]
* [D3.js][d3]

[c3]: https://c3js.org/ 
[d3]: https://d3js.org/ 
[Kafka]: https://kafka.apache.org/
[etcd]: https://github.com/etcd-io/etcd
[RabbitMQ]: https://www.rabbitmq.com/
[redis]: https://redis.io/
[mockapi]: https://www.mockapi.io/
[AMQP]: https://en.wikipedia.org/wiki/Advanced_Message_Queuing_Protocol
