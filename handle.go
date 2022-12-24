package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

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
