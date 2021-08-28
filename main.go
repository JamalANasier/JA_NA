package main

import (
	"fmt"
	"log"

	"github.com/tarantool/go-tarantool"
)

func main() {

	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "guest",
	})

	if err != nil {
		log.Fatalf("Connection refused")
	} else {
		fmt.Println("Oke")
	}

	defer conn.Close()

	// Your logic for interacting with the database
}
