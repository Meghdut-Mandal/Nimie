package main

import (
	"Nimie_alpha/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.Method + " " + r.RequestURI + " ")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

var RegisterUserRoutes = func(router *mux.Router) {
	router.HandleFunc("/user/register", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/status/create", controllers.CreateStatus).Methods("POST")
}

func main() {
	r := mux.NewRouter()
	RegisterUserRoutes(r)
	r.Use(loggingMiddleware)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
