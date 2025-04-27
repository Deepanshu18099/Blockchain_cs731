// this will be sign in page for hyperledger page, signing using email and password and will return jwt token
// tailwind css will be used for styling, basic structure
import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { useEffect } from 'react';
import { useAuth } from "./Authcontext";  // import

const Signin = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState({ api: "" });
  const navigate = useNavigate();
  const { login } = useAuth();   // get login function

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    const apiurl = process.env.REACT_APP_API_URL;
    try {
      const response = await axios.post(`${apiurl}ledger/login`, { email, password });

      if (response.status !== 200) {
        setErrors({ api: "Error signing in" });
        return;
      }

      console.log("Response from API:", response.data);

      const { token, role, userid } = response.data;

      login(token, userid, role);  // update context + localStorage

      navigate("/home");
    } catch (error) {
      console.error("Error signing in:", error);
      setErrors({ api: "Invalid credentials" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-6 rounded shadow-md w-80">
        <h2 className="text-lg font-semibold mb-4">Sign In</h2>
        <p className="mb-4">Please fill in the form to sign in.</p>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="email" className="block text-sm font-medium text-gray-700">
              Email
            </label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring focus:ring-blue-500"
              disabled={loading}
            />
          </div>
          <div className="mb-4">
            <label htmlFor="password" className="block text-sm font-medium text-gray-700">
              Password
            </label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring focus:ring-blue-500"
              disabled={loading}
            />
          </div>
          {errors.api && (
            <p className="text-red-500 text-sm mb-4">{errors.api}</p>
          )}
          <button
            type="submit"
            className={`w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600 ${
              loading ? "opacity-50 cursor-not-allowed" : ""
            }`}
            disabled={loading}
          >
            Sign In
          </button>
        </form>
      </div>
    </div>
  );
};

export default Signin;
