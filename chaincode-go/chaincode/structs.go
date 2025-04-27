package chaincode

type TicketDetails struct {
	TicketID               string        `json:"TicketID"`
	UserID                 string        `json:"Traveller"`
	ProviderID 			   string         `json:"Provider"`
	DateofTravel           string        `json:"DateofTravel"`
	Source                 string        `json:"Source"`
	Destination            string        `json:"Destination"`
	ModeofTravel           string        `json:"ModeofTravel"`
	TransportID            string        `json:"TransportID"` 
	SeatNumber             int32         `json:"SeatNumber"`
	Price                  float64       `json:"Price"`
	DateofBooking          string        `json:"DateofBooking"`
	DateofUpdate           string        `json:"DateofUpdate"`
	PaymentStatus          bool          `json:"PaymentStatus"` 
    IsActive               bool          `json:"IsActive"`
    DepartureTime          string        `json:"DepartureTime"`
	ArrivalTime            string        `json:"ArrivalTime"`
	JourneyDuration        string        `json:"JourneyDuration"`    
	Status                 string        `json:"Status"`
}

type TransportDetails struct {
	ID                     string                     `json:"ID"`
	Source                 string                     `json:"Source"`
	Destination            string                     `json:"Destination"`
	DepartureTime          string                     `json:"DepartureTime"`
	ArrivalTime            string                     `json:"ArrivalTime"`
	BasePrice              float64                    `json:"BasePrice"`
	Rating                 float64                    `json:"Rating"`
	RatingCount		       int32                      `json:"RatingCount"`
	Capacity               int32                      `json:"Capacity"`
	ModeofTravel           string                     `json:"ModeofTravel"`
	JourneyDuration        string                     `json:"JourneyDuration"`
	DateofTravel           []string                   `json:"DateofTravel"`
	SeatMap                map[string][]int32         `json:"AvailableSeats"`
	// Travellers             map[string][]string        `json:"Travellers"`
	ProviderID             string                     `json:"ProviderID"`
}

type User struct {
	Name            string         `json:"Name"`
	Email           string         `json:"Email"`
	Phone           string         `json:"Phone"`
	// PastTravels     []string       `json:"PastTravels"`
	// UpcomingTravels []string       `json:"UpcomingTravels"`  //list of ticketIDs
	BankBalance     float64        `json:"BankBalance"`
	Travels         []string       `json:"Travels"`
	IsAnonymous     bool           `json:"IsAnonymous"`
	PaymentID       []string       `json:"PaymentID"`
}

type Provider struct {
	Name             string             `json:"Name"`
	Email            string             `json:"Email"`
	Phone            string             `json:"Phone"`
	Services         []string           `json:"Services"` // list of transportIDs
	BankBalance      float64            `json:"BankBalance"`
	PaymentID        []string            `json:"PaymentID"`
}

type PaymentDetail struct {
	PaymentID         string       `json:"PaymentID"`
	From              string       `json:"From"`
	To                string       `json:"To"`
	Amount            float64       `json:"Amount"`
	PaymentTime       string       `json:"Payment Time"`
}
