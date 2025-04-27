package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	// "math"
	// "strconv"
	// "strings"
	// "strconv"
	// "sort"
	// "math"
	// "math/rand"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func UserToProviderPayment(ctx contractapi.TransactionContextInterface, userID, providerID string, amount float64) error {

	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("error %s occured", err)
	}
	if userJSON == nil {
		return fmt.Errorf("error: user %s doesn't exist", userID)
	}
	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return fmt.Errorf("error: failed to unmarhsal user %s", userID)
	}

	if user.BankBalance < amount {
		return fmt.Errorf("error: insufficient balance in the user's account")
	}

	providerJSON, err := ctx.GetStub().GetState(providerID)
	if err != nil {
		return fmt.Errorf("error %s occured", err)
	}
	if providerJSON == nil {
		return fmt.Errorf("error: provider %s doesn't exist", providerID)
	}

	var provider Provider
	err = json.Unmarshal(providerJSON, &provider)
	if err != nil {
		return fmt.Errorf("error: failed to unmarshal provider %s", providerID)
	}

	user.BankBalance -= amount
	provider.BankBalance += amount

	paymentID := "payment-" + time.Now().Format("2006-01-02 15:04:05")

	payment := PaymentDetail{
		PaymentID:   paymentID,
		From:        userID,
		To:          providerID,
		Amount:      amount,
		PaymentTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	paymentJSON, _ := json.Marshal(payment)
	ctx.GetStub().PutState(paymentID, paymentJSON)

	user.PaymentID = append(user.PaymentID, paymentID)
	provider.PaymentID = append(provider.PaymentID, paymentID)

	updatedUserJSON, _ := json.Marshal(user)
	err = ctx.GetStub().PutState(userID, updatedUserJSON)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	updatedProviderJSON, _ := json.Marshal(provider)
	return ctx.GetStub().PutState(providerID, updatedProviderJSON)
}

func ProviderToUserPayment(ctx contractapi.TransactionContextInterface, providerID, userID string, amount float64) error {

	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("error %s occured", err)
	}
	if userJSON == nil {
		return fmt.Errorf("error: user %s doesn't exist", userID)
	}
	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return fmt.Errorf("error: failed to unmarhsal user %s", userID)
	}

	providerJSON, err := ctx.GetStub().GetState(providerID)
	if err != nil {
		return fmt.Errorf("error %s occured", err)
	}
	if providerJSON == nil {
		return fmt.Errorf("error: provider %s doesn't exist", providerID)
	}

	var provider Provider
	err = json.Unmarshal(providerJSON, &provider)
	if err != nil {
		return fmt.Errorf("error: failed to unmarshal provider %s", providerID)
	}

	if(provider.BankBalance<amount){
		return fmt.Errorf("error: insufficient balance ")
	}

	user.BankBalance += amount
	provider.BankBalance -= amount

	paymentID := "payment-" + time.Now().Format("2006-01-02")

	payment := PaymentDetail{
		PaymentID:   paymentID,
		From:        userID,
		To:          providerID,
		Amount:      amount,
		PaymentTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	paymentJSON, _ := json.Marshal(payment)
	ctx.GetStub().PutState(paymentID, paymentJSON)

	user.PaymentID = append(user.PaymentID, paymentID)
	provider.PaymentID = append(provider.PaymentID, paymentID)

	updatedUserJSON, _ := json.Marshal(user)
	ctx.GetStub().PutState(userID, updatedUserJSON)

	updatedProviderJSON, _ := json.Marshal(provider)
	ctx.GetStub().PutState(providerID, updatedProviderJSON)

	return nil
}

