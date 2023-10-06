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
		w.Write([]byte("Unauthorized"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			data := s.GetAllMyInfo()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
			return
		}

		info, ok := s.GetMyInfo(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	case http.MethodPost:
		var info MyInfo
		err := json.NewDecoder(r.Body).Decode(&info)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
		s.AddMyInfo(info)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	case http.MethodPut:
		id := r.URL.Query().Get("id")
		var info MyInfo
		err := json.NewDecoder(r.Body).Decode(&info)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
		if s.UpdateMyInfo(id, info) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		s.DeleteMyInfo(id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	}
}

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		log.Println("Shutting down server...")
		os.Exit(0)
	}()
	s := &MyServer{
		data: make(map[string]MyInfo),
	}
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", s))
}
