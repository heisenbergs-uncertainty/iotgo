import Link from 'next/link';
import {api} from '@/lib/api';
import {Device} from '@/types/api';
import ProtectedPage from "@/components/ProtectedRoute";

export default async function DevicesPage() {
    let devices: Device[] = await api.getDevices()

    return (
        <ProtectedPage>
            <div className="py-8 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto">
                <div className="mb-8">
                    <h1 className="text-4xl sm:text-5xl font-bold text-[--color-text-primary] mb-4">
                        Manage Devices
                    </h1>
                    <p className="text-lg text-[--color-text-secondary]">
                        View and manage your IoT devices from a central dashboard.
                    </p>
                </div>

                <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] rounded-lg border-0 overflow-hidden">
                    <div className="p-6 border-b border-gray-700/30">
                        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                            <div>
                                <h2 className="text-2xl font-semibold text-[--color-text-primary] flex items-center">
                                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 mr-2 text-[--color-text-accent]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z" />
                                    </svg>
                                    Devices ({devices.length})
                                </h2>
                                <p className="text-[--color-text-secondary] mt-1 hidden sm:block">
                                    Manage your connected IoT devices
                                </p>
                            </div>

                            <Link
                                href="/devices/new"
                                className="bg-[--color-text-accent] text-white font-semibold py-2 px-4 rounded-lg hover:bg-opacity-90 transition-all duration-300 inline-flex items-center shadow-lg shadow-[--color-text-accent]/20"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                                    <path fillRule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clipRule="evenodd" />
                                </svg>
                                Add Device
                            </Link>
                        </div>
                    </div>

                    {devices.length > 0 ? (
                        <div className="overflow-x-auto">
                            <table className="w-full text-left text-[--color-text-primary]">
                                <thead>
                                <tr className="bg-[--color-background-primary]/30">
                                    <th className="py-3 px-4 font-medium">Name</th>
                                    <th className="py-3 px-4 font-medium">Manufacturer</th>
                                    <th className="py-3 px-4 font-medium">Type</th>
                                    <th className="py-3 px-4 font-medium">Building</th>
                                    <th className="py-3 px-4 font-medium text-right">Actions</th>
                                </tr>
                                </thead>
                                <tbody>
                                {devices.map((device, index) => (
                                    <tr key={device.Id} className={`border-t border-gray-700/30 hover:bg-[--color-background-primary]/10 transition-colors`}>
                                        <td className="py-3 px-4 font-medium">{device.Name}</td>
                                        <td className="py-3 px-4 text-[--color-text-secondary]">{device.Manufacturer}</td>
                                        <td className="py-3 px-4">
                                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-[--color-background-primary]/50 text-[--color-text-primary]">
                                                    {device.Type}
                                                </span>
                                        </td>
                                        <td className="py-3 px-4 text-[--color-text-secondary]">{device.Building}</td>
                                        <td className="py-3 px-4 text-right space-x-1">
                                            <Link
                                                href={`/devices/${device.Id}`}
                                                className="inline-flex items-center justify-center h-8 w-8 rounded-lg bg-[--color-background-primary]/30 text-[--color-text-primary] hover:bg-[--color-text-accent]/10 hover:text-[--color-text-accent] transition-colors"
                                                title="View details"
                                            >
                                                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                                                    <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                                                    <path fillRule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clipRule="evenodd" />
                                                </svg>
                                            </Link>
                                            <Link
                                                href={`/devices/${device.Id}/edit`}
                                                className="inline-flex items-center justify-center h-8 w-8 rounded-lg bg-[--color-background-primary]/30 text-[--color-text-primary] hover:bg-[--color-text-accent]/10 hover:text-[--color-text-accent] transition-colors"
                                                title="Edit device"
                                            >
                                                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                                                    <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                                                </svg>
                                            </Link>
                                            <button/*
                                                onClick={async () => {
                                                    if (confirm('Are you sure you want to delete this device?')) {
                                                        try {
                                                            await api.deleteDevice(device.Id);
                                                            window.location.reload();
                                                        } catch (error) {
                                                            alert('Failed to delete device: ' + error.message);
                                                        }
                                                    }
                                                }}
*/

                                                                                                className="inline-flex items-center justify-center h-8 w-8 rounded-lg bg-[--color-background-primary]/30 text-[--color-text-primary] hover:bg-[--color-danger]/10 hover:text-[--color-danger] transition-colors"
                                                title="Delete device"
                                            >
                                                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                                                    <path fillRule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clipRule="evenodd" />
                                                </svg>
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                                </tbody>
                            </table>
                        </div>
                    ) : (
                        <div className="text-center py-16 px-4">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-16 w-16 mx-auto text-[--color-text-secondary]/30 mb-4" viewBox="0 0 20 20" fill="currentColor">
                                <path d="M5 3a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2V5a2 2 0 00-2-2H5zM5 11a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2v-2a2 2 0 00-2-2H5zM11 5a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V5zM14 11a1 1 0 011 1v1h1a1 1 0 110 2h-1v1a1 1 0 11-2 0v-1h-1a1 1 0 110-2h1v-1a1 1 0 011-1z" />
                            </svg>
                            <h3 className="text-xl font-bold text-[--color-text-primary] mb-2">No devices found</h3>
                            <p className="text-[--color-text-secondary] mb-6 max-w-md mx-auto">
                                You haven't added any devices to your IoT platform yet. Get started by adding your first device.
                            </p>
                            <Link
                                href="/devices/new"
                                className="bg-[--color-text-accent] text-white font-semibold py-2.5 px-5 rounded-lg hover:bg-opacity-90 transition-all duration-300 inline-flex items-center shadow-lg shadow-[--color-text-accent]/20"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                                    <path fillRule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clipRule="evenodd" />
                                </svg>
                                Add Your First Device
                            </Link>
                        </div>
                    )}
                </div>
            </div>
        </ProtectedPage>
    );
}