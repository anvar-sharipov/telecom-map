import { NavLink } from 'react-router-dom';

// 🔹 MobileLink (для мобильного меню)
type Props = {
  to: string;
  text: string;
  setOpen: (v: boolean) => void;
};

const MobileLink = ({ to, text, setOpen }: Props) => {
  return (
    <NavLink to={to} onClick={() => setOpen(false)} className="py-2 text-lg">
      {text}
    </NavLink>
  );
};

export default MobileLink;
