package controllers

import (
	"deepanshu18099/blockchain_ledger_backend/chaincode"
	"deepanshu18099/blockchain_ledger_backend/utils"
	"log"
	"net/http"

	// int to str
	"strconv"

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
	log.Println("AddMoneyToUser function called")
	claims, ok := utils.Authcheck(c)
	log.Println("AddMoneyToUser function called", claims, ok)

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
	// get amount from the request body
	type AddMoneyRequest struct {
		Amount int `json:"amount"`
	}
	var Thisreq AddMoneyRequest
	if err := c.ShouldBindJSON(&Thisreq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	argss := []string{}
	argss = append(argss, email)
	argss = append(argss, strconv.Itoa(Thisreq.Amount))

	// now prepare to send the request to the chaincode
	// Call the chaincode function to create the user on the ledger
	args := chaincode.BuildChaincodeArgs(argss, "AddBalance")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("AddMoneyToUser function called", output)

	// Decode the output
	result := utils.Cleancode2(c, output)

	if result == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode output"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Money added to user",
		"updatedbalance": result["BankBalance"],
		"transaction_id": result["transaction_id"],
	})

}

func GetTransports(c *gin.Context) {
	/*
		Input: Auth token, Source, Destination, Date, Mode of transport
		will be called by user
		Output: List of transports available
	*/

	claims, ok := utils.Authcheck(c)

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
	source := c.Param("source")
	destination := c.Param("destination")
	date := c.Param("date")
	mode := c.Param("mode")

	log.Println("GetTransports function called", source, destination, date, mode)

	// // now prepare to send the request to the chaincode
	args := chaincode.BuildChaincodeArgs([]string{source, destination, date, mode}, "GetAvailableTransports")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// log.Println("GetTransports function called", output)
	// Decode the output
	result := utils.Cleancode2(c, output)
	log.Println("GetTransports done", result)
	if result == nil {
		log.Println("GetTransports function called", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode output"})
		return
	}

	// now return success, and transport details
	c.JSON(http.StatusOK, gin.H{
		"message":        "Transport details fetched successfully",
		"transports":     result["availableTransports"],
		"transaction_id": result["transactionID"],
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
	claims, ok := utils.Authcheck(c)
	// get transportID, date from the request body
	// var transportID string
	// var date string
	// var seatNumber int

	type BookTicketRequest struct {
		TransportID string `json:"transportID"`
		Date        string `json:"date"`
		SeatNumber  string `json:"seatnumber"`
	}
	var Thisreq BookTicketRequest

	// log.Println("contents of c are", c)

	if err := c.ShouldBindJSON(&Thisreq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("BookTicket function called", Thisreq)

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
	args := chaincode.BuildChaincodeArgs([]string{email, Thisreq.TransportID, Thisreq.Date, Thisreq.SeatNumber}, "BookTicket")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("BookTicket function called", output)
	// Decode the output
	// ........................
	clean_output := utils.Cleancode2(c, output)
	if clean_output == nil {
		log.Println("BookTicket function called", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode output"})
		return
	}
	log.Println("BookTicket function called", clean_output)

	// now return success, and updated balance of the user
	c.JSON(http.StatusOK, gin.H{
		"message":        "Ticket booked successfully",
		"updatedbalance": clean_output["BankBalance"],
		"transaction_id": clean_output["transactionID"],
	})

}

func AddTransport(c *gin.Context) {
	/*
		Input: Auth token, Source, Destination, Date, Price, ticketcount
		To create a transport on the ledger by provider
		Output: success message, Transport ID
	*/
	// using the authcheck function
	claims, ok := utils.Authcheck(c)
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

	log.Println("AddTransport function called", email)
	// get Source, Destination, Date, Price, ticketcount from the request body
	// var Source string
	// var Destination string
	// var Date string
	// var Price int
	// var ticketcount int
	type AddTransportRequest struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		Mode        string `json:"mode"`
		StartDate   string `json:"startdate"`
		EndDate     string `json:"enddate"`
		Price       string `json:"price"`
		SeatCount   string `json:"seatcount"`
	}
	log.Println("AddTransport function called", email)
	var Thisreq AddTransportRequest
	if err := c.ShouldBindJSON(&Thisreq); err != nil {
		log.Println("AddTransport function called", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("AddTransport function called", Thisreq)
	Source := Thisreq.Source
	Destination := Thisreq.Destination
	Mode := Thisreq.Mode
	StartDate := Thisreq.StartDate
	EndDate := Thisreq.EndDate
	Price := Thisreq.Price
	SeatCount := Thisreq.SeatCount

	// generate a transportID unique for the transport

	// using the uuid package
	transportID := uuid.New().String()
	log.Println("AddTransport function called", transportID)
	log.Println("AddTransport function called", email, Source, Destination, StartDate, EndDate, Price, SeatCount, Mode)

	// set time to 10 am to 4 pm
	sample_dep_time := "10:00"
	sample_arrival_time := "16:00"
	total_time := "6h"
	// now prepare to send the request to the chaincode with the transportID
	args := chaincode.BuildChaincodeArgs([]string{email, transportID, Source, Destination, sample_dep_time, sample_arrival_time, total_time, Mode, SeatCount, Price, StartDate, EndDate}, "AddTransportService")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("AddTransport function called", output)

	// Decode the output
	clean_output := utils.Cleancode2(c, output)
	if clean_output == nil {
		log.Println("AddTransport function called", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode output"})
		return
	}
	log.Println("AddTransport function called", clean_output)
	// check if the output is success
	log.Println("AddTraction called", clean_output["transportID"])
	// check if the transportID is present in the output
	if clean_output["transportID"] == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add transport"})
		return
	}

	// now return success, and transportID
	c.JSON(http.StatusOK, gin.H{
		"message":        "Transport added successfully",
		"transport_id":   transportID,
		"transaction_id": clean_output["transaction_id"],
	})
}

func GetTransportStatus(c *gin.Context) {
	/*
		Input: Auth token, transportID
		Will return the current status of the transport, vacancy, and Net Income of each date
		Output: success message, transport details
	*/

	claims, ok := utils.Authcheck(c)
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
		"message": "Transport details fetched successfully",
		"Income":  "10000",
		"vacancy": "10",
	})

}

/*now in home, I want option in navbar to view previous bookings, which will result in new page listing user tickets, from ticket id in "userbookings" from this api response

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    const apiurl = process.env.REACT_APP_API_URL;
    try {
      const response = await axios.post(${apiurl}ledger/login, { email, password });

      if (response.status !== 200) {
        setErrors({ api: "Error signing in" });
        return;
      }
      console.log(response)
      const { token, userid, role, balance } = response.data

      login(token, userid, role, balance)

      navigate("/home");
    } catch (error) {
      console.error("Error signing in:", error);
      setErrors({ api: "Invalid credentials" });
    } finally {
      setLoading(false);
    }
  };

for each ticket in the list call GetDetailTicket api to get tickets in this form
		TicketID:        ticket.TicketID,
		DateofTravel:    ticket.DateofTravel,
		Source:          ticket.Source,
		Destination:     ticket.Destination,
		ModeofTravel:    ticket.ModeofTravel,
		TransportID:     ticket.TransportID,
		SeatNumber:      ticket.SeatNumber,
		Price:           ticket.Price,
		ArrivalTime:     ticket.ArrivalTime,
		DepartureTime:   ticket.DepartureTime,
		JourneyDuration: ticket.JourneyDuration,
		DateofBooking:   ticket.DateofBooking,
		DateofUpdate:    ticket.DateofUpdate,
		PaymentStatus:   ticket.PaymentStatus,
		IsActive:        ticket.IsActive,
		Status:          ticket.Status,
*/

func GetTickets(c *gin.Context) {
	/*
		Input: Token
		GIVE List of tickets by that user if user, else listings by provider
		Output: List of tickets(IDS) with details
	*/

	// Get the user ID from the token
	claims, ok := utils.Authcheck(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email, ok := claims["email"].(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	log.Println("GetTickets function called", email)
	// now prepare to send the request to the chaincode
	args := chaincode.BuildChaincodeArgs([]string{email}, "GetDetailUser")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("GetTickets function called", output)
	// Decode the output
	result := utils.Cleancode2(c, output)
	if result == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode output"})
		return
	}
	log.Println("GetTickets function called", result)

	// now return success, and transport details
	c.JSON(http.StatusOK, gin.H{
		"message":        "User tickets fetched successfully",
		"tickets":        result["Travels"],
		"transaction_id": result["transactionID"],
	})
}

func GetAllTransports(c *gin.Context) {
	/*
		Input: Auth token
		Output: List of all transports available of provider
	*/
	// use authmiddleware to check if token is valid and get claims
	// using the authcheck function
	claims, ok := utils.Authcheck(c)
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
	log.Println("GetAllTransports function called", email)
	// now prepare to send the request to the chaincode
	args := chaincode.BuildChaincodeArgs([]string{email}, "GetAllTransports")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("GetAllTransports function called", output)
	// Decode the output
	result := utils.Cleancode2(c, output)
	if result == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode output"})
		return
	}
	log.Println("GetAllTransports function called", result)

	c.JSON(http.StatusOK, gin.H{
		"message":        "All transports fetched successfully",
		"transports":     result["transports"],
		"transaction_id": result["transactionID"],
	})

}

func GetTicketByID(c *gin.Context) {
	/*
		Input: Token
		GIVE List of tickets by that user if user, else listings by provider
		Output: List of tickets(IDS) with details
	*/
	// Get the user ID from the token
	claims, ok := utils.Authcheck(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email, ok := claims["email"].(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	// get the ticket id from the url
	ticketID := c.Param("id")
	log.Println("GetTicketByID function called", email, ticketID)
	// now prepare to send the request to the chaincode
	args := chaincode.BuildChaincodeArgs([]string{ticketID}, "GetDetailTicket")
	output, err := chaincode.RunPeerCommand(args)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("GetTicketByID function called", output)
	// Decode the output
	result := utils.Cleancode2(c, output)
	if result == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode output"})
		return
	}
	log.Println("GetTicketByID function called", result)
	// now return success, and transport details
	c.JSON(http.StatusOK, gin.H{
		"message": "User tickets fetched successfully",
		"ticket":  result,
	})
}

func DeleteTicket(c *gin.Context) {
	/*
		Input: Token
		GIVE id of ticket to be deleted
		Output: success message
	*/
	// Get the user ID from the token
	claims, ok := utils.Authcheck(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	// get the ticket id from the url
	ticketID := c.Param("id")
	log.Println("DeleteTicket function called", email, ticketID)
	// now prepare to send the request to the chaincode, send email and ticketID
	args := chaincode.BuildChaincodeArgs([]string{email, ticketID}, "CancelTicket")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("DeleteTicket function called", output)
	// Decode the output
	result := utils.Cleancode2(c, output)
	if result == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode output"})
		return
	}
	log.Println("DeleteTicket function called", result)
	// now return success, and transport details
	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket deleted successfully",
	})
}

func UpdateTicket(c *gin.Context) {
	/*
		Input: Token, new params of ticket
		Will update the params given to new state
		Output: success message
	*/
	// Get the user ID from the token
	claims, ok := utils.Authcheck(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	// get the ticket id from the url
	ticketID := c.Param("id")

	// newTicketDate, ok := claims["newTicketDate"].(string)
	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "TicketDate error"})
	// 	return
	// }

	// newSeatNum, ok := claims["newSeatNum"].(string)
	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "SeatNum error"})
	// 	return
	// }

	type UpdateTicketRequest struct {
		DateofTravel string `json:"DateofTravel"`
		SeatNumber   string `json:"SeatNumber"`
	}

	var Thisreq UpdateTicketRequest
	if err := c.ShouldBindJSON(&Thisreq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTicketDate := Thisreq.DateofTravel
	newSeatNum := Thisreq.SeatNumber

	log.Println("UpdateTicket function called", email, ticketID)
	// now prepare to send the request to the chaincode, send email and ticketID
	args := chaincode.BuildChaincodeArgs([]string{ticketID, newTicketDate, newSeatNum}, "UpdateTicket")
	output, err := chaincode.RunPeerCommand(args)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("UpdateTicket function called", string(output))

	result := utils.Cleancode2(c, output)
	if result == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode output"})
		return
	}
	log.Println("UpdateTicket function called", result)

	// Decode the output
	// now return success, and transport details
	c.JSON(http.StatusOK, gin.H{
		"message":     "Ticket updated successfully",
		"NewTicketID": result["NewTicketID"],
	})
}
