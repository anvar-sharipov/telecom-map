import { useEffect, useState } from 'react';
import { useDispatch } from 'react-redux';
import { setAuth, clearAuth } from '../features/auth/authSlice';

const API_URL = import.meta.env.VITE_API_URL;

export default function useCheckAuth() {
  const dispatch = useDispatch();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const res = await fetch(`${API_URL}/auth/me`, {
          credentials: 'include',
        });
        if (res.ok) {
          const data = await res.json();
          dispatch(setAuth({ isAuth: true, userId: data.user_id }));
        } else {
          dispatch(clearAuth());
        }
      } catch {
        dispatch(clearAuth());
      } finally {
        setLoading(false);
      }
    };
    checkAuth();
  }, [dispatch]);

  return { loading };
}
