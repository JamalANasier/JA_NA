package main

import (
	"fmt"
	"log"

	"github.com/jamalanasier/JA_NA/setupdb"
	"github.com/tarantool/go-tarantool"
)

func main() {

	setupdb.SetupClichouse()

	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "guest",
	})

	if err != nil {
		log.Fatalf("Connection refused")
		fmt.Println(err.Error())
	} else {
		fmt.Println("Oke")
	}

	resp, err := conn.Select("players", "primary", 0, 1, tarantool.IterEq, []interface{}{uint(15)})
	if err != nil {
		log.Fatalf("Cannot select table")
		fmt.Println(err.Error())
	} else {
		fmt.Println("There is players table there")
		log.Println("Code", resp.Code)
		log.Println("Data", resp.Data)
	}

	//setupdb.SetupTarantool()

	defer conn.Close()

	// Your logic for interacting with the database
}
