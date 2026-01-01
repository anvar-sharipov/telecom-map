import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch } from 'react-redux';
import { showNotification } from '../components/Notifications/notificationSlice';
import Input from '../components/UI/Input';
import PasswordInput from '../components/UI/PasswordInput/PasswordInput';
import ConfirmPassword from '../components/UI/PasswordInput/ConfirmPassword';
import { Loader2 } from 'lucide-react';
import Button from '../components/UI/Button/Button';
import { motion } from 'framer-motion';
// import Snowfall from 'react-snowfall';

const API_URL = import.meta.env.VITE_API_URL;

export default function Register() {
  const [username, setUsername] = useState('');
  const [fullname, setFullname] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [shake, setShake] = useState(false);
  const [loading, setLoading] = useState(false);
  const { t } = useTranslation('auth');
  const dispatch = useDispatch();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (password !== confirmPassword) {
      setShake(true);
      setTimeout(() => setShake(false), 400);
      dispatch(
        showNotification({
          message: t('passwords do not match'),
          type: 'error',
        }),
      );
      return;
    }
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, fullname, password, confirm_password: confirmPassword }),
      });

      const data = await res.json();

      if (!res.ok) {
        if (!loading) {
          dispatch(showNotification({ message: t(data.error), type: 'error' }));
        }
        return;
      }

      dispatch(showNotification({ message: t(data.message), type: 'success' }));
      localStorage.setItem('telecom_map_token', data.token);
    } catch (err) {
      console.error(err);
      dispatch(showNotification({ message: t('something went wrong'), type: 'error' }));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      {/* <Snowfall /> */}
      <form
        autoComplete="off"
        onSubmit={handleSubmit}
        className="w-full max-w-sm p-6 space-y-4 bg-white border border-gray-200 rounded-lg shadow-md dark:bg-zinc-800 dark:border-zinc-700"
      >
        <h2 className="text-2xl font-semibold">{t('register')}</h2>

        <Input
          label={t('username')}
          id="username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder={t('enter username')}
          autoComplete="username"
        />

        <Input
          label={t('fullname')}
          id="fullname"
          autoComplete="username"
          placeholder={t('enter full name')}
          value={fullname}
          onChange={(e) => setFullname(e.target.value)}
        />
        <motion.div
          animate={shake ? { x: [0, -8, 8, -8, 8, 0] } : { x: 0 }}
          transition={{ duration: 0.5 }}
        >
          <PasswordInput
            label={t('password')}
            id="password"
            autoComplete="new-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </motion.div>

        <motion.div
          animate={shake ? { x: [0, -8, 8, -8, 8, 0] } : { x: 0 }}
          transition={{ duration: 0.5 }}
        >
          <ConfirmPassword
            password={password}
            value={confirmPassword}
            onChange={setConfirmPassword}
          />
        </motion.div>

        <Button
          type="submit"
          variant={loading ? 'secondary' : 'primary'}
          icon={loading ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : undefined}
          className="w-full"
        >
          {t('register')}
        </Button>
      </form>
    </div>
  );
}
