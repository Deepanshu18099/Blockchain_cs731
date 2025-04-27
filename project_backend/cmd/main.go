package main

import (
	"deepanshu18099/blockchain_ledger_backend/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)


// capital letter starting fields are exportable fields.
// these are public fields.
// the later json thing is go specific, and is for golang



func main() {
	// load env variables from .env file located at ../.env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/tickets", controllers.GetTickets)
	router.GET("/tickets/:id", controllers.GetTicketByID)
	router.POST("/tickets", controllers.CreateTicket)
	router.PUT("/tickets/:id", controllers.UpdateTicket)
	router.DELETE("/tickets/:id", controllers.DeleteTicket)
	router.POST("/ledger/createuser", controllers.CreateLedgerUser)
	router.POST("/ledger/login", controllers.Login)
	router.POST("/ledger/addMoney", controllers.AddMoneyToUser)
	router.POST("/Addtransport", controllers.AddTransport)
	router.GET("/Gettransports", controllers.GetTransports)
	router.GET("/Gettransport/:id", controllers.GetTransportStatus)
	router.Run("localhost:8080")

}
