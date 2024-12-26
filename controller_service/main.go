package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type ServiceStatus struct {
	ID     string `json:"id"`
	UpTime string `json:"uptime"`
}

var (
	services = make(map[string]ServiceStatus)
	mu       sync.Mutex
)

func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(services)
		log.Println("Status request received")
	})

	http.HandleFunc("/scale", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Action string `json:"action"`
			Count  int    `json:"count"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			log.Printf("Invalid request: %v", err)
			return
		}

		switch request.Action {
		case "increase":
			scaleUp(request.Count)
		case "decrease":
			scaleDown(request.Count)
		default:
			http.Error(w, "Invalid action", http.StatusBadRequest)
			log.Printf("Invalid action: %s", request.Action)
		}
	})

	go monitorServices()

	log.Println("Controller Service is starting on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func monitorServices() {
	for {
		mu.Lock()
		for id, status := range services {
			// Обновление статуса сервисов
			status.UpTime = time.Now().String()
			services[id] = status
		}
		mu.Unlock()
		time.Sleep(10 * time.Second) // Исправлено
	}
}

func scaleUp(count int) {
	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < count; i++ {
		id := fmt.Sprintf("service-%d", len(services)+1)
		services[id] = ServiceStatus{
			ID:     id,
			UpTime: time.Now().String(),
		}
		log.Printf("Service %s started", id)
	}
}

func scaleDown(count int) {
	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < count; i++ {
		if len(services) == 0 {
			break
		}
		for id := range services {
			delete(services, id)
			log.Printf("Service %s stopped", id)
			break
		}
	}
}
