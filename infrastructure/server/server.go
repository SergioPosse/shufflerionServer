package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func StartServer(mux *http.ServeMux) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	fmt.Printf("🚀 Server is starting on port %s...\n", port)

	errChan := make(chan error, 1)

	go func() {
		errChan <- server.ListenAndServe()
	}()

	fmt.Println("✅ Server started successfully!")

	if err := <-errChan; err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server error: %v", err)
	}
}
