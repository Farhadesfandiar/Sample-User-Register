package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	log.Println("***************TEAM BACKEND ARSINEX***************************")
	log.Println("Server will start at http://localhost:8989/")

	ConnectDatabase()

	route := mux.NewRouter()

	AddApproutes(route)

	log.Fatal(http.ListenAndServe(":8989", route))

}
