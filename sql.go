package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"regexp"
	"strings"
)

//create db connection
func getSourceDB(connection Server) *sql.DB {
	db, err := sql.Open("mysql", connection.user+":"+connection.password+"@tcp("+connection.hostname+":3306)/"+connection.database)
	if err != nil {
		log.Println(err)
	}
	return db
}

func getFedDB(connection Server) *sql.DB {
	db, err := sql.Open("mysql", connection.feduser+":"+connection.fedpassword+"@tcp("+connection.fedhostname+":3306)/"+connection.federateddb)
	if err != nil {
		log.Println(err)
	}
	return db
}

//get array containing list of tables from specified database
func getTables(connection Server) []string {
	db := getSourceDB(connection)

	res, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Println(err)
		defer db.Close()
	}
	tables := []string{}
	var ntable string

	for res.Next() {
		res.Scan(&ntable)
		tables = append(tables, ntable+"\n")
	}
	defer db.Close()
	return tables
}

func getSchema(connection Server, table string) string {
	db := getSourceDB(connection)

	var name string
	var create string
	err := db.QueryRow("SHOW CREATE TABLE "+table).Scan(&name, &create)
	if err != nil {
		log.Println(err)
		defer db.Close()
	}
	defer db.Close()
	return create
}

func getFedSchema(connection Server, table string) string {
	db := getFedDB(connection)

	var name string
	var create string
	err := db.QueryRow("SHOW CREATE TABLE "+table).Scan(&name, &create)
	if err != nil {
		log.Println(err)
		defer db.Close()
	}
	defer db.Close()
	return create
}

func federatedConvert(connection Server, schema string) string {
	r, err := regexp.Compile(" AUTO_INCREMENT=[0-9]+")
	if err != nil {
		log.Println(err)
	}
	newengine := "ENGINE=FEDERATED CONNECTION='" + connection.servername + "'"
	converted := strings.Replace(schema, "ENGINE=InnoDB", newengine, -1)
	converted = strings.Replace(converted, "ENGINE=MyISAM", newengine, -1)
	converted = r.ReplaceAllString(converted, "")
	return converted
}

func tableExists(connection Server, table string) bool {
	var count int
	db := getFedDB(connection)
	err := db.QueryRow("SELECT count(*) FROM information_schema.TABLES where (TABLE_SCHEMA = '" + connection.federateddb + "') AND (TABLE_NAME = '" + table + "')").Scan(&count)
	if err != nil {
		log.Println(err)
		defer db.Close()
	}
	defer db.Close()
	if count > 0 {
		return true
	} else {
		return false
	}
}

func createTable(connection Server, table string, schema string) int {
	db := getFedDB(connection)
	_, err := db.Exec(schema)
	if err != nil {
		log.Println("an error occurred")
		fmt.Println(err)
		defer db.Close()
		return 1
	}
	defer db.Close()
	return 0
}

func dropTable(connection Server, table string) int {
	db := getFedDB(connection)
	_, err := db.Exec("DROP TABLE IF EXISTS " + table)
	if err != nil {
		fmt.Println(err)
		defer db.Close()
		return 1
	}
	defer db.Close()
	return 0
}

func schemaCompare(schema string, destschema string) bool {
	r, err := regexp.Compile("(?m)^.*ENGINE=.*$")
	if err != nil {
		log.Println("error compiling regex")
	}
	schema = r.ReplaceAllString(schema, "")
	destschema = r.ReplaceAllString(destschema, "")
	if schema == destschema {
		return true
	} else {
		return false
	}
}
