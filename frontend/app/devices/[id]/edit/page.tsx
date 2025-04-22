'use client';

import {useEffect, useState} from 'react';
import {useParams, useRouter} from 'next/navigation';
import {api} from '@/lib/api';
import {Device} from '@/types/api';
import Link from "next/link";

export default function EditDevicePage() {
    const [form, setForm] = useState<Partial<Device>>({});
    const [error, setError] = useState('');
    const router = useRouter();
    const params = useParams();
    const id = Number(params.id);

    useEffect(() => {
        async function fetchDevice() {
            try {
                const device = await api.getDevice(id);
                setForm(device);
            } catch (err) {
                setError('Failed to load device');
            }
        }
        fetchDevice();
    }, [id]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        try {
            await api.updateDevice(id, form);
            router.push('/devices');
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to update device');
        }
    };

    if (!form.Id) return <div>Loading...</div>;

    return (
        <div className="py-8">
            <h1 className="text-4xl sm:text-5xl font-bold text-white mb-4">
                Edit Device
            </h1>
            <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] rounded-lg p-6 max-w-lg">
                {error && <p className="text-[--color-danger] mb-4">{error}</p>}
                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label htmlFor="name" className="block text-[--color-text-primary] mb-2">
                            Name
                        </label>
                        <input
                            type="text"
                            id="name"
                            value={form.Name || ''}
                            onChange={(e) => setForm({ ...form, Name: e.target.value })}
                            className="w-full px-4 py-2 bg-gray-800 text-[--color-text-primary] rounded border border-gray-700 focus:outline-none focus:border-[--color-primary]"
                            required
                        />
                    </div>
                    <div className="mb-4">
                        <label htmlFor="manufacturer" className="block text-[--color-text-primary] mb-2">
                            Manufacturer
                        </label>
                        <input
                            type="text"
                            id="manufacturer"
                            value={form.Manufacturer || ''}
                            onChange={(e) => setForm({ ...form, Manufacturer: e.target.value })}
                            className="w-full px-4 py-2 bg-gray-800 text-[--color-text-primary] rounded border border-gray-700 focus:outline-none focus:border-[--color-primary]"
                            required
                        />
                    </div>
                    <div className="mb-4">
                        <label htmlFor="type" className="block text-[--color-text-primary] mb-2">
                            Type
                        </label>
                        <input
                            type="text"
                            id="type"
                            value={form.Type || ''}
                            onChange={(e) => setForm({ ...form, Type: e.target.value })}
                            className="w-full px-4 py-2 bg-gray-800 text-[--color-text-primary] rounded border border-gray-700 focus:outline-none focus:border-[--color-primary]"
                            required
                        />
                    </div>
                    <div className="mb-6">
                        <label htmlFor="building" className="block text-[--color-text-primary] mb-2">
                            Building
                        </label>
                        <input
                            type="text"
                            id="building"
                            value={form.Building || ''}
                            onChange={(e) => setForm({ ...form, Building: e.target.value })}
                            className="w-full px-4 py-2 bg-gray-800 text-[--color-text-primary] rounded border border-gray-700 focus:outline-none focus:border-[--color-primary]"
                            required
                        />
                    </div>
                    <div className="flex space-x-4">
                        <button
                            type="submit"
                            className="bg-[--color-primary] text-white font-semibold py-2 px-4 rounded hover:bg-blue-600 transition-colors"
                        >
                            Update Device
                        </button>
                        <Link
                            href="/devices"
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