import { useEffect, useState } from "react";

import Layout from "../components/layout";
import api from "../lib/api";
import {
    type AuthResponse,
    clearAuthData,
    getAdminUser,
    isAuthenticated,
    saveAuthData,
} from "../lib/auth";

// Event creation form data types
interface TicketType {
    name: string;
    price: number;
    available_tickets: number;
}

interface Coupon {
    code: string;
    discount_percentage: number;
    valid_from: string;
    valid_until: string;
    usage_limit: number;
}

interface EventFormData {
    title: string;
    description: string;
    start_time: string;
    end_time: string;
    location: string;
    admin_id: string;
    banner_url?: string;
    icon_url?: string;
    total_seats: number;
    available_seats: number;
    status?: string;
    ticket_types: TicketType[];
    coupons: Coupon[];
}

interface RecentEvent {
    id: string;
    title: string;
    description: string;
    start_time: string;
    end_time: string;
    location: string;
    status: string;
    slug: string;
    total_seats: number;
    available_seats: number;
    created_at: string;
}

// Event Creation Form Component
interface EventCreationFormProps {
    onSubmit: (data: EventFormData) => void;
    loading: boolean;
    onCancel: () => void;
}

function EventCreationForm(
    { onSubmit, loading, onCancel }: EventCreationFormProps,
) {
    const [formData, setFormData] = useState<EventFormData>({
        title: "",
        description: "",
        start_time: "",
        end_time: "",
        location: "",
        admin_id: "",
        banner_url: "",
        icon_url: "",
        total_seats: 100,
        available_seats: 100,
        status: "CREATED",
        ticket_types: [{
            name: "General Admission",
            price: 0,
            available_tickets: 100,
        }],
        coupons: [],
    });

    // Set admin_id when component mounts
    useEffect(() => {
        const adminUser = getAdminUser();
        if (adminUser?.id) {
            setFormData((prev) => ({
                ...prev,
                admin_id: adminUser.id,
                available_seats: prev.total_seats,
            }));
        }
    }, []);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        // Convert datetime-local to ISO 8601 format
        const submitData = {
            ...formData,
            start_time: formData.start_time
                ? new Date(formData.start_time).toISOString()
                : "",
            end_time: formData.end_time
                ? new Date(formData.end_time).toISOString()
                : "",
        };

        onSubmit(submitData);
    };

    const updateField = (
        field: keyof EventFormData,
        value: string | number | TicketType[] | Coupon[],
    ) => {
        setFormData((prev) => ({ ...prev, [field]: value }));
    };

    const addTicketType = () => {
        setFormData((prev) => ({
            ...prev,
            ticket_types: [...prev.ticket_types, {
                name: "",
                price: 0,
                available_tickets: 50,
            }],
        }));
    };

    const updateTicketType = (
        index: number,
        field: keyof TicketType,
        value: string | number,
    ) => {
        setFormData((prev) => ({
            ...prev,
            ticket_types: prev.ticket_types.map((ticket, i) =>
                i === index ? { ...ticket, [field]: value } : ticket
            ),
        }));
    };

    const removeTicketType = (index: number) => {
        setFormData((prev) => ({
            ...prev,
            ticket_types: prev.ticket_types.filter((_, i) => i !== index),
        }));
    };

    const addCoupon = () => {
        setFormData((prev) => ({
            ...prev,
            coupons: [...prev.coupons, {
                code: "",
                discount_percentage: 10,
                valid_from: new Date().toISOString().slice(0, 16),
                valid_until: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000)
                    .toISOString().slice(0, 16),
                usage_limit: 100,
            }],
        }));
    };

    const updateCoupon = (
        index: number,
        field: keyof Coupon,
        value: string | number,
    ) => {
        setFormData((prev) => ({
            ...prev,
            coupons: prev.coupons.map((coupon, i) =>
                i === index ? { ...coupon, [field]: value } : coupon
            ),
        }));
    };

    const removeCoupon = (index: number) => {
        setFormData((prev) => ({
            ...prev,
            coupons: prev.coupons.filter((_, i) => i !== index),
        }));
    };

    return (
        <div>
            {/* Quick action button */}
            <div className="mb-4">
                <button
                    type="button"
                    onClick={() => window.open("/www/create-event", "_blank")}
                    className="w-full text-sm text-center btn-secondary"
                >
                    🚀 Open in New Page (Full Form)
                </button>
            </div>

            <form
                onSubmit={handleSubmit}
                className="space-y-4 max-h-96 overflow-y-auto"
            >
                {/* Basic Event Info */}
                <div className="gap-4 grid grid-cols-1">
                    <div>
                        <label className="block mb-1 font-medium text-gray-700 text-sm">
                            Event Title
                        </label>
                        <input
                            type="text"
                            value={formData.title}
                            onChange={(e) =>
                                updateField("title", e.target.value)}
                            className="w-full input-field"
                            required
                        />
                    </div>
                </div>

                <div>
                    <label className="block mb-1 font-medium text-gray-700 text-sm">
                        Description
                    </label>
                    <textarea
                        value={formData.description}
                        onChange={(e) =>
                            updateField("description", e.target.value)}
                        className="w-full input-field"
                        rows={3}
                        required
                    />
                </div>

                <div className="gap-4 grid grid-cols-1 md:grid-cols-2">
                    <div>
                        <label className="block mb-1 font-medium text-gray-700 text-sm">
                            Start Time
                        </label>
                        <input
                            type="datetime-local"
                            value={formData.start_time}
                            onChange={(e) =>
                                updateField("start_time", e.target.value)}
                            className="w-full input-field"
                            required
                        />
                    </div>
                    <div>
                        <label className="block mb-1 font-medium text-gray-700 text-sm">
                            End Time
                        </label>
                        <input
                            type="datetime-local"
                            value={formData.end_time}
                            onChange={(e) =>
                                updateField("end_time", e.target.value)}
                            className="w-full input-field"
                            required
                        />
                    </div>
                </div>

                <div className="gap-4 grid grid-cols-1 md:grid-cols-2">
                    <div>
                        <label className="block mb-1 font-medium text-gray-700 text-sm">
                            Location
                        </label>
                        <input
                            type="text"
                            value={formData.location}
                            onChange={(e) =>
                                updateField("location", e.target.value)}
                            className="w-full input-field"
                            required
                        />
                    </div>
                    <div>
                        <label className="block mb-1 font-medium text-gray-700 text-sm">
                            Total Seats
                        </label>
                        <input
                            type="number"
                            value={formData.total_seats}
                            onChange={(e) => {
                                const newTotalSeats = parseInt(e.target.value);
                                setFormData((prev) => ({
                                    ...prev,
                                    total_seats: newTotalSeats,
                                    available_seats: newTotalSeats,
                                }));
                            }}
                            className="w-full input-field"
                            min="1"
                            required
                        />
                    </div>
                </div>

                {/* Ticket Types */}
                <div>
                    <div className="flex justify-between items-center mb-2">
                        <label className="block font-medium text-gray-700 text-sm">
                            Ticket Types
                        </label>
                        <button
                            type="button"
                            onClick={addTicketType}
                            className="text-blue-600 hover:text-blue-800 text-sm"
                        >
                            + Add Ticket Type
                        </button>
                    </div>

                    {formData.ticket_types.map((ticket, index) => (
                        <div
                            key={index}
                            className="mb-2 p-3 border border-gray-200 rounded"
                        >
                            {formData.ticket_types.length > 1 && (
                                <div className="flex justify-end mb-2">
                                    <button
                                        type="button"
                                        onClick={() => removeTicketType(index)}
                                        className="text-red-600 hover:text-red-800 text-sm"
                                    >
                                        Remove
                                    </button>
                                </div>
                            )}
                            <div className="gap-3 grid grid-cols-1 md:grid-cols-3">
                                <div>
                                    <label className="block mb-1 font-medium text-gray-600 text-xs">
                                        Ticket Name
                                    </label>
                                    <input
                                        type="text"
                                        placeholder="e.g., General Admission"
                                        value={ticket.name}
                                        onChange={(e) =>
                                            updateTicketType(
                                                index,
                                                "name",
                                                e.target.value,
                                            )}
                                        className="text-sm input-field"
                                        required
                                    />
                                </div>
                                <div>
                                    <label className="block mb-1 font-medium text-gray-600 text-xs">
                                        Price ($)
                                    </label>
                                    <input
                                        type="number"
                                        placeholder="0.00"
                                        value={ticket.price}
                                        onChange={(e) =>
                                            updateTicketType(
                                                index,
                                                "price",
                                                parseFloat(e.target.value),
                                            )}
                                        className="text-sm input-field"
                                        min="0"
                                        step="0.01"
                                        required
                                    />
                                </div>
                                <div>
                                    <label className="block mb-1 font-medium text-gray-600 text-xs">
                                        Available Tickets
                                    </label>
                                    <input
                                        type="number"
                                        placeholder="50"
                                        value={ticket.available_tickets}
                                        onChange={(e) =>
                                            updateTicketType(
                                                index,
                                                "available_tickets",
                                                parseInt(e.target.value),
                                            )}
                                        className="text-sm input-field"
                                        min="1"
                                        required
                                    />
                                </div>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Coupons */}
                <div>
                    <div className="flex justify-between items-center mb-2">
                        <label className="block font-medium text-gray-700 text-sm">
                            Coupons (Optional)
                        </label>
                        <button
                            type="button"
                            onClick={addCoupon}
                            className="text-blue-600 hover:text-blue-800 text-sm"
                        >
                            + Add Coupon
                        </button>
                    </div>

                    {formData.coupons.map((coupon, index) => (
                        <div
                            key={index}
                            className="mb-3 p-4 border border-gray-200 rounded"
                        >
                            <div className="flex justify-end mb-3">
                                <button
                                    type="button"
                                    onClick={() => removeCoupon(index)}
                                    className="text-red-600 hover:text-red-800 text-sm"
                                >
                                    Remove
                                </button>
                            </div>

                            <div className="gap-3 grid grid-cols-1 md:grid-cols-2 mb-3">
                                <div>
                                    <label className="block mb-1 font-medium text-gray-600 text-xs">
                                        Coupon Code
                                    </label>
                                    <input
                                        type="text"
                                        placeholder="e.g., EARLYBIRD2025"
                                        value={coupon.code}
                                        onChange={(e) =>
                                            updateCoupon(
                                                index,
                                                "code",
                                                e.target.value.toUpperCase(),
                                            )}
                                        className="text-sm input-field"
                                        required
                                    />
                                </div>
                                <div>
                                    <label className="block mb-1 font-medium text-gray-600 text-xs">
                                        Discount (%)
                                    </label>
                                    <input
                                        type="number"
                                        placeholder="10"
                                        value={coupon.discount_percentage}
                                        onChange={(e) =>
                                            updateCoupon(
                                                index,
                                                "discount_percentage",
                                                parseFloat(e.target.value),
                                            )}
                                        className="text-sm input-field"
                                        min="1"
                                        max="100"
                                        step="0.1"
                                        required
                                    />
                                </div>
                            </div>

                            <div className="gap-3 grid grid-cols-1 md:grid-cols-3">
                                <div>
                                    <label className="block mb-1 font-medium text-gray-600 text-xs">
                                        Valid From
                                    </label>
                                    <input
                                        type="datetime-local"
                                        value={coupon.valid_from}
                                        onChange={(e) =>
                                            updateCoupon(
                                                index,
                                                "valid_from",
                                                e.target.value,
                                            )}
                                        className="text-sm input-field"
                                        required
                                    />
                                </div>
                                <div>
                                    <label className="block mb-1 font-medium text-gray-600 text-xs">
                                        Valid Until
                                    </label>
                                    <input
                                        type="datetime-local"
                                        value={coupon.valid_until}
                                        onChange={(e) =>
                                            updateCoupon(
                                                index,
                                                "valid_until",
                                                e.target.value,
                                            )}
                                        className="text-sm input-field"
                                        required
                                    />
                                </div>
                                <div>
                                    <label className="block mb-1 font-medium text-gray-600 text-xs">
                                        Usage Limit
                                    </label>
                                    <input
                                        type="number"
                                        placeholder="100"
                                        value={coupon.usage_limit}
                                        onChange={(e) =>
                                            updateCoupon(
                                                index,
                                                "usage_limit",
                                                parseInt(e.target.value),
                                            )}
                                        className="text-sm input-field"
                                        min="1"
                                        required
                                    />
                                </div>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Form Actions */}
                <div className="flex gap-2 pt-4">
                    <button
                        type="submit"
                        disabled={loading}
                        className={`flex-1 btn-primary ${
                            loading ? "opacity-50 cursor-not-allowed" : ""
                        }`}
                    >
                        {loading
                            ? (
                                <div className="flex justify-center items-center">
                                    <div className="mr-2 border-white border-b-2 rounded-full w-4 h-4 animate-spin">
                                    </div>
                                    Creating...
                                </div>
                            )
                            : (
                                "Create Event"
                            )}
                    </button>
                    <button
                        type="button"
                        onClick={onCancel}
                        className="btn-secondary"
                        disabled={loading}
                    >
                        Cancel
                    </button>
                </div>
            </form>
        </div>
    );
}

export default function AdminPage() {
    const [authenticated, setAuthenticated] = useState(false);
    const [loading, setLoading] = useState(true);
    const [loginLoading, setLoginLoading] = useState(false);
    const [credentials, setCredentials] = useState({
        email: "",
        password: "",
    });
    const [error, setError] = useState("");
    const [showEventForm, setShowEventForm] = useState(false);
    const [eventSubmitting, setEventSubmitting] = useState(false);
    const [recentEvents, setRecentEvents] = useState<RecentEvent[]>([]);
    const [eventsLoading, setEventsLoading] = useState(false);

    // Fetch recent events
    const fetchRecentEvents = async () => {
        setEventsLoading(true);
        try {
            const response = await api.get("/events/recent?limit=5");
            if (response.data?.data?.events) {
                setRecentEvents(response.data.data.events);
            }
        } catch (error) {
            console.error("Failed to fetch recent events:", error);
        } finally {
            setEventsLoading(false);
        }
    };

    // Check authentication status on component mount
    useEffect(() => {
        const checkAuth = async () => {
            if (isAuthenticated()) {
                setAuthenticated(true);
                await fetchRecentEvents();
            }
            setLoading(false);
        };

        checkAuth();
    }, []);

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        setLoginLoading(true);

        try {
            const response = await api.post("/admin/login", credentials);
            const data: AuthResponse = response.data;

            // Save authentication data to localStorage and cookies
            saveAuthData(data);
            setAuthenticated(true);

            console.log("✅ Admin login successful:", data.user.name);
        } catch (err: unknown) {
            console.error("❌ Admin login failed:", err);
            let errorMessage = "Login failed. Please try again.";

            if (typeof err === "object" && err !== null) {
                const error = err as {
                    response?: { data?: { message?: string } };
                    message?: string;
                };
                if (error.response?.data?.message) {
                    errorMessage = error.response.data.message;
                } else if (error.message) {
                    errorMessage = error.message;
                }
            }

            setError(errorMessage);
        } finally {
            setLoginLoading(false);
        }
    };

    const handleLogout = () => {
        clearAuthData();
        setAuthenticated(false);
        setCredentials({ email: "", password: "" });
    };

    const handleEventSubmit = async (eventData: EventFormData) => {
        setEventSubmitting(true);
        try {
            const response = await api.post("/events/create", eventData);
            console.log("✅ Event created successfully:", response.data);
            setShowEventForm(false);
            // Refresh recent events after successful creation
            await fetchRecentEvents();
            alert("Event created successfully!");
        } catch (err: unknown) {
            console.error("❌ Event creation failed:", err);
            let errorMessage = "Failed to create event. Please try again.";

            if (typeof err === "object" && err !== null) {
                const error = err as {
                    response?: { data?: { message?: string } };
                    message?: string;
                };
                if (error.response?.data?.message) {
                    errorMessage = error.response.data.message;
                } else if (error.message) {
                    errorMessage = error.message;
                }
            }

            alert(errorMessage);
        } finally {
            setEventSubmitting(false);
        }
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setCredentials((prev) => ({
            ...prev,
            [name]: value,
        }));
        // Clear error when user starts typing
        if (error) setError("");
    };

    // Show loading spinner while checking authentication
    if (loading) {
        return (
            <Layout showNavbar={false} showFooter={false}>
                <div className="flex justify-center items-center bg-gray-50 min-h-screen">
                    <div className="text-center">
                        <div className="mx-auto border-b-2 border-blue-600 rounded-full w-12 h-12 animate-spin">
                        </div>
                        <p className="mt-4 text-gray-600">
                            Loading admin panel...
                        </p>
                    </div>
                </div>
            </Layout>
        );
    }

    if (!authenticated) {
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
                                    htmlFor="email"
                                    className="block font-medium text-gray-700 text-sm"
                                >
                                    Email
                                </label>
                                <input
                                    type="email"
                                    id="email"
                                    name="email"
                                    value={credentials.email}
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
                                className={`w-full btn-primary ${
                                    loginLoading
                                        ? "opacity-50 cursor-not-allowed"
                                        : ""
                                }`}
                                disabled={loginLoading}
                            >
                                {loginLoading
                                    ? (
                                        <div className="flex justify-center items-center">
                                            <div className="mr-2 border-white border-b-2 rounded-full w-4 h-4 animate-spin">
                                            </div>
                                            Signing In...
                                        </div>
                                    )
                                    : (
                                        "Sign In"
                                    )}
                            </button>
                        </form>

                        <div className="mt-6 text-gray-500 text-sm text-center">
                            <p>Use your admin email and password</p>
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
                            onClick={handleLogout}
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

                    {/* Event Creation Section */}
                    <div className="gap-6 grid grid-cols-1 lg:grid-cols-2 mb-8">
                        <div className="card">
                            <div className="p-6">
                                <div className="flex justify-between items-center mb-4">
                                    <h3 className="font-bold text-gray-900 text-lg">
                                        {showEventForm
                                            ? "Create New Event"
                                            : "Quick Actions"}
                                    </h3>
                                    {showEventForm && (
                                        <button
                                            onClick={() =>
                                                setShowEventForm(false)}
                                            className="text-gray-500 hover:text-gray-700"
                                        >
                                            <svg
                                                className="w-5 h-5"
                                                fill="none"
                                                stroke="currentColor"
                                                viewBox="0 0 24 24"
                                            >
                                                <path
                                                    strokeLinecap="round"
                                                    strokeLinejoin="round"
                                                    strokeWidth={2}
                                                    d="M6 18L18 6M6 6l12 12"
                                                />
                                            </svg>
                                        </button>
                                    )}
                                </div>

                                {!showEventForm
                                    ? (
                                        <div className="space-y-3">
                                            <button
                                                onClick={() =>
                                                    setShowEventForm(true)}
                                                className="w-full text-left btn-primary"
                                            >
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
                                    )
                                    : (
                                        <EventCreationForm
                                            onSubmit={handleEventSubmit}
                                            loading={eventSubmitting}
                                            onCancel={() =>
                                                setShowEventForm(false)}
                                        />
                                    )}
                            </div>
                        </div>

                        <div className="card">
                            <div className="p-6">
                                <div className="flex justify-between items-center mb-4">
                                    <h3 className="font-bold text-gray-900 text-lg">
                                        Recent Events
                                    </h3>
                                    <button
                                        onClick={fetchRecentEvents}
                                        className="text-blue-600 hover:text-blue-800 text-sm"
                                        disabled={eventsLoading}
                                    >
                                        {eventsLoading
                                            ? "Loading..."
                                            : "Refresh"}
                                    </button>
                                </div>

                                {eventsLoading
                                    ? (
                                        <div className="flex justify-center items-center py-8">
                                            <div className="border-b-2 border-blue-600 rounded-full w-6 h-6 animate-spin">
                                            </div>
                                        </div>
                                    )
                                    : recentEvents.length > 0
                                    ? (
                                        <div className="space-y-3">
                                            {recentEvents.map((event) => (
                                                <div
                                                    key={event.id}
                                                    className="flex justify-between items-start py-2 border-gray-100 border-b last:border-b-0"
                                                >
                                                    <div className="flex-1">
                                                        <h4 className="font-medium text-gray-900 text-sm">
                                                            {event.title}
                                                        </h4>
                                                        <p className="mt-1 text-gray-600 text-xs">
                                                            {new Date(
                                                                event
                                                                    .start_time,
                                                            ).toLocaleDateString()}
                                                            {" "}
                                                            • {event.location}
                                                        </p>
                                                        <div className="flex items-center gap-2 mt-1">
                                                            <span
                                                                className={`inline-flex px-2 py-1 rounded-full text-xs font-semibold ${
                                                                    event
                                                                            .status ===
                                                                            "active"
                                                                        ? "bg-green-100 text-green-800"
                                                                        : event
                                                                                .status ===
                                                                                "draft"
                                                                        ? "bg-yellow-100 text-yellow-800"
                                                                        : "bg-gray-100 text-gray-800"
                                                                }`}
                                                            >
                                                                {event.status}
                                                            </span>
                                                            <span className="text-gray-500 text-xs">
                                                                {event
                                                                    .available_seats}/{event
                                                                    .total_seats}
                                                                {" "}
                                                                seats
                                                            </span>
                                                        </div>
                                                    </div>
                                                    <span className="text-gray-500 text-xs">
                                                        {new Date(
                                                            event.created_at,
                                                        ).toLocaleDateString()}
                                                    </span>
                                                </div>
                                            ))}
                                        </div>
                                    )
                                    : (
                                        <div className="py-8 text-gray-500 text-center">
                                            <p className="text-sm">
                                                No recent events found
                                            </p>
                                            <p className="mt-1 text-xs">
                                                Create your first event to get
                                                started!
                                            </p>
                                        </div>
                                    )}
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
