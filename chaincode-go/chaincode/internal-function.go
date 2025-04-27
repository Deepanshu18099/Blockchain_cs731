package chaincode

import (
	"encoding/json"
	"fmt"
	"math"

	// "math/rand"
	// "strconv"
	"time"
	// "strconv"
	// "sort"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

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

func (s *SmartContract) AddBalance(ctx contractapi.TransactionContextInterface, email string, amount float64) ([]byte, error) {
	var username = email
	exists, err := s.detailExists(ctx, username)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("the user %s does not exist", email)
	}

	detailJSON, err := ctx.GetStub().GetState(username)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if detailJSON == nil {
		return nil, fmt.Errorf("details of %s do not exist", email)
	}
	var asset User
	err = json.Unmarshal(detailJSON, &asset)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the data: %v", err)
	}

	asset.BankBalance += amount

	updatedUserJSON, err := json.Marshal(asset)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated user data: %v", err)
	}

	err = ctx.GetStub().PutState(username, updatedUserJSON)

	// getting the transaction ID and adding it to returned JSON
	txID := ctx.GetStub().GetTxID()
	updatedUserJSON, err = json.Marshal(struct {
		User
		TransactionID string `json:"transaction_id"`
	}{
		User:          asset,
		TransactionID: txID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to put updated user data: %v", err)
	}
	return updatedUserJSON, nil
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

func (s *SmartContract) MakeUserPublic(ctx contractapi.TransactionContextInterface, email string) error {
	var username = email

	detailJSON, err := ctx.GetStub().GetState(username)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if detailJSON == nil {
		return fmt.Errorf("details of %s do not exist", email)
	}

	var asset User
	err = json.Unmarshal(detailJSON, &asset)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the data: %v", err)
	}

	asset.IsAnonymous = false

	updatedUserJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal updated user data: %v", err)
	}

	return ctx.GetStub().PutState(username, updatedUserJSON)
}
