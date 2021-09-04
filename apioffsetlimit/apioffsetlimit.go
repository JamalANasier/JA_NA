package apioffsetlimit

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber"
)

var (
	Connect *sql.DB
	Err     error
)

type Alllogs struct {
	Limit  string `json:"limit"`
	Offset string `json:"offset"`
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

func GetLogs(c *fiber.Ctx) {
	p := new(Alllogs)

	if err := c.BodyParser(p); err != nil {
		log.Fatal(err)
	}

	Connect, Err = sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if Err != nil {
		log.Printf("Failed to connect")
		log.Fatal(Err)
	}
	var (
		strQuery string = "SELECT current_time, user_agent, ip_address, data_before, data_after FROM player_log ORDER BY current_time"
	)

	myslice := []string{strQuery, "LIMIT", p.Limit, "OFFSET", p.Offset}
	resultQuery := strings.Join(myslice, " ")
	rows, err := Connect.Query(resultQuery)
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
