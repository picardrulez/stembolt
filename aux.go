//these are test functions are created to test the sql connection functionality.  I don't want to delete them yet in case we need them in the future
package main

import (
	"fmt"
	"net/http"
)

func schemaHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	authenticated := checkAuth(auth)
	if authenticated {
		database := r.URL.Query().Get("database")
		table := r.URL.Query().Get("table")
		connection := getConnection(database)
		schema := getSchema(connection, table)
		fmt.Fprintf(w, schema)
	} else {
		fmt.Fprintf(w, "authentication failure")
	}
}

func schemaConvert(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	authenticated := checkAuth(auth)
	if authenticated {
		database := r.URL.Query().Get("database")
		table := r.URL.Query().Get("table")
		connection := getConnection(database)
		schema := getSchema(connection, table)
		converted := federatedConvert(connection, schema)
		fmt.Fprintf(w, schema)
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, converted)
	} else {
		fmt.Fprintf(w, "authentication failure")
	}
}

func existHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	authenticated := checkAuth(auth)
	if authenticated {
		database := r.URL.Query().Get("database")
		table := r.URL.Query().Get("table")
		connection := getConnection(database)
		tableCheck := tableExists(connection, table)
		if tableCheck {
			fmt.Fprintf(w, "table exists")
		} else {
			fmt.Fprintf(w, "table does not exist")
		}
	} else {
		fmt.Fprintf(w, "authentication failure")
	}
}

func dropTableHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	authenticated := checkAuth(auth)
	if authenticated {
		database := r.URL.Query().Get("database")
		table := r.URL.Query().Get("table")
		connection := getConnection(database)
		dropRes := dropTable(connection, table)
		if dropRes > 0 {
			fmt.Fprintf(w, "an error occurred dropping table")
		} else {
			fmt.Fprintf(w, "table dropped")
		}
	} else {
		fmt.Fprintf(w, "authentication failure")
	}
}

func createTableHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	authenticated := checkAuth(auth)
	if authenticated {
		database := r.URL.Query().Get("database")
		table := r.URL.Query().Get("table")
		connection := getConnection(database)
		schema := getSchema(connection, table)
		schema = federatedConvert(connection, schema)
		createRes := createTable(connection, table, schema)
		if createRes > 0 {
			fmt.Fprintf(w, "an error occurred creating table")
		} else {
			fmt.Fprintf(w, "table created")
		}
	} else {
		fmt.Fprintf(w, "authentication failure")
	}
}
