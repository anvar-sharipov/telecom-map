import Modal from '../../../components/UI/Modal/Modal';
import ModalHeader from '../../../components/UI/Modal/ModalHeader';
import UserForm from './UserForm';
// import { apiCreateUser } from '../../../api/users';
const API_URL = import.meta.env.VITE_API_URL;
import { store } from '../../../app/store';

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
            const payload = {
              ...values,
              group_ids: values.group_ids.map(Number),
            };

            console.log('payload', payload);

            const res = await fetch(`${API_URL}/admin/users`, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${store.getState().auth.accessToken}`,
              },
              credentials: 'include',
              body: JSON.stringify(payload),
            });

            if (!res.ok) {
              const text = await res.text();
              throw new Error(text);
            }

            onSuccess();
            onClose();
          }}
          onCancel={onClose}
        />
      </div>
    </Modal>
  );
};

export default AddUserModal;
