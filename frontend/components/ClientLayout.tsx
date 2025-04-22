"use client";

import Navbar from "@/components/Navbar";
import {Tooltip} from "react-tooltip";
import {ReactNode} from "react";
import {AuthProvider} from "@/contexts/AuthContext";
import {ThemeProvider} from "@/contexts/ThemeContext";

interface ClientLayoutProps {
    children: ReactNode;
}

export default function ClientLayout({children}: ClientLayoutProps) {
    return (
        <ThemeProvider>
            <AuthProvider>
                <Navbar/>
                <main className="container mx-auto pt-20 px-4">{children}</main>
                <footer
                    className="bg-[--color-background-secondary] text-[--color-text-secondary] py-5 text-center mt-20">
                    <p>© 2025 IoTGo Platform. All rights reserved.</p>
                </footer>
                <Tooltip id="devices-tooltip"/>
                <Tooltip id="admin-tooltip"/>
                <Tooltip id="users-tooltip"/>
                <Tooltip id="login-tooltip"/>
            </AuthProvider>
        </ThemeProvider>
    );
}
