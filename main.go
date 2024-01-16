package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Transaction struct {
	Date       time.Time
	Narratives []string
	Type       string
	Credit     float64
	Debit      float64
	Currency   string
}

type Balance struct {
	Currency string  `json:"currency"`
	Total    float64 `json:"total"`
}

func main() {
	router := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.POST("/parse", handleCSVUpload)
	router.Run(":" + port)
}

func handleCSVUpload(c *gin.Context) {
	// Check if the request contains a file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Check if the uploaded file is a CSV file
	if !strings.HasSuffix(header.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format, only CSV files are allowed"})
		return
	}

	// Extract the date from the form data
	dateString := c.PostForm("date")

	// Parse the date from the form data
	date, err := time.Parse("02/01/2006", dateString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Parse the CSV file
	transactions, err := parseCSV(file, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing CSV file"})
		return
	}

	balances := fetchBalance(transactions)

	c.JSON(http.StatusOK, balances)
}

func parseCSV(reader io.Reader, dateFilter time.Time) ([]Transaction, error) {
	var transactions []Transaction

	// Create a CSV reader
	csvReader := csv.NewReader(reader)

	// Read the CSV headers
	headers, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	// Map headers to indices for easy access
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header] = i
	}

	// Read the CSV records
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			logrus.Info("file end reached")
			break
		} else if err != nil {
			logrus.Errorf("error reading file %v", err)
			return nil, err
		}

		// Parse each record into a Transaction struct
		transaction := Transaction{
			Date:     parseDate(record[headerMap["Date"]]),
			Type:     record[headerMap["Type"]],
			Credit:   parseFloat(record[headerMap["Credit"]]),
			Debit:    parseFloat(record[headerMap["Debit"]]),
			Currency: record[headerMap["Currency"]],
		}

		// Extract narratives from the record
		transaction.Narratives = make([]string, 5)
		for i := 0; i < 5; i++ {
			transaction.Narratives[i] = record[headerMap[fmt.Sprintf("Narrative %d", i+1)]]
		}

		// Check if the transaction date matches the provided date filter
		if transaction.Date.Equal(dateFilter) && isPayment(transaction) {
			// Append the transaction to the slice
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

func isPayment(transaction Transaction) bool {
	// Use a regular expression to match the payment reference pattern
	regex := regexp.MustCompile(`PAY\d{6}[A-Z]{2}`)
	for _, narrative := range transaction.Narratives {
		if regex.MatchString(narrative) {
			return true
		}
	}
	return false
}

func fetchBalance(transactions []Transaction) []Balance {
	// Calculate account balances per currency
	balances := make(map[string]Balance)

	for _, transaction := range transactions {
		currency := transaction.Currency
		if _, ok := balances[currency]; !ok {
			balances[currency] = Balance{Currency: currency}
		}
		balance := balances[currency]
		balance.Total += transaction.Debit - transaction.Credit
		balances[currency] = balance
	}

	// Respond with the updated balances in JSON format
	response := make([]Balance, 0, len(balances))
	for _, balance := range balances {
		response = append(response, balance)
	}
	return response
}

func parseDate(dateString string) time.Time {
	date, err := time.Parse("02/01/2006", dateString)
	if err != nil {
		return time.Time{}
	}
	return date
}

func parseFloat(value string) float64 {
	result, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", ""), 64)
	if err != nil {
		return 0.0
	}
	return result
}
