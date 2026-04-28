import { Navigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { credentials } = useAuth();

  if (!credentials) {
    return <Navigate to="/admin/login" replace />;
  }

  return <>{children}</>;
}
