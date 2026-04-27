import { createContext, useContext, useState, useCallback } from 'react';
import type { ReactNode } from 'react';

interface AuthState {
  username: string;
  password: string;
}

interface AuthContextValue {
  credentials: AuthState | null;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => void;
  authHeader: () => string | null;
}

const AuthContext = createContext<AuthContextValue | null>(null);

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

interface AuthProviderProps {
  children: ReactNode;
  initialCredentials?: { username: string; password: string };
}

export function AuthProvider({ children, initialCredentials }: AuthProviderProps) {
  const [credentials, setCredentials] = useState<AuthState | null>(initialCredentials ?? null);

  const login = useCallback(async (username: string, password: string): Promise<boolean> => {
    const encoded = btoa(`${username}:${password}`);
    try {
      const res = await fetch(`${API_BASE}/api/v1/admin/posts`, {
        headers: { Authorization: `Basic ${encoded}` },
      });
      if (res.ok) {
        setCredentials({ username, password });
        return true;
      }
      return false;
    } catch {
      return false;
    }
  }, []);

  const logout = useCallback(() => {
    setCredentials(null);
  }, []);

  const authHeader = useCallback((): string | null => {
    if (!credentials) return null;
    return `Basic ${btoa(`${credentials.username}:${credentials.password}`)}`;
  }, [credentials]);

  return (
    <AuthContext.Provider value={{ credentials, login, logout, authHeader }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return ctx;
}
