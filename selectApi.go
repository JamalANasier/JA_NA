package main

import (
	"fmt"
	"log"
	"time"


	"github.com/jmoiron/sqlx"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/tarantool/go-tarantool"
)

type  players struct {
	gorm.Model
	Id  uint16 `json:"Id"`
	Name String `json:"Name"`
	Age int    `json:"Age"`
}

type  players_log struct {
	gorm.Model
	currentTime  DateTime `json:"currentTime"`
	userAgent String `json:"userAgent"` 
	ipAddress String `json:"ipAddress"`
	dataBefore String `json:"dataBefore"`
	dataAfter String `json:"dataAfter"`
}

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/getplayers", book.getPlayersId)
	app.Post("/api/v1/getplayers_log", book.getPlayersLog)
}

func initTarantool() {
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
	conn, err := tarantool.Connect(server, opts)
	if err != nil {
		return nil, err
	}
	resp, err = client.Select(spaceNo, indexNo, 0, 1, tarantool.IterEq, []interface{}{uint(15)})
	return conn, nil
}

func initClickHouse() {
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

	var items []struct {

		currentTime time.Time `db:"currentTime"`
		userAgent string    `db:"userAgent"`
		ipAddress string    `db:"ipAddress"`
		dataBefore string    `db:"dataBefore"`
		dataAfter string    `db:"dataAfter"`
	}

	checkErr(connect.Select(&items, "SELECT currentTime, userAgent, ipAddress, dataBefore, dataAfter FROM players_log"))

	for _, item := range items {
		log.Printf("currentTime: %d, userAgent: %v, browser: %s, categories: %v, action_time: %s", item.currentTime, item.userAgent, item.ipAddress, item.dataBefore, item.dataAfter)
	}
}

func main() {
	app := fiber.New()
	initDatabase()

	setupRoutes(app)
	app.Listen(3000)

	defer database.DBConn.Close()
}
