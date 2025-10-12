import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Layout from "../components/layout";
import api, { type Event, eventAPI } from "../lib/api";
import {
    type AuthResponse,
    clearAuthData,
    getAdminUser,
    isAuthenticated,
    saveAuthData,
} from "../lib/auth";

// Dashboard Stats Interface
interface DashboardStats {
    totalEvents: number;
    totalTicketsSold: number;
    totalRevenue: number;
    upcomingEvents: number;
}

// Event Management State Interface
interface EventFilters {
    search: string;
    sortBy: "title" | "created_at" | "start_time";
    sortOrder: "asc" | "desc";
    page: number;
    limit: number;
}

// Login Form Component
function LoginForm({ onLoginSuccess }: { onLoginSuccess: () => void }) {
    const [credentials, setCredentials] = useState({
        email: "",
        password: "",
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError("");

        try {
            const response = await api.post("/admin/login", credentials);
            const authData: AuthResponse = response.data;

            saveAuthData(authData);
            onLoginSuccess();
        } catch (error: unknown) {
            const errorMessage = error instanceof Error && "response" in error
                ? (error as any).response?.data?.error || "Login failed"
                : "Login failed";
            setError(errorMessage);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex justify-center items-center bg-gray-50 min-h-screen">
            <div className="bg-white shadow-lg mx-4 p-8 border w-full max-w-md">
                <div className="mb-8 text-center">
                    <h1 className="mb-2 font-bold text-black text-3xl">
                        Admin Login
                    </h1>
                    <p className="text-gray-600">
                        Sign in to access the admin panel
                    </p>
                </div>

                {error && (
                    <div className="bg-red-50 mb-6 p-4 border border-red-200 rounded-md">
                        <p className="text-red-600 text-sm">{error}</p>
                    </div>
                )}

                <form onSubmit={handleSubmit} className="space-y-6">
                    <div>
                        <label className="block mb-2 font-medium text-black text-sm">
                            Username
                        </label>
                        <input
                            type="text"
                            required
                            className="px-4 py-3 border border-gray-300 focus:border-black rounded-md focus:outline-none w-full"
                            value={credentials.email}
                            onChange={(e) =>
                                setCredentials((prev) => ({
                                    ...prev,
                                    email: e.target.value,
                                }))}
                        />
                    </div>

                    <div>
                        <label className="block mb-2 font-medium text-black text-sm">
                            Password
                        </label>
                        <input
                            type="password"
                            required
                            className="px-4 py-3 border border-gray-300 focus:border-black rounded-md focus:outline-none w-full"
                            value={credentials.password}
                            onChange={(e) =>
                                setCredentials((prev) => ({
                                    ...prev,
                                    password: e.target.value,
                                }))}
                        />
                    </div>

                    <button
                        type="submit"
                        disabled={loading}
                        className="bg-black hover:bg-gray-800 disabled:opacity-50 px-6 py-3 rounded-md w-full font-medium text-white transition-colors duration-200"
                    >
                        {loading ? "Signing in..." : "Sign In"}
                    </button>
                </form>
            </div>
        </div>
    );
}

// Admin Dashboard Component
function AdminDashboard() {
    const [activeTab, setActiveTab] = useState<
        "dashboard" | "events" | "settings"
    >("dashboard");
    const [stats, setStats] = useState<DashboardStats>({
        totalEvents: 0,
        totalTicketsSold: 0,
        totalRevenue: 0,
        upcomingEvents: 0,
    });
    const [recentEvents, setRecentEvents] = useState<Event[]>([]);
    const [loading, setLoading] = useState(true);

    // Event management state
    const [allEvents, setAllEvents] = useState<Event[]>([]);
    const [eventsLoading, setEventsLoading] = useState(false);
    const [totalEvents, setTotalEvents] = useState(0);
    const [eventFilters, setEventFilters] = useState<EventFilters>({
        search: "",
        sortBy: "created_at",
        sortOrder: "desc",
        page: 1,
        limit: 10,
    });

    useEffect(() => {
        fetchDashboardData();
    }, []);

    const fetchDashboardData = async () => {
        try {
            setLoading(true);

            // Fetch recent events from API
            const response = await eventAPI.getRecentEvents(10);
            if (
                response.data && response.data.data && response.data.data.events
            ) {
                const events = response.data.data.events;
                setRecentEvents(events);

                // Calculate stats from events
                const now = new Date();
                const upcomingCount = events.filter((event: Event) =>
                    new Date(event.start_time) > now
                ).length;

                setStats({
                    totalEvents: events.length,
                    totalTicketsSold: events.reduce(
                        (sum: number, event: Event) =>
                            sum + (event.total_seats - event.available_seats),
                        0,
                    ),
                    totalRevenue: events.reduce(
                        (sum: number, event: Event) =>
                            sum +
                            (event.total_seats - event.available_seats) * 50,
                        0,
                    ), // Mock revenue calculation
                    upcomingEvents: upcomingCount,
                });
            }
        } catch (error) {
            console.error("Error fetching dashboard data:", error);
        } finally {
            setLoading(false);
        }
    };

    const fetchAllEvents = useCallback(async () => {
        console.log("Fetching all events...", { eventFilters });
        try {
            setEventsLoading(true);

            const offset = (eventFilters.page - 1) * eventFilters.limit;
            const response = await eventAPI.getAllEventsWithFilters({
                limit: eventFilters.limit,
                offset: offset,
                // Note: search, sortBy, sortOrder not yet supported by backend
            });

            console.log("API Response:", response.data);

            if (response.data) {
                // The backend returns events directly in the 'data' field
                const events = Array.isArray(response.data.data)
                    ? response.data.data
                    : [];
                console.log("Setting events:", events.length, "events");
                setAllEvents(events);
                setTotalEvents(response.data.pagination?.total_count || 0);
            }
        } catch (error) {
            console.error("Error fetching all events:", error);
        } finally {
            setEventsLoading(false);
        }
    }, [eventFilters]);

    useEffect(() => {
        console.log("useEffect triggered:", { activeTab });
        if (activeTab === "events") {
            console.log("Active tab is events, fetching...");
            fetchAllEvents();
        }
    }, [activeTab, fetchAllEvents]);

    const handleLogout = () => {
        clearAuthData();
        window.location.reload();
    };

    const handleSearchChange = (search: string) => {
        setEventFilters((prev) => ({ ...prev, search, page: 1 }));
    };

    const handleSortChange = (
        sortBy: "title" | "created_at" | "start_time",
    ) => {
        setEventFilters((prev) => ({
            ...prev,
            sortBy,
            sortOrder: prev.sortBy === sortBy && prev.sortOrder === "asc"
                ? "desc"
                : "asc",
            page: 1,
        }));
    };

    const handlePageChange = (page: number) => {
        setEventFilters((prev) => ({ ...prev, page }));
    };

    const getTotalPages = () => {
        return Math.ceil(totalEvents / eventFilters.limit);
    };

    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: "USD",
        }).format(amount);
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    const getStatusColor = (status: string) => {
        switch (status.toLowerCase()) {
            case "active":
            case "published":
                return "bg-green-100 text-green-800";
            case "draft":
                return "bg-yellow-100 text-yellow-800";
            case "cancelled":
                return "bg-red-100 text-red-800";
            default:
                return "bg-gray-100 text-gray-800";
        }
    };

    const adminUser = getAdminUser();

    return (
        <Layout>
            <div className="mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl">
                {/* Header */}
                <div className="flex justify-between items-center mb-8">
                    <div>
                        <h1 className="font-bold text-black text-3xl">
                            Admin Panel
                        </h1>
                        <p className="mt-2 text-gray-600">
                            Welcome back,{" "}
                            {adminUser?.name || adminUser?.username || "Admin"}!
                        </p>
                    </div>
                    <button
                        onClick={handleLogout}
                        className="btn-secondary"
                    >
                        Logout
                    </button>
                </div>

                {/* Navigation Tabs */}
                <div className="mb-8 border-gray-200 border-b">
                    <nav className="flex space-x-8">
                        {[
                            { id: "dashboard", label: "Dashboard", icon: "📊" },
                            { id: "events", label: "Events", icon: "🎫" },
                            { id: "settings", label: "Settings", icon: "⚙️" },
                        ].map((tab) => (
                            <button
                                key={tab.id}
                                onClick={() =>
                                    setActiveTab(
                                        tab.id as
                                            | "dashboard"
                                            | "events"
                                            | "settings",
                                    )}
                                className={`py-2 px-1 border-b-2 font-medium text-sm whitespace-nowrap ${
                                    activeTab === tab.id
                                        ? "border-black text-black"
                                        : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                                }`}
                            >
                                {tab.icon} {tab.label}
                            </button>
                        ))}
                    </nav>
                </div>

                {/* Tab Content */}
                {activeTab === "dashboard" && (
                    <div>
                        {/* Stats Cards */}
                        {loading
                            ? (
                                <div className="gap-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 mb-8">
                                    {[1, 2, 3, 4].map((i) => (
                                        <div
                                            key={i}
                                            className="card loading"
                                        >
                                            <div className="bg-gray-200 mb-2 h-4">
                                            </div>
                                            <div className="bg-gray-200 h-6">
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )
                            : (
                                <div className="gap-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 mb-8">
                                    <div className="card">
                                        <div className="mb-2 font-medium text-gray-500 text-sm uppercase tracking-wide">
                                            Total Events
                                        </div>
                                        <div className="font-bold text-black text-3xl">
                                            {stats.totalEvents}
                                        </div>
                                    </div>
                                    <div className="card">
                                        <div className="mb-2 font-medium text-gray-500 text-sm uppercase tracking-wide">
                                            Tickets Sold
                                        </div>
                                        <div className="font-bold text-black text-3xl">
                                            {stats.totalTicketsSold
                                                .toLocaleString()}
                                        </div>
                                    </div>
                                    <div className="card">
                                        <div className="mb-2 font-medium text-gray-500 text-sm uppercase tracking-wide">
                                            Total Revenue
                                        </div>
                                        <div className="font-bold text-black text-3xl">
                                            {formatCurrency(stats.totalRevenue)}
                                        </div>
                                    </div>
                                    <div className="card">
                                        <div className="mb-2 font-medium text-gray-500 text-sm uppercase tracking-wide">
                                            Upcoming Events
                                        </div>
                                        <div className="font-bold text-black text-3xl">
                                            {stats.upcomingEvents}
                                        </div>
                                    </div>
                                </div>
                            )}

                        {/* Recent Events */}
                        <div className="bg-white shadow-sm card">
                            <div className="flex justify-between items-center px-6 py-4 border-b border-b-gray-500">
                                <h2 className="font-bold text-black text-xl">
                                    Recent Events
                                </h2>
                                <Link
                                    to="/create-event"
                                    className="text-sm btn-primary"
                                >
                                    Create Event
                                </Link>
                            </div>

                            {loading
                                ? (
                                    <div className="p-6">
                                        <div className="space-y-4">
                                            {[1, 2, 3].map((i) => (
                                                <div
                                                    key={i}
                                                    className="bg-gray-200 h-16 animate-pulse"
                                                />
                                            ))}
                                        </div>
                                    </div>
                                )
                                : recentEvents.length === 0
                                ? (
                                    <div className="p-6 text-center">
                                        <p className="mb-4 text-gray-500">
                                            No events found
                                        </p>
                                        <Link
                                            to="/create-event"
                                            className="btn-primary"
                                        >
                                            Create Your First Event
                                        </Link>
                                    </div>
                                )
                                : (
                                    <div className="divide-y divide-gray-300">
                                        {recentEvents.map((event) => (
                                            <div
                                                key={event.id}
                                                className="hover:bg-gray-50 px-6 py-4 transition-colors duration-200"
                                            >
                                                <div className="flex justify-between items-center">
                                                    <div className="flex-1">
                                                        <div className="flex items-center space-x-3">
                                                            <h3 className="font-semibold text-black">
                                                                {event.title}
                                                            </h3>
                                                            <span
                                                                className={`px-2 py-1 text-xs font-medium uppercase tracking-wide rounded ${
                                                                    getStatusColor(
                                                                        event
                                                                            .status,
                                                                    )
                                                                }`}
                                                            >
                                                                {event.status}
                                                            </span>
                                                        </div>
                                                        <div className="flex items-center space-x-4 mt-1 text-gray-500 text-sm">
                                                            <span>
                                                                📅 {formatDate(
                                                                    event
                                                                        .start_time,
                                                                )}
                                                            </span>
                                                            <span>
                                                                📍{" "}
                                                                {event.location}
                                                            </span>
                                                            <span>
                                                                🎫 {event
                                                                    .total_seats -
                                                                    event
                                                                        .available_seats}
                                                                {" "}
                                                                / {event
                                                                    .total_seats}
                                                                {" "}
                                                                sold
                                                            </span>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                )}
                        </div>
                    </div>
                )}

                {activeTab === "events" && (
                    <div className="space-y-6">
                        <div className="flex justify-between items-center">
                            <h2 className="font-bold text-black text-2xl">
                                Event Management
                            </h2>
                            <Link to="/create-event" className="btn-primary">
                                Create New Event
                            </Link>
                        </div>

                        <div className="card">
                            <div className="px-6 py-4 border-b">
                                <div className="flex justify-between items-center">
                                    <h3 className="font-semibold text-black text-lg">
                                        All Events
                                    </h3>
                                    <div className="text-gray-500 text-sm">
                                        {totalEvents} total events
                                    </div>
                                </div>

                                {/* Search and Filters */}
                                <div className="flex sm:flex-row flex-col gap-4 mt-4">
                                    {/* Search Input */}
                                    <div className="flex-1">
                                        <input
                                            type="text"
                                            placeholder="Search events..."
                                            value={eventFilters.search}
                                            onChange={(e) =>
                                                handleSearchChange(
                                                    e.target.value,
                                                )}
                                            className="input-field"
                                        />
                                    </div>

                                    {/* Sort Dropdown */}
                                    <div className="flex gap-2">
                                        <select
                                            value={eventFilters.sortBy}
                                            onChange={(e) =>
                                                handleSortChange(
                                                    e.target.value as
                                                        | "title"
                                                        | "created_at"
                                                        | "start_time",
                                                )}
                                            className="input-field"
                                        >
                                            <option value="created_at">
                                                Sort by Created
                                            </option>
                                            <option value="title">
                                                Sort by Name
                                            </option>
                                            <option value="start_time">
                                                Sort by Start Date
                                            </option>
                                        </select>

                                        <button
                                            onClick={() =>
                                                setEventFilters((prev) => ({
                                                    ...prev,
                                                    sortOrder:
                                                        prev.sortOrder === "asc"
                                                            ? "desc"
                                                            : "asc",
                                                }))}
                                            className="hover:bg-gray-50 px-3 py-2 border border-gray-300 rounded-md transition-colors"
                                            title={`Sort ${
                                                eventFilters.sortOrder === "asc"
                                                    ? "Descending"
                                                    : "Ascending"
                                            }`}
                                        >
                                            {eventFilters.sortOrder === "asc"
                                                ? "↑"
                                                : "↓"}
                                        </button>
                                    </div>
                                </div>
                            </div>

                            <div className="bg-yellow-50 p-4 border-b">
                                <div className="text-sm">
                                    Debug: eventsLoading={String(
                                        eventsLoading,
                                    )}, allEvents.length={allEvents.length},
                                    totalEvents={totalEvents}
                                </div>
                            </div>
                            {eventsLoading
                                ? (
                                    <div className="p-8 text-center">
                                        <div className="text-gray-500">
                                            Loading events...
                                        </div>
                                    </div>
                                )
                                : allEvents.length === 0
                                ? (
                                    <div className="p-8 text-center">
                                        <div className="text-gray-500">
                                            No events found
                                        </div>
                                        <div className="mt-2 text-gray-400 text-xs">
                                            totalEvents:{" "}
                                            {totalEvents}, eventsLoading:{" "}
                                            {String(eventsLoading)}
                                        </div>
                                    </div>
                                )
                                : (
                                    <div className="divide-y divide-gray-300">
                                        {allEvents.map((event) => (
                                            <div
                                                key={event.id}
                                                className="hover:bg-gray-50 px-6 py-4 transition-colors duration-200"
                                            >
                                                <div className="flex justify-between items-center">
                                                    <div className="flex-1">
                                                        <div className="flex items-center space-x-3">
                                                            <h4 className="font-semibold text-black">
                                                                {event.title}
                                                            </h4>
                                                            <span
                                                                className={`px-2 py-1 text-xs font-medium uppercase tracking-wide rounded ${
                                                                    getStatusColor(
                                                                        event
                                                                            .status,
                                                                    )
                                                                }`}
                                                            >
                                                                {event.status}
                                                            </span>
                                                        </div>
                                                        <p className="mt-1 text-gray-600 text-sm line-clamp-2">
                                                            {event.description}
                                                        </p>
                                                        <div className="flex items-center space-x-4 mt-2 text-gray-500 text-sm">
                                                            <span>
                                                                📅 {formatDate(
                                                                    event
                                                                        .start_time,
                                                                )}
                                                            </span>
                                                            <span>
                                                                📍{" "}
                                                                {event.location}
                                                            </span>
                                                            <span>
                                                                🎫 {event
                                                                    .available_seats}
                                                                {" "}
                                                                / {event
                                                                    .total_seats}
                                                                {" "}
                                                                available
                                                            </span>
                                                        </div>
                                                    </div>
                                                    <div className="flex space-x-2">
                                                        <Link
                                                            to={`/events/${event.slug}`}
                                                            className="text-sm btn-secondary"
                                                        >
                                                            View
                                                        </Link>
                                                        <Link
                                                            to={`/edit-event/${event.slug}`}
                                                            className="text-sm btn-primary"
                                                        >
                                                            Edit
                                                        </Link>
                                                    </div>
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                )}

                            {/* Pagination */}
                            {!eventsLoading && allEvents.length > 0 &&
                                getTotalPages() > 1 && (
                                <div className="bg-gray-50 px-6 py-4 border-t">
                                    <div className="flex justify-between items-center">
                                        <div className="text-gray-600 text-sm">
                                            Showing {((eventFilters.page - 1) *
                                                eventFilters.limit) + 1} to{" "}
                                            {Math.min(
                                                eventFilters.page *
                                                    eventFilters.limit,
                                                totalEvents,
                                            )} of {totalEvents} events
                                        </div>

                                        <div className="flex space-x-2">
                                            <button
                                                onClick={() =>
                                                    handlePageChange(
                                                        eventFilters.page - 1,
                                                    )}
                                                disabled={eventFilters.page ===
                                                    1}
                                                className="hover:bg-white disabled:opacity-50 px-3 py-1 border border-gray-300 rounded transition-colors disabled:cursor-not-allowed"
                                            >
                                                Previous
                                            </button>

                                            {Array.from({
                                                length: getTotalPages(),
                                            }, (_, i) => i + 1).map((page) => (
                                                <button
                                                    key={page}
                                                    onClick={() =>
                                                        handlePageChange(page)}
                                                    className={`px-3 py-1 border border-gray-300 rounded transition-colors ${
                                                        page ===
                                                                eventFilters
                                                                    .page
                                                            ? "bg-black text-white border-black"
                                                            : "hover:bg-white"
                                                    }`}
                                                >
                                                    {page}
                                                </button>
                                            ))}

                                            <button
                                                onClick={() =>
                                                    handlePageChange(
                                                        eventFilters.page + 1,
                                                    )}
                                                disabled={eventFilters.page ===
                                                    getTotalPages()}
                                                className="hover:bg-white disabled:opacity-50 px-3 py-1 border border-gray-300 rounded transition-colors disabled:cursor-not-allowed"
                                            >
                                                Next
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                )}

                {activeTab === "settings" && (
                    <div className="space-y-6">
                        <h2 className="font-bold text-black text-2xl">
                            Admin Settings
                        </h2>

                        <div className="gap-6 grid grid-cols-1 lg:grid-cols-2">
                            <div className="card">
                                <h3 className="mb-4 font-semibold text-black text-lg">
                                    Profile Information
                                </h3>
                                <div className="space-y-4">
                                    <div>
                                        <label className="block mb-1 font-medium text-gray-700 text-sm">
                                            Name
                                        </label>
                                        <p className="text-gray-900">
                                            {adminUser?.name || "Not set"}
                                        </p>
                                    </div>
                                    <div>
                                        <label className="block mb-1 font-medium text-gray-700 text-sm">
                                            Username
                                        </label>
                                        <p className="text-gray-900">
                                            {adminUser?.username}
                                        </p>
                                    </div>
                                    <div>
                                        <label className="block mb-1 font-medium text-gray-700 text-sm">
                                            Email
                                        </label>
                                        <p className="text-gray-900">
                                            {adminUser?.email}
                                        </p>
                                    </div>
                                    <div>
                                        <label className="block mb-1 font-medium text-gray-700 text-sm">
                                            Role
                                        </label>
                                        <p className="text-gray-900">
                                            {adminUser?.role}
                                        </p>
                                    </div>
                                </div>
                            </div>

                            <div className="card">
                                <h3 className="mb-4 font-semibold text-black text-lg">
                                    Quick Actions
                                </h3>
                                <div className="space-y-3">
                                    <Link
                                        to="/create-event"
                                        className="block w-full text-center btn-primary"
                                    >
                                        Create New Event
                                    </Link>
                                    <Link
                                        to="/events"
                                        className="block w-full text-center btn-secondary"
                                    >
                                        View All Events
                                    </Link>
                                    <Link
                                        to="/dashboard"
                                        className="block w-full text-center btn-secondary"
                                    >
                                        Public Dashboard
                                    </Link>
                                    <button
                                        onClick={handleLogout}
                                        className="w-full btn-danger"
                                    >
                                        Logout
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </Layout>
    );
}

// Main Admin Page Component
export default function AdminPage() {
    const [authenticated, setAuthenticated] = useState(false);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const checkAuth = async () => {
            try {
                setAuthenticated(isAuthenticated());
            } catch (error) {
                console.error("Auth check failed:", error);
                setAuthenticated(false);
            } finally {
                setLoading(false);
            }
        };

        checkAuth();
    }, []);

    const handleLoginSuccess = () => {
        setAuthenticated(true);
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="text-center">
                    <div className="mx-auto mb-4 border-t-2 border-black rounded-full w-8 h-8 animate-spin">
                    </div>
                    <p className="text-gray-600">Loading...</p>
                </div>
            </div>
        );
    }

    return authenticated
        ? <AdminDashboard />
        : <LoginForm onLoginSuccess={handleLoginSuccess} />;
}
