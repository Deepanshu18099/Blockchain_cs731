package controllers

import (
	"net/http"
	"deepanshu18099/blockchain_ledger_backend/models"
	"github.com/gin-gonic/gin"

)


var tickets = []models.Ticket1{
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

func GetTickets(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tickets)
}
func GetTicketByID(c *gin.Context) {
	id := c.Param("id")
	for _, ticket := range tickets {
		if ticket.ID == id {
			c.JSON(http.StatusOK, ticket)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "ticket not found"})
}
func CreateTicket(c *gin.Context) {
	var newTicket models.Ticket1
	if err := c.ShouldBindJSON(&newTicket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tickets = append(tickets, newTicket)
	c.JSON(http.StatusCreated, newTicket)
}

func UpdateTicket(c *gin.Context) {
	id := c.Param("id")
	var updatedTicket models.Ticket1
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

func DeleteTicket(c *gin.Context) {
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
