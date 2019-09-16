package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

var VERSION = "v0.1.3"
var LOGFILE = "/var/log/stembolt"

func main() {

	//initialize logging
	var _, err = os.Stat(LOGFILE)
	if os.IsNotExist(err) {
		var file, err = os.Create(LOGFILE)
		log.Println("error statting logfile")
		checkError(err)
		defer file.Close()
	}
	f, err := os.OpenFile(LOGFILE, os.O_WRONLY|os.O_APPEND, 0644)
	checkError(err)
	defer f.Close()
	log.SetOutput(f)

	var config = ReadConfig()

	//user flags
	bind := flag.String("bind", config.Bind, "port to bind to")

	//setup http handlers
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/tables", tableHandler)
	http.HandleFunc("/api/schema", schemaHandler)
	http.HandleFunc("/api/convert", schemaConvert)
	http.HandleFunc("/api/tableexist", existHandler)
	http.HandleFunc("/api/updateTable", updateTableHandler)
	http.HandleFunc("/api/updateDatabase", updateDatabaseHandler)
	http.HandleFunc("/api/createTable", createTableHandler)
	http.HandleFunc("/api/dropTable", dropTableHandler)
	fmt.Println("serving stembolt " + VERSION + " on port " + *bind)
	http.ListenAndServe(":"+*bind, nil)
}

func checkError(err error) {
	if err != nil {
		log.Println(err.Error)
	}
}
