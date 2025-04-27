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
	
	bookingDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("error: invalid date format, expected YYYY-MM-DD")
	}

	currentDate := time.Now()
	if bookingDate.Before(currentDate.Truncate(24 * time.Hour)) {
		return "", fmt.Errorf("error: cannot book ticket for a past date: %s", date)
	}

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

	// if transport.AvailableSeats <= 0 {
	// 	return "", fmt.Errorf("no available seats on transport %s", transportID)
	// }

	seats, exists := transport.SeatMap[date]
	if !exists || len(seats) == 0 {
		return "", fmt.Errorf("error: no available seats on transport %s for date %s", transportID, date)
	}

	currentPrice, err := calculateDynamicPrice(ctx, transportID, date)
	if err != nil {
		return "", fmt.Errorf("error while finding the price %s", err)
	}

	if user.BankBalance < currentPrice {
		return "", fmt.Errorf("error insufficient balance! required: %.2f, available: %.2f", currentPrice, user.BankBalance)
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

	// user.BankBalance -= currentPrice
	providerID := transport.ProviderID

	s.UserToProviderPayment(ctx,userID,providerID, currentPrice)

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
		DateofBooking:   time.Now().Format("2006-01-02 15:04:05"), // Use actual timestamp in real implementation
		PaymentStatus:   true,                                     // Payment is verified
		IsActive:        true,
		DepartureTime:   transport.DepartureTime,
		ArrivalTime:     transport.ArrivalTime,
		JourneyDuration: transport.JourneyDuration,
	}

	user.UpcomingTravels = append(user.UpcomingTravels, ticketID)
	updatedUserJSON, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated user data: %v", err)
	}
	err = ctx.GetStub().PutState(userID, updatedUserJSON)
	if err != nil {
		return "", fmt.Errorf("failed to update user state: %v", err)
	}

	updatedTransportJSON, err := json.Marshal(transport)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated transport data: %v", err)
	}
	err = ctx.GetStub().PutState(transportID, updatedTransportJSON)

	if err != nil {
		return "", fmt.Errorf("failed to update transport state: %v", err)
	}

	ticketJSON, err := json.Marshal(ticket)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ticket data: %v", err)
	}
	err = ctx.GetStub().PutState(ticketID, ticketJSON)
	if err != nil {
		return "", fmt.Errorf("failed to save ticket: %v", err)
	}
	return fmt.Sprintf("ticket booked successfully! Ticket ID: %s", ticketID), nil
}