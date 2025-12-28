import './i18n';
import './styles/index.css';
import App from './app/App.tsx';

// Подключил Redux к React. Теперь любой компонент может использовать Redux.
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import { store } from './app/store';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <Provider store={store}>
    <App />
  </Provider>,
);
