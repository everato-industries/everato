import { useState } from "react";
import Layout from "../components/layout";

export default function AdminPage() {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [credentials, setCredentials] = useState({
        username: "",
        password: "",
    });
    const [error, setError] = useState("");

    const handleLogin = (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        // Simple admin authentication - replace with actual authentication
        if (
            credentials.username === "admin" &&
            credentials.password === "admin123"
        ) {
            setIsAuthenticated(true);
        } else {
            setError("Invalid admin credentials");
        }
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setCredentials((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

    if (!isAuthenticated) {
        return (
            <Layout showNavbar={false} showFooter={false}>
                <div className="flex justify-center items-center bg-gray-50 min-h-screen">
                    <div className="bg-white shadow-lg p-8 rounded-lg w-full max-w-md">
                        <div className="mb-8 text-center">
                            <h1 className="font-bold text-gray-900 text-2xl">
                                Admin Panel
                            </h1>
                            <p className="mt-2 text-gray-600">
                                Please sign in to continue
                            </p>
                        </div>

                        <form onSubmit={handleLogin} className="space-y-6">
                            {error && (
                                <div className="bg-red-50 p-3 border border-red-200 rounded text-red-600 text-sm">
                                    {error}
                                </div>
                            )}

                            <div>
                                <label
                                    htmlFor="username"
                                    className="block font-medium text-gray-700 text-sm"
                                >
                                    Username
                                </label>
                                <input
                                    type="text"
                                    id="username"
                                    name="username"
                                    value={credentials.username}
                                    onChange={handleInputChange}
                                    className="mt-1 w-full input-field"
                                    required
                                />
                            </div>

                            <div>
                                <label
                                    htmlFor="password"
                                    className="block font-medium text-gray-700 text-sm"
                                >
                                    Password
                                </label>
                                <input
                                    type="password"
                                    id="password"
                                    name="password"
                                    value={credentials.password}
                                    onChange={handleInputChange}
                                    className="mt-1 w-full input-field"
                                    required
                                />
                            </div>

                            <button
                                type="submit"
                                className="w-full btn-primary"
                            >
                                Sign In
                            </button>
                        </form>

                        <div className="mt-6 text-gray-500 text-sm text-center">
                            <p>Default credentials: admin / admin123</p>
                        </div>
                    </div>
                </div>
            </Layout>
        );
    }

    return (
        <Layout>
            <div className="bg-gray-50 min-h-screen">
                <div className="mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl">
                    {/* Header */}
                    <div className="flex justify-between items-center mb-8">
                        <div>
                            <h1 className="font-bold text-gray-900 text-3xl">
                                Admin Dashboard
                            </h1>
                            <p className="mt-2 text-gray-600">
                                Manage your events and platform
                            </p>
                        </div>
                        <button
                            onClick={() => setIsAuthenticated(false)}
                            className="btn-secondary"
                        >
                            Logout
                        </button>
                    </div>

                    {/* Stats Cards */}
                    <div className="gap-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 mb-8">
                        <div className="card">
                            <div className="p-6">
                                <div className="flex justify-between items-center">
                                    <div>
                                        <p className="font-medium text-gray-600 text-sm">
                                            Total Events
                                        </p>
                                        <p className="font-bold text-gray-900 text-2xl">
                                            127
                                        </p>
                                    </div>
                                    <div className="bg-blue-100 p-3 rounded-lg">
                                        <svg
                                            className="w-6 h-6 text-blue-600"
                                            fill="none"
                                            stroke="currentColor"
                                            viewBox="0 0 24 24"
                                        >
                                            <path
                                                strokeLinecap="round"
                                                strokeLinejoin="round"
                                                strokeWidth={2}
                                                d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                                            />
                                        </svg>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="card">
                            <div className="p-6">
                                <div className="flex justify-between items-center">
                                    <div>
                                        <p className="font-medium text-gray-600 text-sm">
                                            Total Users
                                        </p>
                                        <p className="font-bold text-gray-900 text-2xl">
                                            1,428
                                        </p>
                                    </div>
                                    <div className="bg-green-100 p-3 rounded-lg">
                                        <svg
                                            className="w-6 h-6 text-green-600"
                                            fill="none"
                                            stroke="currentColor"
                                            viewBox="0 0 24 24"
                                        >
                                            <path
                                                strokeLinecap="round"
                                                strokeLinejoin="round"
                                                strokeWidth={2}
                                                d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"
                                            />
                                        </svg>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="card">
                            <div className="p-6">
                                <div className="flex justify-between items-center">
                                    <div>
                                        <p className="font-medium text-gray-600 text-sm">
                                            Revenue
                                        </p>
                                        <p className="font-bold text-gray-900 text-2xl">
                                            $24,890
                                        </p>
                                    </div>
                                    <div className="bg-yellow-100 p-3 rounded-lg">
                                        <svg
                                            className="w-6 h-6 text-yellow-600"
                                            fill="none"
                                            stroke="currentColor"
                                            viewBox="0 0 24 24"
                                        >
                                            <path
                                                strokeLinecap="round"
                                                strokeLinejoin="round"
                                                strokeWidth={2}
                                                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"
                                            />
                                        </svg>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="card">
                            <div className="p-6">
                                <div className="flex justify-between items-center">
                                    <div>
                                        <p className="font-medium text-gray-600 text-sm">
                                            Active Events
                                        </p>
                                        <p className="font-bold text-gray-900 text-2xl">
                                            23
                                        </p>
                                    </div>
                                    <div className="bg-purple-100 p-3 rounded-lg">
                                        <svg
                                            className="w-6 h-6 text-purple-600"
                                            fill="none"
                                            stroke="currentColor"
                                            viewBox="0 0 24 24"
                                        >
                                            <path
                                                strokeLinecap="round"
                                                strokeLinejoin="round"
                                                strokeWidth={2}
                                                d="M13 10V3L4 14h7v7l9-11h-7z"
                                            />
                                        </svg>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Quick Actions */}
                    <div className="gap-6 grid grid-cols-1 lg:grid-cols-2 mb-8">
                        <div className="card">
                            <div className="p-6">
                                <h3 className="mb-4 font-bold text-gray-900 text-lg">
                                    Quick Actions
                                </h3>
                                <div className="space-y-3">
                                    <button className="w-full text-left btn-primary">
                                        Create New Event
                                    </button>
                                    <button className="w-full text-left btn-secondary">
                                        Manage Users
                                    </button>
                                    <button className="w-full text-left btn-secondary">
                                        View Reports
                                    </button>
                                    <button className="w-full text-left btn-secondary">
                                        System Settings
                                    </button>
                                </div>
                            </div>
                        </div>

                        <div className="card">
                            <div className="p-6">
                                <h3 className="mb-4 font-bold text-gray-900 text-lg">
                                    Recent Activity
                                </h3>
                                <div className="space-y-3">
                                    <div className="flex justify-between items-center py-2 border-gray-100 border-b">
                                        <span className="text-gray-900 text-sm">
                                            New event created
                                        </span>
                                        <span className="text-gray-500 text-xs">
                                            2 hours ago
                                        </span>
                                    </div>
                                    <div className="flex justify-between items-center py-2 border-gray-100 border-b">
                                        <span className="text-gray-900 text-sm">
                                            User registration spike
                                        </span>
                                        <span className="text-gray-500 text-xs">
                                            4 hours ago
                                        </span>
                                    </div>
                                    <div className="flex justify-between items-center py-2 border-gray-100 border-b">
                                        <span className="text-gray-900 text-sm">
                                            Payment processed
                                        </span>
                                        <span className="text-gray-500 text-xs">
                                            6 hours ago
                                        </span>
                                    </div>
                                    <div className="flex justify-between items-center py-2">
                                        <span className="text-gray-900 text-sm">
                                            Event published
                                        </span>
                                        <span className="text-gray-500 text-xs">
                                            1 day ago
                                        </span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Recent Events Table */}
                    <div className="card">
                        <div className="p-6">
                            <h3 className="mb-4 font-bold text-gray-900 text-lg">
                                Recent Events
                            </h3>
                            <div className="overflow-x-auto">
                                <table className="min-w-full">
                                    <thead>
                                        <tr className="border-gray-200 border-b">
                                            <th className="px-4 py-3 font-medium text-gray-700 text-sm text-left">
                                                Event
                                            </th>
                                            <th className="px-4 py-3 font-medium text-gray-700 text-sm text-left">
                                                Date
                                            </th>
                                            <th className="px-4 py-3 font-medium text-gray-700 text-sm text-left">
                                                Tickets Sold
                                            </th>
                                            <th className="px-4 py-3 font-medium text-gray-700 text-sm text-left">
                                                Revenue
                                            </th>
                                            <th className="px-4 py-3 font-medium text-gray-700 text-sm text-left">
                                                Status
                                            </th>
                                            <th className="px-4 py-3 font-medium text-gray-700 text-sm text-left">
                                                Actions
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-gray-200">
                                        <tr>
                                            <td className="px-4 py-3 text-gray-900 text-sm">
                                                Tech Conference 2025
                                            </td>
                                            <td className="px-4 py-3 text-gray-600 text-sm">
                                                Oct 15, 2025
                                            </td>
                                            <td className="px-4 py-3 text-gray-600 text-sm">
                                                150/1000
                                            </td>
                                            <td className="px-4 py-3 text-gray-600 text-sm">
                                                $44,850
                                            </td>
                                            <td className="px-4 py-3 text-sm">
                                                <span className="bg-green-100 px-2 py-1 rounded-full text-green-800 text-xs">
                                                    Active
                                                </span>
                                            </td>
                                            <td className="px-4 py-3 text-sm">
                                                <button className="text-blue-600 hover:text-blue-800">
                                                    Edit
                                                </button>
                                            </td>
                                        </tr>
                                        <tr>
                                            <td className="px-4 py-3 text-gray-900 text-sm">
                                                Music Festival 2025
                                            </td>
                                            <td className="px-4 py-3 text-gray-600 text-sm">
                                                Nov 20, 2025
                                            </td>
                                            <td className="px-4 py-3 text-gray-600 text-sm">
                                                500/2000
                                            </td>
                                            <td className="px-4 py-3 text-gray-600 text-sm">
                                                $75,000
                                            </td>
                                            <td className="px-4 py-3 text-sm">
                                                <span className="bg-yellow-100 px-2 py-1 rounded-full text-yellow-800 text-xs">
                                                    Draft
                                                </span>
                                            </td>
                                            <td className="px-4 py-3 text-sm">
                                                <button className="text-blue-600 hover:text-blue-800">
                                                    Edit
                                                </button>
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </Layout>
    );
}
