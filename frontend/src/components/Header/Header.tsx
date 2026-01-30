import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Menu, X, Home, User, LogIn, LogOut } from 'lucide-react';
import { LanguageSwitcher } from './LanguageSwitcher';
import { ThemeToggle } from './ThemeToggle';
import NavItem from './NavItem';
import MobileLink from './MobileLink';
import { useAppDispatch, useAppSelector } from '../../app/hooks';
import type { RootState } from '../../app/store';
import { clearAuth } from '../../features/auth/authSlice';
// import Button from '../UI/Button/Button';
import { showNotification } from '../../features/notifications/notificationSlice';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { mobileMenuSections } from './mobileMenuSections';

// admin panel

const API_URL = import.meta.env.VITE_API_URL;

const listVariants = {
  open: {
    transition: { staggerChildren: 0.05 },
  },
  closed: {},
};

const itemVariants = {
  open: { opacity: 1, y: 0 },
  closed: { opacity: 0, y: -8 },
};

type SectionProps = {
  title: string;
};

type MobileSettingRowProps = {
  label: string;
  children: React.ReactNode;
};

const MobileSettingRow = ({ label, children }: MobileSettingRowProps) => (
  <motion.div
    variants={itemVariants}
    className="flex items-center justify-between px-4 py-3 rounded-xl bg-gray-50 dark:bg-zinc-800"
  >
    <span className="text-sm font-medium">{label}</span>
    {children}
  </motion.div>
);

const MobileSectionTitle = ({ title }: SectionProps) => (
  <div className="px-4 pt-3 pb-1">
    <p className="text-xs font-semibold tracking-wide text-gray-500 uppercase">{title}</p>
  </div>
);

const Header = () => {
  const [open, setOpen] = useState(false);
  const dispatch = useAppDispatch();
  const { t } = useTranslation('auth');
  const { t: tCommon } = useTranslation('common');
  const navigate = useNavigate();

  const { isAuth, loading, user } = useAppSelector((state: RootState) => state.auth);
  console.log('user', user);

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
      // dispatch(
      //   showNotification({
      //     message: t('logout successfull'),
      //     type: 'success',
      //   }),
      // );
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
    <header className="font-bold bg-white border-b border-gray-200 dark:border-zinc-700 dark:bg-zinc-800 text-[12px] md:text-sm">
      <div className="flex items-center justify-between px-4 py-3 mx-auto max-w-7xl">
        {/* Левая часть */}
        <nav className="items-center hidden gap-6 md:flex">
          <NavItem to="/" text="Главная" icon={<Home size={20} />} />
          {isAuth && <NavItem to="/admin-panel" text={tCommon('Admin panel')} icon="" />}

          {!loading &&
            // <Loader2 className="w-4 h-4 animate-spin" />
            (!isAuth ? (
              <>
                <NavItem to="/register" text={t('register')} icon={<User size={20} />} />
                <NavItem to="/login" text={t('login')} icon={<LogIn size={20} />} />
              </>
            ) : (
              // <>
              //   <button
              //     onClick={handleLogout}
              //     className="flex items-center w-full gap-3 px-4 py-3 text-sm font-medium text-red-600 rounded-xl hover:bg-red-50 dark:hover:bg-zinc-800"
              //   >
              //     logout
              //   </button>
              // </>
              <motion.div variants={itemVariants} className="px-2">
                <button
                  onClick={() => {
                    handleLogout();
                  }}
                  className="flex items-center w-full gap-3 px-4 py-3 text-sm font-medium text-red-600 rounded-xl hover:bg-red-50 dark:hover:bg-zinc-800"
                >
                  <LogOut className="w-5 h-5" />
                  Выход
                </button>
              </motion.div>
            ))}
        </nav>

        {/* Правая часть */}
        <div className="flex items-center gap-3 ml-auto">
          <div className="hidden md:block">
            <LanguageSwitcher />
          </div>

          <div className="hidden md:block">
            <ThemeToggle />
          </div>

          <div>{user?.username}</div>

          {/* Бургер */}
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
            <motion.div
              variants={listVariants}
              initial="closed"
              animate="open"
              exit="closed"
              className="flex flex-col gap-1 p-2"
            >
              {/* NAV SECTIONS */}
              {mobileMenuSections.map((section) => (
                <div key={section.title}>
                  <MobileSectionTitle title={section.title} />

                  {section.items.map((item) => {
                    if ((item.auth === 'auth' && !isAuth) || (item.auth === 'guest' && isAuth)) {
                      return null;
                    }

                    if (item.action === 'logout') {
                      return (
                        <motion.div key={item.text} variants={itemVariants} className="px-2">
                          <button
                            onClick={() => {
                              handleLogout();
                              setOpen(false);
                            }}
                            className="flex items-center w-full gap-3 px-4 py-3 text-sm font-medium text-red-600 rounded-xl hover:bg-red-50 dark:hover:bg-zinc-800"
                          >
                            <item.icon className="w-5 h-5" />
                            {tCommon(item.text)}
                          </button>
                        </motion.div>
                      );
                    }

                    return (
                      <motion.div key={item.text} variants={itemVariants}>
                        <MobileLink
                          to={item.to}
                          text={tCommon(item.text)}
                          icon={item.icon}
                          setOpen={setOpen}
                        />
                      </motion.div>
                    );
                  })}

                  <div className="my-2 border-t dark:border-zinc-700" />
                </div>
              ))}

              {/* SETTINGS SECTION (ОДИН РАЗ) */}
              <div>
                <MobileSectionTitle title="Настройки" />

                <div className="flex flex-col gap-2 px-2">
                  <MobileSettingRow label="Язык">
                    <LanguageSwitcher />
                  </MobileSettingRow>

                  <MobileSettingRow label="Тема">
                    <ThemeToggle />
                  </MobileSettingRow>
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  );
};

export default Header;
