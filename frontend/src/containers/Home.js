// Home.js
import React, { useState } from "react";
import { useAuth } from "./Authcontext"; 
import { Link } from "react-router-dom";

function Home() {
  const { userId, role, balance } = useAuth();

  const sampleData = {
    flight: [
      { id: 1, source: "Delhi", destination: "Mumbai", price: 5000, date: "2023-10-01" },
      { id: 2, source: "Bangalore", destination: "Chennai", price: 3000, date: "2023-10-02" }
    ],
    train: [
      { id: 3, source: "Delhi", destination: "Kanpur", price: 800, date: "2023-10-03" },
      { id: 4, source: "Pune", destination: "Hyderabad", price: 1200, date: "2023-10-04" }
    ],
    bus: [
      { id: 5, source: "Mumbai", destination: "Goa", price: 600, date: "2023-10-05" },
      { id: 6, source: "Chennai", destination: "Pondicherry", price: 500, date: "2023-10-06" }
    ]
  };

  const [mode, setMode] = useState("flight");
  const [source, setSource] = useState("");
  const [destination, setDestination] = useState("");
  const [date, setDate] = useState("");

  const handleBooking = () => {
    alert(`Booking a ${mode} ticket from ${source} to ${destination} on ${date}`);
  };

  const travelOptions = sampleData[mode] || [];

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 p-6">
      <div className="bg-white p-6 rounded shadow-md w-full max-w-md">
        <h2 className="text-xl font-bold mb-4 text-center">Welcome to MyTravel.com</h2>
        
        {userId ? (
          <div className="mb-6">
            <p className="text-sm">User ID: {userId}</p>
            <p className="text-sm">Role: {role}</p>
            <p className="text-sm mb-2">Balance: {balance}</p>
          </div>
        ) : (
          <div className="mb-6">
            <p className="mb-2">Please sign up or sign in to continue.</p>
            <div className="flex justify-center gap-4">
              <Link to="/signup" className="text-blue-500 hover:underline">
                Sign Up
              </Link>
              <Link to="/signin" className="text-blue-500 hover:underline">
                Sign In
              </Link>
            </div>
          </div>
        )}

        <div className="mb-4">
          <label className="block font-semibold mb-1">Mode of Travel</label>
          <select
            value={mode}
            onChange={(e) => setMode(e.target.value)}
            className="w-full border rounded p-2"
          >
            <option value="flight">Flight</option>
            <option value="train">Train</option>
            <option value="bus">Bus</option>
          </select>
        </div>

        <div className="mb-4">
          <label className="block font-semibold mb-1">Source</label>
          <input
            type="text"
            value={source}
            onChange={(e) => setSource(e.target.value)}
            placeholder="Enter source"
            className="w-full border rounded p-2"
          />
        </div>

        <div className="mb-4">
          <label className="block font-semibold mb-1">Destination</label>
          <input
            type="text"
            value={destination}
            onChange={(e) => setDestination(e.target.value)}
            placeholder="Enter destination"
            className="w-full border rounded p-2"
          />
        </div>

        <div className="mb-4">
          <label className="block font-semibold mb-1">Date of Travel</label>
          <input
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            className="w-full border rounded p-2"
          />
        </div>

        <button
          onClick={handleBooking}
          className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
        >
          Fetch Available Options
        </button>

        <div className="mt-6">
          <h3 className="text-lg font-semibold mb-2">Available Options</h3>
          {travelOptions.map((option) => (
            <div key={option.id} className="border p-3 rounded mb-2">
              <p><strong>From:</strong> {option.source}</p>
              <p><strong>To:</strong> {option.destination}</p>
              <p><strong>Price:</strong> ₹{option.price}</p>
              <p><strong>Date:</strong> {option.date}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default Home;
