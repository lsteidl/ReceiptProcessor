package main

import (
	"github.com/gin-gonic/gin"
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
