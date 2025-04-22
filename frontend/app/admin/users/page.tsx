import Link from 'next/link';
import {api} from '@/lib/api';
import {User} from '@/types/api';
import {redirect} from 'next/navigation';

export default async function UsersPage() {
    let users: User[] = [];

    try {
        users = await api.getUsers({ limit: 10 });
    } catch (error: any) {
        if (error.message.includes('Unauthorized')) {
            redirect('/login');
        }
        if (error.message.includes('Admin role required')) {
            redirect('/');
        }
        console.error('Failed to fetch users:', error);
    }

    return (
        <div className="py-8">
            <div className="mb-12">
                <h1 className="text-4xl sm:text-5xl font-bold text-white mb-4">
                    Manage Users
                </h1>
                <p className="text-lg text-[--color-text-primary]">
                    View and manage user accounts.
                </p>
            </div>
            <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] rounded-lg p-6">
                <div className="flex justify-between items-center mb-4">
                    <h2 className="text-2xl font-semibold text-white">Users</h2>
                    <Link
                        href="/admin/users/new"
                        className="bg-[--color-primary] text-white font-semibold py-2 px-4 rounded hover:bg-blue-600 transition-colors"
                    >
                        Add User
                    </Link>
                </div>
                {users.length > 0 ? (
                    <div className="overflow-x-auto">
                        <table className="w-full text-left text-[--color-text-primary]">
                            <thead>
                            <tr className="border-b border-gray-700">
                                <th className="py-2 px-4">Username</th>
                                <th className="py-2 px-4">Role</th>
                                <th className="py-2 px-4">Actions</th>
                            </tr>
                            </thead>
                            <tbody>
                            {users.map((user) => (
                                <tr key={user.Id} className="border-b border-gray-700">
                                    <td className="py-2 px-4">{user.Username}</td>
                                    <td className="py-2 px-4">{user.Role}</td>
                                    <td className="py-2 px-4 space-x-2">
                                        <Link
                                            href={`/admin/users/${user.Id}/edit`}
                                            className="text-[--color-text-accent] hover:underline"
                                        >
                                            Edit
                                        </Link>
                                        <button
                                            onClick={async () => {
                                                if (confirm('Are you sure you want to delete this user?')) {
                                                    try {
                                                        await api.deleteUser(user.Id);
                                                        window.location.reload();
                                                    } catch (error: any) {
                                                        alert('Failed to delete user: ' + error.message);
                                                    }
                                                }
                                            }}
                                            className="text-[--color-danger] hover:underline"
                                        >
                                            Delete
                                        </button>
                                    </td>
                                </tr>
                            ))}
                            </tbody>
                        </table>
                    </div>
                ) : (
                    <p className="text-[--color-text-primary]">No users found.</p>
                )}
            </div>
        </div>
    );
}