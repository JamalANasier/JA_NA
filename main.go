package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go"
)

func main() {
	var err error
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?username=&compress=true&debug=true")
	checkErr(err)
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	}

	_, err = connect.Exec(`
		CREATE TABLE IF NOT EXISTS players_log (
			currentTime  DateTime,
			userAgent String,
			ipAddress String,
			dataBefore String,
			dataAfter String,
		) engine=Memory
	`)

	checkErr(err)
	tx, err := connect.Begin()
	checkErr(err)
	checkErr(tx.Commit())
}

func errCheck(err string) {
	if err != nil {
		log.Fatal(err)
	}
}
