package chaincode

import (
	"encoding/json"
	"fmt"

	// "time"
	// "math"
	// "strconv"
	// "math/rand"
	"sort"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func getNextDate(date string) string {
	// Assuming date is in the format "YYYY-MM-DD"
	year, month, day := 0, 0, 0
	fmt.Sscanf(date, "%d-%d-%d", &year, &month, &day)
	day++
	if day > 30 {
		day = 1
		month++
		if month > 12 {
			month = 1
			year++
		}
	}
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func (s *SmartContract) AddTransportService(ctx contractapi.TransactionContextInterface, providerID, vehicleNumber, source, dest, dept, arrt, totalt, mode string, cap int32, basep float64, startDate, endDate string) ([]byte, error) {
	providerJSON, err := ctx.GetStub().GetState(providerID)
	if err != nil {
		return nil, fmt.Errorf("error occured while fetching the provider details %s", err)
	}
	if providerJSON == nil {
		return nil, fmt.Errorf("the provider %s does not exist", providerID)
	}

	var transportID = "transport-" + vehicleNumber + "-" + source + "-" + dest

	transportJSON, err := ctx.GetStub().GetState(transportID)
	if err != nil {
		return nil, fmt.Errorf("error occured while fetching the transport details %s", err)
	}
	if transportJSON != nil {
		return nil, fmt.Errorf("the transport %s already exists", transportID)
	}

	arr := make([]int32, cap) /* for maintaining the available seat numbers in the transport*/
	tempArr := make([]string, cap)

	for i := 0; i < int(cap); i++ {
		arr[i] = int32(i + 1)
		tempArr[i] = ""
	}

	travellerMap := make(map[string][]string)
	seatMap := make(map[string][]int32) /*date wise availabe seat numbers maintained*/

	dates := []string{}
	// Assuming startDate and endDate are in the format "YYYY-MM-DD", add the dates in array
	for {
		dates = append(dates, startDate)
		if startDate == endDate {
			break
		}
		startDate = getNextDate(startDate)
	}

	for _, date := range dates {
		// Copying is necessary, otherwise all map values will point to the same slice
		sort.Slice(arr, func(i, j int) bool {
			return arr[i] < arr[j]
		})
		seatsCopy := make([]int32, len(arr))
		copy(seatsCopy, arr)
		seatMap[date] = seatsCopy

		travellerCopy := make([]string, len(tempArr))
		copy(travellerCopy, tempArr)
		travellerMap[date] = travellerCopy
	}

	/*seat map is created*/

	/*maintaining the travellers for a transportation*/

	newTransport := TransportDetails{
		ID:              transportID,
		Source:          source,
		Destination:     dest,
		DepartureTime:   dept,
		ArrivalTime:     arrt,
		BasePrice:       basep,
		Rating:          3.00, /*all new travels will have a rating of 3 at start*/
		Capacity:        cap,
		ModeofTravel:    mode,
		JourneyDuration: totalt,
		DateofTravel:    dates,
		SeatMap:         seatMap,
		// Travellers: travellerMap,
		ProviderID: providerID,
	}

	transpotJSON, err := json.Marshal(newTransport)
	if err != nil {
		return nil, fmt.Errorf("error in serializing the transport data %s", err)
	}

	err = ctx.GetStub().PutState(transportID, transpotJSON)
	if err != nil {
		return nil, fmt.Errorf("error in putting the transport data %s", err)
	}

	var updateProvider Provider
	err = json.Unmarshal(providerJSON, &updateProvider)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshalling the provider data %s", err)
	}
	updateProvider.Services = append(updateProvider.Services, transportID)

	updateProviderJSON, err := json.Marshal(updateProvider)
	if err != nil {
		return nil, fmt.Errorf("error in serializing the provider data %s", err)
	}

	/*
		should return these fields
			// check if the output is success
		if clean_output["status"] != "success" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add transport"})
			return
		}
		// check if the transportID is present in the output
		if clean_output["transportID"] == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add transport"})
			return
		}

		// now return success, and transportID
		c.JSON(http.StatusOK, gin.H{
			"message":        "Transport added successfully",
			"transport_id":   transportID,
			"transaction_id": clean_output["transaction_id"],
	*/
	err = ctx.GetStub().PutState(providerID, updateProviderJSON)
	if err != nil {
		return nil, fmt.Errorf("error in putting the provider data %s", err)
	}
	// return the transportID and transactionID
	returnitem := map[string]string{
		"transportID":   transportID,
		"transactionID": ctx.GetStub().GetTxID(),
		"message":       "Transport added successfully",
	}

	returnJSON, err := json.Marshal(returnitem)
	if err != nil {
		return nil, fmt.Errorf("error in serializing the return data %s", err)
	}
	return returnJSON, nil
}

func (s *SmartContract) GetAvailableTransports(ctx contractapi.TransactionContextInterface, src, dest, date, mode string) ([]byte, error) {
	/*
		Input: source, destination, date, mode of travel
		Output: list of available transports, in byte array of arrayy of maps
	*/
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var availableTransports []*TransportDetails
	fmt.Println("helloooooooooooooooooooooooooooooooooooooooo")

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(queryResponse.Key, "transport") {
			continue
		}

		var transport TransportDetails
		err = json.Unmarshal(queryResponse.Value, &transport)
		if err != nil {
			// Optional: log error or continue silently
			continue
		}

		ticketPrice, err := calculateDynamicPrice(ctx, transport.ID, date)
		if err != nil {
			return nil, fmt.Errorf("error in getting the payment of the transport")
		}
		fmt.Println("tranport_details", transport)
		// printing dates top 10
		// fmt.Println("ID", transport.ID)
		fmt.Println("size of dates and start date, enddate", len(transport.DateofTravel))
		if len(transport.DateofTravel) > 0 {
			fmt.Println("dates", transport.DateofTravel[0], transport.DateofTravel[len(transport.DateofTravel)-1])
		}

		// fmt.Println("top 10 seats",
		fmt.Println("ticket_price", ticketPrice)

		if transport.Source == src && transport.Destination == dest && transport.ModeofTravel == mode {
			// i := 0
			for _, availableDate := range transport.DateofTravel {
				// i++
				// if i > 40 {
				// 	break
				// }
				// fmt.Println(availableDate, date)
				if availableDate == date {
					requiredTransport := &TransportDetails{
						ID:              transport.ID,
						Source:          transport.Source,
						Destination:     transport.Destination,
						DepartureTime:   transport.DepartureTime,
						ArrivalTime:     transport.ArrivalTime,
						BasePrice:       ticketPrice,
						Rating:          transport.Rating,
						Capacity:        transport.Capacity,
						ModeofTravel:    transport.ModeofTravel,
						JourneyDuration: transport.JourneyDuration,
						DateofTravel:    []string{date},
						SeatMap:         map[string][]int32{date: transport.SeatMap[date]},
					}
					fmt.Println("yeeeeeerahaaa", requiredTransport)
					availableTransports = append(availableTransports, requiredTransport)
					// break
				}
			}
		}

	}

	returnitem := map[string]interface{}{
		"availableTransports": availableTransports,
		"transactionID":       ctx.GetStub().GetTxID(),
		"message":             "Available transports fetched successfully",
	}
	// now serializing the data to byte array
	availableTransportsJSON, err := json.Marshal(returnitem)
	if err != nil {
		return nil, fmt.Errorf("error in serializing the data %s", err)
	}
	return availableTransportsJSON, nil
}
