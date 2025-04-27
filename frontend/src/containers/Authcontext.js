import React, { createContext, useContext, useState, useEffect } from "react";

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [token, setToken] = useState(null);
  const [userId, setUserId] = useState(null);
  const [role, setRole] = useState(null);

  // Load from localStorage on mount
  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    const storedUserId = localStorage.getItem("userid");
    const storedRole = localStorage.getItem("role");

    if (storedToken && storedUserId && storedRole) {
      setToken(storedToken);
      setUserId(storedUserId);
      setRole(storedRole);
    }
  }, []);

  const login = (token, userId, role) => {
    setToken(token);
    setUserId(userId);
    setRole(role);

    localStorage.setItem("token", token);
    localStorage.setItem("userid", userId);
    localStorage.setItem("role", role);
  };

  const logout = () => {
    setToken(null);
    setUserId(null);
    setRole(null);

    localStorage.removeItem("token");
    localStorage.removeItem("userid");
    localStorage.removeItem("role");
  };

  return (
    <AuthContext.Provider value={{ token, userId, role, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
