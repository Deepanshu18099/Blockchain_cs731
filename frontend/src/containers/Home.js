// tailwind css will be used for styling, basic structure
import React from 'react';


// simple dashboard page showing the application details
import { Link } from "react-router-dom";
import { useAuth } from "./Authcontext";  // import

function Home() {
  const { userId, role } = useAuth();  // get values from context

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-6 rounded shadow-md w-80">
        <h2 className="text-lg font-semibold mb-4">Welcome to the App</h2>
        {userId ? (
          <div>
            <p className="mb-4">User ID: {userId}</p>
            <p className="mb-4">Role: {role}</p>
          </div>
        ) : (
          <div>
            <p className="mb-4">Please sign up or sign in to continue.</p>
            <Link to="/signup" className="text-blue-500 hover:underline">
              Sign Up
            </Link>
            <Link to="/signin" className="text-blue-500 hover:underline ml-2">
              Sign In
            </Link>
          </div>
        )}
      </div>
    </div>
  );
}

export default Home;
