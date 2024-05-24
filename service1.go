package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if licenseKey == "" {
		log.Fatal("NEW_RELIC_LICENSE_KEY environment variable not set")
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Service1"),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(nrgin.Middleware(app))

	r.GET("/service1", func(c *gin.Context) {
		txn := nrgin.Transaction(c)

		// Adding custom attributes
		txn.AddAttribute("customAttribute1", "value1")
		txn.AddAttribute("customAttribute2", "value2")

		client := &http.Client{}
		req, _ := http.NewRequest("GET", "http://localhost:8081/service2", nil)
		txn.InsertDistributedTraceHeaders(req.Header)

		resp, err := client.Do(req)
		if err != nil {
			txn.NoticeError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			txn.NoticeError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}

		txn.AddAttribute("service2ResponseStatusCode", resp.StatusCode)

		// Create a channel to communicate the result
		resultCh := make(chan []string, 1)
		errorCh := make(chan error, 1)

		// Start myFunc as a goroutine and pass the transaction and body
		go func(txn *newrelic.Transaction, body []byte, resultCh chan []string, errorCh chan error) {
			goroutineTxn := txn.NewGoroutine()
			defer goroutineTxn.End()

			seg := goroutineTxn.StartSegment("myFunc2")
			seg.AddAttribute("SegAttribute", "SegValue")
			defer seg.End()

			titles, err := myFunc2(goroutineTxn, body)
			if err != nil {
				errorCh <- err
				return
			}
			resultCh <- titles
		}(txn, body, resultCh, errorCh)

		select {
		case titles := <-resultCh:
			c.JSON(http.StatusOK, gin.H{"status": "success", "titles": titles})
		case err := <-errorCh:
			txn.NoticeError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
	})

	r.Run(":8080")
}

func myFunc2(txn *newrelic.Transaction, body []byte) ([]string, error) {
	var data map[string][]map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	titles := []string{}
	for _, book := range data["books"] {
		titles = append(titles, book["title"])
	}

	return titles, nil
}
