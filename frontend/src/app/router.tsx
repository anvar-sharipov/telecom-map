import { createBrowserRouter } from 'react-router-dom';
import Home from '../pages/Home';
import Register from '../pages/Register';
import Layout from './Layout';
import Login from '../pages/Login';
import ProtectedRoute from '../routes/ProtectedRoute';

export const router = createBrowserRouter([
  {
    element: <Layout />, // 👈 общий layout
    children: [
      {
        path: '/login',
        element: <Login />,
      },
      {
        path: '/register',
        element: <Register />,
      },

      // 🔒 ЗАЩИЩЁННЫЕ РОУТЫ
      {
        element: <ProtectedRoute />,
        children: [
          {
            path: '/',
            element: <Home />,
          },
          // сюда потом:
          // {
          //   path: '/profile',
          //   element: <Profile />,
          // },
        ],
      },
    ],
  },
]);
