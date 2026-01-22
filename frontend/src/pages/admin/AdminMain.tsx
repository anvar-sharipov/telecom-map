import { useState } from 'react';
import UsersTable from './users/UserTable';

const tables = [
  { key: 'users', label: 'Users' },
  { key: 'groups', label: 'Groups' },
  { key: 'user_groups', label: 'User Groups' },
];

const AdminMain = () => {
  const [activeTable, setActiveTable] = useState('users');

  return (
    <div className="flex w-full min-h-[calc(100vh-64px)]">
      {/* LEFT SIDEBAR */}
      <aside className="w-64 border-r bg-indigo-50 dark:bg-zinc-900">
        <div className="p-4 text-lg font-semibold bg-indigo-400 border-b dark:bg-zinc-900">
          Панель Администратора
        </div>

        <ul className="p-2 space-y-1">
          {tables.map((table) => (
            <li
              key={table.key}
              onClick={() => setActiveTable(table.key)}
              className={`
                px-3 py-2 rounded cursor-pointer
                ${activeTable === table.key ? 'text-blue-700 font-medium' : 'hover:bg-gray-100 dark:hover:bg-gray-700'}
              `}
            >
              {table.label}
            </li>
          ))}
        </ul>
      </aside>

      {/* RIGHT CONTENT */}
      <main className="flex-1 p-6">
        <h1 className="mb-4 text-2xl font-semibold">{activeTable}</h1>

        <div className="p-4 border rounded">{activeTable === 'users' && <UsersTable />}</div>
      </main>
    </div>
  );
};

export default AdminMain;
