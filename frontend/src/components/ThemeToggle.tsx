// import { useEffect } from 'react';
// import { useAppDispatch, useAppSelector } from '../app/hooks';
// import { toggleTheme } from '../features/theme/themeSlice';

// export function ThemeToggle() {
//   const theme = useAppSelector((state) => state.theme.theme);
//   const dispatch = useAppDispatch();

//   useEffect(() => {
//     document.documentElement.classList.toggle('dark', theme === 'dark');
//     localStorage.setItem('theme', theme);
//   }, [theme]);

//   return (
//     <button
//       onClick={() => dispatch(toggleTheme())}
//       className="px-3 py-1 rounded border text-sm
//                  bg-gray-200 dark:bg-gray-800
//                  text-black dark:text-white"
//     >
//       {theme === 'light' ? '🌙 Dark' : '☀️ Light'}
//     </button>
//   );
// }

import { useEffect, useState } from 'react';

export function ThemeToggle() {
  const [theme, setTheme] = useState<'light' | 'dark'>(
    (localStorage.getItem('theme') as 'light' | 'dark') || 'light',
  );

  useEffect(() => {
    document.documentElement.classList.toggle('dark', theme === 'dark');
    localStorage.setItem('theme', theme);
  }, [theme]);

  const toggleTheme = () => {
    setTheme((prev) => (prev === 'light' ? 'dark' : 'light'));
  };

  return (
    <button
      onClick={toggleTheme}
      className="px-3 py-1 rounded border text-sm
                 bg-gray-200 dark:bg-gray-800
                 text-black dark:text-white"
    >
      {theme === 'light' ? '🌙 Dark' : '☀️ Light'}
    </button>
  );
}
