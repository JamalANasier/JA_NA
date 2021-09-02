package setupdb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/tarantool/go-tarantool"
)

func SetupClichouse() {
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	}

	_, err = connect.Exec(`
	CREATE TABLE IF NOT EXISTS player_log (
		current_time Date,
		user_agent   String,
		ip_address   String,
		data_before  String,
		data_after   String
		) engine=Memory
	`)

	/*
		_, err = connect.Exec(`
			CREATE TABLE IF NOT EXISTS example (
				country_code FixedString(2),
				os_id        UInt8,
				browser_id   UInt8,
				categories   Array(Int16),
				action_day   Date,
				action_time  DateTime
				) engine=Memory
			`)
		*/
	
		if err != nil {
		log.Fatal(err)
		}
	/*	
		va	r (
				tx, _   = connect.Begin()
		s	tmt, _ = tx.Prepare("INSERT INTO player_log (current_time, user_agent, ip_address, data_before, data_after) VALUES (?, ?, ?, ?, ?)")
		)
	de	fer stmt.Close()
	
		for i := 0; i < 10; i++ {
			if _, err := stmt.Exec(
					time.Now(),
			"RU",
					"ipAppdress",
					"dataBefore",
					"dataAfter",
		); err != nil {
					log.Fatal(err)
				}
		}	
		
			if err := tx.Commit(); err != nil {
			log.Fatal(err)
			}
	
			rows, err := connect.Query("SELECT current_time, user_agent, ip_address, data_before, data_after FROM player_log")
			if err != nil {
			log.Fatal(err)
	}
		de	fer rows.Close()
	
			for rows.Next() {
			var (
				currentTime                                 time.Time
					userAgent, ipAddress, dataBefore, dataAfter string
				)
				if err := rows.Scan(&currentTime, &userAgent, &ipAddress, &dataBefore, &dataAfter); err != nil {
				log.Fatal(err)
			}
				log.Printf("currentTime: %s, userAgent: %s, ipAddress: %s, dataBefore: %s, dataAfter: %s", currentTime, userAgent, ipAddress, dataBefore, dataAfter)
			}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		/*if _, err := connect.Exec("DROP TABLE example"); err != nil {
			log.Fatal(err)
		}*/

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
