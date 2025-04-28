// <Route path="/details/:ticketId" element={<Confirmticket />} />

import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from './Authcontext.js';
import axios from 'axios';
import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';

// const sampleData = [
//   { id: 1, transportName: "Indigo Flight 6E-205", source: "Delhi", destination: "Mumbai", date: "2025-05-10", time: "10:00 AM", price: 5000, availableSeats: [1, 2, 3, 4, 5] },
//   { id: 2, transportName: "Rajdhani Express", source: "Delhi", destination: "Kanpur", date: "2025-05-12", time: "5:00 PM", price: 1200, availableSeats: [6, 7, 8, 9] },
//   { id: 3, transportName: "Volvo Bus A/C", source: "Mumbai", destination: "Goa", date: "2025-05-15", time: "9:00 PM", price: 800, availableSeats: [10, 11, 12] }
// ];

const ConfirmTicket = () => {
  const { transportid } = useParams();
  const navigate = useNavigate();
  const { userId } = useAuth();

  // useEffect(() => {
  //   // Fetch the transport details based on id (simulate for now)
  //   console.log("hii")
  //   console.log(transportid)
  //   // const transport = sampleData.find(item => item.id === parseInt(transportid));
  //   console.log("hdddii")
  //   // console.log(transport.transportName)

  //   setTransportDetails(transport);
  // }, [transportid]);

                      
  // load the data from state(navigation from prev page) onClick={() => navigate(`/details/${option.id}`, { state: { option } })}
  const location = useLocation();
  const { option } = location.state || {};

  const [transportDetails, setTransportDetails] = useState(null);
  const [selectedSeat, setSelectedSeat] = useState("");

  useEffect(() => {
    if (option) {
      const transport = {
        date: option["date"],
        time: option["starttime"],
        source: option["source"],
        destination: option["destination"],
        price: option["price"],
        availableSeats: option["seatmap"],
        seatleft: option["seatleft"],
        ReachingTime: option["endtime"],
        transportid: option["id"],
      };
      setTransportDetails(transport);
    }
  }, [option]); // <-- Run only when `option` changes

  const handleConfirmBooking = async () => {
    if (!selectedSeat) {
      alert("Please select a seat before confirming!");
      return;
    }
  
    const apiurl = process.env.REACT_APP_API_URL;
    const token = localStorage.getItem("token");
  
    try {
      console.log(selectedSeat);
      const response = await axios.post(
        `${apiurl}tickets`,
        {
          transportid: transportDetails["transportid"],
          date: transportDetails["date"],
          seatnumber: selectedSeat,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      console.log("Booking response:", response.data);

      if (response.status === 200) {
        alert("Booking confirmed successfully!");
        /*
        New balance should be updated
        	c.JSON(http.StatusOK, gin.H{
		"message":        "Ticket booked successfully",
		"updatedbalance": clean_output["BankBalance"],
		"transaction_id": clean_output["transactionID"],
	})
        */
        // Update the user's balance in local storage or context
        const updatedBalance = response.data.updatedbalance;
        localStorage.setItem("balance", updatedBalance);
        navigate("/home");
      } else {
        alert("Error confirming booking. Please try again.");
      }
    } catch (error) {
      console.error("Error confirming booking:", error);
      alert("Error confirming booking. Please try again.");
    }
  };

  if (!option) {
    console.log("No data found!");
    return <div>No data found!</div>;
  }

  console.log("Transport Details:", transportDetails);
  if (!transportDetails) {
    return <div className="flex justify-center items-center min-h-screen">Loading...</div>;
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 p-6">
      <div className="bg-white p-6 rounded shadow-md w-full max-w-md">
        <h2 className="text-2xl font-bold mb-4 text-center">Confirm Your Ticket</h2>

        {/* showing the data */}
        <div className="mb-4">
          <h3 className="text-lg font-semibold mb-2">Transport Details</h3>
          <p><strong>TransportId :</strong> {transportDetails.transportid}</p>
          <p><strong>Source:</strong> {transportDetails.source}</p>
          <p><strong>Destination:</strong> {transportDetails.destination}</p>
          <p><strong>Date:</strong> {transportDetails.date}</p>
          <p><strong>Time:</strong> {transportDetails.time}:{transportDetails.ReachingTime}</p>
          <p><strong>Price:</strong> ₹{transportDetails.price}</p>
          <p><strong>Available Seats:</strong> {transportDetails.seatleft}</p>
        </div>

        <div className="mb-4">
          <label className="block font-semibold mb-2">Select Seat</label>
            <select
                value={selectedSeat}
                onChange={(e) => setSelectedSeat(e.target.value)}
                className="block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring focus:ring-blue-500"
            >
                <option value="">Select a seat</option>
                {transportDetails.availableSeats.map(seat => (
                    <option key={seat} value={seat}>Seat {seat}</option>
                ))}
            </select>
        </div>
        <button
          onClick={handleConfirmBooking}
          className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
        >
            Confirm Booking
        </button>
        </div>
        </div>
    );
}
export default ConfirmTicket;
