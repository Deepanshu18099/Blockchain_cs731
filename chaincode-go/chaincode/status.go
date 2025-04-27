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
	_ = json.Unmarshal(userJSON, &user)

	/*
	type User struct {
		Name            string         `json:"Name"`
		Email           string         `json:"Email"`
		Phone           string         `json:"Phone"`
		PastTravels     []string       `json:"PastTravels"`
		UpcomingTravels []string       `json:"UpcomingTravels"`  //list of ticketIDs
		BankBalance     float64        `json:"BankBalance"`
		IsAnonymous     bool           `json:"IsAnonymous"`
		PaymentID       []string       `json:"PaymentID"`
	}
	*/
	return User{
		Email: "",
		Name:  "",
		// Plus other safe fields if you want
		// PastTravels:     user.PastTravels,
		// UpcomingTravels: user.UpcomingTravels,
		BankBalance: 	 user.BankBalance,
		PaymentID:       user.PaymentID,
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

	/*
	type Provider struct {
		Name             string             `json:"Name"`
		Email            string             `json:"Email"`
		Phone            string             `json:"Phone"`
		Services         []string           `json:"Services"` // list of transportIDs
		Verified         bool               `json:"Verified"`
		BankBalance      float64            `json:"BankBalance"`
		PaymentID        []string            `json:"PaymentID"`
	}
	*/
	_ = json.Unmarshal(providerJSON, &provider)
	
	return Provider{
		Services: provider.Services,
		BankBalance: provider.BankBalance,
		PaymentID: provider.PaymentID,
	},nil
}

func (s *SmartContract) GetDetailTicket(ctx contractapi.TransactionContextInterface, ticketID string) (TicketDetails, error) {
	ticketJSON,_ := ctx.GetStub().GetState(ticketID)
	var ticket TicketDetails
	_ = json.Unmarshal(ticketJSON,&ticket)
	/*
	type TicketDetails struct {
		TicketID               string        `json:"TicketID"`
		UserID                 string        `json:"OwnerEmail"`
		DateofTravel           string        `json:"DateofTravel"`
		Source                 string        `json:"Source"`
		Destination            string        `json:"Destination"`
		ModeofTravel           string        `json:"ModeofTravel"`
		TransportID            string        `json:"TransportID"` 
		SeatNumber             int32         `json:"SeatNumber"`
		Price                  float64       `json:"Price"`
		DateofBooking          string        `json:"DateofBooking"`
		DateofUpdate           string        `json:"DateofUpdate"`
		PaymentStatus          bool          `json:"PaymentStatus"` 
		IsActive               bool          `json:"IsActive"`
		DepartureTime          string        `json:"DepartureTime"`
		ArrivalTime            string        `json:"ArrivalTime"`
		JourneyDuration        string        `json:"JourneyDuration"`    
		Status                 string        `json:"Status"`
	}		
	*/
	return TicketDetails{
		TicketID: ticket.TicketID,
		DateofTravel: ticket.DateofTravel,
		Source: ticket.Source,
		Destination: ticket.Destination,
		ModeofTravel: ticket.ModeofTravel,
		TransportID: ticket.TransportID,
		SeatNumber: ticket.SeatNumber,
		Price: ticket.Price,
		ArrivalTime: ticket.ArrivalTime,
		DepartureTime: ticket.DepartureTime,
		JourneyDuration: ticket.JourneyDuration,
		DateofBooking: ticket.DateofBooking,
		DateofUpdate: ticket.DateofUpdate,
		PaymentStatus: ticket.PaymentStatus,
		IsActive: ticket.IsActive,
		Status: ticket.Status,
	},nil
}

//##########################################################################################################
//<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
