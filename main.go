package main

import (
	"context"
	"fmt"

	"github.com/viciious/go-tarantool"
)

func main() {
	opts := tarantool.Options{User: "guest"}
	conn, err := tarantool.Connect("127.0.0.1:3301", &opts)
	if err != nil {
		fmt.Printf("Connection refused: %s\n", err.Error())
		return
	}

	query := &tarantool.Insert{Space: "examples", Tuple: []interface{}{uint64(99999), "BB"}}
	resp := conn.Exec(context.Background(), query)

	if resp.Error != nil {
		fmt.Println("Insert failed", resp.Error)
	} else {
		fmt.Println(fmt.Sprintf("Insert succeeded: %#v", resp.Data))
	}

	conn.Close()
}
