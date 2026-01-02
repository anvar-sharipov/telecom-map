import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Menu, X, Home, User, LogIn } from 'lucide-react';
import { LanguageSwitcher } from './LanguageSwitcher';
import { ThemeToggle } from './ThemeToggle';
import NavItem from './NavItem';
import MobileLink from './MobileLink';

const Header = () => {
  const [open, setOpen] = useState(false);

  return (
    <header className="bg-white border-b border-gray-200 dark:border-zinc-700 dark:bg-zinc-800">
      <div className="flex items-center justify-between px-4 py-3 mx-auto max-w-7xl">
        {/* Левая часть */}
        <nav className="items-center hidden gap-6 md:flex">
          <NavItem to="/" text="Главная" icon={<Home size={20} />} />
          <NavItem to="/register" text="Регистрация" icon={<User size={20} />} />
          <NavItem to="/login" text="Логин" icon={<LogIn size={20} />} />
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
              <MobileLink to="register" text="Регистрация" setOpen={setOpen} />
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  );
};

export default Header;
