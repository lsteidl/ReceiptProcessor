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

/*
	Calculates total number of points for each category

@param idFromRequest string indicating which receipt to process
@return int value of total rewared points
*/
func calculatePoints(idFromRequest string) int {
	// Handle situation of invalid ID requested, maybe earlier than this line
	if receipts[idFromRequest].Retailer == "" {
		fmt.Println("Bad Request, Receipt not Found")
		return -101
	}
	fmt.Printf("Retailer is %s", receipts[idFromRequest].Retailer)
	retailerPoints := 0 // Points earned for Retailer Name
	roundPoints := 0    // Points earned for Round receipt total
	multiplePoints := 0 // Points earned for receipt total being a multiple of (0.25)
	itemPoints := 0     // Points earned for every 2 items
	descPoints := 0     // Points earned from item description
	datePoints := 0     // Points earned from purchase date
	timePoints := 0     // Points earned from purchase time

	// One point for every alphanumeric character in the retailer name.
	retailer := receipts[idFromRequest].Retailer
	split := strings.Split(retailer, "") // convert Retailer Name into an array of strings, 1 letter per index
	fmt.Printf("first letter: %s\n", split[0])
	for i := 0; i < len(split); i++ { // iterate through each letter of the Retailer Name
		isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(split[i]) // MustCompile parses a regular expression and returns a Regexp object, MatchString compares this object to "split[i]"
		if isAlphanumeric {
			fmt.Printf("%s is AlphaNumerical\n", split[i])
			retailerPoints++
		} else {
			fmt.Printf("%s is NOT AlphaNumerical\n", split[i])
		}
	}
	fmt.Printf("Total  Retailer Points: %d\n", retailerPoints)
	// 50 points if the total is a round dollar amount with no cents.
	fmt.Printf("Receipt Total is %s\n", receipts[idFromRequest].Total)
	receiptTotal := receipts[idFromRequest].Total
	isRoundDollar, err := regexp.MatchString(`^[0-9]*.[0]*$`, receiptTotal) // determine if receiptTotal matches the definition of "round dollar amount"
	if isRoundDollar && err == nil {                                        // reward 50 points if string matches
		roundPoints = 50
	}
	// 25 points if the total is a multiple of 0.25.
	isMultiple, err := regexp.MatchString(`^[0-9]*(.|.25|.50|.75)[0]*$`, receiptTotal) // determine if receiptTotal matches the definition of "divisible by 0.25"
	if isMultiple && err == nil {                                                      // reward 25 points if string matches
		multiplePoints = 25
	}
	// 5 points for every two items on the receipt.
	itemCount := len(receipts[idFromRequest].Items)
	itemPoints = (itemCount / 2) * 5

	//If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	//The result is the number of points earned.
	for i := 0; i < len(receipts[idFromRequest].Items); i++ { //iterate through each item on the receipt
		//fmt.Printf("Evaluating item %d...\n", i)
		description := receipts[idFromRequest].Items[i].ShortDescription
		price := receipts[idFromRequest].Items[i].Price
		//fmt.Println(description)
		descTrimmed := strings.TrimSpace(description) // trim leading and trailing spaces
		//fmt.Println(descTrimmed)
		trimmedLength := len(descTrimmed) // get length of trimmed Description
		//fmt.Printf("Trimmed Length is %d\n", trimmedLength)
		if trimmedLength%3 == 0 { // check if trimmed Length is a multiple of 3
			fmt.Printf("Item %d - Trimmed Length is a multiple of 3\n", i)
			//	fmt.Printf("price: %s\n", price)
			priceAsFloat, err := strconv.ParseFloat(price, 64)
			if err == nil {
				// Mulitply the price by 0.2 and round up to the nearest integer
				pointsAsFloat := priceAsFloat * (0.2)           // multiply price by 0.2 to get points value
				pointsAsInt := int(pointsAsFloat)               // convert to integer, dropping everything after the decimal
				decimal := pointsAsFloat - float64(pointsAsInt) // store information after the decimal
				//		fmt.Printf("priceAsFloat: %f\n", priceAsFloat)
				//		fmt.Printf("priceAsInt: %d\n", priceAsInt)
				//		fmt.Printf("decimal: %f\n", decimal)
				if decimal != 0 { //if decimal information was discarded, add 1 to satsify the "round up" specification
					//			fmt.Printf("Rounding up\n")
					pointsAsInt++
					descPoints += pointsAsInt
				} else {
					fmt.Printf("Already round!!!\n")
				}
			} else {
				fmt.Println(err)
			}
			//	fmt.Printf("final modified priceAsInt: %d\n", priceAsInt)
		} else {
			fmt.Printf("Item %d - Trimmed Length is NOT a multiple of 3\n", i)
		}
	}

	//6 points if the day in the purchase date is odd.
	purchaseDate := receipts[idFromRequest].PurchaseDate
	isOdd, err := regexp.MatchString(`^[0-9]*-[0-9]*-[0-3](1|3|5|7|9)$`, purchaseDate)
	if isOdd && err == nil {
		datePoints = 6
	}
	//10 points if the time of purchase is after 2:00pm and before 4:00pm. (14:01 to 15:59)
	purchaseTime := receipts[idFromRequest].PurchaseTime
	isValidTime, err := regexp.MatchString(`^1[4-5]:[0-5][0-9]$`, purchaseTime)
	if isValidTime && err == nil {
		timePoints = 10
	}
	// Add all point categories
	TotalPoints := retailerPoints + roundPoints + multiplePoints + itemPoints + descPoints + datePoints + timePoints
	fmt.Printf("Total Points %d\n", TotalPoints)
	fmt.Println("Points Breakdown...")
	fmt.Printf("Retailer Points %d\n", retailerPoints)
	fmt.Printf("Round Points %d\n", roundPoints)
	fmt.Printf("Multiple Points %d\n", multiplePoints)
	fmt.Printf("Item Points %d\n", itemPoints)
	fmt.Printf("Desc Points %d\n", descPoints)
	fmt.Printf("Date Points %d\n", datePoints)
	fmt.Printf("Time Points %d\n", timePoints)

	return TotalPoints
}

/*
	Processes the requested receipt and sends JSON containing the calculated point value as a response, example: { "points": 32 }

@param pointer to the Context struct from the GET request
*/
func getPoints(c *gin.Context) {
	idFromRequest := c.Param("id")
	points := calculatePoints(idFromRequest)
	response := struct { // anonymous struct for the GET response
		Points int `json:"points"`
	}{
		Points: points,
	}
	c.IndentedJSON(http.StatusCreated, response) // send the calculated point value
}

/*
	Processes the JSON info received in the POST request body.
	Sends a JSON containing a unique identifier as a response, example: {  "id": "ceibu5s3c37hfmrav9q0" }

@param pointer to the Context struct from the POST request
*/
func postReceipt(c *gin.Context) {
	var newReceipt receipt
	// Call BindJSON to bind the JSON data to newReceipt
	if err := c.BindJSON(&newReceipt); err != nil {
		fmt.Println("Error after binding", err)
		return
	}
	// generate unique key to pass back in response
	id := xid.New()                    // convert key to String
	receipts[id.String()] = newReceipt // Add the new receipt to the slice.
	response := struct {               //anonymous struct for POST response
		ID string `json:"id"`
	}{
		ID: id.String(),
	}
	c.IndentedJSON(http.StatusCreated, response)
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