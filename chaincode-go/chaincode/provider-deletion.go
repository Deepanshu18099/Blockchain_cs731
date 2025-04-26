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

func (s*SmartContract)DeleteProvider(ctx contractapi.TransactionContextInterface, providerID string)error{
	providerJSON,_ := ctx.GetStub().GetState(providerID)
	if(providerJSON == nil){
		return fmt.Errorf("error: the provider %s doesn't exist",providerID)
	}
	var provider Provider
	json.Unmarshal(providerJSON,&provider)

	/*for all the transports, refund will be given to all the users with a notification*/
	

	/*some penalty to be imposed on the provider*/

	return nil
}