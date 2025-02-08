import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import { message } from "antd";
import PingTable from "./components/PingTable";
import AuthPage from "./components/AuthPage";
import "antd/dist/reset.css";
import axios from "axios";

const App = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // Проверка аутентификации при загрузке
  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const response = await axios.get("http://localhost:5000/", {
        withCredentials: true
      });
      setIsAuthenticated(response.status === 200);
    } catch (error) {
      setIsAuthenticated(false);
    }
  };

  return (
    <Router>
      <Routes>
        <Route
          path="/"
          element={
            isAuthenticated ? (
              <Navigate to="/dashboard" replace />
            ) : (
              <Navigate to="/login" replace />
            )
          }
        />
        <Route
          path="/login"
          element={
            <AuthPage
              type="login"
              onAuthSuccess={() => setIsAuthenticated(true)}
            />
          }
        />
        <Route
          path="/signup"
          element={
            <AuthPage
              type="signup"
              onAuthSuccess={() => setIsAuthenticated(true)}
            />
          }
        />
        <Route
          path="/dashboard"
          element={
            isAuthenticated ? (
              <div>
                <h1 style={{ textAlign: "center", margin: "20px 0" }}>
                  Статистика пингов
                </h1>
                <PingTable onUnauthorized={() => setIsAuthenticated(false)} />
              </div>
            ) : (
              <Navigate to="/login" replace />
            )
          }
        />
      </Routes>
    </Router>
  );
};

export default App;