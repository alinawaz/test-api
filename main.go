package main

import (
	"fmt"

	db "test-api/database"
	router "test-api/routes"
)

func main() {
	fmt.Println("Connecting & setting up database ...")
	db.Connect()
	// Assigning permissions
	db.Query("UPDATE mysql.user SET Grant_priv='Y', Super_priv='Y' WHERE User='root';")
	db.Query("FLUSH PRIVILEGES;")
	// Creating table [orders]
	db.SchemaAddField("id:int:11:pk:ai")
	db.SchemaAddField("distance:string:55")
	db.SchemaAddField("status:string:55")
	db.SchemaCreate("table", "orders")
	fmt.Println("Ready.")

	router.MakeRoutes()

}
