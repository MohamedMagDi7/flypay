package main

import (
	"github.com/callicoder/go-docker/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	//using gorilla mux router
	router := mux.NewRouter()
	// route /api/payment/transaction to our flypayController handler
	router.HandleFunc("/api/payment/transaction", controllers.GetFlyPays).Methods("GET")
	//use port :8000
	port := "8000"
	//Launch the app, visit localhost:8000/api
	log.Println("Starting Server")
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}
