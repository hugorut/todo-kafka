[![](https://godoc.org/github.com/hugorut/todo-kafka?status.svg)](http://godoc.org/github.com/hugorut/todo-kafka)

[![Build Status](https://travis-ci.org/hugorut/todo-kafka.svg?branch=master)](https://travis-ci.org/hugorut/todo-kafka)

# Todo Kafka

A very simple implementation of the classic todos JSON api which uses kafka + event sourcing to persist data to an in-memory store.

## Getting started

### Prerequisites

You must have a running local version of kafka to run successfully. This can be downloaded via brew for mac:

```sh
brew install kafka
```

and then then run as a service:

```sh
brew services start kafka
```

default port is `9092`

### Running

clone the repo

```
git clone git@github.com:hugorut/todo-kafka.git
```

`cd` into the project

build the program

```
go build
```

Then run

```
./todo-kafka
```

`todo-kafka` takes optional arguments:

| name      	| description                                 	| default        	|
|-----------	|---------------------------------------------	|----------------	|
| address   	| kafka service address                       	| 127.0.0.1:9092 	|
| http-port 	| port that the http service should listen on 	| 8081           	|

### JSON API

#### Create
| parameter     | description           | type   |
 required (Y/N)    |
|-----------    |---------------------- |------  |----------------   |
| name          | the name of the todo  | string | Y                 |

```
curl -X POST -d '{"name":"my important todo"} localhost:8081/todo/create
```

#### Update

| parameter     | description                              | type   | required (Y/N)    |
|-----------    |---------------------------------------   |------  |----------------   |
| id            | the id of the todo you wish to update    | int    | Y                 |
| name          | the updated name of the todo             | string | Y                 |

```
curl -X POST -d '{"id":4,"name":"rename it to something else"} localhost:8081/todo/create
```

#### Delete

| parameter     | description                       | type  | required (Y/N)    |
|-----------    |--------------------------------   |------ |----------------   |
| id            | the id todo you wish to delete    | int   | Y                 |

```
curl -X DELETE -d '{"id":4}' localhost:8081/todo/delete
```

#### List

```
curl  localhost:8081/todo/list
```
