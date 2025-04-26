package chaincode

import (
	"encoding/json"
	"fmt"
	// "strconv"
	// "math"
	// "math/rand"
	// "time"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, email,name,phone string)error {
	var username = email
	exists, err := s.detailExists(ctx, username)

	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the user %s already exists", email)
	}

	// randomNumber := strconv.Itoa(rand.Intn(1000000000))
	// userID := randomNumber+email

  	customer := User{
        Name: name,
        Email: email,
        Phone: phone,
        PastTravels: make([]string, 0), // Initialize empty slice
        UpcomingTravels: make([]string, 0), // Initialize empty slice
        BankBalance: 0.0,
        IsAnonymous: false,                         /*Start with zero balance*/
		PaymentID: make([]string,0),
	}
	
  	userJSON, err := json.Marshal(customer)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(username, userJSON)
}


func (s *SmartContract) CreateProvider(ctx contractapi.TransactionContextInterface, email,name,phone string )error {
	var providerName = email
	exists, err := s.detailExists(ctx, providerName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Provider %s already exists", email)
	}

  	provider := Provider{
  		Name: name,
  		Email: email,
  		Phone: phone,
  		Services: make([]string, 0),
  	}
  	providerJSON, err := json.Marshal(provider)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(providerName, providerJSON)
}