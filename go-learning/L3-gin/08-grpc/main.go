package main

import (
	"log"
)

func main() {
	log.Println("Starting gRPC server...")

	if err := startServer(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
