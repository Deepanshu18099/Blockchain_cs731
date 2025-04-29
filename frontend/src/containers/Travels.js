import React, { useEffect, useState } from "react";
import axios from "axios";
import { useAuth } from "./Authcontext.js";

const Travels = () => {
  const { token } = useAuth();
  const [travels, setTravels] = useState([]);
  const [loading, setLoading] = useState(true);

  const apiurl = process.env.REACT_APP_API_URL;

  useEffect(() => {
    const fetchTravels = async () => {
      try {
        const response = await axios.get(`${apiurl}GetTransports`, {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        });

        const { ownedtransports } = response.data;
        console.log("Owned transports:", ownedtransports);

        const travelDetails = await Promise.all(
          ownedtransports.map(async (transportId) => {
            const detailRes = await axios.get(`${apiurl}GetDetailTravel/${transportId}`, {
              headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "application/json",
              },
            });
            return detailRes.data;
          })
        );

        console.log("Detailed travels:", travelDetails);
        setTravels(travelDetails);
      } catch (error) {
        console.error("Error fetching travels:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchTravels();
  }, [apiurl, token]);

  const handleDeleteTravel = async (transportId) => {
    if (!window.confirm("Are you sure you want to delete this travel? All future tickets will be cancelled.")) {
      return;
    }
    try {
      const response = await axios.delete(`${apiurl}DeleteTravel/${transportId}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });

      if (response.status === 200) {
        alert("Travel deleted and future tickets cancelled successfully!");
        setTravels((prev) => prev.filter((travel) => travel.TransportID !== transportId));
      } else {
        alert("Failed to delete travel.");
      }
    } catch (error) {
      console.error("Error deleting travel:", error);
      alert("Failed to delete travel.");
    }
  };

  if (loading) return <div className="text-center mt-8">Loading your travels...</div>;

  if (travels.length === 0) return <div className="text-center mt-8">No travels found.</div>;

  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold mb-4 text-center">My Travels</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {travels.map((travel) => (
          <div key={travel.TransportID} className="border rounded-lg shadow-md p-4 bg-white">
            <p><strong>Transport ID:</strong> {travel.TransportID}</p>
            <p><strong>Travel Date:</strong> {travel.DateofTravel}</p>
            <p><strong>Source:</strong> {travel.Source}</p>
            <p><strong>Destination:</strong> {travel.Destination}</p>
            <p><strong>Mode:</strong> {travel.ModeofTravel}</p>
            <p><strong>Seats Available:</strong> {travel.SeatsAvailable}</p>
            <p><strong>Price:</strong> ₹{travel.Price}</p>
            <div className="flex gap-4 mt-4">
              <button
                onClick={() => handleDeleteTravel(travel.TransportID)}
                className="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600"
              >
                Delete Travel
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Travels;
