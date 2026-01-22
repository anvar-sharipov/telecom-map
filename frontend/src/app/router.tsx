import { createBrowserRouter } from 'react-router-dom';
import Home from '../pages/Home';
import Register from '../pages/Register';
import Layout from './Layout';
import Login from '../pages/Login';
import ProtectedRoute from '../routes/ProtectedRoute';

// admin panel
import AdminMain from '../pages/admin/AdminMain';

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
          {
            path: '/admin-panel',
            element: <AdminMain />,
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
