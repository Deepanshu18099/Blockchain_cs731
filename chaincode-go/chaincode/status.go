package chaincode
import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s*SmartContract)GetDetailUser(ctx contractapi.TransactionContextInterface, userID string)(User,error){
	userJSON,err := ctx.GetStub().GetState(userID)
	if err != nil{
		return User{}, fmt.Errorf("error occured while fetching the users details %s",err)
	}
	if userJSON == nil{
		return User{},fmt.Errorf("the user %s does not exist",userID)
	}
	var user User

	err = json.Unmarshal(userJSON, &user)
	if err != nil{
		return user, fmt.Errorf("error while pointing to the user details")
	}
	return user,nil
}

func (s*SmartContract)GetDetailProvider(ctx contractapi.TransactionContextInterface, providerID string)(Provider,error){
	providerJSON,err := ctx.GetStub().GetState(providerID)
	if err != nil{
		return Provider{}, fmt.Errorf("error occured while fetching the users details %s",err)
	}
	if providerJSON == nil{
		return Provider{}, fmt.Errorf("the provider %s does not exist",providerID)
	}
	var provider Provider

	err = json.Unmarshal(providerJSON, &provider)
	if err != nil{
		return provider, fmt.Errorf("error while pointing to the provider details")
	}
	return provider,nil
}

func (s*SmartContract)GetDetailTicket(ctx contractapi.TransactionContextInterface, ticketID string)(TicketDetails,error){
	ticketJSON,err := ctx.GetStub().GetState(ticketID)
	var ticket TicketDetails

	if err != nil{
		return ticket, fmt.Errorf("error while fetching the ticket details")
	}
	if ticketJSON == nil{
		return ticket, fmt.Errorf("error: Ticket with ticket id = %s does not exist",ticketID)
	}
	err = json.Unmarshal(ticketJSON,&ticket)
	if err != nil{
		return ticket, fmt.Errorf("error while pointing to the ticket details")
	}

	return ticket, nil
}

//##########################################################################################################
//<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>