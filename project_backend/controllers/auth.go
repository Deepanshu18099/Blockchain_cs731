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

	"log"
	"math/rand"
	"net/http"
	"strconv"
	"os"
	"fmt"
	
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
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
		userID = "USR_" +  strconv.Itoa(rand.Intn(1000000))
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

	// Call the chaincode function to create the user on the ledger
	args := chaincode.BuildChaincodeArgs(user1)
	output, err := chaincode.RunPeerCommand(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	

	// register user id, password, and email in the database
	_, err = collection.InsertOne(c, bson.M{
		"email":    user.Email,
		"name":     user.Name,
		"phone":    user.Phone,
		"role":     user.Role,
		"password": user.Password,
		"user_id":  userID,
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
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists in the database
	db := database.ConnectDB()
	collection := db.Database("blockchain_ledger").Collection("users")
	err := collection.FindOne(c, bson.M{"email": user.Email, "password": user.Password}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
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