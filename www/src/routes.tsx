import { Route, Routes } from "react-router-dom";
import HomePage from "./pages/home";
import LoginPage from "./pages/auth/login";
import RegisterPage from "./pages/auth/register";
import EventsPage from "./pages/events";
import EventDetailPage from "./pages/event-detail";
import DashboardPage from "./pages/dashboard";
import AdminPage from "./pages/admin";

export default function AppRoutes() {
    return (
        <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/events" element={<EventsPage />} />
            <Route path="/events/:slug" element={<EventDetailPage />} />
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/auth/login" element={<LoginPage />} />
            <Route path="/auth/register" element={<RegisterPage />} />
            <Route path="/admin" element={<AdminPage />} />

            {/* Placeholder routes for other pages mentioned in navigation */}
            <Route
                path="/about"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            About Page - Coming Soon
                        </h1>
                    </div>
                }
            />
            <Route
                path="/contact"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            Contact Page - Coming Soon
                        </h1>
                    </div>
                }
            />
            <Route
                path="/careers"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            Careers Page - Coming Soon
                        </h1>
                    </div>
                }
            />
            <Route
                path="/press"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            Press Page - Coming Soon
                        </h1>
                    </div>
                }
            />
            <Route
                path="/help"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            Help Center - Coming Soon
                        </h1>
                    </div>
                }
            />
            <Route
                path="/privacy"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            Privacy Policy - Coming Soon
                        </h1>
                    </div>
                }
            />
            <Route
                path="/terms"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            Terms of Service - Coming Soon
                        </h1>
                    </div>
                }
            />
            <Route
                path="/faq"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            FAQ - Coming Soon
                        </h1>
                    </div>
                }
            />
            <Route
                path="/create-event"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            Create Event - Coming Soon
                        </h1>
                    </div>
                }
            />
            <Route
                path="/organizer"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            For Organizers - Coming Soon
                        </h1>
                    </div>
                }
            />
            <Route
                path="/pricing"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            Pricing - Coming Soon
                        </h1>
                    </div>
                }
            />
            <Route
                path="/checkout"
                element={
                    <div className="p-8 text-center">
                        <h1 className="font-bold text-2xl">
                            Checkout - Coming Soon
                        </h1>
                    </div>
                }
            />

            <Route
                path="*"
                element={
                    <div className="flex justify-center items-center min-h-screen">
                        <div className="text-center">
                            <h1 className="mb-4 font-bold text-black text-6xl">
                                404
                            </h1>
                            <p className="mb-6 text-gray-600 text-xl">
                                Page not found
                            </p>
                            <a href="/" className="btn-primary">
                                Go Home
                            </a>
                        </div>
                    </div>
                }
            />
        </Routes>
    );
}
