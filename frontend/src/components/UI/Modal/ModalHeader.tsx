import { X } from 'lucide-react';

type Props = {
  title: string;
  onClose: () => void;
};

const ModalHeader = ({ title, onClose }: Props) => {
  return (
    <div className="flex items-center justify-between px-6 py-4 border-b">
      <h3 className="text-lg font-semibold">{title}</h3>

      <button onClick={onClose} className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-600">
        <X size={18} />
      </button>
    </div>
  );
};

export default ModalHeader;
