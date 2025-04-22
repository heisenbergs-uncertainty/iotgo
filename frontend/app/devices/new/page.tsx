'use client';

import {useState} from 'react';
import {useRouter} from 'next/navigation';
import {api} from '@/lib/api';
import Link from "next/link";
import ProtectedPage from "@/components/ProtectedRoute";
import {allowedDeviceTypes, allowedLocations} from "@/lib/const";
import {AllowedDeviceType, AllowedLocation} from "@/types/devices";

export default function NewDevicePage() {
    const [form, setForm] = useState({ Name: '', Manufacturer: '', Type: '', Building: '' });
    const [errors, setErrors] = useState({ Name: '', Manufacturer: '', Type: '', Building: '', general: '' });
    const [isSubmitting, setIsSubmitting] = useState(false);
    const router = useRouter();

    const validateForm = () => {
        const newErrors = { Name: '', Manufacturer: '', Type: '', Building: '', general: '' };
        if (!form.Name) newErrors.Name = 'Device name is required.';
        if (!form.Manufacturer) newErrors.Manufacturer = 'Manufacturer is required.';
        if (!form.Type) newErrors.Type = 'Device type is required.';
        if (!form.Building) newErrors.Building = 'Building location is required.';
        setErrors(newErrors);
        return !Object.values(newErrors).some((error) => error);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!validateForm()) return;

        setErrors({ Name: '', Manufacturer: '', Type: '', Building: '', general: '' });
        setIsSubmitting(true);
        try {
            await api.createDevice(form);
            router.push('/devices');
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'Failed to create device.';
            setErrors((prev) => ({ ...prev, general: errorMessage }));
            setIsSubmitting(false);
        }
    };

    return (
        <ProtectedPage>
            <div className="py-8 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto">
                <div className="mb-8">
                    <Link href="/devices" className="inline-flex items-center text-[--color-text-secondary] hover:text-[--color-text-accent] mr-4 transition-colors">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                            <path fillRule="evenodd" d="M9.707 16.707a1 1 0 01-1.414 0l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l4.293 4.293a1 1 0 010 1.414z" clipRule="evenodd" />
                        </svg>
                        Back to Devices
                    </Link>
                    <h1 className="text-4xl font-bold text-[--color-text-primary] mb-3">Add New Device</h1>
                    <p className="text-lg text-[--color-text-secondary]">Register a new IoT device to start monitoring it.</p>
                </div>

                <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] rounded-lg border-0 overflow-hidden max-w-2xl">
                    <div className="p-6">
                        {errors.general && (
                            <div className="mb-6 p-3 rounded-lg bg-red-500/10 border border-red-500/20">
                                <p className="text-[--color-danger] text-sm">{errors.general}</p>
                            </div>
                        )}

                        <form onSubmit={handleSubmit} className="space-y-5">
                            <div>
                                <label htmlFor="name" className="block text-[--color-text-primary] text-sm font-medium mb-2">
                                    Device Name <span className="text-[--color-danger]">*</span>
                                </label>
                                <input
                                    type="text"
                                    id="name"
                                    value={form.Name}
                                    onChange={(e) => setForm({ ...form, Name: e.target.value })}
                                    className="w-full py-2 px-3 bg-[--color-background-primary] text-[--color-text-primary] rounded-lg border border-gray-700/30 focus:outline-none focus:ring-2 focus:ring-[--color-text-accent] focus:border-transparent transition-all"
                                    placeholder="Enter device name"
                                />
                                {errors.Name && <p className="text-[--color-danger] text-sm mt-1">{errors.Name}</p>}
                            </div>

                            <div>
                                <label htmlFor="manufacturer" className="block text-[--color-text-primary] text-sm font-medium mb-2">
                                    Manufacturer <span className="text-[--color-danger]">*</span>
                                </label>
                                <input
                                    type="text"
                                    id="manufacturer"
                                    value={form.Manufacturer}
                                    onChange={(e) => setForm({ ...form, Manufacturer: e.target.value })}
                                    className="w-full py-2 px-3 bg-[--color-background-primary] text-[--color-text-primary] rounded-lg border border-gray-700/30 focus:outline-none focus:ring-2 focus:ring-[--color-text-accent] focus:border-transparent transition-all"
                                    placeholder="Enter manufacturer name"
                                />
                                {errors.Manufacturer && <p className="text-[--color-danger] text-sm mt-1">{errors.Manufacturer}</p>}
                            </div>

                            <div>
                                <label htmlFor="type" className="block text-[--color-text-primary] text-sm font-medium mb-2">
                                    Device Type <span className="text-[--color-danger]">*</span>
                                </label>
                                <select
                                    id="type"
                                    value={form.Type}
                                    onChange={(e) => setForm({ ...form, Type: e.target.value as AllowedDeviceType })}
                                    className="w-full py-2 px-3 bg-[--color-background-primary] text-[--color-text-primary] rounded-lg border border-gray-700/30 focus:outline-none focus:ring-2 focus:ring-[--color-text-accent] focus:border-transparent transition-all"
                                >
                                    <option value="" disabled>Select device type</option>
                                    {allowedDeviceTypes.map((key, val) => (
                                        <option key={val} value={key}>{key}</option>
                                    ))}
                                </select>
                                {errors.Type && <p className="text-[--color-danger] text-sm mt-1">{errors.Type}</p>}                                {errors.Type && <p className="text-[--color-danger] text-sm mt-1">{errors.Type}</p>}
                            </div>

                            <div>
                                <label htmlFor="building" className="block text-[--color-text-primary] text-sm font-medium mb-2">
                                    Building <span className="text-[--color-danger]">*</span>
                                </label>
                                <select
                                    id="building"
                                    value={form.Building}
                                    onChange={(e) => setForm({ ...form, Building: e.target.value as AllowedLocation})}
                                    className="w-full py-2 px-3 bg-[--color-background-primary] text-[--color-text-primary] rounded-lg border border-gray-700/30 focus:outline-none focus:ring-2 focus:ring-[--color-text-accent] focus:border-transparent transition-all"
                                >
                                    <option value="" disabled>Select building location</option>
                                    {allowedLocations.map((key, val) => (
                                        <option key={val} value={key}>{key}</option>
                                    ))}

                                </select>
                                {errors.Building && <p className="text-[--color-danger] text-sm mt-1">{errors.Building}</p>}
                            </div>

                            <div className="pt-4 border-t border-gray-700/30 flex items-center justify-end space-x-4">
                                <Link href="/devices" className="px-4 py-2 bg-[--color-background-primary] text-[--color-text-primary] rounded-lg border border-gray-700/30 hover:bg-[--color-background-primary]/80 transition-colors">
                                    Cancel
                                </Link>
                                <button
                                    type="submit"
                                    disabled={isSubmitting}
                                    className="bg-[--color-text-accent] text-white font-medium py-2 px-4 rounded-lg hover:bg-opacity-90 focus:outline-none focus:ring-2 focus:ring-[--color-text-accent] focus:ring-opacity-50 transition-all duration-300 flex items-center shadow-lg shadow-[--color-text-accent]/20 disabled:opacity-70"
                                >
                                    {isSubmitting ? 'Creating...' : 'Create Device'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </ProtectedPage>
    );
}