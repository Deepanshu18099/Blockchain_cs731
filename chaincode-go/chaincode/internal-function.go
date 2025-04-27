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

