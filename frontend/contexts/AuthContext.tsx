'use client';

import {createContext, ReactNode, useContext, useEffect, useState} from 'react';
import {api} from '@/lib/api';
import {User} from "@/types/api";
import { useRouter } from 'next/router';

interface AuthContextType {
    user: User | null;
    setUser: (user: User | null) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        // Check session on mount
        validateSession();

    }, []);


    async function validateSession() {
        try {
            const data = await api.getMe();
            if (data.username && data.roles && data.user_id) {
                setUser({ 
                    Username: data.username, 
                    Roles: data.roles,  // Changed from Role to Roles to match backend
                    Id: data.user_id 
                });
            } else {
                setUser(null)
            }
        } catch (error) {
            console.error('Session check failed:', error);
        } finally {
            setIsLoading(false);
        }
    }

  

    return (
        <AuthContext.Provider value={{   user, setUser }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) throw new Error('useAuth must be used within an AuthProvider');
    return context;
};