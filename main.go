package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func createTransactionFileIfNotExists(filename string, attempts int) error {
	file, err := os.Open(filename)

	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		if attempts == 0 {
			return err
		}

		defer file.Close()

		file, _ = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)

		file.WriteString("customer_id,transaction_id,amount,transaction_type,timestamp\n")
		file.WriteString("acc_123,tx_1,100.00,CREDIT,2024-10-27T10:00:00Z\n")
		file.WriteString("acc_456,tx_2,50.00,CREDIT,2024-10-27T10:01:00Z\n")
		file.WriteString("acc_123,tx_3,25.50,DEBIT,2024-10-27T10:02:00Z\n")
		file.WriteString("acc_456,tx_4,10.00,CREDIT,2024-10-27T10:03:00Z\n")
		file.WriteString("acc_789,tx_5,200.00,CREDIT,2024-10-27T10:04:00Z\n")
		file.WriteString("acc_123,tx_6,5.00,DEBIT,2024-10-27T10:05:00Z")

		return createTransactionFileIfNotExists(filename, attempts-1)
	}

	return err
}

func main() {
	argsWithoutProg := os.Args[1:]

	fileName := ""

	if len(argsWithoutProg) > 0 {
		fileName = argsWithoutProg[0]

		extension := filepath.Ext(fileName)

		if !strings.HasSuffix(strings.ToLower(extension), "csv") {
			fmt.Printf("%s is not a CSV", fileName)
			return
		}

		f, err := os.Open(fileName)

		if os.IsNotExist(err) {
			fmt.Printf("%s does not exist", fileName)
			return
		}

		defer f.Close()

	} else {
		fileName = "transactions.csv"

		err := createTransactionFileIfNotExists(fileName, 3)

		if err != nil {
			fmt.Printf("Unable to create or access %s", fileName)
			return
		}
	}

	file, err := os.Open(fileName)

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	defer file.Close()

	csv := csv.NewReader(file)

	// skip header row

	header, err := csv.Read()

	if !slices.Contains(header, "customer_id") || !slices.Contains(header, "amount") || !slices.Contains(header, "transaction_type") {
		fmt.Printf("%s is missing required headers: customer_id, amount, transaction_type", fileName)
	}

	customerIdIndex := slices.Index(header, "customer_id")
	amountIndex := slices.Index(header, "amount")
	transactionTypeIndex := slices.Index(header, "transaction_type")

	maxUsedIndex := math.Max(float64(customerIdIndex), float64(amountIndex))
	maxUsedIndex = math.Max(maxUsedIndex, float64(transactionTypeIndex))

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	lineCount := 0

	balances := make(map[string]float64)

	for {
		record, err := csv.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("%v", err)
			lineCount++
			continue
		}

		if len(record)-1 < int(maxUsedIndex) {
			fmt.Printf("Record on line %d does not have enough columns, skipping\n", lineCount)
			lineCount++
			continue
		}

		customerId := record[customerIdIndex]
		amountString := record[amountIndex]
		if strings.HasPrefix(amountString, "£") {
			amountString = strings.ReplaceAll(amountString, "£", "")
		}

		amount, err := strconv.ParseFloat(amountString, 64)

		if err != nil {
			fmt.Printf("Record on line %d has an amount set as %s which cannot be converted to a number, skipping\n", lineCount, amountString)
			lineCount++
			continue
		}

		transactionType := record[transactionTypeIndex]
		transactionType = strings.ToUpper(transactionType)

		if transactionType != "CREDIT" && transactionType != "DEBIT" {
			fmt.Printf("Record on line %d has a transaction type of %s which is not the allowed values of DEBIT or CREDIT, skipping\n", lineCount, transactionType)
			lineCount++
			continue
		}

		currentBalance, exists := balances[customerId]

		if !exists {
			currentBalance = 0
		}

		if transactionType == "CREDIT" {
			currentBalance += amount
		} else {
			currentBalance -= amount
		}

		balances[customerId] = currentBalance

		lineCount++
	}

	if len(balances) == 0 {
		fmt.Println("No valid balances found")
		return
	}

	keys := []string{}

	for key := range balances {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("%s,%.2f\n", key, balances[key])
	}

}
