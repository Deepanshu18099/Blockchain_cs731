import { useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useEffect } from "react";


const SignUp = () => {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [errors, setErrors] = useState({});
  const navigate = useNavigate();
  const [pvtKey, setPvtKey] = useState("");
  const [publicKey, setPublicKey] = useState("");




  const timeout = 10; // seconds
  // function to handle timer
  const [timeoutleft, setTimeoutleft] = useState();


  // to handle timer useeffect
  useEffect(() => {
    const interval = setInterval(() => {
      setTimeoutleft((prev) => {
        if (prev <= 0) {
          clearInterval(interval);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
    return () => clearInterval(interval);
  }, [timeout]);



  // pvt key and public key will be returned from the api, which user has to save
  // endpoint will be REACT_APP_API_URL + 'ledger/createuser', use env
  const handleSubmit = async (e) => {
    e.preventDefault();
    // print something to test
    console.log("User Name:", name);
    const apiurl = process.env.REACT_APP_API_URL;
    try {

      // send a post request to the api
      const response = await axios.post(`${apiurl}ledger/createuser`, {
        name,
        email,
        phone,
      });

      console.log("Response from API:", response.data);

      
      if (response.status !== 200) {
        setErrors({ api: "Error creating user" });
        return;
      }
      console.log("Private Key:", pvtKey);
      console.log("User created successfully:", response.data);

      setPvtKey(response.data.output.privateKey);
      setPublicKey(response.data.message);
      console.log("Private Key:", pvtKey);
      setTimeoutleft(timeout);

      // wait seconds to save the keys, and navigate to signin
      setTimeout(() => {
        navigate("/signin");
      }, timeout * 1000);

    } catch (error) {
      console.error("Error signing up:", error);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-6 rounded shadow-md w-80">
        <h2 className="text-lg font-semibold mb-4">Sign Up</h2>
        {
          !pvtKey && 
          <>
          <p className="mb-4">Please fill in the form to sign up.</p>

          <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label
              htmlFor="name"
              className="block text-sm font-medium text-gray-700"
              >
              Name
            </label>
            <input
              type="text"
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring focus:ring-blue-500"
              />
          </div>
          <div className="mb-4">
            <label
              htmlFor="email"
              className="block text-sm font-medium text-gray-700"
              >
              Email
            </label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring focus:ring-blue-500"
            />
          </div>
          <div className="mb-4">
            <label
              htmlFor="phone"
              className="block text-sm font-medium text-gray-700"
              >
              Phone
            </label>
            <input
              type="text"
              id="phone"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              required
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring focus:ring-blue-500"
              />
          </div>
          <button
            type="submit"
            className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
            >
            Sign Up
          </button>
        </form>
        </>
        }
        {pvtKey && (
          <>
          <h3>
            Keys generated successfully. Please save the keys safely.
          </h3>
          <p className="mt-4 text-red-600">
            Please save the keys safely. You have {timeoutleft} seconds left to
            save the keys.
            <br />
            Private Key: {pvtKey}
            <br />
            Public Key: {publicKey}
            <br />
            You will be redirected to sign in page in {timeout} seconds.
          </p>
          </>
        )}
        <p className="mt-4 text-blue-500">
          Already have an account?{" "}
          <a href="/signin" className="underline">
            Sign In
          </a>
        </p>
      </div>
    </div>
  );
}

export default SignUp;