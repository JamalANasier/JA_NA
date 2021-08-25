package clickhouse

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ClickHouse/clickhouse-go"
)

func clickhouse() {
	var err error
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?username=&compress=true&debug=true")
	errCheck(err)
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

	errCheck(err)
	tx, err := connect.Begin()
	errCheck(err)
	errCheck(tx.Commit())
}

func errCheck(err string) {
	if err != nil {
		log.Fatal(err)
	}
}
