package apiplayerlog

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber"
	"github.com/tarantool/go-tarantool"
)

var (
	Connect *sql.DB
	Err     error
)

type Player_log struct {
	Id        int    `json:"id" xml:"id" form:"id"`
	UserAgent string `json:"user_agent" xml:"user_agent" form:"user_agent"`
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

func GetPlayerLogs(c *fiber.Ctx) {
	p := new(Player_log)

	if err := c.BodyParser(p); err != nil {
		log.Fatal(err)
	}
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

	resp, err := client.Select("players", "primary", 0, 0, tarantool.IterEq, []interface{}{uint(p.Id)})
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(resp.Data)

	Connect, Err = sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
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
