package tarantool

import (
	"github.com/tarantool/go-tarantool"
	"log"
	"time"
)

func connectTarantool() {
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
	client, err := tarantool.Connect(server, opts)
	if err != nil {
		log.Fatalf("Failed to connect: %s", err.Error())
	}

	resp, err := client.Ping()
	log.Println(resp.Code)
	log.Println(resp.Data)
	log.Println(err)

	// insert new tuple { 10, 1 }
	resp, err = client.Insert(spaceNo, []interface{}{uint(10), 1})
    // or
	resp, err = client.Insert("test", []interface{}{uint(10), 1})
	log.Println("Insert")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// delete tuple with primary key { 10 }
	resp, err = client.Delete(spaceNo, indexNo, []interface{}{uint(10)})
    // or
	resp, err = client.Delete("test", "primary", []interface{}{uint(10)})
	log.Println("Delete")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// replace tuple with { 13, 1 }
	resp, err = client.Replace(spaceNo, []interface{}{uint(13), 1})
    // or
	resp, err = client.Replace("test", []interface{}{uint(13), 1})
	log.Println("Replace")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// update tuple with primary key { 13 }, incrementing second field by 3
	resp, err = client.Update(spaceNo, indexNo, []interface{}{uint(13)}, []interface{}{[]interface{}{"+", 1, 3}})
    // or
	resp, err = client.Update("test", "primary", []interface{}{uint(13)}, []interface{}{[]interface{}{"+", 1, 3}})
	log.Println("Update")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// insert tuple {15, 1} or increment second field by 1
	resp, err = client.Upsert(spaceNo, []interface{}{uint(15), 1}, []interface{}{[]interface{}{"+", 1, 1}})
    // or
	resp, err = client.Upsert("test", []interface{}{uint(15), 1}, []interface{}{[]interface{}{"+", 1, 1}})
	log.Println("Upsert")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// select just one tuple with primay key { 15 }
	resp, err = client.Select(spaceNo, indexNo, 0, 1, tarantool.IterEq, []interface{}{uint(15)})
    // or
	resp, err = client.Select("test", "primary", 0, 1, tarantool.IterEq, []interface{}{uint(15)})
	log.Println("Select")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// select tuples by condition ( primay key > 15 ) with offset 7 limit 5
	// BTREE index supposed
	resp, err = client.Select(spaceNo, indexNo, 7, 5, tarantool.IterGt, []interface{}{uint(15)})
    // or
	resp, err = client.Select("test", "primary", 7, 5, tarantool.IterGt, []interface{}{uint(15)})
	log.Println("Select")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// call function 'func_name' with arguments
	resp, err = client.Call("func_name", []interface{}{1, 2, 3})
	log.Println("Call")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// run raw lua code
	resp, err = client.Eval("return 1 + 2", []interface{}{})
	log.Println("Eval")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)
}
