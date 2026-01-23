import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { motion, AnimatePresence } from 'framer-motion';
import { Menu, X, ChevronLeft } from 'lucide-react';

import UsersTable from './users/UserTable';

const desktopSidebarVariants = {
  open: {
    width: 256,
    transition: { type: 'spring', stiffness: 220, damping: 26 },
  },
  closed: {
    width: 56,
    transition: { type: 'spring', stiffness: 220, damping: 26 },
  },
};

const mobileSidebarVariants = {
  hidden: { x: '-100%' },
  visible: {
    x: 0,
    transition: { type: 'spring', stiffness: 300, damping: 30 },
  },
};

const AdminMain = () => {
  const { t } = useTranslation('common');

  const [activeTable, setActiveTable] = useState('users');
  const [collapsed, setCollapsed] = useState(false); // desktop
  const [mobileOpen, setMobileOpen] = useState(false); // mobile

  const tables = [
    { key: 'users', label: t('Users') },
    { key: 'groups', label: t('Groups') },
    { key: 'user_groups', label: t('User Groups') },
  ];

  return (
    <div className="flex w-full min-h-[calc(100vh-64px)] bg-gray-50 dark:bg-zinc-950">
      {/* MOBILE HAMBURGER */}
      <button
        onClick={() => setMobileOpen(true)}
        className="fixed z-40 p-2 text-white bg-indigo-600 rounded-md shadow-lg md:hidden top-4 left-4"
      >
        <Menu className="w-5 h-5" />
      </button>

      {/* DESKTOP SIDEBAR */}
      <motion.aside
        variants={desktopSidebarVariants}
        animate={collapsed ? 'closed' : 'open'}
        className="relative hidden overflow-hidden border-r md:block bg-indigo-50 dark:bg-zinc-900"
      >
        {/* HEADER */}
        <div className="flex items-center px-4 bg-indigo-400 border-b h-14 dark:bg-zinc-900">
          {!collapsed && (
            <span className="text-sm font-semibold truncate">Панель администратора</span>
          )}
        </div>

        {/* MENU */}
        <ul className="p-2 space-y-1">
          {!collapsed &&
            tables.map((table) => (
              <li
                key={table.key}
                onClick={() => setActiveTable(table.key)}
                className={`
                  px-3 py-2 rounded-md cursor-pointer text-sm
                  transition-colors
                  ${
                    activeTable === table.key
                      ? 'bg-indigo-200 text-indigo-900 dark:text-indigo-200 font-medium dark:bg-zinc-700'
                      : 'hover:bg-indigo-100 dark:hover:bg-zinc-800'
                  }
                `}
              >
                {table.label}
              </li>
            ))}
        </ul>

        {/* DESKTOP TOGGLE */}
        <button
          onClick={() => setCollapsed(!collapsed)}
          className="absolute flex items-center justify-center w-6 h-6 text-white -translate-y-1/2 bg-indigo-600 rounded-full shadow-md dark:bg-indigo-700 -right-0 top-1/2 hover:bg-indigo-700 dark:hover:bg-indigo-600"
        >
          <motion.div animate={{ rotate: collapsed ? 180 : 0 }} transition={{ duration: 0.25 }}>
            <ChevronLeft className="w-4 h-4" />
          </motion.div>
        </button>
      </motion.aside>

      {/* MOBILE OVERLAY SIDEBAR */}
      <AnimatePresence>
        {mobileOpen && (
          <>
            {/* BACKDROP */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 0.5 }}
              exit={{ opacity: 0 }}
              onClick={() => setMobileOpen(false)}
              className="fixed inset-0 z-40 bg-black"
            />

            {/* SIDEBAR */}
            <motion.aside
              variants={mobileSidebarVariants}
              initial="hidden"
              animate="visible"
              exit="hidden"
              className="fixed inset-y-0 left-0 z-50 w-64 shadow-xl bg-indigo-50 dark:bg-zinc-900"
            >
              {/* HEADER */}
              <div className="flex items-center justify-between px-4 border-b h-14">
                <span className="text-sm font-semibold">Панель администратора</span>
                <button onClick={() => setMobileOpen(false)}>
                  <X />
                </button>
              </div>

              {/* MENU */}
              <ul className="p-2 space-y-1">
                {tables.map((table) => (
                  <li
                    key={table.key}
                    onClick={() => {
                      setActiveTable(table.key);
                      setMobileOpen(false);
                    }}
                    className={`
                      px-3 py-2 rounded-md cursor-pointer text-sm
                      transition-colors
                      ${
                        activeTable === table.key
                          ? 'bg-indigo-200 text-indigo-900 font-medium dark:bg-zinc-700 dark:text-indigo-200'
                          : 'hover:bg-indigo-100 dark:hover:bg-zinc-800'
                      }
                    `}
                  >
                    {table.label}
                  </li>
                ))}
              </ul>
            </motion.aside>
          </>
        )}
      </AnimatePresence>

      {/* CONTENT */}
      <main className="flex-1 p-6">{activeTable === 'users' && <UsersTable />}</main>
    </div>
  );
};

export default AdminMain;
