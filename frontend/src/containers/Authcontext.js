import React, { createContext, useContext, useState, useEffect } from "react";

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [token, setToken] = useState(null);
  const [userId, setUserId] = useState(null);
  const [role, setRole] = useState(null);
  const [balance, setBalance] = useState(null);

  // Load from localStorage on mount
  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    const storedUserId = localStorage.getItem("userid");
    const storedRole = localStorage.getItem("role");
    const storedBalance = localStorage.getItem("balance");


    if (storedToken && storedUserId && storedRole) {
      setToken(storedToken);
      setUserId(storedUserId);
      setRole(storedRole);
      setBalance(storedBalance);
    }
  }, []);

  const login = (token, userId, role, balance) => {
    setToken(token);
    setUserId(userId);
    setRole(role);
    setBalance(balance)

    localStorage.setItem("token", token);
    localStorage.setItem("userid", userId);
    localStorage.setItem("role", role);
    localStorage.setItem("balance", balance)
  };

  const logout = () => {
    setToken(null);
    setUserId(null);
    setRole(null);
    setBalance(null);

    localStorage.removeItem("token");
    localStorage.removeItem("userid");
    localStorage.removeItem("role");
    localStorage.removeItem("balance")
  };

  return (
    <AuthContext.Provider value={{ token, userId, role, balance, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
