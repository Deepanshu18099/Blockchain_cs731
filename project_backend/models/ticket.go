package models



type Ticket struct {
	ID          string `gorm:"primaryKey" json:"id"`
	UserID      string `json:"user_id"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Seat        string `json:"seat"`
	Price       float64 `json:"price"`
	Status      string `json:"status"` // e.g. PENDING, CONFIRMED
}
type Ticket1 struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Destination string  `json:"destination"`
	Source      string  `json:"source"`
}
// Used for input validation
type TicketRequest struct {
	UserID      string  `json:"user_id" binding:"required"`
	Source      string  `json:"source" binding:"required"`
	Destination string  `json:"destination" binding:"required"`
	Date        string  `json:"date" binding:"required"`
	Time        string  `json:"time" binding:"required"`
	Seat        string  `json:"seat" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}
