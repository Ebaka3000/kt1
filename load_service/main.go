package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Status request received")
		fmt.Fprintf(w, "Load Service is running on %s", os.Getenv("HOSTNAME"))
	})

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Work request received")
		// Эмуляция рабочей нагрузки
		time.Sleep(2 * time.Second)
		fmt.Fprintf(w, "Work done on %s", os.Getenv("HOSTNAME"))
	})

	log.Println("Load Service is starting on port 8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
