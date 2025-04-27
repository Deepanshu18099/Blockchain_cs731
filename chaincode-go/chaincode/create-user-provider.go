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

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, email,name,phone,id,role string)error {
	var username = email
	exists, err := s.detailExists(ctx, username)

	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the user %s already exists", email)
	}

	// log the details
	fmt.Printf("User Details: %s, %s, %s, %s, %s\n", email, name, phone, id, role)


	if role == "user" {
		user := User{
			Name: name,
			Email: email,
			Phone: phone,
			PastTravels: make([]string, 0),
			UpcomingTravels: make([]string, 0),
			BankBalance: 0.0,
			IsAnonymous: false,
			PaymentID: make([]string, 0),
		}
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return ctx.GetStub().PutState(username, userJSON)
	}
	if role == "provider" {
		provider := Provider{
			Name: name,
			Email: email,
			Phone: phone,
			Services: make([]string, 0),
			Verified: false,
			BankBalance: 0.0,
			PaymentID: make([]string, 0),
		}
		providerJSON, err := json.Marshal(provider)
		if err != nil {
			return err
		}
		return ctx.GetStub().PutState(username, providerJSON)
	}
	return fmt.Errorf("Invalid role")
}


// func (s *SmartContract) CreateProvider(ctx contractapi.TransactionContextInterface, email,name,phone string )error {
// 	var providerName = email
// 	exists, err := s.detailExists(ctx, providerName)
// 	if err != nil {
// 		return err
// 	}
// 	if exists {
// 		return fmt.Errorf("the Provider %s already exists", email)
// 	}

//   	provider := Provider{
//   		Name: name,
//   		Email: email,
//   		Phone: phone,
//   		Services: make([]string, 0),
//   	}
//   	providerJSON, err := json.Marshal(provider)
// 	if err != nil {
// 		return err
// 	}
// 	return ctx.GetStub().PutState(providerName, providerJSON)
// }