package controllers

import (
	"deepanshu18099/blockchain_ledger_backend/chaincode"
	"log"
	"net/http"
	"github.com/google/uuid"


	"github.com/gin-gonic/gin"
)



func AddMoneyToUser(c *gin.Context) {
	/*
	Input: Auth token, amount
	This should update the balance of the user in ledger
	Output: success message, updated balance of user
	*/
	// use authmiddleware to check if token is valid and get claims
	// using the authcheck function
	claims, ok := Authcheck(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	argss := []string{}
	argss = append(argss, email)
	// now prepare to send the request to the chaincode
	// Call the chaincode function to create the user on the ledger
	args := chaincode.BuildChaincodeArgs(argss, "AddMoneyToUser")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("AddMoneyToUser function called", output)

	// ................missing part.......................

	// // check if the output has the updated balance
	// updatedbalance, ok := outputdecoded["updatedbalance"].(string)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated balance"})
	// 	return
	// }
	// now return the updated balance
	c.JSON(http.StatusOK, gin.H{
		"message":        "Money added to user",
		"updatedbalance": "1100",
		"transaction_id": "tx_10101",
	})

}


func GetTransports(c *gin.Context) {
	/*
	Input: Auth token, Source, Destination, Date, Mode of transport
	will be called by user
	Output: List of transports available
	*/
	
	claims, ok := Authcheck(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	email, ok := claims["email"].(string)
	log.Println("GetTransports function called", email)
	if !ok {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	// get Source, Destination, Date, Mode of transport from the request body using a struct
	type GetTransportsRequest struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		Date        string `json:"date"`
		Mode        string `json:"mode"`
	}
	

	var Thisreq GetTransportsRequest


	if err := c.ShouldBindJSON(&Thisreq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// now prepare to send the request to the chaincode
	args := chaincode.BuildChaincodeArgs([]string{Thisreq.Source, Thisreq.Destination, Thisreq.Date, Thisreq.Mode}, "GetTransports")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("GetTransports function called", output)
	// Decode the output
	// .........................


	// now return success, and transport details
	c.JSON(http.StatusOK, gin.H{
		"message":        "Transport details fetched successfully",
		"transports":     "List of transports",
		"transaction_id": "tx_10101",
	})
}



func BookTicket(c *gin.Context) {
	/*
	Input: Auth token, transportID, date, seatNumber
	This should update the ticket of both the users in ledger
	Also update the Balance of user and provider
	Output: success message, updated balance of user
	*/
	// use authmiddleware to check if token is valid and get claims
	// using the authcheck function
	claims, ok := Authcheck(c)
	// get transportID, date from the request body
	// var transportID string
	// var date string
	// var seatNumber int

	type BookTicketRequest struct {
		TransportID string `json:"transportID"`
		Date        string `json:"date"`
		SeatNumber  int    `json:"seatNumber"`
	}
	var Thisreq BookTicketRequest

	if err := c.ShouldBindJSON(&Thisreq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	// now call func (s *SmartContract) BookTicket(ctx contractapi.TransactionContextInterface, userID, transportID, date string, seatNumber int32) (string, error) {
	// now prepare to send the request to the chaincode
	args := chaincode.BuildChaincodeArgs([]string{email, Thisreq.TransportID, Thisreq.Date, string(Thisreq.SeatNumber)}, "BookTicket")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("BookTicket function called", output)
	// Decode the output
	// ........................



	// now return success, and updated balance of the user
	c.JSON(http.StatusOK, gin.H{
		"message":        "Ticket booked successfully",
		"updatedbalance": "1100",
		"transaction_id": "tx_10101",
	})
	
}




func GetUserTickets(c *gin.Context) {
	/*
	Input: Auth token
	Output: List of tickets booked by the user
	*/
	// use authmiddleware to check if token is valid and get claims
	// using the authcheck function
	claims, ok := Authcheck(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	// now prepare to send the request to the chaincode
	args := chaincode.BuildChaincodeArgs([]string{email}, "GetUserTickets")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("GetUserTickets function called", output)
	// Decode the output
	// ...............................


	// now return success, and transport details
	c.JSON(http.StatusOK, gin.H{
		"message":        "User tickets fetched successfully",
		"tickets":        "List of tickets",
		"transaction_id": "tx_10101",
	})
}



func AddTransport(c *gin.Context) {
	/*
	Input: Auth token, Source, Destination, Date, Price, ticketcount
	To create a transport on the ledger by provider
	Output: success message, Transport ID
	*/
	// using the authcheck function
	claims, ok := Authcheck(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	// get Source, Destination, Date, Price, ticketcount from the request body
	// var Source string
	// var Destination string
	// var Date string
	// var Price int
	// var ticketcount int
	type AddTransportRequest struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		Date        string `json:"date"`
		Price       int    `json:"price"`
		TicketCount  int    `json:"ticketcount"`
	}
	var Thisreq AddTransportRequest
	if err := c.ShouldBindJSON(&Thisreq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Source := Thisreq.Source
	Destination := Thisreq.Destination
	Date := Thisreq.Date
	Price := Thisreq.Price
	ticketcount := Thisreq.TicketCount

	
	// generate a transportID unique for the transport
	var transportID string
	// using the uuid package
	transportID = uuid.New().String()


	// now prepare to send the request to the chaincode with the transportID
	args := chaincode.BuildChaincodeArgs([]string{transportID, email, Source, Destination, Date, string(Price), string(ticketcount)}, "AddTransport")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("AddTransport function called", output)
	// Decode the output
	// ...............................


	// now return success, and transportID
	c.JSON(http.StatusOK, gin.H{
		"message":      "Transport added successfully",
		"transport_id": transportID,
		"transaction_id": "tx_10101",
	})
}
	


func GetTransportStatus(c *gin.Context) {
	/*
	Input: Auth token, transportID
	Will return the current status of the transport, vacancy, and Net Income of each date
	Output: success message, transport details
	*/

	claims, ok := Authcheck(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	email, ok := claims["email"].(string)
	log.Println("GetTransportStatus function called", email)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	// get transportID from the request body
	var transportID string
	if err := c.ShouldBindJSON(&transportID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// now prepare to send the request to the chaincode
	args := chaincode.BuildChaincodeArgs([]string{transportID}, "GetTransportStatus")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("GetTransportStatus function called", output)
	// Decode the output
	// ...............................

	// now return success, and transport details
	c.JSON(http.StatusOK, gin.H{
		"message":        "Transport details fetched successfully",
		"Income":         "10000",
		"vacancy":        "10",
	})

}