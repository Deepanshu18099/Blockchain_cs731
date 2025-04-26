package chaincode

import (
	"deepanshu18099/blockchain_ledger_backend/models"
	"fmt"
	"os"
	"os/exec"
	"encoding/json"
)



// BuildChaincodeArgs constructs the peer command
func BuildChaincodeArgs(user models.UserRequest23) []string {
	ccInput := fmt.Sprintf(`{"function":"CreateUser","Args":["%s","%s","%s","%s"]}`,
		user.Email, user.Name, user.Phone, user.UserID)
	// Set environment variables
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
func RunPeerCommand(args []string) (models.UserRequest32, error) {
	cmd := exec.Command("peer", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return models.UserRequest32{}, fmt.Errorf("failed to run peer command: %w", err)
	}
	var sample_outs models.UserRequest32
	err = json.Unmarshal(output, &sample_outs)
	if err != nil {
		return models.UserRequest32{}, fmt.Errorf("failed to unmarshal output: %w", err)
	}
	return sample_outs, nil
}
