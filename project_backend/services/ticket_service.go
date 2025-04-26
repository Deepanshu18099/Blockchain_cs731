package services

import (
	"deepanshu18099/blockchain_ledger_backend/database"
	"deepanshu18099/blockchain_ledger_backend/models"

	"github.com/google/uuid"
)

type TicketService struct{}

func NewTicketService() *TicketService {
	return &TicketService{}
}

func (s *TicketService) CreateTicket(req models.TicketRequest) (models.Ticket, error) {
	ticket := models.Ticket{
		ID:          uuid.New().String(),
		UserID:      req.UserID,
		Source:      req.Source,
		Destination: req.Destination,
		Date:        req.Date,
		Time:        req.Time,
		Seat:        req.Seat,
		Price:       req.Price,
		Status:      "PENDING", // Can be updated after payment
	}

	result := database.DB.Create(&ticket)
	return ticket, result.Error
}

func (s *TicketService) GetTicketsByUser(userID string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := database.DB.Where("user_id = ?", userID).Find(&tickets).Error
	return tickets, err
}

func (s *TicketService) DeleteTicket(ticketID string) error {
	return database.DB.Delete(&models.Ticket{}, "id = ?", ticketID).Error
}
