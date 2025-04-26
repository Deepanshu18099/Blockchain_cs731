package main

import (
	"deepanshu18099/blockchain_ledger_backend/controllers"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const Testing_Mode bool = true

// Sample outputs for testing
var sample_outs = map[string]interface{}{
	"CreateUser": map[string]string{
		"publicKey":  "SamplePublicKey_ABC123XYZ",
		"privateKey": "SamplePrivateKey_DEF456UVW",
		"txID":       "SampleTransactionID_789GHI",
	},
}

type UserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

const HyperledgerPath = "/home/maverick/Downloads/courses8/cs731/Hyperledger/fabric"

// Chaincode settings (hardcoded for simplicity)
const (
	OrdererAddress     = "localhost:7050"
	OrdererTLSHostname = "orderer.example.com"
	CAFile             = HyperledgerPath + "/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
	ChannelName        = "ticketsystem"
	ChaincodeName      = "keyvalchaincode"
	Peer0Org1Address   = "localhost:7051"
	Peer0Org1TLSCert   = HyperledgerPath + "/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
	Peer0Org2Address   = "localhost:9051"
	Peer0Org2TLSCert   = HyperledgerPath + "/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
)

// BuildChaincodeArgs constructs the peer command
func BuildChaincodeArgs(user UserRequest) []string {
	ccInput := fmt.Sprintf(`{"function":"CreateUser","Args":["%s","%s","%s"]}`,
		user.Email, user.Name, user.Phone)

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
func RunPeerCommand(args []string) (string, error) {
	cmd := exec.Command("peer", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Handler function
func createLedgerUser(c *gin.Context) {
	var user UserRequest
	// if in testing mode, just return success and output from list of sample outs above for each function

	// print something to test

	if Testing_Mode {
		c.JSON(http.StatusOK, gin.H{
			"message": "User created on the ledger",
			"output":  sample_outs["CreateUser"],
		})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	args := BuildChaincodeArgs(user)
	output, err := RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to invoke chaincode",
			"error":   err.Error(),
			"output":  output,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created on the ledger",
		"output":  output,
	})
}

// capital letter starting fields are exportable fields.
// these are public fields.
// the later json thing is go specific, and is for golang
type ticket struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Destination string  `json:"destination"`
	Source      string  `json:"source"`
}

var tickets = []ticket{
	{
		ID: "1", Name: "Ticket 1", Price: 100.0, Destination: "New York", Source: "Los Angeles",
	},
	{
		ID: "2", Name: "Ticket 2", Price: 200.0, Destination: "Chicago", Source: "San Francisco",
	},
	{
		ID: "3", Name: "Ticket 3", Price: 300.0, Destination: "Miami", Source: "Seattle",
	},
}

func getTickets(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tickets)
}
func getTicketByID(c *gin.Context) {
	id := c.Param("id")
	for _, ticket := range tickets {
		if ticket.ID == id {
			c.JSON(http.StatusOK, ticket)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "ticket not found"})
}
func createTicket(c *gin.Context) {
	var newTicket ticket
	if err := c.ShouldBindJSON(&newTicket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tickets = append(tickets, newTicket)
	c.JSON(http.StatusCreated, newTicket)
}

func updateTicket(c *gin.Context) {
	id := c.Param("id")
	var updatedTicket ticket
	if err := c.ShouldBindJSON(&updatedTicket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, ticket := range tickets {
		if ticket.ID == id {
			tickets[i] = updatedTicket
			c.JSON(http.StatusOK, updatedTicket)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "ticket not found"})
}

func deleteTicket(c *gin.Context) {
	id := c.Param("id")
	for i, ticket := range tickets {
		if ticket.ID == id {
			tickets = append(tickets[:i], tickets[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "ticket deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "ticket not found"})
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/tickets", getTickets)
	router.GET("/tickets/:id", getTicketByID)
	router.POST("/tickets", createTicket)
	router.PUT("/tickets/:id", updateTicket)
	router.DELETE("/tickets/:id", deleteTicket)
	router.POST("/ledger/createuser", createLedgerUser)

	router.Run("localhost:8080")

	controllers.Prinn()
}
