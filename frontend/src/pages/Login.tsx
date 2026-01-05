import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import Input from '../components/UI/Input';
import PasswordInput from '../components/UI/PasswordInput/PasswordInput';
import { motion } from 'framer-motion';
import { Loader2 } from 'lucide-react';
import Button from '../components/UI/Button/Button';
import { useDispatch, useSelector } from 'react-redux';
import { showNotification } from '../components/Notifications/notificationSlice';
import clsx from 'clsx';
import { useNavigate } from 'react-router-dom';
import { setAuth } from '../features/auth/authSlice';
import DebugAuth from '../components/Debug/DebugAuth';

const API_URL = import.meta.env.VITE_API_URL;

const Login = () => {
  const { t } = useTranslation('auth');
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [shakePassword, setShakePassword] = useState(false);
  const [shakeUsername, setShakeUsername] = useState(false);
  const [usernameError, setUsernameError] = useState(false);
  const [passwordError, setPasswordError] = useState(false);

  const shakeFunc = (filed: 'username' | 'password') => {
    if (filed.includes('password')) {
      setShakePassword(true);
      setPasswordError(true);
      setTimeout(() => {
        setShakePassword(false);
      }, 400);
      return;
    }
    if (filed.includes('username')) {
      setShakeUsername(true);
      setUsernameError(true);
      setTimeout(() => {
        setShakeUsername(false);
      }, 400);
      return;
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    // console.log('username', username);
    // console.log('password', password);

    if (username === '') {
      shakeFunc('username');
      dispatch(
        showNotification({
          message: t('username cant be empty'),
          type: 'error',
        }),
      );
      return;
    }

    if (password === '') {
      shakeFunc('password');
      dispatch(
        showNotification({
          message: t('password cant be empty'),
          type: 'error',
        }),
      );
      return;
    }

    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ username, password }),
      });

      const data = await res.json();

      if (!res.ok) {
        if (data.error.includes('username')) shakeFunc('username');
        if (data.error.includes('password')) shakeFunc('password');
        dispatch(
          showNotification({
            message: t(data.error),
            type: 'error',
          }),
        );
        return;
      }
      console.log('data login', data);
      dispatch(setAuth(data));

      // (Временно) refresh_token можно сохранить в localStorage, ⚠️ Потом мы это уберём в cookie.
      // ❌ УДАЛЯЕМ
      // localStorage.setItem('refresh_token_telecom_map', data.refresh_token);
      // localStorage.setItem('telecom_map_token', data.token);
      // teper FRONTEND — вообще без refresh token, refresh teper tolko w HTPSCookie

      dispatch(
        showNotification({
          message: t('login successful'),
          type: 'success',
        }),
      );
      navigate('/');
      console.log('data', data);
    } catch (err) {
      console.log('err login == ', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    document.title = t('login');
  }, [t]);
  return (
    <div className="flex items-center justify-center">
      <DebugAuth />
      <form
        autoComplete="off"
        onSubmit={handleSubmit}
        className="w-full max-w-sm p-6 space-y-4 bg-white border border-gray-200 rounded-lg shadow-md dark:bg-zinc-800 dark:border-zinc-700"
      >
        <h2 className="text-2xl font-semibold">{t('login')}</h2>

        <motion.div
          animate={shakeUsername ? { x: [0, -8, 8, -8, 8, 0] } : { x: 0 }}
          transition={{ duration: 0.5 }}
        >
          <Input
            label={t('username')}
            id="username"
            value={username}
            onChange={(e) => {
              setUsername(e.target.value);
              setUsernameError(false);
            }}
            placeholder={t('enter username')}
            autoComplete="username"
            className={clsx(usernameError && 'border-red-500 dark:border-red-500')}
          />
        </motion.div>

        <motion.div
          animate={shakePassword ? { x: [0, -8, 8, -8, 8, 0] } : { x: 0 }}
          transition={{ duration: 0.5 }}
        >
          <PasswordInput
            label={t('password')}
            id="password"
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
              setPasswordError(false);
            }}
            className={clsx(passwordError && 'border-red-500 dark:border-red-500')}
          />
        </motion.div>

        <Button
          type="submit"
          variant={loading ? 'secondary' : 'primary'}
          icon={loading ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : undefined}
          className="w-full"
        >
          {t('login')}
        </Button>
      </form>
    </div>
  );
};

export default Login;
