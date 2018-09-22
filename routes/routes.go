package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	controllers "test-api/server"
)

// MakeRoutes for the application
func MakeRoutes() {
	router := mux.NewRouter()
	router.Path("/orders").
		Queries("page", "{page}").
		Queries("limit", "{limit}").
		HandlerFunc(controllers.Index).
		Methods("GET")
	router.HandleFunc("/orders/{page}/{limit}", controllers.Index).Methods("GET")
	router.HandleFunc("/order", controllers.Store).Methods("POST")
	router.HandleFunc("/order/{id}", controllers.Update).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", router))
}
