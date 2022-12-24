package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// receipt represents data about a submitted receipt.
type receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []item `json:"items"`
	Total        string `json:"total"`
}

// item represents data for each item within a receipt
type item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// receipts map, holds receipt submission data
var receipts = map[string]receipt{}

/*
	Calculates total number of points for each category
	Calls functions from calculate.go

@param idFromRequest string indicating which receipt to process from receipts map
@return int value of total rewarded points
*/
func calculatePoints(idFromRequest string) int {
	// Retrieve information needed for point calculation
	retailer := receipts[idFromRequest].Retailer
	receiptTotal := receipts[idFromRequest].Total
	itemCount := len(receipts[idFromRequest].Items)
	items := receipts[idFromRequest].Items
	purchaseDate := receipts[idFromRequest].PurchaseDate
	purchaseTime := receipts[idFromRequest].PurchaseTime

	// Add all point categories
	totalPoints := getRetailerPoints(retailer) + getRoundPoints(receiptTotal) + getMultiplePoints(receiptTotal) +
		getItemPoints(itemCount) + getDescPoints(items) + getDatePoints(purchaseDate) + getTimePoints(purchaseTime)

	return totalPoints
}

/*
	Processes the requested receipt
	Responds to request with JSON containing the calculated point value as a response, example: { "points": 32 }

@param pointer to the Context struct from the GET request
*/
func getPoints(c *gin.Context) {
	idFromRequest := c.Param("id")
	if receipts[idFromRequest].Retailer == "" { // Handle situation of unknown ID request, return error code 404
		fmt.Println("Bad Request, Receipt not Found")
		c.IndentedJSON(http.StatusNotFound, "Bad Request, Receipt not Found")
	} else {
		points := calculatePoints(idFromRequest)
		response := struct { // anonymous struct for the GET response
			Points int `json:"points"`
		}{
			Points: points,
		}
		c.IndentedJSON(http.StatusOK, response) // send the calculated point value
	}
}

/*
	Processes the JSON info received in the POST request body.
	Responds to request with JSON containing a unique identifier as a response, example: {  "id": "ceibu5s3c37hfmrav9q0" }

@param pointer to the Context struct from the POST request
*/
func postReceipt(c *gin.Context) {
	var newReceipt receipt
	err := c.BindJSON(&newReceipt)
	if err != nil { // bind the JSON data to newReceipt
		fmt.Println("Error after binding JSON data in postReceipt()", err)
		c.IndentedJSON(http.StatusBadRequest, "Error processing data")
	} else {
		id := xid.New()                    // generate unique key to pass back in response
		receipts[id.String()] = newReceipt // Add the new receipt to the map.
		response := struct {               //anonymous struct for POST response
			ID string `json:"id"`
		}{
			ID: id.String(),
		}
		c.IndentedJSON(http.StatusCreated, response)
	}
}

/*
Provides 2 endpoints, Process Receipts (POST) and Get Points (GET)
Utilizes the Gin Web Framework
Listening and serving on localhost:8080
*/
func main() {
	router := gin.Default() // initializes Gin router
	router.POST("/receipts/process", postReceipt)
	router.GET("/receipts/:id/points", getPoints)

	router.Run("localhost:8080") // Attach the router to an Http server
}
