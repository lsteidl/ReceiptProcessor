package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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
	Awards points for retailer name
	One point for every alphanumeric character in the retailer name.

@param idFromRequest string indicating which receipt to process from receipts map
@return int value of rewarded points
*/
func getRetailerPoints(retailer string) int {
	retailerPoints := 0
	split := strings.Split(retailer, "") // convert Retailer Name into a slice of strings, 1 letter per index
	for i := 0; i < len(split); i++ {    // iterate through each letter of the Retailer Name
		isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(split[i]) // MustCompile parses a regular expression and returns a Regexp object, MatchString compares this object to "split[i]"
		if isAlphanumeric {
			retailerPoints++
		}
	}
	return retailerPoints
}

/*
	Awards points for round dollar price total
	50 points if the total is a round dollar amount with no cents.

@param idFromRequest string indicating which receipt to process from receipts map
@return int value of rewarded points
*/
func getRoundPoints(receiptTotal string) int {
	roundPoints := 0
	isRoundDollar, err := regexp.MatchString(`^([0-9]*[.]*[0]*)$`, receiptTotal) // determine if receiptTotal matches the definition of "round dollar amount"
	if isRoundDollar && err == nil {                                             // reward 50 points if string matches
		roundPoints = 50
	}
	return roundPoints

}

/*
	Awards points for dollar price total divisibility
	25 points if the total is a multiple of 0.25.

@param idFromRequest string indicating which receipt to process from receipts map
@return int value of rewarded points
*/
func getMultiplePoints(receiptTotal string) int {
	multiplePoints := 0
	isMultiple, err := regexp.MatchString(`^[0-9]*(.|.25|.50|.75)[0]*$`, receiptTotal) // determine if receiptTotal matches the definition of "divisible by 0.25"
	if isMultiple && err == nil {                                                      // reward 25 points if string matches
		multiplePoints = 25
	}
	return multiplePoints
}

/*
	Awards points for item count
	5 points for every two items on the receipt.

@param idFromRequest string indicating which receipt to process from receipts map
@return int value of rewarded points
*/
func getItemPoints(itemCount int) int {
	itemPoints := 0
	itemPoints = (itemCount / 2) * 5

	return itemPoints
}

/*
	Evaluates a single description and
	If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	The result is the number of points earned.

@param description string untrimmed item description
@param price string item price
@return int value of rewarded points
*/
func evaluateDescription(description string, price string) int {
	descPoints := 0
	descTrimmed := strings.TrimSpace(description) // trim leading and trailing spaces
	trimmedLength := len(descTrimmed)             // get length of trimmed Description
	if trimmedLength%3 == 0 {                     // check if trimmed Length is a multiple of 3
		priceAsFloat, err := strconv.ParseFloat(price, 64)
		if err == nil {
			// Mulitply the price by 0.2 and round up to the nearest integer
			pointsAsFloat := priceAsFloat * (0.2)           // multiply price by 0.2 to get points value
			pointsAsInt := int(pointsAsFloat)               // convert to integer, dropping everything after the decimal
			decimal := pointsAsFloat - float64(pointsAsInt) // store information after the decimal
			if decimal != 0 {                               // if decimal information was discarded, add 1 to satsify the "round up" specification
				pointsAsInt++
				descPoints += pointsAsInt
			}
		} else {
			fmt.Println(err)
		}
	}
	return descPoints

}

/*
	Awards points for each qualified item description
	If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	The result is the number of points earned.

@param idFromRequest string indicating which receipt to process from receipts map
@return int value of rewarded points
*/
func getDescPoints(items []item) int {
	descPoints := 0                   // Points earned from item description
	for i := 0; i < len(items); i++ { //iterate through each item on the receipt
		description := items[i].ShortDescription
		price := items[i].Price
		descPoints += evaluateDescription(description, price)
	}
	return descPoints
}

/*
	Awards points for receipt date.
	6 points if the day in the purchase date is odd. Assumes date format from examples (YYYY-MM-DD)

@param idFromRequest string indicating which receipt to process from receipts map
@return int value of rewarded points
*/
func getDatePoints(purchaseDate string) int {
	datePoints := 0
	isOdd, err := regexp.MatchString(`^[0-9]*-[0-9]*-[0-3](1|3|5|7|9)$`, purchaseDate)
	if isOdd && err == nil {
		datePoints = 6
	}
	return datePoints
}

/*
	Awards points for receipt time.
	10 points if the time of purchase is after 2:00pm and before 4:00pm. (14:00 to 15:59)

@param idFromRequest string indicating which receipt to process from receipts map
@return int value of rewarded points
*/
func getTimePoints(purchaseTime string) int {
	timePoints := 0
	isValidTime, err := regexp.MatchString(`^1[4-5]:[0-5][0-9]$`, purchaseTime)
	if isValidTime && err == nil {
		timePoints = 10
	}
	return timePoints
}

/*
	Calculates total number of points for each category

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
