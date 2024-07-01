import React, { useState } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate, Link } from 'react-router-dom';
import BudgetPage from './components/BugdetPage';
import TransactionAnalysisPage from './components/TransactionAnalysisPage';
import LoginPage from './components/LoginPage';
import { AuthProvider, useAuth } from './components/AuthContext';
import './App.css';
import RegisterPage from './components/Register';

function Navigation() {
  const { isAuthenticated, setIsAuthenticated, logout } = useAuth();

  if (!isAuthenticated) return null;

  return (
    <nav>
      <ul>
        <li><Link to="/">Income & Budget</Link></li>
        <li><Link to="/transactions">Transaction Analysis</Link></li>
        <li><button onClick={() => logout()}>Logout</button></li>
      </ul>
    </nav>
  );
}

function ProtectedRoute({ children }) {
  const { isAuthenticated } = useAuth();
  console.log(isAuthenticated)
  return isAuthenticated ? children : <Navigate to="/login" replace />;
}

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="App">
          <Navigation />
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route path="/" element={
              <ProtectedRoute>
                <BudgetPage />
              </ProtectedRoute>
            } />
            <Route path="/transactions" element={
              <ProtectedRoute>
                <TransactionAnalysisPage />
              </ProtectedRoute>
            } />
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;