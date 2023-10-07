/*
Package main is the main entry point for the server application.

The server manages MyInfo entities and exposes HTTP APIs for creating, retrieving,
updating, and deleting these entities.

This package is maintained by Srikanth Bhandary.
*/
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

const myAppSecret = "myAppSecret12254"

// MyInfo represents information about an entity.
type MyInfo struct {
	ID   string `json:"number"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// MyServer is responsible for managing MyInfo entities.
type MyServer struct {
	m    sync.Mutex
	data map[string]MyInfo
}

// AddMyInfo adds a new MyInfo entity or updates an existing one.
func (s *MyServer) AddMyInfo(info MyInfo) {
	s.m.Lock()
	defer s.m.Unlock()
	s.data[info.ID] = info
}

// GetMyInfo retrieves a MyInfo entity by its ID.
func (s *MyServer) GetMyInfo(id string) (*MyInfo, bool) {
	s.m.Lock()
	defer s.m.Unlock()
	info, ok := s.data[id]
	return &info, ok
}

// GetAllMyInfo retrieves all MyInfo entities.
func (s *MyServer) GetAllMyInfo() []MyInfo {
	s.m.Lock()
	defer s.m.Unlock()
	var allInfo []MyInfo
	for _, info := range s.data {
		allInfo = append(allInfo, info)
	}
	return allInfo
}

// DeleteMyInfo deletes a MyInfo entity by its ID.
func (s *MyServer) DeleteMyInfo(id string) {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.data, id)
}

// UpdateMyInfo updates an existing MyInfo entity by its ID.
func (s *MyServer) UpdateMyInfo(id string, info MyInfo) bool {
	s.m.Lock()
	defer s.m.Unlock()
	_, ok := s.data[id]
	if ok {
		s.data[id] = info
	}
	return ok
}

// ServeHTTP handles incoming HTTP requests and dispatches them based on the HTTP method.
func (s *MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("api-key") != myAppSecret {
		log.Println(r.Header.Get("api-key"))
		w.WriteHeader(http.StatusUnauthorized)
		response := map[string]string{"error": "Unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			data := s.GetAllMyInfo()
			json.NewEncoder(w).Encode(data)
			return
		}

		info, ok := s.GetMyInfo(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			response := map[string]string{"error": "Not Found"}
			json.NewEncoder(w).Encode(response)
			return
		}
		json.NewEncoder(w).Encode(info)
	case http.MethodPost:
		var info MyInfo
		err := json.NewDecoder(r.Body).Decode(&info)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]string{"error": "Bad Request"}
			json.NewEncoder(w).Encode(response)
			return
		}
		s.AddMyInfo(info)
		w.WriteHeader(http.StatusCreated)
		response := map[string]string{"message": "Created"}
		json.NewEncoder(w).Encode(response)
	case http.MethodPut:
		id := r.URL.Query().Get("id")
		var info MyInfo
		err := json.NewDecoder(r.Body).Decode(&info)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]string{"error": "Bad Request"}
			json.NewEncoder(w).Encode(response)
			return
		}
		if s.UpdateMyInfo(id, info) {
			w.WriteHeader(http.StatusOK)
			response := map[string]string{"message": "OK"}
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusNotFound)
			response := map[string]string{"error": "Not Found"}
			json.NewEncoder(w).Encode(response)
		}
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		s.DeleteMyInfo(id)
		w.WriteHeader(http.StatusOK)
		response := map[string]string{"message": "OK"}
		json.NewEncoder(w).Encode(response)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Method Not Allowed"}
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	// Initialize the server with some default data
	initialData := map[string]MyInfo{
		"1": {ID: "1", Name: "Alice", Age: 30},
		"2": {ID: "2", Name: "Bob", Age: 35},
	}
	s := &MyServer{
		data: initialData,
	}

	// Set up a signal channel to handle interrupts for graceful shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Start the server
	go func() {
		<-ch
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", s))
}
