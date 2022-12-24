package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

/*
	Validates JSON Date, Time, Total and each Item's Price

@param pointer to the Context struct from the GET request
*/
func validate(newReceipt receipt) (bool, string) {
	purchaseDate := newReceipt.PurchaseDate
	purchaseTime := newReceipt.PurchaseTime
	purchaseTotal := newReceipt.Total
	items := newReceipt.Items
	isValid := true
	validateResult := "Error: "
	isValidDate, _ := regexp.MatchString(`^\d{4}\-(0?[1-9]|1[012])\-(0?[1-9]|[12][0-9]|3[01])$`, purchaseDate) // only YYYY-MM-DD or YYYY-M-D are valid
	isValidTime, _ := regexp.MatchString(`^[0-2][0-3]:[0-5][0-9]$`, purchaseTime)                              // only HH:MM is valid
	isValidTotal, _ := regexp.MatchString(`(^[0-9]*[.][0-9]*$)|(^[0-9]*$)`, purchaseTotal)                     // only valid forms are 0.0 and 0
	regExp := regexp.MustCompile(`(^[0-9]*[.][0-9]*$)|(^[0-9]*$)`)
	for i := 0; i < len(items); i++ {
		isValidPrice := regExp.MatchString(items[i].Price) // only valid forms are 0.0 and 0
		if !isValidPrice {
			validateResult += ("Invalid Price for " + items[i].ShortDescription + ", ")
			isValid = false
		}
		isValidPrice = true
	}
	if !isValidDate { // create string containing error information                                                                // determine if receiptTotal has valid format)
		validateResult += "Invalid Purchase Date, "
		isValid = false
	}
	if !isValidTime {
		validateResult += "Invalid Purchase Time, "
		isValid = false
	}
	if !isValidTotal {
		validateResult += "Invalid Purchase Total"
		isValid = false
	}

	return isValid, validateResult
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
		id := xid.New()                                 // generate unique key to pass back in response
		isValid, validateResult := validate(newReceipt) // validate Date, Time, Price
		if isValid {
			receipts[id.String()] = newReceipt // Add the new receipt to the map.
			response := struct {               //anonymous struct for POST response to valid data
				ID string `json:"id"`
			}{
				ID: id.String(),
			}
			c.IndentedJSON(http.StatusCreated, response)
		} else {
			response := struct { //anonymous struct for response to invalid data
				Message string `json:"message"`
			}{
				Message: validateResult,
			}
			c.IndentedJSON(http.StatusBadRequest, response)

		}
	}
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
