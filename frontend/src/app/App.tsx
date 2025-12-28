import { RouterProvider } from 'react-router-dom';
import { router } from './router';

export default function App() {
  return <RouterProvider router={router} />;
}

// import Register from '../pages/Register';

// function App() {
//   return (
//     <div>
//       <Register />;
//     </div>
//   );
// }

// export default App;
