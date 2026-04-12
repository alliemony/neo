import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Layout } from "../components/layout/Layout";
import { useAuth } from "../hooks/useAuth";
import { adminGetPosts } from "../services/api";

const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:8080";

export function AdminLogin() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { login, authMode, authLoading, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  // If already authenticated, redirect to admin
  if (isAuthenticated) {
    navigate("/admin");
    return null;
  }

  if (authLoading) {
    return (
      <Layout>
        <div className="max-w-sm mx-auto py-12 text-center">
          <p className="text-text-secondary">Loading…</p>
        </div>
      </Layout>
    );
  }

  // OAuth mode: show "Sign in with GitHub" button
  if (authMode === "oauth") {
    return (
      <Layout>
        <div className="max-w-sm mx-auto py-12 text-center">
          <h1 className="font-heading text-2xl font-bold mb-8">Admin Login</h1>
          <a
            href={`${API_BASE}/api/v1/auth/login`}
            className="inline-flex items-center gap-2 px-6 py-3 border-2 border-border bg-[#24292f] text-white font-heading text-sm hover:opacity-90 no-underline"
          >
            <svg width="20" height="20" viewBox="0 0 16 16" fill="currentColor">
              <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z" />
            </svg>
            Sign in with GitHub
          </a>
        </div>
      </Layout>
    );
  }

  // Basic auth mode: show username/password form
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      await adminGetPosts({ username, password });
      login(username, password);
      navigate("/admin");
    } catch {
      setError("Invalid credentials. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <div className="max-w-sm mx-auto py-12">
        <h1 className="font-heading text-2xl font-bold mb-6">Admin Login</h1>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-heading mb-1">Username</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full border-2 border-border bg-surface px-3 py-2 text-sm font-body text-text-primary focus:outline-none focus:border-accent"
              autoComplete="username"
            />
          </div>
          <div>
            <label className="block text-sm font-heading mb-1">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full border-2 border-border bg-surface px-3 py-2 text-sm font-body text-text-primary focus:outline-none focus:border-accent"
              autoComplete="current-password"
            />
          </div>
          {error && <p className="text-accent text-sm">{error}</p>}
          <button
            type="submit"
            disabled={loading}
            className="w-full border-2 border-border bg-accent text-white px-4 py-2 text-sm font-heading hover:opacity-90 disabled:opacity-50"
          >
            {loading ? "Logging in…" : "Login"}
          </button>
        </form>
      </div>
    </Layout>
  );
}
