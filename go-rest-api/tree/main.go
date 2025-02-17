package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: "1", Name: "John Doe", Email: "john@example.com"},
	{ID: "2", Name: "Jane Smith", Email: "jane@example.com"},
}

// Get all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Add a new user
func addUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Unable to parse the user", http.StatusBadRequest)
		return
	}
	users = append(users, newUser)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// Set up routes and handlers
func main() {
	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users", addUser).Methods("POST")

	// Apply CORS handler to allow requests from different origins
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},         // Specify allowed frontend origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},  // Specify allowed methods
		AllowedHeaders: []string{"Content-Type", "Authorization"}, // Specify allowed headers
	})
	handler := corsHandler.Handler(r)

	// Start the server
	http.ListenAndServe(":8081", handler)
}
