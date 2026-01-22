import type { UserDTO } from '../../../types/user';

const mockUsers: UserDTO[] = [
  {
    id: 1,
    username: 'admin',
    full_name: 'Admin User',
    is_active: true,
    created_at: '2024-01-10T12:00:00Z',
    groups: [{ id: 1, name: 'admin', description: 'Administrators', is_active: true }],
  },
  {
    id: 2,
    username: 'john',
    full_name: 'John Doe',
    is_active: false,
    created_at: '2024-01-12T15:30:00Z',
    groups: [{ id: 2, name: 'manager', description: 'Managers', is_active: true }],
  },
];

const UsersTable: React.FC = () => {
  return (
    <div className="bg-white shadow rounded-xl">
      <div className="px-6 py-4 border-b">
        <h2 className="text-lg font-semibold">Users</h2>
      </div>

      <table className="w-full text-sm">
        <thead className="text-gray-600 bg-gray-100">
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
          {mockUsers.map((user) => (
            <tr key={user.id} className="transition border-t hover:bg-gray-50">
              <td className="px-4 py-3">{user.id}</td>
              <td className="px-4 py-3 font-medium">{user.username}</td>
              <td className="px-4 py-3">{user.full_name}</td>

              <td className="px-4 py-3">
                <div className="flex flex-wrap gap-1">
                  {user.groups.map((group) => (
                    <span
                      key={group.id}
                      className="px-2 py-1 text-xs text-blue-700 bg-blue-100 rounded"
                    >
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
  );
};

export default UsersTable;
