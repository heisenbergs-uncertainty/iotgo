'use client';

import {useState} from 'react';
import {useRouter} from 'next/navigation';
import {api} from '@/lib/api';
import Link from "next/link";

export default function NewUserPage() {
    const [form, setForm] = useState({ Username: '', Password: '', Role: 'user' });
    const [error, setError] = useState('');
    const router = useRouter();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        try {
            await api.createUser(form);
            router.push('/admin/users');
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to create user');
        }
    };

    return (
        <div className="py-8">
            <h1 className="text-4xl sm:text-5xl font-bold text-white mb-4">
                Add New User
            </h1>
            <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] rounded-lg p-6 max-w-lg">
                {error && <p className="text-[--color-danger] mb-4">{error}</p>}
                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label htmlFor="username" className="block text-[--color-text-primary] mb-2">
                            Username
                        </label>
                        <input
                            type="text"
                            id="username"
                            value={form.Username}
                            onChange={(e) => setForm({ ...form, Username: e.target.value })}
                            className="w-full px-4 py-2 bg-gray-800 text-[--color-text-primary] rounded border border-gray-700 focus:outline-none focus:border-[--color-primary]"
                            required
                        />
                    </div>
                    <div className="mb-4">
                        <label htmlFor="password" className="block text-[--color-text-primary] mb-2">
                            Password
                        </label>
                        <input
                            type="password"
                            id="password"
                            value={form.Password}
                            onChange={(e) => setForm({ ...form, Password: e.target.value })}
                            className="w-full px-4 py-2 bg-gray-800 text-[--color-text-primary] rounded border border-gray-700 focus:outline-none focus:border-[--color-primary]"
                            required
                        />
                    </div>
                    <div className="mb-6">
                        <label htmlFor="role" className="block text-[--color-text-primary] mb-2">
                            Role
                        </label>
                        <select
                            id="role"
                            value={form.Role}
                            onChange={(e) => setForm({ ...form, Role: e.target.value })}
                            className="w-full px-4 py-2 bg-gray-800 text-[--color-text-primary] rounded border border-gray-700 focus:outline-none focus:border-[--color-primary]"
                        >
                            <option value="user">User</option>
                            <option value="admin">Admin</option>
                        </select>
                    </div>
                    <div className="flex space-x-4">
                        <button
                            type="submit"
                            className="bg-[--color-primary] text-white font-semibold py-2 px-4 rounded hover:bg-blue-600 transition-colors"
                        >
                            Create User
                        </button>
                        <Link
                            href="/admin/users"
                            className="bg-[--color-secondary] text-white font-semibold py-2 px-4 rounded hover:bg-gray-600 transition-colors"
                        >
                            Cancel
                        </Link>
                    </div>
                </form>
            </div>
        </div>
    );
}