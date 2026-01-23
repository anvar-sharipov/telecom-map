import { useEffect, useState } from 'react';
import type { UserDTO } from '../../../types/user';
import LoadingSpin from '../../../components/UI/LoadingSpin';
import { useTranslation } from 'react-i18next';
import { store } from '../../../app/store';
import { setAuth, clearAuth } from '../../../features/auth/authSlice';
import AddUserModal from './AddUserModal';
import { Plus } from 'lucide-react';

const API_URL = import.meta.env.VITE_API_URL;

const UsersTable: React.FC = () => {
  const [users, setUsers] = useState<UserDTO[]>([]);
  const [loading, setLoading] = useState(true);
  const [addOpen, setAddOpen] = useState(false);
  const { t } = useTranslation('common');

  const fetchUsers = async () => {
    setLoading(true);

    try {
      let token = store.getState().auth.accessToken;

      // 1️⃣ первый запрос
      let res = await fetch(`${API_URL}/admin/users`, {
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
        res = await fetch(`${API_URL}/admin/users`, {
          method: 'GET',
          credentials: 'include',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
      }

      if (!res.ok) {
        throw new Error('failed to fetch users');
      }

      const data = await res.json();
      setUsers(data);
    } catch (err) {
      console.error('cant get users', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  if (loading) {
    return <LoadingSpin />;
  }

  return (
    <div className="shadow rounded-xl">
      <div className="flex items-center justify-between px-6 py-4 border-b">
        <h2 className="text-lg font-semibold">{t('Users')}</h2>

        <button
          onClick={() => setAddOpen(true)}
          className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700"
        >
          <Plus size={16} />
          Add user
        </button>
      </div>

      {/* DESKTOP TABLE */}
      <div className="hidden md:block">
        <table className="w-full">
          <thead>
            <tr>
              <th className="px-4 py-3 text-left">ID</th>
              <th className="px-4 py-3 text-left">Username</th>
              <th className="px-4 py-3 text-left">Full name</th>
              <th className="px-4 py-3 text-left">Groups</th>
              <th className="px-4 py-3 text-center">Active</th>
              <th className="px-4 py-3 text-left">Created</th>
            </tr>
          </thead>

          <tbody>
            {users.map((user) => (
              <tr key={user.id} className="border-t hover:bg-gray-50 dark:hover:bg-gray-800">
                <td className="px-4 py-3">{user.id}</td>
                <td className="px-4 py-3 font-medium">{user.username}</td>
                <td className="px-4 py-3">{user.full_name}</td>

                <td className="px-4 py-3">
                  <div className="flex flex-wrap gap-1">
                    {user.groups.map((group) => (
                      <span key={group.id} className="px-2 py-1 text-blue-700 bg-blue-100 rounded">
                        {group.name}
                      </span>
                    ))}
                  </div>
                </td>

                <td className="px-4 py-3 text-center">
                  {user.is_active ? (
                    <span className="font-medium text-green-600">Yes</span>
                  ) : (
                    <span className="font-medium text-red-500">No</span>
                  )}
                </td>

                <td className="px-4 py-3 text-gray-500">
                  {new Date(user.created_at).toLocaleDateString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* MOBILE CARDS */}
      <div className="p-4 space-y-3 md:hidden">
        {users.map((user) => (
          <div key={user.id} className="p-4 shadow dark:border dark:border-gray-700 rounded-xl">
            <div className="flex justify-between">
              <div>
                <p className="font-semibold">{user.username}</p>
                <p className="text-gray-500">{user.full_name}</p>
              </div>

              <span
                className={`text-xs px-2 py-1 rounded ${
                  user.is_active
                    ? 'bg-green-100 text-green-700 dark:bg-green-700 dark:text-green-100'
                    : 'bg-red-100 text-red-600 dark:bg-red-600 dark:text-red-100'
                }`}
              >
                {user.is_active ? 'Active' : 'Inactive'}
              </span>
            </div>

            <div className="flex flex-wrap gap-1 mt-3">
              {user.groups.map((group) => (
                <span key={group.id} className="px-2 py-1 text-blue-700 bg-blue-100 rounded">
                  {group.name}
                </span>
              ))}
            </div>

            <div className="mt-3 text-gray-500">
              Created: {new Date(user.created_at).toLocaleDateString()}
            </div>
          </div>
        ))}
      </div>
      <AddUserModal open={addOpen} onClose={() => setAddOpen(false)} onSuccess={fetchUsers} />
    </div>
  );
};

export default UsersTable;
