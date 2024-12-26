package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		// Логика запуска реплики сервиса
		cmd := exec.Command("go", "run", "c:/Users/ilyal/Desktop/667/load_service/main.go")
		err := cmd.Start()
		if err != nil {
			fmt.Fprintf(w, "Error starting service: %v", err)
			log.Printf("Error starting service: %v", err)
			return
		}
		fmt.Fprintf(w, "Service started")
		log.Println("Service started")
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		// Логика остановки реплики сервиса
		cmd := exec.Command("pkill", "-f", "c:/Users/ilyal/Desktop/667/load_service/main.go")
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(w, "Error stopping service: %v", err)
			log.Printf("Error stopping service: %v", err)
			return
		}
		fmt.Fprintf(w, "Service stopped")
		log.Println("Service stopped")
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Agent Service is running")
		log.Println("Status request received")
	})

	log.Println("Agent Service is starting on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
