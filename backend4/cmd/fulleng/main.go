package main

import (
	"log"

	"github.com/katerinasolntsye/fulleng/internal/app"
)

func main() {
	application := app.NewApp()

	if err := application.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// curl -H "Content-Type: application/json" -d "{\"username\": \"johndoe\", \"password\": \"mysecurepassword\"}" http://localhost:8000/signup
