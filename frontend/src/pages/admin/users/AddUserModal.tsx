import Modal from '../../../components/UI/Modal/Modal';
import ModalHeader from '../../../components/UI/Modal/ModalHeader';
import UserForm from './UserForm';
// import { apiCreateUser } from '../../../api/users';

type Props = {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
};

const AddUserModal = ({ open, onClose, onSuccess }: Props) => {
  return (
    <Modal open={open} onClose={onClose} closeOnBackdrop width="md">
      <ModalHeader title="Add user" onClose={onClose} />

      <div className="p-6">
        <UserForm
          mode="create"
          onSubmit={async (values) => {
            // await apiCreateUser(values);
            onSuccess(); // 🔥 обновляем таблицу
            onClose(); // 🔥 закрываем модалку
          }}
          onCancel={onClose}
        />
      </div>
    </Modal>
  );
};

export default AddUserModal;
