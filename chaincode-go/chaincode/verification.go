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

// func (s*SmartContract)VerifyProvider(ctx contractapi.TransactionContextInterface, providerID string)(string, error){
// 	providerJSON,err := ctx.GetStub().GetState(providerID)
// 	if err != nil{
// 		return "",fmt.Errorf("error in getting the providers details from the hyperledger")
// 	}
// 	if providerJSON == nil{
// 		return "The provider doesn't  exist",fmt.Errorf("error")
// 	}

// 	var provider Provider

// 	err = json.Unmarshal(providerJSON, &provider)
// 	if err != nil{
// 		return "", fmt.Errorf("error while pointing to the provider details")
// 	}
// 	provider.Verified = true
// 	updatedProviderJSON,_ := json.Marshal(provider)
// 	ctx.GetStub().PutState(providerID, updatedProviderJSON)

// 	return "The provider verified successfully",nil
// }

func (s *SmartContract) VerifyProvider(ctx contractapi.TransactionContextInterface, providerID, transportID string) (string,error) {
    providerJSON, _ := ctx.GetStub().GetState(providerID)
    var provider Provider
   json.Unmarshal(providerJSON, &provider)

	found := false
    for _, service := range provider.Services {
        if service == transportID {
            found = true
            break
        }
    }
	if(!found){
		return "verification not successful",nil
	} else{
		return "verification successful",nil
	}
}

func (s*SmartContract)VerifyTicket(ctx contractapi.TransactionContextInterface,  ticketID,providerID, userID string)(string, error){
	ticketJSON,_ := ctx.GetStub().GetState(ticketID)
	if ticketJSON == nil{
		return "ticket doesn't exist", fmt.Errorf("error")
	}
	var ticket TicketDetails
	json.Unmarshal(ticketJSON,&ticket)
	if(ticket.UserID == userID && ticket.ProviderID == providerID){
		ticket.Status = "verified"
	}else{
		ticket.Status = "not verified"
	}
	updatedTicketJSON,_ := json.Marshal(ticket)
	ctx.GetStub().PutState(ticketID,updatedTicketJSON)
	return "Ticket verified successfully",nil
}

func (s*SmartContract)VerifyTransaction(ctx contractapi.TransactionContextInterface, paymentID, from, to string)(string, error){
	paymentJSON,_ := ctx.GetStub().GetState(paymentID)
	
	if(paymentJSON == nil){
		return "Error: ", fmt.Errorf("error: payment detsils for the id: %s doesn't exist", paymentID)
	}

	var payment PaymentDetail

	_ = json.Unmarshal(paymentJSON, &payment)
	
	if(payment.From == from && payment.To == to){
		return "Transaction verification successfull",nil
	}
	return "Error: Not Verified",fmt.Errorf("error: identities not mached")
}