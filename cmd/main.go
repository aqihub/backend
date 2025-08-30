package main

import (
	"go-tropic-thunder/pkg/routes"
	"log"
	"net/http"
)

func main() {
	// Load environment variables
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}

	router := routes.SetupRouter()

	log.Println("Starting server on :3000...")
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
