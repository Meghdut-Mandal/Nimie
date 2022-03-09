package main

import (
	"github.com/Meghdut-Mandal/Nimie/controllers/grpc_api"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//port := os.Getenv("PORT")

	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc_api.NimieApiServerImpl{}

	grpcServer := grpc.NewServer()

	log.Printf("Server started on port 8080")

	grpc_api.RegisterNimieApiServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
