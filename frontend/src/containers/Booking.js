import React, { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useParams } from "react-router-dom";
import { useLocation } from "react-router-dom";


/*
Ticket structure in go:
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
const Bookings = () => {
  const [tickets, setTickets] = useState([]);
  const navigate = useNavigate();
  const apiurl = process.env.REACT_APP_API_URL;
  const token = localStorage.getItem("token");

  useEffect(() => {
    const fetchBookings = async () => {
      const apiurl = process.env.REACT_APP_API_URL;
      const token = localStorage.getItem("token");
      if (!token) {
        alert("Not authorized");
        navigate("/signin");
        return;
      }

      try {
        // 1. Get the user bookings
        const response = await axios.get(`${apiurl}tickets`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        const userBookings = response.data.tickets; // assuming this is an array of TicketIDs

        // it will be string from user.Travels = append(user.Travels, ticketID)
        console.log("User Bookings:", userBookings);

        // 2. Now fetch ticket details for each TicketID
        const ticketDetailsPromises = userBookings.map((ticketId) =>
          axios.get(`${apiurl}tickets/${ticketId}`, {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          })
        );

        const ticketsResponse = await Promise.all(ticketDetailsPromises);
        console.log("Tickets Response:", ticketsResponse);
        const ticketData = ticketsResponse.map((res) => res.data.ticket);

        setTickets(ticketData);
      } catch (error) {
        console.error("Error fetching bookings:", error);
        alert("Failed to fetch bookings");
      }
    };

    fetchBookings();
  }, [navigate]);

  
  const handleDelete = async (ticketId) => {
    if (!window.confirm("Are you sure you want to delete this ticket?")) {
      return;
    }
    try {
      // .DELETE("/tickets/:id
      const response = await axios.delete(`${apiurl}tickets/${ticketId}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });

      if (response.status === 200) {
        alert("Ticket deleted successfully!");
        // Refresh the list
        setTickets((prev) => prev.filter((ticket) => ticket.TicketID !== ticketId));
      } else {
        alert("Failed to delete ticket.");
      }
    } catch (error) {
      console.error("Error deleting ticket:", error);
      alert("Failed to delete ticket.");
    }
  };

  const handleUpdate = async (ticketId) => {
    // new page is needed
    // ask for newTicketDate, and newSeatNumber
    const newTicketDate = prompt("Enter new travel date (YYYY-MM-DD):");
    const newSeatNumber = prompt("Enter new seat number:");
    if (!newTicketDate || !newSeatNumber) {
      alert("Please provide both new travel date and seat number.");
      return;
    }
    try {
      const response = await axios.put(
        `${apiurl}tickets/${ticketId}`,
        {
          DateofTravel: newTicketDate,
          SeatNumber: newSeatNumber,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );

      if (response.status === 200) {
        alert("Ticket updated successfully!");
        // Refresh the list
        var NewTicketID = response.data.NewTicketID;
        setTickets((prev) =>
          prev.map((ticket) =>
            ticket.TicketID === ticketId ? { ...ticket, DateofTravel: newTicketDate, SeatNumber: newSeatNumber, TicketID: NewTicketID } : ticket
          )
        );
      } else {
        alert("Failed to update ticket.");
      }
    } catch (error) {
      console.error("Error updating ticket:", error);
      alert("Failed to update ticket.");
    }
    
  };

  if (tickets.length === 0) return <div className="text-center mt-8">No bookings found.</div>;



  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold mb-4 text-center">My Bookings</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {tickets.map((ticket) => (
          <div key={ticket.TicketID} className="border rounded-lg shadow-md p-4 bg-white">
            <p><strong>Ticket ID:</strong> {ticket.TicketID}</p>
            <p><strong>Travel Date:</strong> {ticket.DateofTravel}</p>
            <p><strong>Source:</strong> {ticket.Source}</p>
            <p><strong>Destination:</strong> {ticket.Destination}</p>
            <p><strong>Mode:</strong> {ticket.ModeofTravel}</p>
            <p><strong>Seat:</strong> {ticket.SeatNumber}</p>
            <p><strong>Price:</strong> ₹{ticket.Price}</p>
            <p><strong>Status:</strong> {ticket.Status}</p>
            <div className="flex gap-4 mt-4">
              <button
                onClick={() => handleDelete(ticket.TicketID)}
                className="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600"
              >
                Delete
              </button>
              <button
                onClick={() => handleUpdate(ticket.TicketID)}
                className="bg-blue-500 text-white px-3 py-1 rounded hover:bg-blue-600"
              >
                Update
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Bookings;
