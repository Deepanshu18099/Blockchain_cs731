// this will be sign in page for hyperledger page, signing using userid, pvt key and recieve token from api
// tailwind css will be used for styling, basic structure
import React, { useState } from 'react';
import axios from 'axios';


const Signin = () => {
  const [userId, setUserId] = useState('');
  const [pvtKey, setPvtKey] = useState('');
  const [token, setToken] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post('http://localhost:5000/signin', {
        userId,
        pvtKey,
      });
      setToken(response.data.token);
    } catch (error) {
      console.error('Error signing in:', error);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-6 rounded shadow-md w-80">
        <h2 className="text-lg font-semibold mb-4">Sign In</h2>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="userId" className="block text-sm font-medium text-gray-700">User ID</label>
            <input
              type="text"
              id="userId"
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              required
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring focus:ring-blue-500"
            />
          </div>
          <div className="mb-4">
            <label htmlFor="pvtKey" className="block text-sm font-medium text-gray-700">Private Key</label>
            <input
              type="password"
              id="pvtKey"
              value={pvtKey}
              onChange={(e) => setPvtKey(e.target.value)}
              required
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring focus:ring-blue-500"
            />
          </div>
          <button type="submit" className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600">Sign In</button>
        </form>
        {token && <p className="mt-4 text-green-600">Token: {token}</p>}
      </div>
    </div>
    );
}
export default Signin;
