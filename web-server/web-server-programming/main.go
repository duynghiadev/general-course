package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"web-server-programming/database"
)

func main() {
	dbClient, err := database.NewDBClient()
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	fmt.Println("Connected to MongoDB!")

	// HTTP handlers
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// Create new user
			var user database.User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := dbClient.CreateUser(user); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})

		case http.MethodGet:
			// Check if username parameter exists
			username := r.URL.Query().Get("username")
			if username != "" {
				// Get single user
				user, err := dbClient.GetUser(username)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if user == nil {
					http.Error(w, "user not found", http.StatusNotFound)
					return
				}
				json.NewEncoder(w).Encode(user)
				return
			}

			// Get all users
			users, err := dbClient.GetAllUsers()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(users)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the server
	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
