package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Somebody has to teach you officers the difference between a warp matrix flux capacitor and a self-sealing stem bolt."+"\n")
	fmt.Fprintf(w, "stembolt "+VERSION)
}

func tableHandler(w http.ResponseWriter, r *http.Request) {
	database := r.URL.Query().Get("database")
	connection := getConnection(database)
	tables := getTables(connection)

	for _, table := range tables {
		fmt.Fprintf(w, table)
	}
}

func updateTableHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	authenticated := checkAuth(auth)
	if authenticated {
		database := r.URL.Query().Get("database")
		table := r.URL.Query().Get("table")
		connection := getConnection(database)
		schema := getSchema(connection, table)
		schema = federatedConvert(connection, schema)
		tableExistsCheck := tableExists(connection, table)
		if tableExistsCheck != true {
			createres := createTable(connection, table, schema)
			if createres > 0 {
				fmt.Fprintf(w, "problem creating table\n")
			} else {
				fmt.Fprintf(w, "table created\n")
			}
		} else {
			destschema := getFedSchema(connection, table)
			schemasMatch := schemaCompare(schema, destschema)
			if schemasMatch {
				fmt.Fprintf(w, "Table "+table+" already matches production\n")
			} else {
				dropRes := dropTable(connection, table)
				if dropRes > 0 {
					log.Println("an error occurred dropping table \"" + table + "\"\n")
					fmt.Fprintf(w, "an error occurred dropping table \""+table+"\"\n")
				}
				createRes := createTable(connection, table, schema)
				if createRes > 0 {
					log.Println("an error occurred creating table \"" + table + "\"\n")
					fmt.Fprintf(w, "an error occurred creating table \""+table+"\"\n")
				}
				fmt.Fprintf(w, "Federated table has been updated from prod\n")
			}
		}
	} else {
		fmt.Fprintf(w, "authentication failure")
	}
}

func updateDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	authenticated := checkAuth(auth)
	if authenticated {
		database := r.URL.Query().Get("database")
		connection := getConnection(database)
		tables := getTables(connection)
		for _, table := range tables {
			table = strings.Replace(table, "\n", "", -1)
			schema := getSchema(connection, table)
			schema = federatedConvert(connection, schema)
			tableExistsCheck := tableExists(connection, table)
			if tableExistsCheck != true {
				fmt.Fprintf(w, "Creating table "+table+"\n")
				createres := createTable(connection, table, schema)
				if createres > 0 {
					log.Println("problem creating table" + table + "\n")
				}
			} else {
				destschema := getFedSchema(connection, table)
				schemasMatch := schemaCompare(schema, destschema)
				if schemasMatch {
				} else {
					log.Println(schema + "\n")
					log.Println(destschema + "\n")
					fmt.Fprintf(w, "updating table "+table+"\n")
					dropRes := dropTable(connection, table)
					if dropRes > 0 {
						log.Println("an error occurred dropping table\"" + table + "\"")
					}
					createRes := createTable(connection, table, schema)
					if createRes > 0 {
						log.Println("an error occured creating table \"" + table + "\"")
					}
				}
			}
		}
		fmt.Fprintf(w, "All tables on "+database+" are up to date")
	} else {
		fmt.Fprintf(w, "authentication failure")
	}
}
