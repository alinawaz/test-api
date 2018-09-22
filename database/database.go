package database

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Global Variables
var configLoaded = false
var server = "localhost"
var database = ""
var username = "root"
var password = "123"
var connection *sql.DB
var fieldString = ""

// Loading the config .env file for
// database connection
func config() {
	if configLoaded {
		return
	}
	dat, err := ioutil.ReadFile("./.env")
	if err != nil {
		log.Println(err)
	}
	s := strings.Split(string(dat), "\n")
	for _, element := range s {
		line := strings.TrimSpace(element)
		if line != "" {
			token := strings.Split(string(line), "=")
			if token[0] == "DB_SERVER" {
				server = token[1]
			}
			if token[0] == "DB_NAME" {
				database = token[1]
			}
			if token[0] == "DB_USERNAME" {
				username = token[1]
			}
			if token[0] == "DB_PASSWORD" {
				password = token[1]
			}
			configLoaded = true
		}
	}
}

// Opening the database connection
func Connect() {
	var err error
	config()

	connection, err = sql.Open("mysql", username+":"+password+"@"+server+"/"+database)
	if err != nil {
		log.Println(err.Error())
	}
}

// Close the database connection
func Close() {
	connection.Close()
}

// Query the mysql database with RAW query
func Query(queryString string) {
	_, err := connection.Exec(queryString)
	if err != nil {
		log.Println(err)
	}
}

// Select results from table as map
func Select(queryString string) []map[string]string {

	var newRows = []map[string]string{}

	rows, err := connection.Query(queryString)
	if err != nil {
		log.Println(err)
	} else {

		columns, err := rows.Columns()
		if err != nil {
			log.Println(err)
		}
		values := make([]sql.RawBytes, len(columns))

		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				log.Println(err)
			}
			var value string
			var newRow = make(map[string]string)
			for i, col := range values {
				if col == nil {
					value = "NULL"
				} else {
					value = string(col)
				}
				newRow[columns[i]] = value
			}
			newRows = append(newRows, newRow)
		}
	}
	return newRows
}

// Fill model map with values
func Fill(modelMap map[string]string, values []string) map[string]string {
	var counter = 0
	for k := range modelMap {
		log.Println(k)
		if k == "id" {
			//ignored as it's PK
		} else {
			log.Println("key=", k, "value=", values[counter])
			modelMap[k] = values[counter]
			counter++
		}
	}
	return modelMap
}

// Insert records via map
func Insert(tableName string, m map[string]string) {
	var columns = ""
	var values = ""
	for k, v := range m {
		if k == "id" {
			//ignored as it's PK
		} else {
			if columns == "" {
				columns += k
			} else {
				columns += "," + k
			}
			if values == "" {
				values += "'" + v + "'"
			} else {
				values += ",'" + v + "'"
			}
		}
	}
	var tempQuery = "INSERT INTO " + tableName + " (" + columns + ") VALUES (" + values + ");"
	Query(tempQuery)
}

// For creating database or tables
func SchemaCreate(what string, name string) {
	if what == "database" {
		var queryString = "CREATE DATABASE IF NOT EXISTS " + name + ";"
		Query(queryString)
	}
	if what == "table" {
		var queryString = "CREATE TABLE IF NOT EXISTS " + name + " (" + fieldString + "); "
		Query(queryString)
		fieldString = ""
	}
}

// Making database table's schema
func SchemaAddField(command string) {
	tokens := strings.Split(command, ":")
	var fieldName = tokens[0]
	var fieldType = tokens[1]
	var fieldLength = tokens[2]
	var isPK = ""
	var ai = ""
	if len(tokens) > 3 {
		isPK = tokens[3]
	}
	if len(tokens) > 4 {
		ai = tokens[4]
	}

	if fieldString != "" {
		fieldString += ","
	}

	if fieldType == "int" {
		fieldString += fieldName + " INT(" + fieldLength + ")"
		if ai != "" {
			fieldString += " AUTO_INCREMENT"
		}
		if isPK != "" {
			fieldString += " PRIMARY KEY"
		}
	} else if fieldType == "string" {
		fieldString += fieldName + " VARCHAR(" + fieldLength + ")"
	} else if fieldType == "decimal" {
		fieldString += fieldName + " DECIMAL(" + fieldLength + ")"
	}
}
