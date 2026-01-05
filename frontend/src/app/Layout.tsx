import { Outlet, NavLink } from 'react-router-dom';
import Header from '../components/Header/Header';
import Notifications from '../components/Notifications/Notifications';
import useCheckAuth from '../hooks/useCheckAuth';

export default function Layout() {
  useCheckAuth();
  return (
    <div className="min-h-screen text-gray-900 transition-colors bg-gray-50 dark:bg-zinc-900 dark:text-gray-100">
      {/* 🔔 Уведомления */}
      <Notifications />

      <Header />

      <main className="flex min-h-[calc(100vh-64px)] items-center justify-center">
        <Outlet />
      </main>
    </div>
  );
}
