import { Outlet } from 'react-router-dom';
import { LanguageSwitcher } from './LanguageSwitcher';
import { ThemeToggle } from '../components/ThemeToggle';

export default function Layout() {
  return (
    <div className="min-h-screen bg-white dark:bg-gray-900/50 transition-colors">
      <header className="p-4 flex justify-between items-center border-b dark:border-gray-700">
        <LanguageSwitcher />
        <ThemeToggle />
      </header>

      <main className="p-4">
        <Outlet />
      </main>
    </div>
  );
}
