package main

import (
	"database/sql"
	"ffCommon/log/log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Open doesn't open a connection
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:13429)/ff_game")
	if err != nil {
		log.FatalLogger.Println(err)
		return
	}
	defer db.Close()
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(2)

	// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
	err = db.Ping()
	if err != nil {
		log.FatalLogger.Println(err)
		return
	}

	// sql string
	sqlTruncate := "truncate table account"               // truncate table account to clear data
	sqlInsert := "INSERT INTO account VALUES( ?, ? )"     // ? = placeholder
	sqlQuery := "SELECT name FROM account WHERE uuid = ?" // ? = placeholder

	// Prepare statement for inserting data
	stmtInsert, err := db.Prepare(sqlInsert)
	if err != nil {
		log.FatalLogger.Println(err.Error())
		return
	}
	defer stmtInsert.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtQuery, err := db.Prepare(sqlQuery)
	if err != nil {
		log.FatalLogger.Println(err.Error())
		return
	}
	defer stmtQuery.Close()

	// truncate table account directly
	_, err = db.Query(sqlTruncate)
	if err != nil {
		log.FatalLogger.Println(err.Error())
		return
	}

	for i := 0; i < 1000; i++ {
		_, err = stmtInsert.Exec(i, strconv.Itoa(i)) // Insert tuples (i, i^2)
		if err != nil {
			log.FatalLogger.Println(err.Error())
		}
	}

	// // Insert
	// for i := 0; i < 25; i++ {
	// 	_, err = stmtInsert.Exec(i*i, strconv.Itoa(i)) // Insert tuples (i, i^2)
	// 	if err != nil {
	// 		log.FatalLogger.Println(err.Error())
	// 	}
	// }

	// var name string

	// // Query the uuid of 1
	// err = stmtQuery.QueryRow(1).Scan(&name) // WHERE uuid = 1
	// if err != nil {
	// 	log.FatalLogger.Println(err.Error())
	// }
	// log.RunLogger.Printf("The name of uuid 1 is: %s\n", name)

	// // Query the uuid of 2
	// rows, err := stmtQuery.Query(2) // WHERE uuid = 2
	// if err != nil {
	// 	log.FatalLogger.Println(err.Error())
	// }
	// defer rows.Close()

	// columns, err := rows.Columns()
	// if err != nil {
	// 	log.FatalLogger.Println(err.Error())
	// }

	// log.RunLogger.Printf("len columns is %d\n", len(columns))
	// if len(columns) > 0 {
	// 	log.RunLogger.Printf("columns is %v\n", columns)

	// 	for rows.Next() {
	// 		err = rows.Scan(&name)
	// 		if err != nil {
	// 			log.FatalLogger.Println(err.Error())
	// 		}
	// 		log.RunLogger.Printf("The name of uuid 2 is: %s\n", name)
	// 	}
	// }

	// // Query the uuid of 2 will trigger panic
	// err = stmtQuery.QueryRow(2).Scan(&name) // WHERE uuid = 2
	// if err != nil {
	// 	log.FatalLogger.Println(err.Error())
	// }
	// log.RunLogger.Printf("The name of uuid 2 is: %s\n", name)

	// gen, _ := uuid.NewGeneratorSafe(0)
	// for i := 0; i < 10; i++ {
	// 	go util.SafeGo(query, i, stmtQuery, gen)
	// }
}

// func query(params ...interface{}) {
// 	index, _ := params[0].(int)
// 	stmtQuery, _ := params[1].(*sql.Stmt)
// 	gen, _ := params[2].(*uuid.GeneratorSafe)
// 	for i := 0; i < 100; i++ {
// 		uuid := gen.Gen()
// 	}
// }

// func query(params ...interface{}) {
// 	index, _ := params[0].(int)
// 	stmtQuery, _ := params[1].(*sql.Stmt)
// 	gen, _ := params[2].(*uuid.GeneratorSafe)
// 	for i := 0; i < 100; i++ {
// 		uuid := gen.Gen()
// 	}
// }

func insert(params ...interface{}) {

}
