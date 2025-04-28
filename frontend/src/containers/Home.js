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
      { id: "1", source: "Delhi", destination: "Mumbai", price: 5000, date: "2023-10-01", starttime: "10:00", endtime: "12:00", seatleft: 50, providerId: "provider1" },
      { id: "2", source: "Bangalore", destination: "Chennai", price: 3000, date: "2023-10-02", starttime: "14:00", endtime: "15:30", seatleft: 30, providerId: "provider2" }
    ],
    train: [
      { id: "3", source: "Delhi", destination: "Kanpur", price: 800, date: "2023-10-03", starttime: "08:00", endtime: "10:00", seatleft: 100, providerId: "provider3" },
      { id: "4", source: "Pune", destination: "Hyderabad", price: 1200, date: "2023-10-04", starttime: "18:00", endtime: "20:30", seatleft: 20, providerId: "provider4" }
    ],
    bus: [
      { id: "5", source: "Mumbai", destination: "Goa", price: 600, date: "2023-10-05", starttime: "22:00", endtime: "06:00", seatleft: 40, providerId: "provider5" },
      { id: "6", source: "Chennai", destination: "Pondicherry", price: 500, date: "2023-10-06", starttime: "09:00", endtime: "11:00", seatleft: 60, providerId: "provider6" }
    ]
  };

  const { userId, role } = useAuth();
  const [mode, setMode] = useState("flight");
  const [source, setSource] = useState("");
  const [destination, setDestination] = useState("");
  const [date, setDate] = useState("");
  const [showDropdown, setShowDropdown] = useState(false);
  const [addMoneyAmount, setAddMoneyAmount] = useState("");
  const [transportId, setTransportId] = useState("");
  const [price, setPrice] = useState("");
  const [startdate, setStartDate] = useState("");
  const [enddate, setEndDate] = useState("");
  const [seatcount, setSeatCount] = useState(0);
  const [travelOptions, setTravelOptions] = useState([]);
  const token = localStorage.getItem("token");
  var curr_balance = localStorage.getItem("balance")

  const [balance, setBalance] = useState(Number(curr_balance)); // Initialize balance state


  // navigate
  const navigate = useNavigate();


  const handleAddTransport = async () => {
    if (!source || !destination || !startdate || !enddate || !price) {
      alert("Please fill in all fields.");
      return;
    }
    try {
      const response = await axios.post(
        "http://localhost:8080/Addtransport",
        {
          mode,
          source,
          destination,
          startdate,
          enddate,
          price,
          seatcount
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );
      if (response.status === 200) {
        alert("Transport option added successfully!");
        setSource("");
        setDestination("");
        setStartDate("");
        setEndDate("");
        setPrice("");

        // get transport id and show it
        setTransportId(response.data.transport_id);
        // adding visiblity of transport id box
        const transportIdBox = document.querySelector(".transport-id-box");
        transportIdBox.classList.remove("hidden");
        transportIdBox.classList.add("visible");
        
        setTimeout(() => {
          setTransportId("");
          // now hide the transport id box
          transportIdBox.classList.remove("visible");
          transportIdBox.classList.add("hidden");
          
        }
        , 10000); // clear transport id after 10 seconds

        // visible show transport id box for 10 seconds
      } else {
        alert("Failed to add transport option.");
      }
    } catch (error) {
      console.error(error);
      alert("Failed to add transport option.");
    }
  }
  // set travel options to sample data in case of no data
  if (travelOptions.length === 0) {
    setTravelOptions(sampleData);
  }


  const set_data = (data) => {
    // now make the object from parsing the data
    const newData = data.map((item) => ({
      id: item.ID,
      source: item.Source,
      destination: item.Destination,
      price: item.BasePrice,
      date: item.DateofTravel,
      starttime: item.DepartureTime,
      endtime: item.ArrivalTime,
      seatleft: item.Capacity,
      providerId: item.ProviderID
    }));
    setTravelOptions(newData);
  }

  const handleBooking = async() => {
    // alert(`Booking a ${mode} ticket from ${source} to ${destination} on ${date}`);
    if (!source || !destination || !date) {
      alert("Please fill in all fields.");
      return;
    }
    // make get api call to get the data of selected option
    try{
      const response = await axios.get(
        `http://localhost:8080/Gettransports/${mode}/${source}/${destination}/${date}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );
      // console.log(response);
      console.log(response.data);
      if (response.status === 200) {
        console.log("Data fetched successfully");
        alert("Data fetched successfully");
        // navigate(`/details/${response.data.id}`);
      } else {
        alert("Failed to fetch data.");
      }
      // console.log(response);
      // now response.data.transports will be a list of transports
      /*
      structure of response.data:
      newTransport := TransportDetails{
        ID:              transportID,
        Source:          source,
        Destination:     dest,
        DepartureTime:   dept,
        ArrivalTime:     arrt,
        BasePrice:       basep,
        Rating:          3.00, all new travels will have a rating of 3 at start
        Capacity:        cap,
        ModeofTravel:    mode,
        JourneyDuration: totalt,
        DateofTravel:    dates,
        SeatMap:         seatMap,
        // Travellers: travellerMap,
        ProviderID: providerID,
      */

      // Now using the response data, we can show the available options
      set_data(response.data.transports);
    }
    catch (error) {
      console.error(error);
      alert("Failed to fetch data.");
    }
  };

  
  const handleSignOut = () => {
    localStorage.removeItem("token");
    navigate("/signin");
  };

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
        setBalance((prevBalance) => Number(prevBalance) + Number(addMoneyAmount));
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
          {/* for role of user */}
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
      {
        role === "user" && (
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
        )
      }

{role === "provider" && (
  <div className="flex flex-col items-center justify-center p-6">
    <div className="bg-white p-8 rounded-2xl shadow-lg w-full max-w-md">
      <h2 className="text-2xl font-bold mb-6 text-center text-gray-800">
        Welcome to MyTravel.com
      </h2>

      {userId ? (
        <>
          {/* Add Balance Section */}
          <div className="mb-8">
            <h3 className="text-lg font-semibold mb-4 text-gray-700">Add Balance</h3>
            <div className="flex flex-col gap-4">
              <input
                type="number"
                value={addMoneyAmount}
                onChange={(e) => setAddMoneyAmount(e.target.value)}
                placeholder="Enter amount to add"
                className="w-full border rounded-lg p-3 focus:outline-none focus:ring-2 focus:ring-green-400"
              />
              <button
                onClick={handleAddMoney}
                className="bg-green-500 text-white py-2 rounded-lg hover:bg-green-600 transition"
              >
                Add Money
              </button>
            </div>
          </div>

          {/* Add Transport Section */}
          <div className="mb-8">
            <h3 className="text-lg font-semibold mb-4 text-gray-700">Add Transport Option</h3>
            <div className="flex flex-col gap-4">
              {/* area to show transport id after submission */}

              <div className="transport-id-box hidden mb-4">
                <p className="text-green-600">You new Transport ID: {transportId}</p>
                <p className="text-gray-600">Note this in 10 seconds</p>
              </div>

              <select
                value={mode}
                onChange={(e) => setMode(e.target.value)}
                className="w-full border rounded-lg p-3 focus:outline-none focus:ring-2 focus:ring-blue-400"
              >
                <option value="flight">Flight</option>
                <option value="train">Train</option>
                <option value="bus">Bus</option>
              </select>

              <input
                type="text"
                value={source}
                onChange={(e) => setSource(e.target.value)}
                placeholder="Enter source"
                className="w-full border rounded-lg p-3 focus:outline-none focus:ring-2 focus:ring-blue-400"
              />

              <input
                type="text"
                value={destination}
                onChange={(e) => setDestination(e.target.value)}
                placeholder="Enter destination"
                className="w-full border rounded-lg p-3 focus:outline-none focus:ring-2 focus:ring-blue-400"
              />

              <label className="block font-semibold mb-1">Start Date</label>
              <input
                type="date"
                value={startdate}
                onChange={(e) => setStartDate(e.target.value)}
                className="w-full border rounded-lg p-3 focus:outline-none focus:ring-2 focus:ring-blue-400"
              />

              <label className="block font-semibold mb-1">End Date</label>
              <input
                type="date"
                value={enddate}
                onChange={(e) => setEndDate(e.target.value)}
                className="w-full border rounded-lg p-3 focus:outline-none focus:ring-2 focus:ring-blue-400"
              />


              <label className="block font-semibold mb-1">Set the Base Price</label>
              <input
                type="number"
                value={price}
                onChange={(e) => setPrice(e.target.value)}
                placeholder="Enter price"
                className="w-full border rounded-lg p-3 focus:outline-none focus:ring-2 focus:ring-blue-400"
              />

              <label className="block font-semibold mb-1">Set the Seat Count</label>
              <input
                type="number"
                value={seatcount}
                onChange={(e) => setSeatCount(e.target.value)}
                placeholder="Enter seat count"
                className="w-full border rounded-lg p-3 focus:outline-none focus:ring-2 focus:ring-blue-400"
              />

              <button
                onClick={handleAddTransport}
                className="bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600 transition"
              >
                Submit Transport
              </button>
            </div>
          </div>
        </>
      ) : (
        <div className="mb-8">
          <p className="mb-4 text-center text-gray-600">Please sign up or sign in to continue.</p>
          <div className="flex justify-center gap-6">
            <Link to="/signup" className="text-blue-500 hover:underline">
              Sign Up
            </Link>
            <Link to="/signin" className="text-blue-500 hover:underline">
              Sign In
            </Link>
          </div>
        </div>
      )}
    </div>
  </div>
)}

        {/* ending divs*/}
      {/* Footer */}
      <div className="bg-gray-800 text-white py-4 text-center">
        <p>&copy; 2023 MyTravel.com. All rights reserved.</p>
      </div>
    </div>
  );
}


export default Home;
