import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  type ReactNode,
} from "react";

const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:8080";

interface OAuthUser {
  username: string;
  avatar_url: string;
  provider: string;
}

interface AuthContextValue {
  // Legacy basic auth (kept for basic mode)
  credentials: { username: string; password: string } | null;
  // OAuth user (set when mode is oauth and user is authenticated)
  user: OAuthUser | null;
  // Auth mode: "basic" or "oauth"
  authMode: string;
  // Whether the auth mode has been loaded
  authLoading: boolean;
  login: (username: string, password: string) => void;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [credentials, setCredentials] = useState<{
    username: string;
    password: string;
  } | null>(null);
  const [user, setUser] = useState<OAuthUser | null>(null);
  const [authMode, setAuthMode] = useState("basic");
  const [authLoading, setAuthLoading] = useState(true);

  // Detect auth mode and check OAuth session on mount
  useEffect(() => {
    fetch(`${API_BASE}/api/v1/auth/mode`)
      .then((res) => (res.ok ? res.json() : { mode: "basic" }))
      .then((data) => {
        setAuthMode(data.mode);
        if (data.mode === "oauth") {
          // Check if we have an active session
          return fetch(`${API_BASE}/api/v1/auth/me`, {
            credentials: "include",
          })
            .then((res) => (res.ok ? res.json() : null))
            .then((userData) => {
              if (userData) setUser(userData);
            });
        }
      })
      .catch(() => {})
      .finally(() => setAuthLoading(false));
  }, []);

  const login = useCallback((username: string, password: string) => {
    setCredentials({ username, password });
  }, []);

  const logout = useCallback(() => {
    setCredentials(null);
    setUser(null);
    if (authMode === "oauth") {
      fetch(`${API_BASE}/api/v1/auth/logout`, {
        method: "POST",
        credentials: "include",
      }).catch(() => {});
    }
  }, [authMode]);

  const isAuthenticated =
    authMode === "oauth" ? user !== null : credentials !== null;

  return (
    <AuthContext.Provider
      value={{
        credentials,
        user,
        authMode,
        authLoading,
        login,
        logout,
        isAuthenticated,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return ctx;
}
