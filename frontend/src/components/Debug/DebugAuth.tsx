import { useSelector } from 'react-redux';
import type { RootState } from '../../app/store';

const DebugAuth = () => {
  const auth = useSelector((state: RootState) => state.auth);

  console.log('AUTH STATE == ', auth);

  return null;
};

export default DebugAuth;
