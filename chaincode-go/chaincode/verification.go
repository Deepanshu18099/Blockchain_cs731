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

func (s *SmartContract) VerifyProvider(ctx contractapi.TransactionContextInterface, providerID string) error {
    // Get invoker's X.509 certificate
    clientID, err := ctx.GetClientIdentity().GetID()
    if err != nil {
        return fmt.Errorf("failed to get client identity: %v", err)
    }
	if clientID == ""{
		return fmt.Errorf("client ID doesn't exist")
	}

    // Verify invoker is from Regulatory Authority's MSP
    err = ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "regulator")
    if err != nil {
        return fmt.Errorf("unauthorized verification attempt")
    }

    // Existing verification logic
    providerJSON, _ := ctx.GetStub().GetState(providerID)
    var provider Provider
    json.Unmarshal(providerJSON, &provider)
    
    provider.Verified = true
    updatedProviderJSON, _ := json.Marshal(provider)
    return ctx.GetStub().PutState(providerID, updatedProviderJSON)
}

func (s*SmartContract)VerifyTicket(ctx contractapi.TransactionContextInterface, ticketID string)(string, error){
	ticketJSON,_ := ctx.GetStub().GetState(ticketID)
	if ticketJSON == nil{
		return "ticket doesn't exist", fmt.Errorf("error")
	}
	var ticket TicketDetails
	json.Unmarshal(ticketJSON,&ticket)
	ticket.Status = "verified"
	updatedTicketJSON,_ := json.Marshal(ticket)
	ctx.GetStub().PutState(ticketID,updatedTicketJSON)
	return "Ticket verified successfully",nil
}

func (s*SmartContract)VerifyTransaction(ctx contractapi.TransactionContextInterface, paymentID string)(string, error){
	paymentJSON,_ := ctx.GetStub().GetState(paymentID)
	if paymentJSON == nil{
		return "Error: payment with given ID doesn't exist", fmt.Errorf("error: payment with id= %s doesn't exist",paymentID)
	}
	return "Transaction verification successfull",nil
}