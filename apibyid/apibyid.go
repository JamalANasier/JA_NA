package apibyid

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber"
	"github.com/jamalanasier/JA_NA/apioffsetlimit"
	"github.com/jamalanasier/JA_NA/apiplayerlog"
	"github.com/tarantool/go-tarantool"
)

type Players_log struct {
	Id          int      `json:"id,string,omitempty"`
	Name        string   `json:"name"`
	Age         int      `json:"age,string,omitempty"`
	CurrentTime Datetime `json:"currentTime"`
	UserAgent   string   `json:"userAgent"`
	IpAddress   string   `json:"ipAddress"`
	DataBefore  string   `json:"dataBefore"`
	DataAfter   string   `json:"dataAfter"`
}

type Players struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Datetime struct {
	time.Time
}

func (t *Datetime) UnmarshalJSON(input []byte) error {
	strInput := strings.Trim(string(input), `"`)
	newTime, err := time.Parse(time.RFC3339, strInput)
	if err != nil {
		return err
	}

	t.Time = newTime
	return nil
}

func GetBooks(c *fiber.Ctx) {
	Connect, Err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if Err != nil {
		log.Printf("Failed to connect")
		log.Fatal(Err)
	}
	rows, err := Connect.Query("SELECT current_time, user_agent, ip_address, data_before, data_after FROM player_log")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
	c.Send(string(jsonData))
}

func NewBook(c *fiber.Ctx) {
	p := new(Players_log)
	err := c.BodyParser(p)

	if err != nil {
		log.Fatal(err)
	}

	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := tarantool.Connect("127.0.0.1:3303", tarantool.Opts{
		User: "guest",
	})

	if err != nil {
		log.Fatalf("Connection refused")
		fmt.Println(err.Error())
	} else {
		fmt.Println(conn)
		_, err := conn.Upsert("players", []interface{}{int(p.Id), "Jamal", 25}, []interface{}{[]interface{}{"+", 1, 1}})
		if err != nil {
			log.Fatal(err)
			return
		}

		var (
			tx, _   = connect.Begin()
			stmt, _ = tx.Prepare("INSERT INTO player_log (current_time, user_agent, ip_address, data_before, data_after) VALUES (?, ?, ?, ?, ?)")
		)
		defer stmt.Close()
		if _, err := stmt.Exec(
			p.CurrentTime,
			p.UserAgent,
			p.IpAddress,
			p.DataBefore,
			p.DataBefore,
		); err != nil {
			log.Fatal(err)
			return
		}
	}
}

// New Player by Id
func NewPlayer(c *fiber.Ctx) {
	p := new(Players_log)
	err := c.BodyParser(p)

	if err != nil {
		log.Fatal(err)
	}

	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := tarantool.Connect("127.0.0.1:3303", tarantool.Opts{
		User: "guest",
	})

	if err != nil {
		log.Fatalf("Connection refused")
		fmt.Println(err.Error())
	}

	resp, err := conn.Upsert("players", []interface{}{p.Id, p.Name, p.Age, 1}, []interface{}{[]interface{}{"+", 1, 1}})

	if err != nil {
		log.Fatal("tarantool exception", resp)
		return
	}

	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO player_log (current_time, user_agent, ip_address, data_before, data_after) VALUES (?, ?, ?, ?, ?)")
	)
	defer stmt.Close()
	const layoutUS = "2006-01-02"

	if _, err := stmt.Exec(
		p.CurrentTime.Format(layoutUS),
		p.UserAgent,
		p.IpAddress,
		p.DataBefore,
		p.DataAfter,
	); err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	c.Send("Upsert ok")
}

func helloWorld(c *fiber.Ctx) {
	c.Send("Hello, World!")
}

func SetupRoutes(app *fiber.App) {
	app.Get("/", helloWorld)

	app.Get("/api/v1/getlogs", apioffsetlimit.GetLogs)
	app.Get("/api/v1/getplayerlogs", apiplayerlog.GetPlayerLogs)
	app.Get("/api/v1/books", GetBooks)
	app.Post("/api/v1/newplayer", NewPlayer)
	app.Post("/api/v1/book", NewBook)
}
