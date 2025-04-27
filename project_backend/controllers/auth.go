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
	"deepanshu18099/blockchain_ledger_backend/chaincode"
	"deepanshu18099/blockchain_ledger_backend/database"
	"deepanshu18099/blockchain_ledger_backend/models"

	// "encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	// "regexp"
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
	log.Println("Loading environment variables...")

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

	// type *mongo.Client
	db := database.ConnectDB()
	collection := db.Database("blockchain_ledger").Collection("users")
	// check if user already exists
	var existingUser models.User
	err = collection.FindOne(c, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// log.Println("Loading env11ironment variables...")

	// creating random user id
	var userID string
	for {
		// generate a random user id
		log.Println(rand.Intn(1000000))
		userID = "USR_" + strconv.Itoa(rand.Intn(1000000))
		// check if user id already exists
		err = collection.FindOne(c, bson.M{"user_id": userID}).Decode(&existingUser)
		if err != mongo.ErrNoDocuments {
			// user id already exists, generate a new one
			continue
		}
		break
	}

	// log.Println("Loading envi22ronment variables...")

	var user1 models.UserRequest23
	user1.UserID = userID
	user1.Email = user.Email
	user1.Name = user.Name
	user1.Phone = user.Phone
	user1.Role = user.Role

	// args list of strings
	argss := []string{}
	argss = append(argss, user1.Email)
	argss = append(argss, user1.Name)
	argss = append(argss, user1.Phone)
	argss = append(argss, user1.UserID)
	argss = append(argss, user1.Role)

	// now prepare to send the request to the chaincode
	// Call the chaincode function to create the user on the ledger
	args := chaincode.BuildChaincodeArgs(argss, "CreateEntity")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// register user id, password, and email in the database
	_, err = collection.InsertOne(c, bson.M{
		"email":    user.Email,
		"username": user.Name,
		"phone":    user.Phone,
		"role":     user.Role,
		"password": user.Password,
		"userid":   userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	// return success
	c.JSON(http.StatusOK, gin.H{
		"message": "User created on the ledger",
		"output":  output,
		"userid":  userID,
	})
}

// will check email, password and returns access token for further use
func Login(c *gin.Context) {
	var user models.Signin
	log.Println("Signin function called")
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Signin function called")

	// Check if the user exists in the database
	db := database.ConnectDB()
	collection := db.Database("blockchain_ledger").Collection("users")

	var existingUser models.User
	log.Println("Signin function called")
	err := collection.FindOne(c, bson.M{"email": user.Email, "password": user.Password}).Decode(&existingUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	log.Println(existingUser)

	log.Println("Signin function called")

	// now getting user details from the chaincode
	argss := []string{}
	argss = append(argss, existingUser.Email)
	// now prepare to send the request to the chaincode
	args := chaincode.BuildChaincodeArgs(argss, "GetDetailUser")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Decode the output
	// outputDecoded := make(map[string]interface{})
	log.Println("Signin function called", output)
	// err = json.Unmarshal(output, &outputDecoded)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode output"})
	// 	return
	// }
	// log.Println("Signin function called", outputDecoded)

	// // check if the output has the updated balance
	// updatedbalance, ok := outputDecoded["BankBalance"].(string)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated balance"})
	// 	return
	// }

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   user.Email,
		"role":    existingUser.Role,
		"userid":  existingUser.UserID,
		"balance": "1000",
		"exp":     jwt.TimeFunc().Add(time.Hour * 24).Unix(), // Token expiration time
	})

	log.Println("Signin function called", existingUser.UserID)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// Return the token to the client
	c.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"role":    existingUser.Role,
		"userid":  existingUser.UserID,
		"balance": "1000",
	})

}

// checking authorization of a request
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// will be called by all the services to check if the token is valid and get the claims
func Authcheck(c *gin.Context) (jwt.MapClaims, bool) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		c.Abort()
		return nil, false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return nil, false
	}

	return claims, true
}

func AddMoneyToUser(c *gin.Context) {
	// use authmiddleware to check if token is valid and get claims
	// using the authcheck function
	claims, ok := Authcheck(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	argss := []string{}
	argss = append(argss, email)
	// now prepare to send the request to the chaincode
	// Call the chaincode function to create the user on the ledger
	args := chaincode.BuildChaincodeArgs(argss, "AddMoneyToUser")
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("AddMoneyToUser function called", output)

	// ................missing part.......................

	// // check if the output has the updated balance
	// updatedbalance, ok := outputdecoded["updatedbalance"].(string)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated balance"})
	// 	return
	// }
	// now return the updated balance
	c.JSON(http.StatusOK, gin.H{
		"message":        "Money added to user",
		"updatedbalance": "900",
		"transaction_id": "tx_10101",
	})

}
