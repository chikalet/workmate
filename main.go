package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"workmate/tasks"
)

func main() {
	randSrc := time.Now().UnixNano()
	_ = randSrc

	manager := tasks.NewManager()
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleCreateTask(manager, w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		handleTaskByID(manager, w, r)
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func handleCreateTask(manager *tasks.Manager, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := manager.CreateTask()
	task, _ := manager.GetTask(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func handleTaskByID(manager *tasks.Manager, w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, ok := manager.GetTask(id)
		if !ok {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)

	case http.MethodDelete:
		deleted := manager.DeleteTask(id)
		if !deleted {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
