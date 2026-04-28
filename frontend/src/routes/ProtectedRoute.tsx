import type { ReactNode } from "react";
import { Navigate } from "react-router-dom";
import { Layout } from "../components/layout/Layout";
import { useAuth } from "../hooks/useAuth";

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const { authLoading, isAuthenticated } = useAuth();

  if (authLoading) {
    return (
      <Layout>
        <div className="max-w-sm mx-auto py-12 text-center">
          <p className="text-text-secondary">Loading…</p>
        </div>
      </Layout>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/admin/login" replace />;
  }

  return <>{children}</>;
}
