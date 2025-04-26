package controllers

import (
	"net/http"
	"deepanshu18099/blockchain_ledger_backend/models"
	"deepanshu18099/blockchain_ledger_backend/services"
	"github.com/gin-gonic/gin"
)

// TicketController struct to group related functions
type TicketController struct {
	TicketService *services.TicketService
}

// NewTicketController initializes and returns the controller
func NewTicketController(ticketService *services.TicketService) *TicketController {
	return &TicketController{
		TicketService: ticketService,
	}
}

// CreateTicket godoc
// POST /api/tickets
func (tc *TicketController) CreateTicket(c *gin.Context) {
	var req models.TicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := tc.TicketService.CreateTicket(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// GetUserTickets godoc
// GET /api/tickets/user/:userID
func (tc *TicketController) GetUserTickets(c *gin.Context) {
	userID := c.Param("userID")
	tickets, err := tc.TicketService.GetTicketsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

// DeleteTicket godoc
// DELETE /api/tickets/:ticketID
func (tc *TicketController) DeleteTicket(c *gin.Context) {
	ticketID := c.Param("ticketID")
	err := tc.TicketService.DeleteTicket(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found or cannot be deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}
