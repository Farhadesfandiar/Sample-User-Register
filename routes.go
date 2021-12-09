package main

import (
	"log"

	"github.com/gorilla/mux"
)

func AddApproutes(route *mux.Router) {

	log.Println("Loadeding Routes...")

	route.HandleFunc("/", RenderHome)

	route.HandleFunc("/login", RenderLogin)

	route.HandleFunc("/register", RenderRegister)

	route.HandleFunc("/signin", SignInUser).Methods("POST")

	route.HandleFunc("/signup", SignUpUser).Methods("POST")

	route.HandleFunc("/userDetails", GetUserDetails).Methods("GET")

	log.Println("Routes are Loaded.")
	log.Println("Routes sucessfully loaded on localhost port 8989")
	log.Println("Navigate to registration page if you are a new user :8989/register")
	log.Println("***********************************************************************")
}
