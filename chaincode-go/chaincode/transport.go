package chaincode

import (
	"encoding/json"
	"fmt"
	// "time"
	// "math"
	// "strconv"
	"strings"
	"sort"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)


func (s *SmartContract) AddTransportService(ctx contractapi.TransactionContextInterface, providerID,vehicleNumber, source,dest,dept,arrt,totalt,mode string, cap int32,basep float64,dates []string) error {
	providerJSON,err := ctx.GetStub().GetState(providerID)
	if err != nil{
		return fmt.Errorf("error occured while getting details of the provider: %s", err)
	}
	if providerJSON == nil{
		return fmt.Errorf("the provider: %s does not exist",providerID)
	}

	var transportID = "transport:"+vehicleNumber+"-"+source+"-"+dest

	transportJSON,err := ctx.GetStub().GetState(transportID)
	if err != nil{
		return fmt.Errorf("error occured while getting details of the transportation: %s", err)
	}
	if transportJSON != nil{
		return fmt.Errorf("the transportation already exist")
	}

	arr := make([]int32, cap) /* for maintaining the available seat numbers in the transport*/
	tempArr := make([]string,cap)

	for i := 0; i < int(cap); i++ {
		arr[i] = int32(i + 1)
		tempArr[i] = ""
	}

	travellerMap := make(map[string][]string)
	seatMap := make(map[string][]int32) /*date wise availabe seat numbers maintained*/

	for _, date := range dates {
		// Copying is necessary, otherwise all map values will point to the same slice
		sort.Slice(arr, func(i, j int) bool {
			return arr[i] < arr[j]
		})
		seatsCopy := make([]int32, len(arr))
		copy(seatsCopy, arr)
		seatMap[date] = seatsCopy

		travellerCopy := make([]string,len(tempArr))
		copy(travellerCopy, tempArr)
		travellerMap[date] = travellerCopy
	}

	/*seat map is created*/

	/*maintaining the travellers for a transportation*/

	newTransport := TransportDetails{
		ID: transportID,
		Source: source,
		Destination: dest,
		DepartureTime:dept,
		ArrivalTime: arrt,
		BasePrice: basep,
		Rating: 3.00, /*all new travels will have a rating of 3 at start*/
		Capacity:cap,
		ModeofTravel: mode,
		JourneyDuration: totalt,
		DateofTravel: dates,
		SeatMap: seatMap,
		Travellers: travellerMap,
		ProviderID: providerID,
	}

	transpotJSON,err := json.Marshal(newTransport)
	if(err != nil){
		return fmt.Errorf("error occured %s", err)
	}

	err = ctx.GetStub().PutState(transportID, transpotJSON)
	if(err != nil){
		return fmt.Errorf("error occured while storing the transport details, %s",err)
	}

	var updateProvider Provider
	err = json.Unmarshal(providerJSON, &updateProvider)
	if err != nil{
		return fmt.Errorf("error in pointing to the updated provider data")
	}
	updateProvider.Services = append(updateProvider.Services, transportID)

	updateProviderJSON,err := json.Marshal(updateProvider)
	if err != nil{
		return fmt.Errorf("error in updating provider data")
	}
	return ctx.GetStub().PutState(providerID,updateProviderJSON)
}



func (s *SmartContract) GetAvailableTransports(ctx contractapi.TransactionContextInterface, src, dest, date string) ([]*TransportDetails, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var availableTransports []*TransportDetails

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(queryResponse.Key, "transport:") {
			continue
		}

		var transport TransportDetails
		err = json.Unmarshal(queryResponse.Value, &transport)
		if err != nil {
			// Optional: log error or continue silently
			continue
		}

		ticketPrice,err := calculateDynamicPrice(ctx,transport.ID, date)
		if err != nil{
			return nil,fmt.Errorf("error in getting the payment of the transport")
		}
		if transport.Source == src && transport.Destination == dest {
			for _, availableDate := range transport.DateofTravel {
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
					availableTransports = append(availableTransports, requiredTransport)
					break
				}
			}
		}
	}

	return availableTransports, nil
}

// func (s*SmartContract)GetAvailableTransports(ctx contractapi.TransactionContextInterface,src,dest,date string) ([]*TransportDetails,error){
// 	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer resultsIterator.Close()

// 		var availableTransports []*TransportDetails
// 		for resultsIterator.HasNext() {
// 			queryResponse, err := resultsIterator.Next()
// 			if err != nil {
// 				return nil, err
// 			}

// 			var transport TransportDetails
// 			err = json.Unmarshal(queryResponse.Value, &transport)
// 			if err != nil{
// 				return nil,fmt.Errorf("error in pointing the query responses")
// 			}
// 			if transport.Source == src && transport.Destination==dest{
// 				for _, availableDate := range transport.DateofTravel{
// 					if availableDate == date{
// 						requiredTransport := &TransportDetails{
// 							ID: transport.ID,
// 							Source: transport.Source,
// 							Destination: transport.Destination,
// 							DepartureTime: transport.DepartureTime,
// 							ArrivalTime: transport.ArrivalTime,
// 							BasePrice: transport.BasePrice,
// 							Rating: transport.Rating,
// 							Capacity: transport.Capacity,
// 							AvailableSeats: transport.AvailableSeats,
// 							ModeofTravel: transport.ModeofTravel,
// 							JourneyDuration: transport.JourneyDuration,
// 							DateofTravel: []string{date},
// 							SeatMap: map[string][]int32{date: transport.SeatMap[date]},
// 						}
// 						availableTransports = append(availableTransports,requiredTransport)
// 						break
// 					}
// 				}
// 			}
// 		}
// 	return availableTransports,nil;
// }