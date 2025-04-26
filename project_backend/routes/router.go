package routes

import (
	"deepanshu18099/blockchain_ledger_backend/controllers"
	"deepanshu18099/blockchain_ledger_backend/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	ticketService := services.NewTicketService()
	ticketController := controllers.NewTicketController(ticketService)

	api := router.Group("/api")
	{
		api.POST("/tickets", ticketController.CreateTicket)
		api.GET("/tickets/user/:userID", ticketController.GetUserTickets)
		api.DELETE("/tickets/:ticketID", ticketController.DeleteTicket)
	}
}
