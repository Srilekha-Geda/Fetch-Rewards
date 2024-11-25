package main

import (
	"encoding/json"
	"fmt"

	"crypto/sha256"
	"encoding/hex"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Response structure
type Response struct {
	ID string `json:"id"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var receipts = make(map[string]int) // In-memory storage for receipts and their points

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running!"))
	})

	http.HandleFunc("/receipts/process", processReceiptHandler)
	http.HandleFunc("/receipts/", getPointsHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler for processing receipts
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Calculate points
	points := calculatePoints(receipt)

	// Log request for debugging
	log.Println("Processing receipt...")

	// Generate unique ID for the receipt
	id := generateReceiptID(receipt.Retailer, receipt.PurchaseDate, receipt.PurchaseTime, receipt.Total)
	receipts[id] = points

	// Respond with the ID
	response := map[string]string{"id": id}
	// response := Response{
	// 	ID: id,
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateReceiptID(retailer, purchaseDate, purchaseTime, total string) string {
	// Concatenate key fields
	data := retailer + purchaseDate + purchaseTime + total

	// Generate a SHA-256 hash
	hash := sha256.Sum256([]byte(data))

	// Convert hash to a hex string
	return hex.EncodeToString(hash[:])
}

// Handler for getting points
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the receipt ID from the URL
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	if id == "" {
		http.Error(w, "Missing receipt ID", http.StatusBadRequest)
		return
	}

	points, exists := receipts[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Respond with the points
	response := map[string]int{"points": points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Function to calculate points based on the rules
func calculatePoints(receipt Receipt) int {
	points := 0

	// 1. One point for every alphanumeric character in the retailer name
	for _, char := range receipt.Retailer {
		if char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z' || char >= '0' && char <= '9' {
			points++
		}
	}

	// 2. 50 points if the total is a round dollar amount with no cents
	if strings.HasSuffix(receipt.Total, ".00") {
		points += 50
	}

	// 3. 25 points if the total is a multiple of 0.25
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// 4. 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// 5. Points based on item description length
	for _, item := range receipt.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int(math.Ceil(price * 0.2))
			}
		}
	}

	// 6. 6 points if the day in the purchase date is odd
	date, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && date.Day()%2 != 0 {
		points += 6
	}

	// 7. 10 points if the time of purchase is between 2:00pm and 4:00pm
	timeOfPurchase, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil && timeOfPurchase.Hour() == 14 {
		points += 10
	}

	return points
}
