package server

import (
	"fmt"
	"strconv"
	db "test-api/database"
)

var table = "orders"
var fields = map[string]string{
	"id":       "",
	"distance": "",
	"status":   ""}

func get(page string, limit string) []map[string]string {
	pageNumber, _ := strconv.Atoi(page)
	limitNumber, _ := strconv.Atoi(limit)
	if pageNumber > 0 {
		pageNumber--
	}
	var start = pageNumber * limitNumber
	page = strconv.Itoa(start)
	return db.Select("SELECT * FROM " + table + " LIMIT " + page + "," + limit)
}

func create(distance float64, status string) {
	var dataSet = db.Fill(fields, []string{fmt.Sprintf("%.2f", distance), status})
	db.Insert(table, dataSet)
}

func getLast() []map[string]string {
	return db.Select("SELECT * FROM orders ORDER BY id DESC LIMIT 1")
}

func getByID(id string) []map[string]string {
	return db.Select("SELECT * FROM orders WHERE id=" + id)
}

func updateStatus(id string) {
	_ = db.Select("UPDATE orders SET status='TAKEN' WHERE id=" + id)
}
