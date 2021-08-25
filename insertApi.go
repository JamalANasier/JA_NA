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
	app.Get("/api/v1/players", book.GetPlayers)
	app.Post("/api/v1/players_log", book.NewLog)
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
	_, err = conn.Replace(spaceNo, []interface{}{uint(1111), 19, "world"})
	if err != nil {
		conn.Close()
		return nil, err
	}
	_, err = conn.Replace(spaceNo, []interface{}{uint(1112), 21, "werld"})
	if err != nil {
		conn.Close()
		return nil, err
	}
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

	stmt, err := tx.Prepare("INSERT INTO players_log (currentTime, userAgent, ipAddress, dataBefore, dataAfter) VALUES (?, ?, ?, ?, ?, ?)")
	checkErr(err)

	for i := 0; i < 100; i++ {
		if _, err := stmt.Exec(
			time.Now(),
			"userAgent",
			"userIpAddress",
			"dataBefore this",
			"dataAfter this"
		); err != nil {
			log.Fatal(err)
		}
	}
	checkErr(tx.Commit())
}

func main() {
	app := fiber.New()
	initDatabase()

	setupRoutes(app)
	app.Listen(3000)

	defer database.DBConn.Close()
}
