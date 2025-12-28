import { createBrowserRouter } from 'react-router-dom';
import Home from '../pages/Home';
import Register from '../pages/Register';
import Layout from './Layout';

export const router = createBrowserRouter([
  {
    element: <Layout />, // 👈 layout-обертка
    children: [
      {
        path: '/',
        element: <Home />,
      },
      {
        path: '/register',
        element: <Register />,
      },
    ],
  },
]);
