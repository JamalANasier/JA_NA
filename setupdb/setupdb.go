package setupdb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/tarantool/go-tarantool"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Success clickhouse ok")
	}
}

func SetupClichouse() {
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?username=&compress=true&debug=true")
	checkErr(err)
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	} else {
		fmt.Println("Connecting to clickhouse")
	}

	_, err = connect.Exec(`
		CREATE TABLE IF NOT EXISTS example (
			currentTime DateTime,
			userAgent String,
			ipAddress String,
			dataBefore String,
			dataAfter String
		) 
		ENGINE = ReplacingMergeTree(currentTime)
		PARTITION BY tuple()
		ORDER BY (currentTime, userAgent);
	`)
	checkErr(err)
}

func SetupTarantool() {
	//spaceNo := uint32(512)
	//indexNo := uint32(0)

	server := "127.0.0.1:3303"
	opts := tarantool.Opts{
		Timeout:       500 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          "guest",
	}
	client, err := tarantool.Connect(server, opts)
	if err != nil {
		log.Fatalf("Failed to connect: %s", err.Error())
	} else {
		fmt.Println("Connecting as guest")
	}

	resp, err := client.Ping()
	log.Println(resp.Code)
	log.Println(resp.Data)
	log.Println(err)

	resp, err = client.Select("players", "primary", 0, 1, tarantool.IterEq, []interface{}{uint(15)})
	log.Println("Select")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)
}
