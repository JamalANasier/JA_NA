package apibyid

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber"
	"github.com/tarantool/go-tarantool"
	"gorm.io/gorm"
)

var (
	DBConn  *gorm.DB
	Connect *sql.DB
	Err     error
)

type Players_log struct {
	Id          uint     `json:"id,string,omitempty"`
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

// GetBooks ...
func GetBooks(c *fiber.Ctx) {
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

// GetBook ...
func GetBook(c *fiber.Ctx) {
	/*
		id := c.Params("id")
		db := DBConn
		var book Book
		db.Find(&book, id)
		c.JSON(book)
	*/
	c.Send("Hello, World!")
}

// NewBook ...
func NewBook(c *fiber.Ctx) {
	p := new(Players_log)
	err := c.BodyParser(p)

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
			tx, _   = Connect.Begin()
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
			tx, _   = Connect.Begin()
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
	c.Send("Upsert ok")
}

// DeleteBook ...
func DeleteBook(c *fiber.Ctx) {
	/*
		id := c.Params("id")
		db := DBConn
		var book Book
		db.First(&book, id)
		if book.Title == "" {
			c.Status(http.StatusNotFound).Send("No book found with given id")
			return
		}
		db.Delete(&book)
		c.Status(http.StatusNoContent).Send()
	*/
	c.Send("Hello, World!")
}

func helloWorld(c *fiber.Ctx) {
	c.Send("Hello, World!")
}

func SetupRoutes(app *fiber.App) {
	app.Get("/", helloWorld)
	app.Get("/", helloWorld)

	app.Get("/api/v1/books", GetBooks)
	app.Get("/api/v1/book/:id", GetBook)
	app.Post("/api/v1/newplayer", NewPlayer)
	app.Post("/api/v1/book", NewBook)
	app.Delete("/api/v1/book/:id", DeleteBook)
}
