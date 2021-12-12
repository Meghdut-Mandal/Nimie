package main

import (
	"Nimie_alpha/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
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
	router.HandleFunc("/conversation/{conversation_id:[0-9]+}/{user_id:[0-9]+}", controllers.HandleChatConnections)

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	r := mux.NewRouter()
	RegisterUserRoutes(r)
	r.Use(loggingMiddleware)
	http.Handle("/", r)
	go controllers.HandleMessages()
	log.Fatal(http.ListenAndServe("localhost:"+port, r))
}
