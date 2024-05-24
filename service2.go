package main

import (
	"fmt"
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
		newrelic.ConfigAppName("Service2"),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(nrgin.Middleware(app))

	r.GET("/service2", func(c *gin.Context) {
		txn := nrgin.Transaction(c)

		// Adding custom attributes
		txn.AddAttribute("customAttribute1", "value1")
		txn.AddAttribute("customAttribute2", "value2")

		// Sample books data
		books := []map[string]string{}
		for i := 1; i <= 20; i++ {
			books = append(books, map[string]string{"title": "NR GO" + fmt.Sprintf("-%d", i), "author": "Gulab " + "Sidhwani" + fmt.Sprintf("-%d", i)})
		}

		c.JSON(http.StatusOK, gin.H{"books": books})
	})

	r.Run(":8081")
}
