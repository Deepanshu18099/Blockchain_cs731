// Home.js
import React, { useState } from "react";
import { useAuth } from "./Authcontext"; 
import { Link } from "react-router-dom";
import { FaUserCircle } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import axios from "axios";

function Home() {
  
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

  const { userId, role } = useAuth();
  const [mode, setMode] = useState("flight");
  const [source, setSource] = useState("");
  const [destination, setDestination] = useState("");
  const [date, setDate] = useState("");
  const [showDropdown, setShowDropdown] = useState(false);
  const [addMoneyAmount, setAddMoneyAmount] = useState("");
  const [balance, setBalance] = useState(0); // Initialize balance state
  const token = localStorage.getItem("token");


  // navigate
  const navigate = useNavigate();

  const handleBooking = () => {
    alert(`Booking a ${mode} ticket from ${source} to ${destination} on ${date}`);
  };

  
  const handleSignOut = () => {
    localStorage.removeItem("token");
    navigate("/signin");
  };
  const travelOptions = sampleData[mode] || [];

  const handleAddMoney = async () => {
    if (!addMoneyAmount) {
      alert("Please enter an amount to add.");
      return;
    }
    try {
      const response = await axios.post(
        "http://localhost:8080/ledger/addMoney",
        { Amount: Number(addMoneyAmount) },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );
      if (response.status === 200) {
        alert("Money added successfully!");
        setBalance((prevBalance) => prevBalance + Number(addMoneyAmount));
        setAddMoneyAmount("");
      } else {
        alert("Failed to add money.");
      }
    } catch (error) {
      console.error(error);
      alert("Failed to add money.");
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Navbar */}
      <div className="flex justify-between items-center bg-white p-4 shadow-md">
        <div className="text-xl font-bold">MyTravel.com</div>
        <div className="flex items-center gap-6">
          {userId && (
            <>
              <div className="text-sm">
                <p><strong>User ID:</strong> {userId}</p>
                <p><strong>Role:</strong> {role}</p>
                <p><strong>Balance:</strong> ₹{balance}</p>
              </div>
              <div className="relative">
                <FaUserCircle
                  size={32}
                  className="cursor-pointer"
                  onClick={() => setShowDropdown(!showDropdown)}
                />
                {showDropdown && (
                  <div className="absolute right-0 mt-2 w-40 bg-white rounded shadow-md z-10">
                    <button
                      onClick={handleSignOut}
                      className="w-full text-left px-4 py-2 hover:bg-gray-100"
                    >
                      Sign Out
                    </button>
                  </div>
                )}
              </div>
            </>
          )}
        </div>
      </div>

      {/* Main Content */}
      <div className="flex flex-col items-center justify-center p-6">
        <div className="bg-white p-6 rounded shadow-md w-full max-w-md">
          <h2 className="text-xl font-bold mb-4 text-center">Welcome to MyTravel.com</h2>

          {userId ? (
            <div className="mb-6">
              <div className="flex flex-col gap-2">
                <input
                  type="number"
                  value={addMoneyAmount}
                  onChange={(e) => setAddMoneyAmount(e.target.value)}
                  placeholder="Enter amount to add"
                  className="w-full border rounded p-2"
                />
                <button
                  onClick={handleAddMoney}
                  className="bg-green-500 text-white py-2 rounded hover:bg-green-600"
                >
                  Add Money
                </button>
              </div>
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

          {/* Travel Mode Selector */}
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

          {/* Source Input */}
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

          {/* Destination Input */}
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

          {/* Date Picker */}
          <div className="mb-4">
            <label className="block font-semibold mb-1">Date of Travel</label>
            <input
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              className="w-full border rounded p-2"
            />
          </div>

          {/* Fetch Options Button */}
          <button
            onClick={handleBooking}
            className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
          >
            Fetch Available Options
          </button>

          {/* Available Options */}
          <div className="mt-6">
            <h3 className="text-lg font-semibold mb-2">Available Options</h3>
            {travelOptions.map((option) => (
              <div
                key={option.id}
                className="border p-4 mb-2 rounded cursor-pointer hover:shadow"
                onClick={() => navigate(`/details/${option.id}`)}
              >
                <p><strong>Source:</strong> {option.source}</p>
                <p><strong>Destination:</strong> {option.destination}</p>
                <p><strong>Date:</strong> {option.date}</p>
                <p><strong>Price:</strong> ₹{option.price}</p>
                <p className="text-blue-500 hover:underline mt-2">View Details</p>
              </div>       
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}


export default Home;
