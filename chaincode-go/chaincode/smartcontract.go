package chaincode

import (
	"encoding/json"
	"fmt"
	// "time"
	// "math"
	// "strconv"
	// "strings"
	// "strconv"
	// "sort"
	// "math"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Structs are defined in struct.go file
// The below code is working fine, it restricts user to make duplicate account

// create user and provider in the register.go file

// to make a function private, we will have to rename it starting with lowercase
// Also if we have to make this function private then the call to
// this function will be made by an internal function with its userID only
// this is done for testing
// Or we can do it this way also


func (s *SmartContract) AddBalance(ctx contractapi.TransactionContextInterface, email string, amount float64) error {
	var username = email
	exists, err := s.detailExists(ctx, username)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the User %s does not exist", email)
	}

	detailJSON, err := ctx.GetStub().GetState(username)

	if err != nil {
		return fmt.Errorf("failed to read from worls state: %v", err)
	}
	if detailJSON == nil {
		return fmt.Errorf("details of %s does not exist", email)
	}
	var asset User
	err = json.Unmarshal(detailJSON, &asset)
	if err != nil {
		return fmt.Errorf("failed to marshal the data")
	}
	
	asset.BankBalance += amount

	updatedUserJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal updated user data: %v", err)
	}
	return ctx.GetStub().PutState(username, updatedUserJSON)
}

func (s *SmartContract) GetBalance(ctx contractapi.TransactionContextInterface, email string) (float64, error) {
	var username = email
	detailJSON, err := ctx.GetStub().GetState(username)
	if err != nil {
		return 0.0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if detailJSON == nil {
		return 0.0, fmt.Errorf("details of %s do not exist", email)
	}

	var user User
	err = json.Unmarshal(detailJSON, &user)
	if err != nil {
		return 0.0, fmt.Errorf("failed to unmarshal user data: %v", err)
	}
	return user.BankBalance, nil
}

func (s *SmartContract) MakeUserAnonymous(ctx contractapi.TransactionContextInterface, email string) error {
	var username = email

	detailJSON, err := ctx.GetStub().GetState(username)

	if err != nil {
		return fmt.Errorf("failed to read from worls state: %v", err)
	}
	if detailJSON == nil {
		return fmt.Errorf("details of %s does not exist", email)
	}

	var asset User
	err = json.Unmarshal(detailJSON, &asset)
	if err != nil {
		return fmt.Errorf("failed to marshal the data")
	}
	asset.IsAnonymous = true
	updatedUserJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal updated user data: %v", err)
	}
	return ctx.GetStub().PutState(username, updatedUserJSON)
}

//book ticket code in ticket.go file

/*
Errors:
ticket booking ka logic me thoda issue hai, usme travel date select nahi ho rhi hai
struct me jo arrays hai unme struct ki IDs rakh sakte hai.....pura struct rakhne ki need nahi hai

get transportation detail me saare date ki values nahi chahiye hame
*/

/*
func (s*SmartContract)UpdateTicket(ctx contractapi.TransactionContextInterface,id, date string, seat int32)error{
	ticketJSON,err := ctx.GetStub().GetState(id)
	if err != nil{
		return fmt.Errorf("error occured while fetching the ticket details")
	}
	if ticketJSON == nil{
		return fmt.Errorf("ticket doesn't exists with the ticket ID: %s", id)
	}

	var ticket TicketDetails

	err = json.Unmarshal(ticketJSON,&ticket)
	if err != nil{
		return 	fmt.Errorf("error occured while creating the ticket variable")
	}

	ticket.DateofTravel = date
	ticket.SeatNumber = seat

	updatedTicketJSON,err := json.Marshal(ticket)
	if err != nil{
		return fmt.Errorf("error in converting ticket to JSON format")
	}

	// logic to charge for the ticket update is not added yet

	return ctx.GetStub().PutState(id,updatedTicketJSON)
}

func (s*SmartContract)DeleteTicket(ctx contractapi.TransactionContextInterface, userID, ticketID string)error{
	userJSON,err := ctx.GetStub().GetState(userID)
	if err != nil{
		return fmt.Errorf("error occured while fetching the users details %s",err)
	}
	if userJSON == nil{
		return fmt.Errorf("the user %s does not exist",userID)
	}
	var user User

	err = json.Unmarshal(userJSON, &user)
	if err != nil{
		return fmt.Errorf("error occured in pointing user")
	}

	ticketJSON,err := ctx.GetStub().GetState(ticketID)
	if err != nil{
		return fmt.Errorf("error while fetching the ticket details: %s",err)
	}
	if ticketJSON == nil{
		return fmt.Errorf("the ticket with id: %s does not exist",ticketID)
	}

	var ticket TicketDetails

	err = json.Unmarshal(ticketJSON, &ticket)
	if err != nil{
		return fmt.Errorf("error in pointing the ticket details")
	}
	departureTime, err := time.Parse(time.RFC3339, ticket.DepartureTime)
	if err != nil {
		return fmt.Errorf("failed to parse departure time: %s", err)
	}

	timeNow := time.Now().UTC()

	if timeNow.After(departureTime){
		return fmt.Errorf("you cannot delete the ticket after the journey has started")
	}

	transportID := ticket.TransportID
	transportJSON,err := ctx.GetStub().GetState(transportID)
	if err != nil{
		return fmt.Errorf("error while fetching the transport details")
	}
	var transport TransportDetails
	err = json.Unmarshal(transportJSON,&transport)
	if err != nil{
		return fmt.Errorf("error in pointing to the transport details")
	}

	for i, tid := range user.UpcomingTravels {
		if tid == ticketID {
			user.UpcomingTravels = append(user.UpcomingTravels[:i], user.UpcomingTravels[i+1:]...)
			break
		}
	}

	updatedUserJSON,err := json.Marshal(user)
	if err != nil{
		return fmt.Errorf("error occured after deleting the ticket from user upcoming travels")
	}

	err = ctx.GetStub().PutState(userID,updatedUserJSON)
	if err != nil{
		return fmt.Errorf("error in uoadting the user in hyperledger")
	}

	transport.AvailableSeats += 1
	seatNumber := ticket.SeatNumber
	date := ticket.DateofTravel
	// transport.SeatsAvailable = append(seatNumber,transport.SeatsAvailable ...)
	transport.SeatMap[date] = append(transport.SeatMap[date], seatNumber)
	slices.Sort(transport.SeatMap[date])

	//Waiting seats are not implemented now......, can be implemented
	updatedTransportJSON,err := json.Marshal(transport)
	if err != nil{
		return fmt.Errorf("error while creating JSON format of the updated transport details")
	}
	err = ctx.GetStub().PutState(transportID,updatedTransportJSON)
	if err != nil{
		return fmt.Errorf("error while updating the data of transport in Hyperledger")
	}

	return ctx.GetStub().DelState(ticketID)
}


func (s *SmartContract) AddTransportService(ctx contractapi.TransactionContextInterface, email,id, source,dest,dept,arrt,totalt,mode string, cap int32,basep float64,dates []string) error {
	providerJSON,err := ctx.GetStub().GetState(email)
	if err != nil{
		return fmt.Errorf("error occured while getting details of the provider: %s", err)
	}
	if providerJSON == nil{
		return fmt.Errorf("the provider: %s does not exist",email)
	}

	var transportID = id+"-"+source+"-"+dest
	transportJSON,err := ctx.GetStub().GetState(transportID)
	if err != nil{
		return fmt.Errorf("error occured while getting details of the transportation: %s", err)
	}
	if transportJSON != nil{
		return fmt.Errorf("the transportation already exist")
	}

	arr := make([]int32, cap)

	for i := 0; i < int(cap); i++ {
		arr[i] = int32(i + 1)
	}

	seatMap := make(map[string][]int32)

	for _, date := range dates {
		// Important: Copy the seats slice, otherwise all map values will point to the same slice
		seatsCopy := make([]int32, len(arr))
		copy(seatsCopy, arr)
		sort.Slice(arr, func(i, j int) bool {
			return arr[i] < arr[j]
		})
		seatMap[date] = seatsCopy
	}

	newTransport := TransportDetails{
		ID: transportID,
		Source: source,
		Destination: dest,
		DepartureTime:dept,
		ArrivalTime: arrt,
		BasePrice: basep,
		Rating: 3.00,
		Capacity:cap,
		AvailableSeats: cap,
		ModeofTravel: mode,
		JourneyDuration: totalt,
		DateofTravel: dates,
		SeatMap: seatMap,
	}

	transpotJSON,err := json.Marshal(newTransport)
	if(err != nil){
		return fmt.Errorf("error occured %s", err)
	}

	err = ctx.GetStub().PutState(transportID, transpotJSON)
	if(err != nil){
		return fmt.Errorf("error occured while storing the transport details, %s",err)
	}

	var updateProvider Provider
	err = json.Unmarshal(providerJSON, &updateProvider)
	if err != nil{
		return fmt.Errorf("error in pointing to the updated provider data")
	}
	updateProvider.Services = append(updateProvider.Services, transportID)

	updateProviderJSON,err := json.Marshal(updateProvider)
	if err != nil{
		return fmt.Errorf("error in updating provider data")
	}
	return ctx.GetStub().PutState(email,updateProviderJSON)
}

func (s*SmartContract)GetAvailableTransports(ctx contractapi.TransactionContextInterface,src,dest,date string) ([]*TransportDetails,error){
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
		if err != nil {
			return nil, err
		}
		defer resultsIterator.Close()

		var availableTransports []*TransportDetails
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return nil, err
			}

			var transport TransportDetails
			err = json.Unmarshal(queryResponse.Value, &transport)
			if err != nil{
				return nil,fmt.Errorf("error in pointing the query responses")
			}
			if transport.Source == src && transport.Destination==dest{
				for _, availableDate := range transport.DateofTravel{
					if availableDate == date{
						requiredTransport := &TransportDetails{
							ID: transport.ID,
							Source: transport.Source,
							Destination: transport.Destination,
							DepartureTime: transport.DepartureTime,
							ArrivalTime: transport.ArrivalTime,
							BasePrice: transport.BasePrice,
							Rating: transport.Rating,
							Capacity: transport.Capacity,
							AvailableSeats: transport.AvailableSeats,
							ModeofTravel: transport.ModeofTravel,
							JourneyDuration: transport.JourneyDuration,
							DateofTravel: []string{date},
							SeatMap: map[string][]int32{date: transport.SeatMap[date]},
						}
						availableTransports = append(availableTransports,requiredTransport)
						break
					}
				}
			}
		}
	return availableTransports,nil;
}


func (t*SmartContract) detailExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	detailJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return detailJSON != nil, nil
}


func calculateDynamicPrice(ctx contractapi.TransactionContextInterface, transportID string) (float64, error) {
    transportJSON, err := ctx.GetStub().GetState(transportID)
    if err != nil || transportJSON == nil {
        return 0, fmt.Errorf("transportation not found")
    }
    var transport TransportDetails
    json.Unmarshal(transportJSON, &transport)

    // Calculate days until departure
    departureTime, _ := time.Parse(time.RFC3339, transport.DepartureTime)

    availablePercent := float64(transport.AvailableSeats) / float64(transport.Capacity) * 100
    seatFactor := 1.00

    switch{
    case availablePercent <= 30:
    	seatFactor = 1.10
    case availablePercent <=20:
    	seatFactor = 1.20
    case availablePercent <=10:
    	seatFactor = 1.40
    default:
    	seatFactor = 1.00

    }

    now := time.Now().UTC()
	timeToTravel := departureTime.Sub(now).Hours()

	timeFactor := 1.00

	switch {
	case timeToTravel <= 24:
		timeFactor = 1.40
	case timeToTravel <= 48:
		timeFactor = 1.20
	case timeToTravel <= 72:
		timeFactor = 1.10
	default:
		timeFactor = 1.00
	}

    currentPrice := transport.BasePrice * timeFactor * seatFactor
    return math.Round(currentPrice*100)/100, nil
}
*/
