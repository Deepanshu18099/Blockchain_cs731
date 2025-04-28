package main

import (
	"deepanshu18099/blockchain_ledger_backend/controllers"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	// allow requests from all origins, cors error should not occur, also from localhost
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"}, // <-- allow ALL headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// router.Use(cors.Default())
	router.GET("/tickets", controllers.GetTickets)
	router.GET("/tickets/:id", controllers.GetTicketByID)
	router.POST("/tickets", controllers.CreateTicket)
	router.PUT("/tickets/:id", controllers.UpdateTicket)
	router.DELETE("/tickets/:id", controllers.DeleteTicket)
	router.POST("/ledger/createuser", controllers.CreateLedgerUser)
	router.POST("/ledger/login", controllers.Login)
	router.POST("/ledger/addMoney", controllers.AddMoneyToUser)
	router.POST("/Addtransport", controllers.AddTransport)


	// to get all transports available for that selections by the user
	// `http://localhost:8080/Gettransports/${mode}/${source}/${destination}/${date}`,
	 
	router.GET("/Gettransports/:mode/:source/:destination/:date", controllers.GetTransports)
	router.GET("/Gettransport/:id", controllers.GetTransportStatus)
	router.Run(":8080")

}
