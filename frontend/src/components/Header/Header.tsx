import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Menu, X, Home, User, LogIn } from 'lucide-react';
import { LanguageSwitcher } from './LanguageSwitcher';
import { ThemeToggle } from './ThemeToggle';
import NavItem from './NavItem';
import MobileLink from './MobileLink';
import { useSelector, useDispatch } from 'react-redux';
import type { RootState } from '../../app/store';
import { clearAuth } from '../../features/auth/authSlice';
import Button from '../UI/Button/Button';
import { showNotification } from '../Notifications/notificationSlice';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import useCheckAuth from '../../hooks/useCheckAuth';

const API_URL = import.meta.env.VITE_API_URL;

const Header = () => {
  const [open, setOpen] = useState(false);
  const dispatch = useDispatch();
  const isAuth = useSelector((state: RootState) => state.auth.isAuth);
  const { t } = useTranslation('auth');
  const navigate = useNavigate();
  const { loading } = useCheckAuth();

  const handleLogout = async () => {
    try {
      const res = await fetch(`${API_URL}/auth/logout`, {
        method: 'POST',
        credentials: 'include',
      });

      if (res.status !== 200) {
        const errorData = await res.json();
        dispatch(
          showNotification({
            message: t(errorData.error || 'logout error'),
            type: 'error',
          }),
        );
        return;
      }

      dispatch(clearAuth());
      dispatch(
        showNotification({
          message: t('logout successfull'),
          type: 'success',
        }),
      );
      navigate('/login');
    } catch (err) {
      console.log('err logout == ', err);
      dispatch(
        showNotification({
          message: t('logout error'),
          type: 'error',
        }),
      );
    }
  };

  // if (loading)
  //   return (
  //     <header className="h-12 bg-white border-b border-gray-200 dark:border-zinc-700 dark:bg-zinc-800"></header>
  //   );

  return (
    <header className="bg-white border-b border-gray-200 dark:border-zinc-700 dark:bg-zinc-800">
      <div className="flex items-center justify-between px-4 py-3 mx-auto max-w-7xl">
        {/* Левая часть */}
        <nav className="items-center hidden gap-6 md:flex">
          <NavItem to="/" text="Главная" icon={<Home size={20} />} />
          {!loading &&
            // <Loader2 className="w-4 h-4 animate-spin" />
            (!isAuth ? (
              <>
                <NavItem to="/register" text="Регистрация" icon={<User size={20} />} />
                <NavItem to="/login" text="Логин" icon={<LogIn size={20} />} />
              </>
            ) : (
              <>
                <Button onClick={handleLogout} variant="primary" size="sm">
                  logout
                </Button>
              </>
            ))}
        </nav>

        {/* Правая часть */}
        <div className="flex items-center gap-3">
          <LanguageSwitcher />
          <ThemeToggle />

          {/* Бургер (только на мобилке) */}
          <button
            onClick={() => setOpen(!open)}
            className="p-2 rounded md:hidden hover:bg-gray-200 dark:hover:bg-gray-700"
          >
            {open ? <X /> : <Menu />}
          </button>
        </div>
      </div>

      {/* МОБИЛЬНОЕ МЕНЮ */}
      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.25 }}
            className="overflow-hidden border-t md:hidden dark:border-gray-700"
          >
            <div className="flex flex-col gap-3 p-4">
              <MobileLink to="/" text="Главная" setOpen={setOpen} />
              {!loading &&
                // <Loader2 className="w-4 h-4 animate-spin" />
                (!isAuth ? (
                  <>
                    <MobileLink to="/register" text="Регистрация" setOpen={setOpen} />
                    <MobileLink to="/login" text="Логин" setOpen={setOpen} />
                  </>
                ) : (
                  <Button onClick={handleLogout} variant="primary" size="sm">
                    logout
                  </Button>
                ))}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  );
};

export default Header;
