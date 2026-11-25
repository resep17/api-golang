package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Struct untuk data dummy
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	City  string `json:"city"`
}

// Data dummy (in-memory)
var users = []User{
	{ID: 1, Name: "Budi Santoso", Email: "budi@example.com", City: "Jakarta"},
	{ID: 2, Name: "Siti Aminah", Email: "siti@example.com", City: "Bandung"},
	{ID: 3, Name: "Andi Wijaya", Email: "andi@example.com", City: "Surabaya"},
}

func main() {
	mux := http.NewServeMux()

	// Endpoint root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Simple Go API on Vercel",
			"endpoints": "/api/users, /api/users/:id",
		})
	})

	// GET semua users (support filter ?city=)
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cityFilter := r.URL.Query().Get("city")
		result := users

		if cityFilter != "" {
			var filtered []User
			for _, u := range users {
				if strings.EqualFold(u.City, cityFilter) {
					filtered = append(filtered, u)
				}
			}
			result = filtered
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"total":   len(result),
			"data":    result,
		})
	})

	// GET user by ID
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Ambil ID dari path /api/users/1
		idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"success":false,"message":"ID tidak valid"}`, http.StatusBadRequest)
			return
		}

		for _, u := range users {
			if u.ID == id {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": true,
					"data":    u,
				})
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "User tidak ditemukan",
		})
	})

	// Vercel menyediakan PORT via environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // fallback untuk local testing
	}

	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}