package chaincode

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	// "strconv"
	// "sort"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (t *SmartContract) detailExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	detailJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return detailJSON != nil, nil
}

func calculateDynamicPrice(ctx contractapi.TransactionContextInterface, transportID, date string) (float64, error) {
	transportJSON, err := ctx.GetStub().GetState(transportID)
	if err != nil || transportJSON == nil {
		return 0, fmt.Errorf("transportation not found")
	}
	var transport TransportDetails
	json.Unmarshal(transportJSON, &transport)

	// Calculate days until departure
	departureTime, _ := time.Parse(time.RFC3339, transport.DepartureTime)

	availablePercent := float64(len(transport.SeatMap[date])) / float64(transport.Capacity) * 100
	seatFactor := 1.00

	switch {
	case availablePercent <= 30:
		seatFactor = 1.10
	case availablePercent <= 20:
		seatFactor = 1.20
	case availablePercent <= 10:
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
	case timeToTravel <= 72:
		timeFactor = 1.20
	case timeToTravel <= 120:
		timeFactor = 1.10
	default:
		timeFactor = 1.00
	}

	currentPrice := transport.BasePrice * timeFactor * seatFactor
	return math.Round(currentPrice*100) / 100, nil
}

func (s *SmartContract) UserToProviderPayment(ctx contractapi.TransactionContextInterface, userID, providerID string, amount float64) error {
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

	user.BankBalance -= amount
	provider.BankBalance += amount

	randomNumber := strconv.Itoa(rand.Intn(1000000000))
	paymentID := "payment:" + randomNumber + "-" + time.Now().Format("2006-01-02")

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

func (s *SmartContract) ProviderToUserPayment(ctx contractapi.TransactionContextInterface, providerID, userID string, amount float64) error {
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

	user.BankBalance += amount
	provider.BankBalance -= amount

	randomNumber := strconv.Itoa(rand.Intn(1000000000))
	paymentID := "payment:" + randomNumber + "-" + time.Now().Format("2006-01-02")

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
