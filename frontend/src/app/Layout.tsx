import { Outlet, NavLink } from 'react-router-dom';
import Header from '../components/Header/Header';
import Notifications from '../components/Notifications/Notifications';

export default function Layout() {
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
