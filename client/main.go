package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Set the timeout for the request and create the context
	timeout := 300 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create the HTTP request to fetch the current USD exchange rate
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8082/quote", nil)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Timeout exceeded. The execution time was insufficient.")
		}
		panic(err)
	}

	// Make the request and get the response
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		log.Println("Unexpected status code:", response.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Remove special characters from the exchange rate value
	exchangeRateValue := strings.Trim(string(body), "\"\n")

	// Convert the exchange rate value to a floating-point number
	exchangeRate, err := strconv.ParseFloat(exchangeRateValue, 64)
	if err != nil {
		panic(err)
	}

	// Prompt the user for the USD value to convert to BRL
	var usdValue float64
	log.Println("Enter the USD value to convert to BRL:")
	fmt.Scanln(&usdValue)

	// Calculate the total value in BRL
	totalValue := exchangeRate * usdValue
	log.Println("The value of", usdValue, "USD in BRL is", totalValue)

	// Create a file to save the exchange rate
	file, err := os.Create("quote.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write the exchange rate to the file
	_, err = file.WriteString(fmt.Sprintf("The value of %v USD in BRL is %v", usdValue, totalValue))
	if err != nil {
		panic(err)
	}

	log.Println("Exchange rate saved successfully in 'quote.txt' file.")
}
