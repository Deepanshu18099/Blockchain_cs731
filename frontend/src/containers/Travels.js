import React, { useEffect, useState } from "react";
import axios from "axios";
import { useAuth } from "./Authcontext.js";
import { useNavigate } from "react-router-dom";



/*
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
*/
const Travels = () => {
  const { token } = useAuth();
  const [travels, setTravels] = useState([]);
  const [loading, setLoading] = useState(true);

  const apiurl = process.env.REACT_APP_API_URL;
  const navigate = useNavigate();

  useEffect(() => {
    const fetchTravels = async () => {
        if (!token) {
            alert("Not authorized");
            navigate("/signin");
        }
      try {
        setLoading(true);
        const response = await axios.get(`${apiurl}GetTransports`, {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        });
        console.log("Response:", response);

        
        setTravels(response.data.transports);
        console.log("Owned transports:", response.data.transports);
      } catch (error) {
        console.error("Error fetching travels:", error);
      }
    setLoading(false);
    };

    fetchTravels();
  }, [apiurl, token]);

  const handleDeleteTravel = async (transportId) => {
//     if (!window.confirm("Are you sure you want to delete this travel? All future tickets will be cancelled.")) {
//       return;
//     }
//     try {
//       const response = await axios.delete(`${apiurl}DeleteTravel/${transportId}`, {
//         headers: {
//           Authorization: `Bearer ${token}`,
//           "Content-Type": "application/json",
//         },
//       });

//       if (response.status === 200) {
//         alert("Travel deleted and future tickets cancelled successfully!");
//         setTravels((prev) => prev.filter((travel) => travel.TransportID !== transportId));
//       } else {
//         alert("Failed to delete travel.");
//       }
//     } catch (error) {
//       console.error("Error deleting travel:", error);
//       alert("Failed to delete travel.");
//     }
  };

  if (loading) return <div className="text-center mt-8">Loading your travels...</div>;

  if (travels.length === 0) return <div className="text-center mt-8">No travels found.</div>;

  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold mb-4 text-center">My Travels</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {travels.map((travel) => (
          <div key={travel.TransportID} className="bg-white p-4 rounded shadow-md">
            <h3 className="text-lg font-semibold">{travel.Source} to {travel.Destination}</h3>
            <p>Departure: {travel.DepartureTime}</p>
            <p>Arrival: {travel.ArrivalTime}</p>
            <p>Base Price: ${travel.BasePrice}</p>
            <p>Rating: {travel.Rating} ({travel.RatingCount} reviews)</p>
            <p>Capacity: {travel.Capacity}</p>
            {/* first key of seatmap will contain array of seats available */}
            {/* <p>Capacity: {travel.SeatMap[Object.keys(travel.SeatMap)[0]].length} seats available</p> */}
            <p>Mode of Travel: {travel.ModeofTravel}</p>
            {/* <p>Date of travels: {travel.DateofTravel.join(", ")[0]}, {travel.DateofTravel.join(", ")[1]}</p> */}
            <button
              onClick={() => handleDeleteTravel(travel.TransportID)}
              className="mt-2 bg-red-500 text-white px-4 py-2 rounded"
            >
              Delete Travel
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Travels;
