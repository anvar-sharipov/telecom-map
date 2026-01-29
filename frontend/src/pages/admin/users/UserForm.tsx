import { Formik, Form, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';
import { useEffect, useState } from 'react';
import { store } from '../../../app/store';
import { clearAuth, setAuth } from '../../../features/auth/authSlice';
const API_URL = import.meta.env.VITE_API_URL;

export type UserFormMode = 'create' | 'edit';

export interface UserFormValues {
  username: string;
  full_name: string;
  password?: string;
  is_active: boolean;
  group_ids: number[];
}

type Props = {
  mode: UserFormMode;
  initialValues?: Partial<UserFormValues>;
  onSubmit: (values: UserFormValues) => Promise<void>;
  onCancel?: () => void;
};

export default function UserForm({ mode, initialValues, onSubmit, onCancel }: Props) {
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchGroups = async () => {
    setLoading(true);

    try {
      let token = store.getState().auth.accessToken;

      // 1️⃣ первый запрос
      let res = await fetch(`${API_URL}/admin/groups?is_active=true`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      // 2️⃣ access истёк → пробуем refresh
      if (res.status === 401) {
        console.log('access expired → refresh');

        const refreshRes = await fetch(`${API_URL}/auth/refresh`, {
          method: 'POST',
          credentials: 'include',
        });

        // 🔴 только тут решаем logout
        if (refreshRes.status === 401 || refreshRes.status === 403) {
          store.dispatch(clearAuth());
          throw new Error('refresh expired');
        }

        if (!refreshRes.ok) {
          throw new Error('refresh failed');
        }

        const refreshData = await refreshRes.json();

        // 3️⃣ сохраняем новый access
        store.dispatch(setAuth(refreshData));

        token = refreshData.access_token;

        // 4️⃣ повторяем запрос
        res = await fetch(`${API_URL}/admin/groups?is_active=true`, {
          method: 'GET',
          credentials: 'include',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
      }

      if (!res.ok) {
        throw new Error('failed to fetch groups');
      }

      const data = await res.json();
      setGroups(data);
      //   console.log('data', data);
    } catch (err) {
      console.error('cant get groups', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGroups();
  }, []);

  useEffect(() => {
    console.log('groups', groups);
  }, [groups]);

  const isCreate = mode === 'create';

  const validationSchema = Yup.object({
    username: Yup.string().required('Required'),
    full_name: Yup.string().required('Required'),
    ...(isCreate && {
      password: Yup.string().min(6).required('Required'),
    }),
  });

  return (
    <Formik<UserFormValues>
      initialValues={{
        username: '',
        full_name: '',
        password: '',
        is_active: true,
        group_ids: [],
        ...initialValues,
      }}
      validationSchema={validationSchema}
      enableReinitialize
      onSubmit={onSubmit}
    >
      {({ isSubmitting }) => (
        <Form className="space-y-4">
          {/* username */}
          <div>
            <label className="block text-sm">Username</label>
            <Field name="username" disabled={!isCreate} className="input" />
            <ErrorMessage name="username" component="div" className="text-xs text-red-500" />
          </div>

          {/* full name */}
          <div>
            <label className="block text-sm">Full name</label>
            <Field name="full_name" className="input" />
          </div>

          {/* password */}
          {isCreate && (
            <div>
              <label className="block text-sm">Password</label>
              <Field name="password" type="password" className="input" />
            </div>
          )}

          {/* active */}
          <label className="flex items-center gap-2">
            <Field type="checkbox" name="is_active" />
            Active
          </label>

          {/* groups */}
          <div>
            <label className="block mb-1 text-sm">Groups</label>
            <div className="space-y-1">
              {groups?.map((g: any) => (
                <label key={g.ID} className="flex items-center gap-2">
                  <Field type="checkbox" name="group_ids" value={String(g.ID)} />
                  {g.Name}
                </label>
              ))}
            </div>
          </div>

          {/* actions */}
          <div className="flex justify-end gap-2 pt-4">
            {onCancel && (
              <button type="button" onClick={onCancel} className="btn-secondary">
                Cancel
              </button>
            )}
            <button type="submit" disabled={isSubmitting} className="btn-primary">
              {isCreate ? 'Create' : 'Save'}
            </button>
          </div>
        </Form>
      )}
    </Formik>
  );
}

// Использование (Create)
{
  /* <UserForm
  mode="create"
  onSubmit={async (values) => {
    await apiCreateUser(values);
    closeModal();
    fetchUsers();
  }}
/> */
}

// Использование (Edit)
{
  /* <UserForm
  mode="edit"
  initialValues={user}
  onSubmit={async (values) => {
    await apiUpdateUser(user.id, values);
    closeModal();
    fetchUsers();
  }}
/> */
}
