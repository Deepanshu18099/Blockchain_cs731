package chaincode

import (
	"encoding/json"
	"fmt"

	// "math"
	"slices"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *SmartContract) CancelTicket(ctx contractapi.TransactionContextInterface, userID, ticketID string) ([]byte, error) {
	/*get the user's details from the hyperledger*/
	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return nil, fmt.Errorf("error while fetching the user details: %s", err)
	}
	if userJSON == nil {
		return nil, fmt.Errorf("the user with id: %s does not exist", userID)
	}
	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, fmt.Errorf("error in pointing to the user details")
	}

	/*get the ticket details which the user wants to delete*/
	ticketJSON, err := ctx.GetStub().GetState(ticketID)
	if err != nil {
		return nil, fmt.Errorf("error while fetching the ticket details: %s", err)
	}
	if ticketJSON == nil {
		return nil, fmt.Errorf("the ticket with id: %s does not exist", ticketID)
	}

	var ticket TicketDetails
	err = json.Unmarshal(ticketJSON, &ticket)
	if err != nil {
		return nil, fmt.Errorf("error in pointing to the ticket details")
	}

	/*only the owner can cancel the ticket*/
	if ticket.UserID != userID {
		return nil, fmt.Errorf("error: you are not the owner of this ticket")
	}

	departureDateTimeStr := fmt.Sprintf("%s %s", ticket.DateofTravel, ticket.DepartureTime)
	departureDateTime, err := time.Parse("2006-01-02 15:04", departureDateTimeStr)
	if err != nil {
		return nil, fmt.Errorf("error while parsing the departure date and time: %s", err)
	}
	timeNow := time.Now()
	if timeNow.After(departureDateTime) {
		return nil, fmt.Errorf("error: the ticket has already departed")
	}

	transportID := ticket.TransportID
	transportJSON, err := ctx.GetStub().GetState(transportID)
	if err != nil {
		return nil, fmt.Errorf("error while fetching the transport details: %s", err)
	}
	var transport TransportDetails
	err = json.Unmarshal(transportJSON, &transport)
	if err != nil {
		return nil, fmt.Errorf("error in pointing to the transport details")
	}
	var flag = false

	for i, tid := range user.Travels {
		if tid == ticketID {
			user.Travels = append(user.Travels[:i], user.Travels[i+1:]...)
			flag = true
			break
		}
	}

	if !flag {
		return nil, fmt.Errorf("error: the ticket does not exist in the user's travel history")
	}

	// transport.Travellers[ticket.DateofTravel][ticket.SeatNumber-1] = "" /*removing the userID from travellers list*/

	//refund process
	travelDate, _ := time.Parse("2006-01-02", ticket.DateofTravel)
	hoursUntil := time.Until(travelDate).Hours()
	penaltyRate := 0.10 /*default penalty of 10% of the ticket charge*/

	switch {
	case hoursUntil <= 72:
		penaltyRate = 0.20
	case hoursUntil <= 48:
		penaltyRate = 0.30
	case hoursUntil <= 24:
		penaltyRate = 0.50
	case hoursUntil <= 12:
		penaltyRate = 0.75
	case hoursUntil <= 4:
		penaltyRate = 1.00
	default:
		penaltyRate = 0.10
	}

	// user.BankBalance += ticket.Price*(1-penaltyRate)
	refundPrice := ticket.Price * (1 - penaltyRate)

	// ProviderToUserPayment(ctx,transport.ProviderID,userID,refundPrice)
	///////////////////////////////////////////////

	providerJSON, err := ctx.GetStub().GetState(transport.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("error %s occured", err)
	}
	if providerJSON == nil {
		return nil, fmt.Errorf("error: provider %s doesn't exist", transport.ProviderID)
	}
	var provider Provider
	err = json.Unmarshal(providerJSON, &provider)
	if err != nil {
		return nil, fmt.Errorf("error: failed to unmarshal provider %s", transport.ProviderID)
	}

	user.BankBalance += refundPrice
	provider.BankBalance -= refundPrice

	paymentID := "payment-" + time.Now().Format("2006-01-02 15:04:05")
	payment := PaymentDetail{
		PaymentID:   paymentID,
		From:        userID,
		To:          transport.ProviderID,
		Amount:      refundPrice,
		PaymentTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	paymentJSON, _ := json.Marshal(payment)
	ctx.GetStub().PutState(paymentID, paymentJSON)

	user.PaymentID = append(user.PaymentID, paymentID)
	provider.PaymentID = append(provider.PaymentID, paymentID)

	updatedProviderJSON, _ := json.Marshal(provider)
	ctx.GetStub().PutState(transport.ProviderID, updatedProviderJSON)

	///////////////////////////////////////////////////////////////////////

	updatedUserJSON, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("error while updating the user details")
	}

	err = ctx.GetStub().PutState(userID, updatedUserJSON)
	if err != nil {
		return nil, fmt.Errorf("error while updating the user details in hyperledger")
	}

	seatNumber := ticket.SeatNumber
	date := ticket.DateofTravel
	//to check if the seat already present? it should not happen
	seatExists := false
	for _, seat := range transport.SeatMap[date] {
		if seat == seatNumber {
			seatExists = true
			break
		}
	}

	// transport.Travellers[date][seatNumber] = ""

	if !seatExists {
		transport.SeatMap[date] = append(transport.SeatMap[date], seatNumber)
		slices.Sort(transport.SeatMap[date]) // Sort to maintain order
	}

	//Waiting seats are not implemented now......, can be implemented
	updatedTransportJSON, err := json.Marshal(transport)
	if err != nil {
		return nil, fmt.Errorf("error while updating the transport details")
	}
	err = ctx.GetStub().PutState(transportID, updatedTransportJSON)
	if err != nil {
		return nil, fmt.Errorf("error while updating the transport details in hyperledger")
	}

	returnitem := map[string]interface{}{
		"message": "ticket cancelled successfully",
	}
	returnitemJSON, err := json.Marshal(returnitem)
	if err != nil {
		return nil, fmt.Errorf("error while creating JSON format of the return item")
	}
	err = ctx.GetStub().SetEvent("TicketCancelled", returnitemJSON)
	if err != nil {
		return nil, fmt.Errorf("error while setting event: %s", err)
	}
	return returnitemJSON, nil
}
