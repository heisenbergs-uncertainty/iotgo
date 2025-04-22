import Link from 'next/link';
import {api} from '@/lib/api';
import {redirect} from 'next/navigation';

export default async function AdminPage() {
    let deviceCount = 0;
    let userCount = 0;

    try {
        const mainData = await api.getMainData();
        deviceCount = mainData.DeviceCount;
        const users = await api.getUsers({ limit: 1 }); // Just to get count
        userCount = users.length > 0 ? (await api.getUsers()).length : 0;
    } catch (error: any) {
        if (error.message.includes('Unauthorized')) {
            redirect('/login');
        }
        if (error.message.includes('Admin role required')) {
            redirect('/');
        }
        console.error('Failed to fetch admin data:', error);
    }

    return (
        <div className="py-8">
            <h1 className="text-4xl sm:text-5xl font-bold text-white mb-4">
                Admin Dashboard
            </h1>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] rounded-lg p-6">
                    <h2 className="text-2xl font-semibold text-white mb-4">
                        <i className="bi bi-hdd-stack mr-2"></i>Devices
                    </h2>
                    <p className="text-4xl font-bold text-[--color-text-accent]">{deviceCount}</p>
                    <p className="text-[--color-text-secondary] mb-4">Total devices</p>
                    <Link
                        href="/devices"
                        className="bg-[--color-primary] text-white font-semibold py-2 px-4 rounded hover:bg-blue-600 transition-colors"
                    >
                        Manage Devices
                    </Link>
                </div>
                <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] rounded-lg p-6">
                    <h2 className="text-2xl font-semibold text-white mb-4">
                        <i className="bi bi-people mr-2"></i>Users
                    </h2>
                    <p className="text-4xl font-bold text-[--color-text-accent]">{userCount}</p>
                    <p className="text-[--color-text-secondary] mb-4">Total users</p>
                    <Link
                        href="/admin/users"
                        className="bg-[--color-primary] text-white font-semibold py-2 px-4 rounded hover:bg-blue-600 transition-colors"
                    >
                        Manage Users
                    </Link>
                </div>
            </div>
        </div>
    );
}