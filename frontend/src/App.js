// import logo from './logo.svg';
import './App.css';
import Signup from './containers/Singup.js';
import SignIn from './containers/Signin.js';
import Home from './containers/Home.js';
import Confirmticket from './containers/selectconfirmticket.js';
import { Route, Routes } from 'react-router-dom';
import Bookings from "./containers/Booking.js";
import Travels from "./containers/Travels.js"

function App() {
  return (
    // <div className="App">
    //   <header className="App-header">
    //     <img src={logo} className="App-logo" alt="logo" />
    //     <p>
    //       Edit fsf sfs <code>src/App.js</code> and save to reload.
    //     </p>
    //     <a
    //       className="App-link"
    //       href="https://reactjs.org"
    //       target="_blank"
    //       rel="noopener noreferrer"
    //     >
    //       Learn React
    //     </a>
    //   </header>
    // </div>
    <Routes>
      <Route path="/" element={
        <div className="flex items-center justify-center min-h-screen bg-gray-100">
          <div className="bg-white p-6 rounded shadow-md w-80">
            <h2 className="text-lg font-semibold mb-4">Welcome to the App</h2>
            <p className="mb-4">Please sign up or sign in to continue.</p>
            <div className="flex justify-between">
              <a href="/signup" className="text-blue-500 hover:underline">Sign Up</a>
              <a href="/signin" className="text-blue-500 hover:underline">Sign In</a>
            </div>  
          </div>
        </div>
      } />
      <Route path="/signup" element={<Signup />} />
      <Route path="/signin" element={<SignIn />} />
      <Route path="/home" element={<Home />} />
      <Route path="/details/:transportid" element={<Confirmticket />} />
      <Route path="/bookings" element={<Bookings />} />
      <Route path="/travels" element={<Travels />} />

      {/* get path for /detail/id api, send id as a argument which will be used*/}
      {/* Add more routes as needed */}
    </Routes>
    

  );
}

export default App;
