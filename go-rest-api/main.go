package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Transaction struct {
	TransactionID int    `json:"transaction_id"`
	UserName      string `json:"user_name"`
	TxnDone       bool   `json:"txn_done"`
	TxnAmount     int    `json:"txn_amount"`
}

var db *sql.DB

func main() {
	// Set up database connection
	connStr := "user=postgres password=1234 dbname=postgres sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database!")

	// Set up the routes
	http.HandleFunc("/transactions", handleTransactions)     // GET for all transactions, POST to create a transaction
	http.HandleFunc("/transactions/", handleTransactionByID) // GET, PUT, DELETE for a specific transaction by ID

	// Start the server
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// Handle requests to /transactions (GET for all transactions, POST to create a transaction)
func handleTransactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTransactions(w, r)
	case "POST":
		createTransaction(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// Handle requests to /transactions/{id} (GET, PUT, DELETE for a specific transaction)
func handleTransactionByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/transactions/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getTransactionByID(w, r, id)
	case "PUT":
		updateTransaction(w, r, id)
	case "DELETE":
		deleteTransaction(w, r, id)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var txn Transaction
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&txn)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	// Insert the transaction into the database
	err = db.QueryRow("INSERT INTO transactions (user_name, txn_done, txn_amount) VALUES ($1, $2, $3) RETURNING transaction_id", txn.UserName, txn.TxnDone, txn.TxnAmount).Scan(&txn.TransactionID)
	if err != nil {
		http.Error(w, "Failed to insert transaction into database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(txn)
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT transaction_id, user_name, txn_done, txn_amount FROM transactions")
	if err != nil {
		http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var txn Transaction
		if err := rows.Scan(&txn.TransactionID, &txn.UserName, &txn.TxnDone, &txn.TxnAmount); err != nil {
			http.Error(w, "Failed to scan transaction data", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, txn)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func getTransactionByID(w http.ResponseWriter, r *http.Request, id int) {
	var txn Transaction
	err := db.QueryRow("SELECT transaction_id, user_name, txn_done, txn_amount FROM transactions WHERE transaction_id = $1", id).Scan(&txn.TransactionID, &txn.UserName, &txn.TxnDone, &txn.TxnAmount)
	if err == sql.ErrNoRows {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to retrieve transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txn)
}

func updateTransaction(w http.ResponseWriter, r *http.Request, id int) {
	var txn Transaction
	if r.Method != "PUT" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&txn)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE transactions SET user_name = $1, txn_done = $2, txn_amount = $3 WHERE transaction_id = $4", txn.UserName, txn.TxnDone, txn.TxnAmount, id)
	if err != nil {
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction updated successfully"})
}

func deleteTransaction(w http.ResponseWriter, r *http.Request, id int) {
	if r.Method != "DELETE" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	result, err := db.Exec("DELETE FROM transactions WHERE transaction_id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted successfully"})
}
