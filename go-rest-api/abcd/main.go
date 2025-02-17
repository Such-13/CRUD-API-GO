package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Transaction struct {
	TransactionID int    `json:"transaction_id"`
	UserName      string `json:"user_name"`
	TxnDone       bool   `json:"txn_done"`
	TxnAmount     int    `json:"txn_amount"`
}

func main() {
	// Ask for a transaction ID
	fmt.Print("Enter transaction ID: ")
	var transactionID int
	_, err := fmt.Scanf("%d", &transactionID)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	// Define the API URL with the transaction ID
	url := fmt.Sprintf("http://localhost:8081/transactions/%d", transactionID)

	// Make the GET request to the API
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: Received non-OK response code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Parse the JSON data into a Transaction struct
	var txn Transaction
	err = json.Unmarshal(body, &txn)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Output the transaction details
	fmt.Printf("Transaction ID: %d\n", txn.TransactionID)
	fmt.Printf("User: %s\n", txn.UserName)
	fmt.Printf("Transaction Done: %t\n", txn.TxnDone)
	fmt.Printf("Transaction Amount: %d\n", txn.TxnAmount)
}
