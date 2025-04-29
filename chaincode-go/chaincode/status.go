package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// func (s*SmartContract)GetDetailUser(ctx contractapi.TransactionContextInterface, userID string)(User,error){
// 	userJSON,err := ctx.GetStub().GetState(userID)
// 	if err != nil{
// 		return User{}, fmt.Errorf("error occured while fetching the users details %s",err)
// 	}
// 	if userJSON == nil{
// 		return User{},fmt.Errorf("the user %s does not exist",userID)
// 	}
// 	var user User

// 	err = json.Unmarshal(userJSON, &user)
// 	if err != nil{
// 		return user, fmt.Errorf("error while pointing to the user details")
// 	}
// 	return user,nil
// }

func (s *SmartContract) GetDetailUser(ctx contractapi.TransactionContextInterface, userID string) (User, error) {
	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return User{}, fmt.Errorf("error occured while fetching the users details %s", err)
	}
	if userJSON == nil {
		return User{}, fmt.Errorf("the user %s does not exist", userID)
	}
	var user User
	json.Unmarshal(userJSON, &user)

	return User{
		Email: "",
		Name:  "",
		Travels:     user.Travels,
		BankBalance: user.BankBalance,
		PaymentID:   user.PaymentID,
	}, nil
}



func (s *SmartContract) GetDetailProvider(ctx contractapi.TransactionContextInterface, providerID string) (Provider, error) {
	providerJSON, err := ctx.GetStub().GetState(providerID)
	if err != nil {
		return Provider{}, fmt.Errorf("error occured while fetching the users details %s", err)
	}
	if providerJSON == nil {
		return Provider{}, fmt.Errorf("the provider %s does not exist", providerID)
	}
	var provider Provider

	json.Unmarshal(providerJSON, &provider)

	return Provider{
		Services:    provider.Services,
		BankBalance: provider.BankBalance,
		PaymentID:   provider.PaymentID,
	}, nil
}

func (s *SmartContract) GetDetailTicket(ctx contractapi.TransactionContextInterface, ticketID string) (TicketDetails, error) {
	ticketJSON, _ := ctx.GetStub().GetState(ticketID)
	if ticketJSON == nil {
		return TicketDetails{}, fmt.Errorf("error: ticket does not exist")
	}
	var ticket TicketDetails
	json.Unmarshal(ticketJSON, &ticket)

	return TicketDetails{
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
	}, nil
}

//##########################################################################################################
//<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
