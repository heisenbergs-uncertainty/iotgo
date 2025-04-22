'use client';

import {useEffect, useState} from 'react';
import {useAuth} from '@/contexts/AuthContext';
import Link from 'next/link';
import {api} from '@/lib/api';
import {useRouter} from "next/navigation";
import {User} from "@/types/api";

export default function LoginPage() {
    const [isLoading, setIsLoading] = useState(false);
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const { user, setUser, } = useAuth();
    const router = useRouter();
    useEffect(() => {
        setIsLoading(true)
        if(user){
            router.push("/")
        }
        setIsLoading(false)
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        try {
            console.log('Login attempt:', { username });
            const resp = await api.login(username, password);
            const authenticatedUser : User = {
                Roles: resp.roles,
                Username: resp.username,
                Id: resp.user_id,
            }
            console.log('Login successful, navigating to /');
            setUser(authenticatedUser);
            router.push("/")
        } catch (err) {
            console.error('Login error:', err);
            setError(err instanceof Error ? err.message : 'Login failed');
        }
    };

    return (
        <div className="flex items-center justify-center min-h-screen px-4 py-12 bg-[--color-background-primary] transition-colors duration-300">
            <div className="w-full max-w-md">
                {/* Logo/Brand */}
                <div className="text-center mb-8">
                    <h1 className="text-4xl font-bold text-[--color-text-accent]">IoTGo</h1>
                    <p className="text-[--color-text-secondary] mt-2">Your IoT Management Platform</p>
                </div>

                {/* Card */}
                <div className="bg-[--color-background-secondary] rounded-xl shadow-[0_8px_30px_var(--color-card-shadow)] p-8 transition-all duration-300">
                    <h2 className="text-2xl font-bold text-[--color-text-primary] mb-6">Log in to your account</h2>

                    {/* Error display */}
                    {error && (
                        <div className="mb-6 p-3 rounded-lg bg-red-500/10 border border-red-500/20">
                            <p className="text-[--color-danger] text-sm flex items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                                    <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                                </svg>
                                {error}
                            </p>
                        </div>
                    )}

                    <form onSubmit={handleSubmit}>
                        <div className="space-y-5">
                            {/* Username field */}
                            <div>
                                <label htmlFor="username" className="block text-[--color-text-primary] text-sm font-medium mb-2">
                                    Username
                                </label>
                                <div className="relative">
                                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-[--color-text-secondary]" viewBox="0 0 20 20" fill="currentColor">
                                            <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
                                        </svg>
                                    </div>
                                    <input
                                        type="text"
                                        id="username"
                                        value={username}
                                        onChange={(e) => setUsername(e.target.value)}
                                        className="w-full pl-10 pr-3 py-2 bg-[--color-background-primary] text-[--color-text-primary] rounded-lg border border-gray-700/30 focus:outline-none focus:ring-2 focus:ring-[--color-text-accent] focus:border-transparent transition-all"
                                        required
                                        autoComplete="username"
                                        placeholder="Enter your username"
                                    />
                                </div>
                            </div>

                            {/* Password field */}
                            <div>
                                <label htmlFor="password" className="block text-[--color-text-primary] text-sm font-medium mb-2">
                                    Password
                                </label>
                                <div className="relative">
                                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-[--color-text-secondary]" viewBox="0 0 20 20" fill="currentColor">
                                            <path fillRule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clipRule="evenodd" />
                                        </svg>
                                    </div>
                                    <input
                                        type="password"
                                        id="password"
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                        className="w-full pl-10 pr-3 py-2 bg-[--color-background-primary] text-[--color-text-primary] rounded-lg border border-gray-700/30 focus:outline-none focus:ring-2 focus:ring-[--color-text-accent] focus:border-transparent transition-all"
                                        required
                                        autoComplete="current-password"
                                        placeholder="Enter your password"
                                    />
                                </div>
                            </div>

                            {/* Login button */}
                            <div>
                                <button
                                    type="submit"
                                    disabled={isLoading}
                                    className="w-full bg-[--color-text-accent] text-white font-medium py-2.5 px-4 rounded-lg hover:bg-opacity-90 focus:outline-none focus:ring-2 focus:ring-[--color-text-accent] focus:ring-opacity-50 transition-all duration-300 flex items-center justify-center"
                                >
                                    {isLoading ? (
                                        <>
                                            <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                            </svg>
                                            Logging in...
                                        </>
                                    ) : 'Sign in'}
                                </button>
                            </div>
                        </div>
                    </form>
                </div>

                {/* Footer */}
                <div className="mt-8 text-center">
                    <p className="text-[--color-text-secondary] text-sm">
                        Don't have an account? Contact your administrator
                    </p>
                    <div className="mt-4">
                        <Link
                            href="/"
                            className="text-[--color-text-accent] hover:underline text-sm inline-flex items-center"
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
                                <path fillRule="evenodd" d="M9.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L7.414 9H15a1 1 0 110 2H7.414l2.293 2.293a1 1 0 010 1.414z" clipRule="evenodd" />
                            </svg>
                            Back to home
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
}