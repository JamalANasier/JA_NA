# API server test, base on Go Language

This repository is just an example how to create API server in Golang with a real Database. For this repostiory case, I use Tarantool and Clickhouse with go-Fiber.

The `go-tarantool` package has everything necessary for interfacing with
[Tarantool 1.6+](http://tarantool.org/).

The `go-clickhouse` Golang SQL database driver for [Yandex ClickHouse](https://clickhouse.yandex/)

The advantage of integrating Go with Tarantool, which is an application server
plus a DBMS, is that Go programmers can handle databases and perform on-the-fly
recompilations of embedded Lua routines, just as in C, with responses that are
faster than other packages according to public benchmarks.

## Table of contents

* [Installation](#installation)
* [Hello World](#hello-world)
* [API reference](#api-reference)
* [Walking\-through example in Go](#walking-through-example-in-go)
* [Help](#help)
* [Alternative connectors](#alternative-connectors)

## Installation

to further enrich your knowledge of this repository, we encourage you to install some amazing tools :
* [Docker](https://docs.docker.com/language/golang/build-images/)
  build, share and run any app, anywhere - on-prem or in the cloud.
* [Visual Studio Code](https://code.visualstudio.com/)
  Debug code right from the editor. Launch or attach to your running apps and debug with break points, call stacks, and an interactive console.

then clone this repository to start coding and say :
```
$ git clone https://github.com/JamalANasier/JA_NA
```

<h2>Hello World</h2>

Run this repo to test that your app is run

```
curl "http://localhost:3000/
```

<h2>API reference</h2>

When all above are set we can now test again requests to our app all the parameter send and requested are in JSON format. You can further enhance these on the codes

```
    curl "http://localhost:3000/
    curl "http://localhost:3000/api/v1/getlogs" 
	curl "http://localhost:3000/api/v1/getplayerlogs"
	curl "http://localhost:3000/api/v1/newplayer"
```


## Walking-through example in Go

We can now have a closer look at the `main.go` program and make some observations
about what it does.

```go
package main

import (
     "fmt"
     "github.com/tarantool/go-tarantool"
)

func main() {
   opts := tarantool.Opts{User: "guest"}
   conn, err := tarantool.Connect("127.0.0.1:3301", opts)
   // conn, err := tarantool.Connect("/path/to/tarantool.socket", opts)
   if err != nil {
       fmt.Println("Connection refused:", err)
   }
   resp, err := conn.Insert(999, []interface{}{99999, "BB"})
   if err != nil {
     fmt.Println("Error", err)
     fmt.Println("Code", resp.Code)
   }
}
```

**Observation 1:** the line "`github.com/tarantool/go-tarantool`" in the
`import(...)` section brings in all Tarantool-related functions and structures.

**Observation 2:** the line beginning with "`Opts :=`" sets up the options for
`Connect()`. In this example, there is only one thing in the structure, a user
name. The structure can also contain:

* `Pass` (password),
* `Timeout` (maximum number of milliseconds to wait before giving up),
* `Reconnect` (number of seconds to wait before retrying if a connection fails),
* `MaxReconnect` (maximum number of times to retry).

**Observation 3:** the line containing "`tarantool.Connect`" is essential for
beginning any session. There are two parameters:

* a string with `host:port` format, and
* the option structure that was set up earlier.

**Observation 4:** the `err` structure will be `nil` if there is no error,
otherwise it will have a description which can be retrieved with `err.Error()`.

**Observation 5:** the `Insert` request, like almost all requests, is preceded by
"`conn.`" which is the name of the object that was returned by `Connect()`.
There are two parameters:

* a space number (it could just as easily have been a space name), and
* a tuple.

## Help

To contact `go-tarantool` developers on any problems, create an issue at
[tarantool/go-tarantool](http://github.com/tarantool/go-tarantool/issues).

The developers of the [Tarantool server](http://github.com/tarantool/tarantool)
will also be happy to provide advice or receive feedback.

## Usage

```go
package main

import (
	"github.com/tarantool/go-tarantool"
	"log"
	"time"
)

func main() {
	spaceNo := uint32(512)
	indexNo := uint32(0)

	server := "127.0.0.1:3013"
	opts := tarantool.Opts{
		Timeout:       500 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          "test",
		Pass:          "test",
	}
	client, err := tarantool.Connect(server, opts)
	if err != nil {
		log.Fatalf("Failed to connect: %s", err.Error())
	}

	resp, err := client.Ping()
	log.Println(resp.Code)
	log.Println(resp.Data)
	log.Println(err)

	// insert new tuple { 10, 1 }
	resp, err = client.Insert(spaceNo, []interface{}{uint(10), 1})
    // or
	resp, err = client.Insert("test", []interface{}{uint(10), 1})
	log.Println("Insert")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// delete tuple with primary key { 10 }
	resp, err = client.Delete(spaceNo, indexNo, []interface{}{uint(10)})
    // or
	resp, err = client.Delete("test", "primary", []interface{}{uint(10)})
	log.Println("Delete")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// replace tuple with { 13, 1 }
	resp, err = client.Replace(spaceNo, []interface{}{uint(13), 1})
    // or
	resp, err = client.Replace("test", []interface{}{uint(13), 1})
	log.Println("Replace")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// update tuple with primary key { 13 }, incrementing second field by 3
	resp, err = client.Update(spaceNo, indexNo, []interface{}{uint(13)}, []interface{}{[]interface{}{"+", 1, 3}})
    // or
	resp, err = client.Update("test", "primary", []interface{}{uint(13)}, []interface{}{[]interface{}{"+", 1, 3}})
	log.Println("Update")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// insert tuple {15, 1} or increment second field by 1
	resp, err = client.Upsert(spaceNo, []interface{}{uint(15), 1}, []interface{}{[]interface{}{"+", 1, 1}})
    // or
	resp, err = client.Upsert("test", []interface{}{uint(15), 1}, []interface{}{[]interface{}{"+", 1, 1}})
	log.Println("Upsert")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// select just one tuple with primay key { 15 }
	resp, err = client.Select(spaceNo, indexNo, 0, 1, tarantool.IterEq, []interface{}{uint(15)})
    // or
	resp, err = client.Select("test", "primary", 0, 1, tarantool.IterEq, []interface{}{uint(15)})
	log.Println("Select")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// select tuples by condition ( primay key > 15 ) with offset 7 limit 5
	// BTREE index supposed
	resp, err = client.Select(spaceNo, indexNo, 7, 5, tarantool.IterGt, []interface{}{uint(15)})
    // or
	resp, err = client.Select("test", "primary", 7, 5, tarantool.IterGt, []interface{}{uint(15)})
	log.Println("Select")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// call function 'func_name' with arguments
	resp, err = client.Call("func_name", []interface{}{1, 2, 3})
	log.Println("Call")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// run raw lua code
	resp, err = client.Eval("return 1 + 2", []interface{}{})
	log.Println("Eval")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)
}
```

## Schema

```go
    // save Schema to local variable to avoid races
    schema := client.Schema

    // access Space objects by name or id
    space1 := schema.Spaces["some_space"]
    space2 := schema.SpacesById[20] // it's a map
    fmt.Printf("Space %d %s %s\n", space1.Id, space1.Name, space1.Engine)
    fmt.Printf("Space %d %d\n", space1.FieldsCount, space1.Temporary)

    // access index information by name or id
    index1 := space1.Indexes["some_index"]
    index2 := space1.IndexesById[2] // it's a map
    fmt.Printf("Index %d %s\n", index1.Id, index1.Name)

    // access index fields information by index
    indexField1 := index1.Fields[0] // it's a slice
    indexField2 := index1.Fields[1] // it's a slice
    fmt.Printf("IndexFields %s %s\n", indexField1.Name, indexField1.Type)

    // access space fields information by name or id (index)
    spaceField1 := space.Fields["some_field"]
    spaceField2 := space.FieldsById[3]
    fmt.Printf("SpaceField %s %s\n", spaceField1.Name, spaceField1.Type)
```

## Alternative connectors

- https://github.com/viciious/go-tarantool
  Has tools to emulate tarantool, and to being replica for tarantool.