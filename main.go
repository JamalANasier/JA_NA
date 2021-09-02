package main

import (
	"github.com/gofiber/fiber"
	"github.com/jamalanasier/JA_NA/apibyid"
	"github.com/jamalanasier/JA_NA/setupdb"
	_ "github.com/tarantool/go-tarantool"
)

func main() {

	setupdb.SetupClichouse()
	//setupdb.SetupTarantool()
	app := fiber.New()
	apibyid.SetupRoutes(app)
	app.Listen(3000)
	/*
		conn, err := tarantool.Connect("127.0.0.1:3303", tarantool.Opts{
			User: "guest",
		})

		if err != nil {
			log.Fatalf("Connection refused")
			fmt.Println(err.Error())
		} else {
			fmt.Println("Oke")
		}

		resp, err := conn.Select("players", "primary", 0, 1, tarantool.IterEq, []interface{}{uint(15)})
		if err != nil {
			log.Fatalf("Cannot select table")
			fmt.Println(err.Error())
		} else {
			fmt.Println("There is players table there")
			log.Println("Code", resp.Code)
			log.Println("Data", resp.Data)
		}

		//setupdb.SetupTarantool()

		defer conn.Close()

	*/
}
