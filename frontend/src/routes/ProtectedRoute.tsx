import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAppSelector } from '../app/hooks';

export default function ProtectedRoute() {
  const { isAuth, loading } = useAppSelector((state) => state.auth);
  const location = useLocation();

  if (loading) return null;

  if (!isAuth) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  return <Outlet />;
}
