package chaincode

import (
	"encoding/json"
	"fmt"
	// "math"
	"slices"
	"time"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

/*
sanity check, you can't book a ticket for a past date
get the ticket details
get the users detail
get the transport detail

ticket updating penalty should be imposed (first check the user's available balance)

user-->upcoming travels should be updated with the new ticketID
newticket ID should be given
newTicket ID should be updated in the user's upcoming travel
previous date's seat should be now available
new seat should be now allotted to the user and removed from the seat map
*/
func (s *SmartContract) UpdateTicket(ctx contractapi.TransactionContextInterface, ticketID, date string, newSeat int32) error {
	/*the updated date can't be in the past*/
	inputDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("error parsing input date: %v", err)
	}
	currDate := time.Now()
	if inputDate.Before(currDate) {
		return fmt.Errorf("error: the date is gone")
	}
	
	/*getting the details of the booked ticket*/
	ticketJSON, err := ctx.GetStub().GetState(ticketID)
	if err != nil {
		return fmt.Errorf("error occured while fetching the ticket details")
	}
	if ticketJSON == nil {
		return fmt.Errorf("ticket doesn't exists with the ticket ID: %s", ticketID)
	}
	var ticket TicketDetails
	err = json.Unmarshal(ticketJSON, &ticket)
	if err != nil {
		return fmt.Errorf("error occured while creating the ticket variable")
	}


	/*getting the details of the transport*/
	transportID := ticket.TransportID
	transportJSON, err := ctx.GetStub().GetState(transportID)
	if err != nil {
		return fmt.Errorf("error occured while fetching the transoirt details")
	}
	if transportJSON == nil {
		return fmt.Errorf("error the transport associated with this ticket does not exist")
	}
	var transport TransportDetails
	err = json.Unmarshal(transportJSON, &transport)
	if err != nil {
		return fmt.Errorf("error while pointing to the transport details")
	}

	/*getting the user's details to impose the penalty*/
	userID := ticket.UserID
	var user User
	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("error in fetching the user details")
	}
	if userJSON == nil {
		return fmt.Errorf("error: the user doesn't exist")
	}
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return fmt.Errorf("error in pointing to the user JSON")
	}

	/*calculating the penalty to be imposed*/
	previousDate := ticket.DateofTravel
	previousTicketPrice := ticket.Price
	newTicketPrice, _ := calculateDynamicPrice(ctx, transportID, date) /*as we have dynamic price so if the new date costs more than that is to be paid*/
	penaltyPrice,_ := calculateDynamicPrice(ctx,transportID,previousDate)

	penaltyPrice += newTicketPrice-previousTicketPrice /*extra to be paid according to the dynamic price*/

	if user.BankBalance < penaltyPrice {
		return fmt.Errorf("error: the user doesn't have sufficient balance to pay the penaly charge for ticket update")
	}
	
	/*removing the newSeat from the new travel date*/
	var flag = true
	for i, value := range transport.SeatMap[date] {
		if value == newSeat {
			transport.SeatMap[date] = append(transport.SeatMap[date][:i], transport.SeatMap[date][i+1:]...)
			flag = false
			break
		}
	}
	if flag {
		return fmt.Errorf("error: the seat is already booked")
	}
	
	// user.BankBalance -= penaltyPrice /*penalty imposed*/
	s.UserToProviderPayment(ctx,userID, transport.ProviderID,penaltyPrice)

	// previousSeat := ticket.SeatNumber
	transport.SeatMap[previousDate] = append(transport.SeatMap[previousDate], ticket.SeatNumber) /*previously booked seat is now available for booking*/
	slices.Sort(transport.SeatMap[previousDate])
	

	newTicketID := fmt.Sprintf("%s-%s-%d", ticket.TransportID, date, newSeat) /*new ticket ID id provided to the user*/

	for i, value := range user.UpcomingTravels {
		if value == ticketID {
			user.UpcomingTravels = append(user.UpcomingTravels[:i], user.UpcomingTravels[i+1:]...) 
			/*the previous ticket ID is deleted from the users upcoming travels*/
			break
		}
	}
	user.UpcomingTravels = append(user.UpcomingTravels, newTicketID)
	/*new ticket ID id now added to the user's upcoming travels list*/

	updatedUserJSON, err := json.Marshal(user)
	if err != nil || updatedUserJSON == nil {
		return fmt.Errorf("error while making JSON of the updated user")
	}
	err = ctx.GetStub().PutState(userID, updatedUserJSON)

	if err != nil {
		return fmt.Errorf("error while updating the user details in the hyperledger")
	}
	//new seat number updated in the record successfully
	// transport.Travellers[previousDate][previousSeat]="" /*removing the userID from previously booked seat*/

	ticket.DateofTravel = date
	ticket.SeatNumber = newSeat
	ticket.DateofUpdate = time.Now().Format("2006-01-02 15:04:05")

	ticket.TicketID = newTicketID

	updatedTicketJSON, err := json.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("error in converting ticket to JSON format")
	}

	err = ctx.GetStub().DelState(ticketID) /*the previous ticket is deleted*/
	if err != nil {
		return fmt.Errorf("failed to delete old ticket entry: %v", err)
	}
	return ctx.GetStub().PutState(newTicketID, updatedTicketJSON)
}