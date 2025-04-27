package chaincode

import (
	"encoding/json"
	"fmt"
	// "math"
	"slices"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *SmartContract) CancelTicket(ctx contractapi.TransactionContextInterface, userID, ticketID string) error {
	/*get the user's details from the hyperledger*/
	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("error occured while fetching the users details %s", err)
	}
	if userJSON == nil {
		return fmt.Errorf("the user %s does not exist", userID)
	}
	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return fmt.Errorf("error occured in pointing user")
	}

	/*get the ticket details which the user wants to delete*/
	ticketJSON, err := ctx.GetStub().GetState(ticketID)
	if err != nil {
		return fmt.Errorf("error while fetching the ticket details: %s", err)
	}
	if ticketJSON == nil {
		return fmt.Errorf("the ticket with id: %s does not exist", ticketID)
	}

	var ticket TicketDetails
	err = json.Unmarshal(ticketJSON, &ticket)
	if err != nil {
		return fmt.Errorf("error in pointing the ticket details")
	}
	
	/*only the owner can cancel the ticket*/
	if ticket.UserID != userID {
		return fmt.Errorf("error: user %s is not the owner of ticket %s", userID, ticketID)
	}


	departureDateTimeStr := fmt.Sprintf("%s %s", ticket.DateofTravel, ticket.DepartureTime)
	departureDateTime, err := time.Parse("2006-01-02 15:04", departureDateTimeStr)
	if err != nil {
		return fmt.Errorf("failed to parse combined departure date and time: %s", err)
	}
	timeNow := time.Now()
	if timeNow.After(departureDateTime) {
		return fmt.Errorf("error: you cannot delete the ticket after the journey has started")
	}

	transportID := ticket.TransportID
	transportJSON, err := ctx.GetStub().GetState(transportID)
	if err != nil {
		return fmt.Errorf("error while fetching the transport details")
	}
	var transport TransportDetails
	err = json.Unmarshal(transportJSON, &transport)
	if err != nil {
		return fmt.Errorf("error in pointing to the transport details")
	}
	var flag = false

	for i, tid := range user.Travels {
		if tid == ticketID {
			user.Travels = append(user.Travels[:i], user.Travels[i+1:]...)
			flag = true
			break
		}
	}

	if(!flag){
		return fmt.Errorf("error: the ticket is not present in the user's upcoming travels list")
	}
	
	// transport.Travellers[ticket.DateofTravel][ticket.SeatNumber-1] = "" /*removing the userID from travellers list*/

	//refund process
	travelDate, _ := time.Parse("2006-01-02", ticket.DateofTravel)
    hoursUntil := time.Until(travelDate).Hours()
    penaltyRate := 0.10 /*default penalty of 10% of the ticket charge*/

	switch{
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
	refundPrice := ticket.Price*(1-penaltyRate)

	// ProviderToUserPayment(ctx,transport.ProviderID,userID,refundPrice)
	///////////////////////////////////////////////
	
	providerJSON, err := ctx.GetStub().GetState(transport.ProviderID)
	if err != nil {
		return fmt.Errorf("error %s occured", err)
	}
	if providerJSON == nil {
		return fmt.Errorf("error: provider %s doesn't exist", transport.ProviderID)
	}
	var provider Provider
	err = json.Unmarshal(providerJSON, &provider)
	if err != nil {
		return fmt.Errorf("error: failed to unmarshal provider %s", transport.ProviderID)
	}

	user.BankBalance += refundPrice
	provider.BankBalance -= refundPrice

	paymentID := "payment-" + time.Now().Format("2006-01-02")
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
		return fmt.Errorf("error occured after deleting the ticket from user upcoming travels")
	}

	err = ctx.GetStub().PutState(userID, updatedUserJSON)
	if err != nil {
		return fmt.Errorf("error in uoadting the user in hyperledger")
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
		return fmt.Errorf("error while creating JSON format of the updated transport details")
	}
	err = ctx.GetStub().PutState(transportID, updatedTransportJSON)
	if err != nil {
		return fmt.Errorf("error while updating the data of transport in Hyperledger")
	}

	return ctx.GetStub().DelState(ticketID)
}

