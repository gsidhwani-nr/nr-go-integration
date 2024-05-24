package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type Service1Response struct {
	Status string   `json:"status"`
	Titles []string `json:"titles"`
}

func main() {
	licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if licenseKey == "" {
		log.Fatal("NEW_RELIC_LICENSE_KEY environment variable not set")
	}

	// Create a New Relic application
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("ClientApp"),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	serviceid := 0
	for i := 0; i < 2000; i++ { // Loop to make 2000 calls
		// Determine service ID based on loop index
		temp := i % 5
		switch temp {
		case 1:
			serviceid = 1
		case 2:
			serviceid = 2
		case 3:
			serviceid = 3
		case 4:
			serviceid = 4
		case 0:
			serviceid = 5
		}

		// Start a new transaction for each call
		txn := app.StartTransaction(fmt.Sprintf("CallService1-%d", serviceid))

		// Adding custom attributes to the transaction
		txn.AddAttribute("callNumber", i+1)
		txn.AddAttribute("callDescription", "Invoking Service1")

		// Create a new HTTP request
		req, _ := http.NewRequest("GET", "http://localhost:8080/service1", nil)
		txn.InsertDistributedTraceHeaders(req.Header)

		// Make the HTTP request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			txn.NoticeError(err)
			log.Printf("Failed to call Service1: %v", err)
			continue
		}

		// Log the response status
		log.Printf("Response from Service1: %s", resp.Status)

		// Add response status code as an attribute
		txn.AddAttribute("responseStatusCode", resp.StatusCode)

		// Parse JSON response
		var service1Resp Service1Response
		err = json.NewDecoder(resp.Body).Decode(&service1Resp)
		if err != nil {
			txn.NoticeError(err)
			log.Printf("Failed to decode JSON response: %v", err)
			continue
		}

		// Print book titles
		fmt.Printf("Book titles from Service1: %v\n", service1Resp.Titles)

		// End the transaction
		txn.End()

		// Close the response body
		resp.Body.Close()

		// Wait for 2 seconds before the next call
		time.Sleep(2 * time.Second)
	}

	// Shutdown the New Relic application
	app.Shutdown(10 * time.Second)
}
