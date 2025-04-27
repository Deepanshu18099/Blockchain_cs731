package chaincode

import (
	// "deepanshu18099/blockchain_ledger_backend/models"
	// "encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// // BuildChaincodeArgs constructs the peer command
// func BuildChaincodeArgs(user models.UserRequest23, funcname string) []string {
// 	ccInput := fmt.Sprintf(`{"function":"%s","Args":["%s","%s","%s","%s","%s"]}`,
// 		funcname, user.Email, user.Name, user.Phone, user.UserID, user.Role)
// 	// Set environment variables

// argss := []string{}
//
//	argss = append(argss, arg1)
//	argss = append(argss, arg2)....
//
// now making it generalised with func name and args seperated by comma
func BuildChaincodeArgs(args []string, funcname string) []string {
	// all args had to be split into strings in loop seperated by comma in ccinput args
	// or else it will be treated as a single string which is not correct
	var ccInput string
	ccInput = fmt.Sprintf(`{"function":"%s","Args":[`, funcname)
	for i, arg := range args {
		if i == len(args)-1 {
			ccInput += fmt.Sprintf(`"%s"`, arg)
		} else {
			ccInput += fmt.Sprintf(`"%s",`, arg)
		}
	}
	ccInput += `]}`

	OrdererAddress := os.Getenv("OrdererAddress")
	OrdererTLSHostname := os.Getenv("OrdererTLSHostname")
	CAFile := os.Getenv("CAFile")
	ChannelName := os.Getenv("ChannelName")
	ChaincodeName := os.Getenv("ChaincodeName")
	Peer0Org1Address := os.Getenv("Peer0Org1Address")
	Peer0Org1TLSCert := os.Getenv("Peer0Org1TLSCert")
	Peer0Org2Address := os.Getenv("Peer0Org2Address")
	Peer0Org2TLSCert := os.Getenv("Peer0Org2TLSCert")

	return []string{
		"chaincode", "invoke",
		"-o", OrdererAddress,
		"--ordererTLSHostnameOverride", OrdererTLSHostname,
		"--tls",
		"--cafile", CAFile,
		"-C", ChannelName,
		"-n", ChaincodeName,
		"--peerAddresses", Peer0Org1Address,
		"--tlsRootCertFiles", Peer0Org1TLSCert,
		"--peerAddresses", Peer0Org2Address,
		"--tlsRootCertFiles", Peer0Org2TLSCert,
		"-c", ccInput,
	}
}

// RunPeerCommand executes the command and returns output
func RunPeerCommand(args []string) ([]byte, error) {
	cmd := exec.Command("peer", args...)
	log.Printf("Executing command: %s\n", cmd.String())
	// cmd.Dir = "/path/to/your/working/directory" // Set the working directory if needed
	cmd.Env = os.Environ() // Use the current environment variables

	// returnstype: fmt error or ctx.GetStub().PutState(username, providerJSON)
	output, err := cmd.CombinedOutput()

	// log.Printf("Output: %s\n", output)

	if err != nil {
		log.Printf("Error executing command: %s\n", err)
		return nil, err
	}
	log.Printf("Command output: %s\n", string(output))
	return output, nil

}
