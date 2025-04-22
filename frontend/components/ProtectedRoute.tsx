"use client";

import {useAuth} from "@/contexts/AuthContext";
import {redirect, useRouter} from "next/navigation";
import {ReactNode, useEffect} from "react";

interface ProtectedRouteProps {
    children: ReactNode;
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
    const { user } = useAuth();
    const router = useRouter(); 

    useEffect(() => {
        if (user == null) {
            router.push("/login");
        }
    }, [ user, router]);

    // TODO: Replace with a loading indicator
    if (user === null) {
        return <div>Loading...</div>; // Show a loading indicator
    }

    if (!user) {
        return null; // Prevent rendering if user is not authenticated
    }

    return <>{children}</>; // Render the protected content
}