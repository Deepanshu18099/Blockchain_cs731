/*
This is another way to write
a long comment in Go, using
the block comment syntax.
It's less common for general
comments but can be useful
in certain situations.
*/

package controllers

import (
	"deepanshu18099/blockchain_ledger_backend/models"
	"deepanshu18099/blockchain_ledger_backend/chaincode"
	"deepanshu18099/blockchain_ledger_backend/database"
	
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"math/rand"

)

// Sample outputs for testing
var sample_outs = map[string]interface{}{
	"CreateUser": map[string]string{
		"publicKey":  "SamplePublicKey_ABC123XYZ",
		"privateKey": "SamplePrivateKey_DEF456UVW",
		"txID":       "SampleTransactionID_789GHI",
	},
}

// Handler function
func CreateLedgerUser(c *gin.Context) {

	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Testing_Mode := os.Getenv("TestingMode") == "true"

	var user models.UserRequest12
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


	// create random userids(name + random chars) which are not present in the database mongodb



	
	db := database.ConnectDB()
	if err := db.AutoMigrate(&models.UserRequest12{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to migrate database"})
		return
	}


	// Generate a random userID (name + random chars)for the user and return when found a available id
	
	// check if mail is already present in the database
	var existingUser models.UserRequest12
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// generating the id now in while loop
	var UserID string
	for {
		UserID = user.Name + "_" + string(rand.Intn(10000))
		if err := db.Where("user_id = ?", UserID).First(&existingUser).Error; err != nil {
			break
		}
	}

	
	comp_user := models.UserRequest23{
		Email:  user.Email,
		Name:   user.Name,
		Phone:  user.Phone,
		UserID: UserID,
		Role:   user.Role,
	}



	// one success send this to ledger and save in ledger and certificates returned
	args := chaincode.BuildChaincodeArgs(comp_user)
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to invoke chaincode",
			"error":   err.Error(),

			.....
		})
		return
	}
	// Save the user to the database

	c.JSON(http.StatusOK, gin.H{
		"message": "User created on the ledger",
		"output":  output,
	})
}
