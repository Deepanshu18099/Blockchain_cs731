package controllers

import (
	// "net/http"
	// "deepanshu18099/blockchain_ledger_backend/models"
	// "github.com/gin-gonic/gin"
	// "deepanshu18099/blockchain_ledger_backend/utils"

)


/*now in home, I want option in navbar to view previous bookings, which will result in new page listing user tickets, from ticket id in "userbookings" from this api response

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    const apiurl = process.env.REACT_APP_API_URL;
    try {
      const response = await axios.post(${apiurl}ledger/login, { email, password });

      if (response.status !== 200) {
        setErrors({ api: "Error signing in" });
        return;
      }
      console.log(response)
      const { token, userid, role, balance } = response.data

      login(token, userid, role, balance)

      navigate("/home");
    } catch (error) {
      console.error("Error signing in:", error);
      setErrors({ api: "Invalid credentials" });
    } finally {
      setLoading(false);
    }
  };

for each ticket in the list call GetDetailTicket api to get tickets in this form
		TicketID:        ticket.TicketID,
		DateofTravel:    ticket.DateofTravel,
		Source:          ticket.Source,
		Destination:     ticket.Destination,
		ModeofTravel:    ticket.ModeofTravel,
		TransportID:     ticket.TransportID,
		SeatNumber:      ticket.SeatNumber,
		Price:           ticket.Price,
		ArrivalTime:     ticket.ArrivalTime,
		DepartureTime:   ticket.DepartureTime,
		JourneyDuration: ticket.JourneyDuration,
		DateofBooking:   ticket.DateofBooking,
		DateofUpdate:    ticket.DateofUpdate,
		PaymentStatus:   ticket.PaymentStatus,
		IsActive:        ticket.IsActive,
		Status:          ticket.Status,
*/

// func GetTickets(c *gin.Context) {
// 	/*
// 	Input: Token
// 	GIVE List of tickets by that user if user, else listings by provider
// 	Output: List of tickets(IDS) with details
// 	*/

// 	// Get the user ID from the token
// 	claims, err := utils.Authcheck(c)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}
// 	email := claims["email"].(string)

// }
// func GetTicketByID(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "ticket details"})
// }


