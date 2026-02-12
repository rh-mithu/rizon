import React, { createContext, useContext, useState } from 'react';
import { api } from '../../services/api';

type AuthContextType = {
  isAuthenticated: boolean;
  login: () => void;
  logout: () => void;
  requestLoginLink: (email: string) => Promise<void>;
};

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const login = () => setIsAuthenticated(true);
  const logout = () => setIsAuthenticated(false);

  const requestLoginLink = async (email: string) => {
    await api.requestLoginLink(email);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout, requestLoginLink }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
