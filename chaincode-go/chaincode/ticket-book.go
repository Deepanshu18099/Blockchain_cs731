package chaincode

import (
	"encoding/json"
	"fmt"

	// "math"
	// "slices"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *SmartContract) BookTicket(ctx contractapi.TransactionContextInterface, userID, transportID, date string, seatNumber int32) (string, error) {
	/*Ticket in the past can not be booked*/
	bookingDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("error: invalid date format, expected YYYY-MM-DD")
	}
	currentDate := time.Now()
	if bookingDate.Before(currentDate.Truncate(24 * time.Hour)) {
		return "", fmt.Errorf("error: cannot book ticket for a past date: %s", date)
	}

	/*details of the user*/
	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return "", fmt.Errorf("error: failed to read user data: %v", err)
	}
	if userJSON == nil {
		return "", fmt.Errorf("error: user %s does not exist", userID)
	}
	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return "", fmt.Errorf("error: failed to parse user data: %v", err)
	}

	/*transport details from the given transport ID*/
	transportJSON, err := ctx.GetStub().GetState(transportID)
	if err != nil {
		return "", fmt.Errorf("error: failed to read transport data: %v", err)
	}
	if transportJSON == nil {
		return "", fmt.Errorf("error: transport service %s does not exist", transportID)
	}
	var transport TransportDetails
	err = json.Unmarshal(transportJSON, &transport)
	if err != nil {
		return "", fmt.Errorf("error: failed to parse transport data: %v", err)
	}

	/*if the transport does not have any seats then we can not overbook*/
	seats, exists := transport.SeatMap[date]
	if !exists || len(seats) == 0 {
		return "", fmt.Errorf("error: no available seats on transport %s for date %s", transportID, date)
	}

	/*the transportation stystem used dynamic price*/
	currentPrice, err := calculateDynamicPrice(ctx, transportID, date)
	if err != nil {
		return "", fmt.Errorf("error while finding the price %s", err)
	}

	/*checking if the user has sufficient balance to book the ticket*/
	if user.BankBalance < currentPrice {
		return "Fail", fmt.Errorf("error insufficient balance! required: %.2f, available: %.2f", currentPrice, user.BankBalance)
	}

	var flag = true
	for i, value := range transport.SeatMap[date] {
		if value == seatNumber {
			transport.SeatMap[date] = append(transport.SeatMap[date][:i], transport.SeatMap[date][i+1:]...)
			flag = false
			break
		}
	}
	// transport.Travellers[date][seatNumber-1] = userID /*maintaining the userID of the travellers*/
	if flag {
		return "error: ticket can not be booked", fmt.Errorf("the seat is already booked")
	}


	providerID := transport.ProviderID

	providerJSON, err := ctx.GetStub().GetState(providerID)
	if err != nil {
		return "Fail", fmt.Errorf("error %s occured", err)
	}
	if providerJSON == nil {
		return "Fail", fmt.Errorf("error: provider %s doesn't exist", providerID)
	}
	var provider Provider
	err = json.Unmarshal(providerJSON, &provider)
	if err != nil {
		return "Fail", fmt.Errorf("error: failed to unmarshal provider %s", providerID)
	}

	// err = UserToProviderPayment(ctx, userID, providerID, currentPrice)

	///////////////////////////////////////////////
	user.BankBalance -= currentPrice
	provider.BankBalance += currentPrice

	paymentID := "payment-" + time.Now().Format("2006-01-02 15:04:05")

	payment := PaymentDetail{
		PaymentID:   paymentID,
		From:        userID,
		To:          providerID,
		Amount:      currentPrice,
		PaymentTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	paymentJSON, _ := json.Marshal(payment)
	ctx.GetStub().PutState(paymentID, paymentJSON)

	user.PaymentID = append(user.PaymentID, paymentID)
	provider.PaymentID = append(provider.PaymentID, paymentID)

	updatedProviderJSON, _ := json.Marshal(provider)
	ctx.GetStub().PutState(providerID, updatedProviderJSON)
	///////////////////////////////

	ticketID := fmt.Sprintf("%s-%s-%d", transportID, date, seatNumber)

	ticket := TicketDetails{
		TicketID:        ticketID,
		UserID:          userID,
		ProviderID:      providerID,
		DateofTravel:    date,
		Source:          transport.Source,
		Destination:     transport.Destination,
		ModeofTravel:    transport.ModeofTravel,
		TransportID:     transportID, // Using TicketID as PNR
		SeatNumber:      seatNumber,
		Price:           currentPrice,
		DateofBooking:   time.Now().Format("2006-01-02"), // Use actual timestamp in real implementation
		PaymentStatus:   true,                                     // Payment is verified
		IsActive:        true,
		DepartureTime:   transport.DepartureTime,
		ArrivalTime:     transport.ArrivalTime,
		JourneyDuration: transport.JourneyDuration,
	}

	user.Travels = append(user.Travels, ticketID)

	updatedUserJSON, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated user data: %v", err)
	}
	err = ctx.GetStub().PutState(userID, updatedUserJSON)
	if err != nil {
		return "", fmt.Errorf("failed to update user state: %v", err)
	}

	// updatedUserJSON,_ := json.Marshal(user)
	// ctx.GetStub().PutState(userID, updatedUserJSON)

	updatedTransportJSON, err := json.Marshal(transport)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated transport data: %v", err)
	}
	ctx.GetStub().PutState(transportID, updatedTransportJSON)

	ticketJSON,_ := json.Marshal(ticket)
	ctx.GetStub().PutState(ticketID, ticketJSON)
	
	return fmt.Sprintf("ticket booked successfully! Ticket ID: %s", ticketID), nil
}

