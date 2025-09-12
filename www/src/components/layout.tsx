import type { ReactNode } from "react";
import Navbar from "./navbar";
import Footer from "./footer";

interface LayoutProps {
    children: ReactNode;
    showNavbar?: boolean;
    showFooter?: boolean;
}

export default function Layout({
    children,
    showNavbar = true,
    showFooter = true,
}: LayoutProps) {
    return (
        <div className="flex flex-col bg-white min-h-screen">
            {showNavbar && <Navbar />}
            <main className="flex-1">
                {children}
            </main>
            {showFooter && <Footer />}
        </div>
    );
}
