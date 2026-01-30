import { Home, Shield, UserPlus, LogIn, User, LogOut } from 'lucide-react';
// import { useTranslation } from 'react-i18next';

export const mobileMenuSections = [
  {
    title: 'Основное',
    items: [
      {
        to: '/',
        text: 'Главная',
        icon: Home,
        auth: 'any',
      },
    ],
  },
  {
    title: 'Админ',
    items: [
      {
        to: '/admin-panel',
        text: 'Admin panel',
        icon: Shield,
        auth: 'auth',
      },
    ],
  },
  {
    title: 'Аккаунт',
    items: [
      {
        to: '/register',
        text: 'Регистрация',
        icon: UserPlus,
        auth: 'guest',
      },
      {
        to: '/login',
        text: 'Вход',
        icon: LogIn,
        auth: 'guest',
      },
      {
        to: '/profile',
        text: 'Профиль',
        icon: User,
        auth: 'auth',
      },
      {
        to: 'logout',
        text: 'Выход',
        icon: LogOut,
        auth: 'auth',
        action: 'logout',
      },
    ],
  },
];
