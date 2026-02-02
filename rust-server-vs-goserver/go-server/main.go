package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HashResponse struct {
	Hash      string `json:"hash"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
}

func main() {
	http.HandleFunc("/hash", hashHandler)
	http.HandleFunc("/health", healthHandler)

	fmt.Println("Go server starting on :8080")
	http.ListenAndServe(":8080", nil)
}

func hashHandler(w http.ResponseWriter, r *http.Request) {
	input := fmt.Sprintf("input-%d", time.Now().UnixNano())

	// SHA256 hash 100 iterations
	data := []byte(input)
	for i := 0; i < 100; i++ {
		h := sha256.Sum256(data)
		data = h[:]
	}

	response := HashResponse{
		Hash:      hex.EncodeToString(data),
		Timestamp: time.Now().UnixMilli(),
		Source:    "go",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
