import { useEffect } from 'react';
// import { useDispatch } from 'react-redux';
import { useAppDispatch } from '../app/hooks';
import { setAuth, clearAuth, setLoading } from '../features/auth/authSlice';

const API_URL = import.meta.env.VITE_API_URL;

export default function useCheckAuth() {
  const dispatch = useAppDispatch();

  useEffect(() => {
    const checkAuth = async () => {
      dispatch(setLoading(true));

      try {
        // 1️⃣ Проверка refresh cookie
        const res = await fetch(`${API_URL}/auth/me`, {
          credentials: 'include',
        });

        if (res.ok) {
          // const user = await res.json();
          const data = await res.json();
          console.log('data', data);

          dispatch(
            setAuth({
              access_token: null,
              expires_at: null,
              user: data.user,
            }),
          );
          return;
        } else {
          console.log('NONNONONO');
        }

        // 2️⃣ Пробуем refresh
        const refreshRes = await fetch(`${API_URL}/auth/refresh`, {
          method: 'POST',
          credentials: 'include',
        });

        if (!refreshRes.ok) {
          console.log('NONNONONO2');
          throw new Error('refresh failed');
        }

        const data = await refreshRes.json();
        console.log('data === ', data);

        dispatch(
          setAuth({
            access_token: data.access_token,
            expires_at: data.expires_at,
            user: data.user,
          }),
        );
      } catch {
        dispatch(clearAuth());
      }
    };

    checkAuth();
  }, [dispatch]);
}
