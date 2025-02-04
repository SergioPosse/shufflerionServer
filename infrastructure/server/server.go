package server

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer() {

    fmt.Println("Server is starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}