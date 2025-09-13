import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import Layout from "../components/layout";
import api from "../lib/api";
import { getAdminUser, isAuthenticated } from "../lib/auth";

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
    banner_file?: File;
    icon_file?: File;
    banner_url?: string;
    icon_url?: string;
    total_seats: number;
    available_seats: number;
    status?: string;
    ticket_types: TicketType[];
    coupons: Coupon[];
    // Additional fields for comprehensive event creation
    event_type: string;
    category: string;
    tags: string;
    max_tickets_per_user: number;
    booking_start_time: string;
    booking_end_time: string;
    refund_policy: string;
    terms_and_conditions: string;
    contact_email: string;
    contact_phone: string;
    social_links: {
        website?: string;
        facebook?: string;
        twitter?: string;
        instagram?: string;
        linkedin?: string;
    };
    venue_details: {
        venue_name: string;
        address_line1: string;
        address_line2?: string;
        city: string;
        state: string;
        postal_code: string;
        country: string;
        latitude?: number;
        longitude?: number;
    };
    organizer_info: {
        organizer_name: string;
        organizer_email: string;
        organizer_phone: string;
        organization: string;
    };
}

export default function CreateEventPage() {
    const navigate = useNavigate();
    const [authenticated, setAuthenticated] = useState(false);
    const [loading, setLoading] = useState(true);
    const [submitLoading, setSubmitLoading] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const [formData, setFormData] = useState<EventFormData>({
        title: "",
        description: "",
        start_time: "",
        end_time: "",
        location: "",
        admin_id: "",
        total_seats: 100,
        available_seats: 100,
        status: "CREATED",
        ticket_types: [{
            name: "General Admission",
            price: 50,
            available_tickets: 50,
        }],
        coupons: [],
        event_type: "CONFERENCE",
        category: "TECHNOLOGY",
        tags: "",
        max_tickets_per_user: 10,
        booking_start_time: "",
        booking_end_time: "",
        refund_policy: "",
        terms_and_conditions: "",
        contact_email: "",
        contact_phone: "",
        social_links: {},
        venue_details: {
            venue_name: "",
            address_line1: "",
            city: "",
            state: "",
            postal_code: "",
            country: "USA",
        },
        organizer_info: {
            organizer_name: "",
            organizer_email: "",
            organizer_phone: "",
            organization: "",
        },
    });

    // Check authentication status on component mount
    useEffect(() => {
        const checkAuth = async () => {
            try {
                if (!isAuthenticated()) {
                    navigate("/admin");
                    return;
                }

                const adminUser = getAdminUser();
                if (!adminUser) {
                    navigate("/admin");
                    return;
                }

                setFormData((prev) => ({
                    ...prev,
                    admin_id: adminUser.id,
                    contact_email: adminUser.email,
                    organizer_info: {
                        ...prev.organizer_info,
                        organizer_email: adminUser.email,
                        organizer_name: adminUser.username || "",
                    },
                }));
                setAuthenticated(true);
            } catch (error) {
                console.error("Auth check failed:", error);
                navigate("/admin");
            } finally {
                setLoading(false);
            }
        };

        checkAuth();
    }, [navigate]);

    const updateField = (
        field: keyof EventFormData,
        value: string | number | File | TicketType[] | Coupon[] | any,
    ) => {
        setFormData((prev) => ({ ...prev, [field]: value }));
    };

    const updateNestedField = (
        parentField: keyof EventFormData,
        childField: string,
        value: string | number,
    ) => {
        setFormData((prev) => ({
            ...prev,
            [parentField]: {
                ...(prev[parentField] as any),
                [childField]: value,
            },
        }));
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

    const handleFileChange = (
        field: "banner_file" | "icon_file",
        file: File | null,
    ) => {
        if (file) {
            setFormData((prev) => ({
                ...prev,
                [field]: file,
            }));
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitLoading(true);
        setError("");
        setSuccess("");

        try {
            // Convert datetime-local to ISO format
            const submitData = {
                ...formData,
                start_time: new Date(formData.start_time).toISOString(),
                end_time: new Date(formData.end_time).toISOString(),
                booking_start_time: formData.booking_start_time
                    ? new Date(formData.booking_start_time).toISOString()
                    : new Date().toISOString(),
                booking_end_time: formData.booking_end_time
                    ? new Date(formData.booking_end_time).toISOString()
                    : new Date(formData.start_time).toISOString(),
                tags: formData.tags.split(",").map((tag) => tag.trim()).filter(
                    (tag) => tag,
                ),
            };

            // Remove file fields from API submission (would need separate upload endpoint)
            const { banner_file, icon_file, ...apiData } = submitData;

            // For now, use placeholder URLs if files are selected
            if (banner_file) {
                apiData.banner_url = "https://placeholder.com/banner.jpg";
            }
            if (icon_file) {
                apiData.icon_url = "https://placeholder.com/icon.jpg";
            }

            const response = await api.post("/events/create", apiData);

            if (response.data) {
                setSuccess("Event created successfully!");
                // Reset form or redirect
                setTimeout(() => {
                    navigate("/admin");
                }, 2000);
            }
        } catch (error: any) {
            console.error("Event creation failed:", error);
            const errorMessage = error.response?.data?.message ||
                error.response?.data?.error ||
                "Failed to create event. Please try again.";
            setError(errorMessage);
        } finally {
            setSubmitLoading(false);
        }
    };

    if (loading) {
        return (
            <Layout>
                <div className="flex justify-center items-center min-h-screen">
                    <div className="border-b-2 border-blue-600 rounded-full w-8 h-8 animate-spin">
                    </div>
                </div>
            </Layout>
        );
    }

    if (!authenticated) {
        return null;
    }

    return (
        <Layout>
            <div className="bg-gray-50 min-h-screen">
                <div className="mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-6xl">
                    <div className="mb-8">
                        <div className="flex justify-between items-center">
                            <h1 className="font-bold text-gray-900 text-3xl">
                                Create New Event
                            </h1>
                            <button
                                onClick={() => navigate("/admin")}
                                className="btn-secondary"
                            >
                                Back to Dashboard
                            </button>
                        </div>
                        <p className="mt-2 text-gray-600">
                            Create a comprehensive event with all details,
                            ticket types, and promotional coupons.
                        </p>
                    </div>

                    {error && (
                        <div className="bg-red-50 mb-4 p-4 border border-red-200 rounded text-red-700">
                            {error}
                        </div>
                    )}

                    {success && (
                        <div className="bg-green-50 mb-4 p-4 border border-green-200 rounded text-green-700">
                            {success}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-8">
                        {/* Basic Event Information */}
                        <div className="card">
                            <h2 className="mb-4 font-semibold text-gray-900 text-xl">
                                Basic Information
                            </h2>

                            <div className="gap-6 grid grid-cols-1 md:grid-cols-2">
                                <div className="md:col-span-2">
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Event Title *
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.title}
                                        onChange={(e) =>
                                            updateField(
                                                "title",
                                                e.target.value,
                                            )}
                                        className="w-full input-field"
                                        required
                                    />
                                </div>

                                <div className="md:col-span-2">
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Description *
                                    </label>
                                    <textarea
                                        value={formData.description}
                                        onChange={(e) => updateField(
                                            "description",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        rows={4}
                                        required
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Event Type *
                                    </label>
                                    <select
                                        value={formData.event_type}
                                        onChange={(e) =>
                                            updateField(
                                                "event_type",
                                                e.target.value,
                                            )}
                                        className="w-full input-field"
                                        required
                                    >
                                        <option value="CONFERENCE">
                                            Conference
                                        </option>
                                        <option value="WORKSHOP">
                                            Workshop
                                        </option>
                                        <option value="SEMINAR">Seminar</option>
                                        <option value="MEETUP">Meetup</option>
                                        <option value="FESTIVAL">
                                            Festival
                                        </option>
                                        <option value="CONCERT">Concert</option>
                                        <option value="EXHIBITION">
                                            Exhibition
                                        </option>
                                        <option value="OTHER">Other</option>
                                    </select>
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Category *
                                    </label>
                                    <select
                                        value={formData.category}
                                        onChange={(e) =>
                                            updateField(
                                                "category",
                                                e.target.value,
                                            )}
                                        className="w-full input-field"
                                        required
                                    >
                                        <option value="TECHNOLOGY">
                                            Technology
                                        </option>
                                        <option value="BUSINESS">
                                            Business
                                        </option>
                                        <option value="EDUCATION">
                                            Education
                                        </option>
                                        <option value="HEALTH">Health</option>
                                        <option value="ARTS">
                                            Arts & Culture
                                        </option>
                                        <option value="SPORTS">Sports</option>
                                        <option value="ENTERTAINMENT">
                                            Entertainment
                                        </option>
                                        <option value="NETWORKING">
                                            Networking
                                        </option>
                                        <option value="OTHER">Other</option>
                                    </select>
                                </div>

                                <div className="md:col-span-2">
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Tags (comma-separated)
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.tags}
                                        onChange={(e) =>
                                            updateField("tags", e.target.value)}
                                        placeholder="e.g., AI, Machine Learning, Tech Conference"
                                        className="w-full input-field"
                                    />
                                </div>
                            </div>
                        </div>

                        {/* Date & Time Information */}
                        <div className="card">
                            <h2 className="mb-4 font-semibold text-gray-900 text-xl">
                                Date & Time
                            </h2>

                            <div className="gap-6 grid grid-cols-1 md:grid-cols-2">
                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Start Time *
                                    </label>
                                    <input
                                        type="datetime-local"
                                        value={formData.start_time}
                                        onChange={(e) => updateField(
                                            "start_time",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        required
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        End Time *
                                    </label>
                                    <input
                                        type="datetime-local"
                                        value={formData.end_time}
                                        onChange={(e) =>
                                            updateField(
                                                "end_time",
                                                e.target.value,
                                            )}
                                        className="w-full input-field"
                                        required
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Booking Opens
                                    </label>
                                    <input
                                        type="datetime-local"
                                        value={formData.booking_start_time}
                                        onChange={(e) => updateField(
                                            "booking_start_time",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Booking Closes
                                    </label>
                                    <input
                                        type="datetime-local"
                                        value={formData.booking_end_time}
                                        onChange={(e) => updateField(
                                            "booking_end_time",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                    />
                                </div>
                            </div>
                        </div>

                        {/* Media & Branding */}
                        <div className="card">
                            <h2 className="mb-4 font-semibold text-gray-900 text-xl">
                                Media & Branding
                            </h2>

                            <div className="gap-6 grid grid-cols-1 md:grid-cols-2">
                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Event Banner
                                    </label>
                                    <input
                                        type="file"
                                        accept="image/*"
                                        onChange={(e) => handleFileChange(
                                            "banner_file",
                                            e.target.files?.[0] || null,
                                        )}
                                        className="w-full input-field"
                                    />
                                    <p className="mt-1 text-gray-500 text-sm">
                                        Recommended size: 1200x600px (JPG, PNG)
                                    </p>
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Event Icon/Logo
                                    </label>
                                    <input
                                        type="file"
                                        accept="image/*"
                                        onChange={(e) => handleFileChange(
                                            "icon_file",
                                            e.target.files?.[0] || null,
                                        )}
                                        className="w-full input-field"
                                    />
                                    <p className="mt-1 text-gray-500 text-sm">
                                        Recommended size: 300x300px (JPG, PNG)
                                    </p>
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Banner URL (alternative)
                                    </label>
                                    <input
                                        type="url"
                                        value={formData.banner_url || ""}
                                        onChange={(e) => updateField(
                                            "banner_url",
                                            e.target.value,
                                        )}
                                        placeholder="https://example.com/banner.jpg"
                                        className="w-full input-field"
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Icon URL (alternative)
                                    </label>
                                    <input
                                        type="url"
                                        value={formData.icon_url || ""}
                                        onChange={(e) =>
                                            updateField(
                                                "icon_url",
                                                e.target.value,
                                            )}
                                        placeholder="https://example.com/icon.jpg"
                                        className="w-full input-field"
                                    />
                                </div>
                            </div>
                        </div>

                        {/* Venue Information */}
                        <div className="card">
                            <h2 className="mb-4 font-semibold text-gray-900 text-xl">
                                Venue Information
                            </h2>

                            <div className="gap-6 grid grid-cols-1 md:grid-cols-2">
                                <div className="md:col-span-2">
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Venue Name *
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.venue_details
                                            .venue_name}
                                        onChange={(e) => updateNestedField(
                                            "venue_details",
                                            "venue_name",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        required
                                    />
                                </div>

                                <div className="md:col-span-2">
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Address Line 1 *
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.venue_details
                                            .address_line1}
                                        onChange={(e) => updateNestedField(
                                            "venue_details",
                                            "address_line1",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        required
                                    />
                                </div>

                                <div className="md:col-span-2">
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Address Line 2
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.venue_details
                                            .address_line2 || ""}
                                        onChange={(e) => updateNestedField(
                                            "venue_details",
                                            "address_line2",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        City *
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.venue_details.city}
                                        onChange={(e) => updateNestedField(
                                            "venue_details",
                                            "city",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        required
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        State/Province *
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.venue_details.state}
                                        onChange={(e) => updateNestedField(
                                            "venue_details",
                                            "state",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        required
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Postal Code *
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.venue_details
                                            .postal_code}
                                        onChange={(e) => updateNestedField(
                                            "venue_details",
                                            "postal_code",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        required
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Country *
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.venue_details.country}
                                        onChange={(e) => updateNestedField(
                                            "venue_details",
                                            "country",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        required
                                    />
                                </div>

                                <div className="md:col-span-2">
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Full Location Description
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.location}
                                        onChange={(e) =>
                                            updateField(
                                                "location",
                                                e.target.value,
                                            )}
                                        placeholder="e.g., Grand Ballroom, Silicon Valley Convention Center"
                                        className="w-full input-field"
                                    />
                                </div>
                            </div>
                        </div>

                        {/* Capacity & Booking Settings */}
                        <div className="card">
                            <h2 className="mb-4 font-semibold text-gray-900 text-xl">
                                Capacity & Booking Settings
                            </h2>

                            <div className="gap-6 grid grid-cols-1 md:grid-cols-3">
                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Total Seats *
                                    </label>
                                    <input
                                        type="number"
                                        value={formData.total_seats}
                                        onChange={(e) => updateField(
                                            "total_seats",
                                            parseInt(e.target.value),
                                        )}
                                        className="w-full input-field"
                                        min="1"
                                        required
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Available Seats *
                                    </label>
                                    <input
                                        type="number"
                                        value={formData.available_seats}
                                        onChange={(e) => updateField(
                                            "available_seats",
                                            parseInt(e.target.value),
                                        )}
                                        className="w-full input-field"
                                        min="1"
                                        required
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Max Tickets Per User
                                    </label>
                                    <input
                                        type="number"
                                        value={formData.max_tickets_per_user}
                                        onChange={(e) => updateField(
                                            "max_tickets_per_user",
                                            parseInt(e.target.value),
                                        )}
                                        className="w-full input-field"
                                        min="1"
                                    />
                                </div>
                            </div>
                        </div>

                        {/* Ticket Types */}
                        <div className="card">
                            <div className="flex justify-between items-center mb-4">
                                <h2 className="font-semibold text-gray-900 text-xl">
                                    Ticket Types
                                </h2>
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
                                    className="bg-gray-50 mb-4 p-4 rounded-lg"
                                >
                                    {formData.ticket_types.length > 1 && (
                                        <div className="flex justify-end mb-3">
                                            <button
                                                type="button"
                                                onClick={() =>
                                                    removeTicketType(index)}
                                                className="text-red-600 hover:text-red-800 text-sm"
                                            >
                                                Remove
                                            </button>
                                        </div>
                                    )}
                                    <div className="gap-4 grid grid-cols-1 md:grid-cols-3">
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
                                                className="w-full input-field"
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
                                                        parseFloat(
                                                            e.target.value,
                                                        ),
                                                    )}
                                                className="w-full input-field"
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
                                                        parseInt(
                                                            e.target.value,
                                                        ),
                                                    )}
                                                className="w-full input-field"
                                                min="1"
                                                required
                                            />
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>

                        {/* Coupons */}
                        <div className="card">
                            <div className="flex justify-between items-center mb-4">
                                <h2 className="font-semibold text-gray-900 text-xl">
                                    Promotional Coupons
                                </h2>
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
                                    className="bg-gray-50 mb-4 p-4 rounded-lg"
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

                                    <div className="gap-4 grid grid-cols-1 md:grid-cols-2 mb-3">
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
                                                        e.target.value
                                                            .toUpperCase(),
                                                    )}
                                                className="w-full input-field"
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
                                                value={coupon
                                                    .discount_percentage}
                                                onChange={(e) =>
                                                    updateCoupon(
                                                        index,
                                                        "discount_percentage",
                                                        parseFloat(
                                                            e.target.value,
                                                        ),
                                                    )}
                                                className="w-full input-field"
                                                min="1"
                                                max="100"
                                                step="0.1"
                                                required
                                            />
                                        </div>
                                    </div>

                                    <div className="gap-4 grid grid-cols-1 md:grid-cols-3">
                                        <div>
                                            <label className="block mb-1 font-medium text-gray-600 text-xs">
                                                Valid From
                                            </label>
                                            <input
                                                type="datetime-local"
                                                value={coupon.valid_from}
                                                onChange={(e) => updateCoupon(
                                                    index,
                                                    "valid_from",
                                                    e.target.value,
                                                )}
                                                className="w-full input-field"
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
                                                onChange={(e) => updateCoupon(
                                                    index,
                                                    "valid_until",
                                                    e.target.value,
                                                )}
                                                className="w-full input-field"
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
                                                onChange={(e) => updateCoupon(
                                                    index,
                                                    "usage_limit",
                                                    parseInt(e.target.value),
                                                )}
                                                className="w-full input-field"
                                                min="1"
                                                required
                                            />
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>

                        {/* Organizer Information */}
                        <div className="card">
                            <h2 className="mb-4 font-semibold text-gray-900 text-xl">
                                Organizer Information
                            </h2>

                            <div className="gap-6 grid grid-cols-1 md:grid-cols-2">
                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Organizer Name *
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.organizer_info
                                            .organizer_name}
                                        onChange={(e) => updateNestedField(
                                            "organizer_info",
                                            "organizer_name",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        required
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Organization
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.organizer_info
                                            .organization}
                                        onChange={(e) => updateNestedField(
                                            "organizer_info",
                                            "organization",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Organizer Email *
                                    </label>
                                    <input
                                        type="email"
                                        value={formData.organizer_info
                                            .organizer_email}
                                        onChange={(e) => updateNestedField(
                                            "organizer_info",
                                            "organizer_email",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        required
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Organizer Phone
                                    </label>
                                    <input
                                        type="tel"
                                        value={formData.organizer_info
                                            .organizer_phone}
                                        onChange={(e) => updateNestedField(
                                            "organizer_info",
                                            "organizer_phone",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                    />
                                </div>
                            </div>
                        </div>

                        {/* Contact & Social Links */}
                        <div className="card">
                            <h2 className="mb-4 font-semibold text-gray-900 text-xl">
                                Contact & Social Links
                            </h2>

                            <div className="gap-6 grid grid-cols-1 md:grid-cols-2">
                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Contact Email
                                    </label>
                                    <input
                                        type="email"
                                        value={formData.contact_email}
                                        onChange={(e) => updateField(
                                            "contact_email",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Contact Phone
                                    </label>
                                    <input
                                        type="tel"
                                        value={formData.contact_phone}
                                        onChange={(e) => updateField(
                                            "contact_phone",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Website
                                    </label>
                                    <input
                                        type="url"
                                        value={formData.social_links.website ||
                                            ""}
                                        onChange={(e) => updateNestedField(
                                            "social_links",
                                            "website",
                                            e.target.value,
                                        )}
                                        placeholder="https://example.com"
                                        className="w-full input-field"
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Facebook
                                    </label>
                                    <input
                                        type="url"
                                        value={formData.social_links.facebook ||
                                            ""}
                                        onChange={(e) => updateNestedField(
                                            "social_links",
                                            "facebook",
                                            e.target.value,
                                        )}
                                        placeholder="https://facebook.com/page"
                                        className="w-full input-field"
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Twitter
                                    </label>
                                    <input
                                        type="url"
                                        value={formData.social_links.twitter ||
                                            ""}
                                        onChange={(e) => updateNestedField(
                                            "social_links",
                                            "twitter",
                                            e.target.value,
                                        )}
                                        placeholder="https://twitter.com/handle"
                                        className="w-full input-field"
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        LinkedIn
                                    </label>
                                    <input
                                        type="url"
                                        value={formData.social_links.linkedin ||
                                            ""}
                                        onChange={(e) => updateNestedField(
                                            "social_links",
                                            "linkedin",
                                            e.target.value,
                                        )}
                                        placeholder="https://linkedin.com/company/name"
                                        className="w-full input-field"
                                    />
                                </div>
                            </div>
                        </div>

                        {/* Policies & Terms */}
                        <div className="card">
                            <h2 className="mb-4 font-semibold text-gray-900 text-xl">
                                Policies & Terms
                            </h2>

                            <div className="space-y-4">
                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Refund Policy
                                    </label>
                                    <textarea
                                        value={formData.refund_policy}
                                        onChange={(e) => updateField(
                                            "refund_policy",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        rows={3}
                                        placeholder="Describe your refund policy..."
                                    />
                                </div>

                                <div>
                                    <label className="block mb-2 font-medium text-gray-700 text-sm">
                                        Terms and Conditions
                                    </label>
                                    <textarea
                                        value={formData.terms_and_conditions}
                                        onChange={(e) => updateField(
                                            "terms_and_conditions",
                                            e.target.value,
                                        )}
                                        className="w-full input-field"
                                        rows={4}
                                        placeholder="Enter terms and conditions..."
                                    />
                                </div>
                            </div>
                        </div>

                        {/* Form Actions */}
                        <div className="flex justify-end gap-4 pt-6">
                            <button
                                type="button"
                                onClick={() => navigate("/admin")}
                                className="px-6 py-2 btn-secondary"
                                disabled={submitLoading}
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                disabled={submitLoading}
                                className={`px-8 py-2 btn-primary ${
                                    submitLoading
                                        ? "opacity-50 cursor-not-allowed"
                                        : ""
                                }`}
                            >
                                {submitLoading
                                    ? (
                                        <div className="flex justify-center items-center">
                                            <div className="mr-2 border-white border-b-2 rounded-full w-4 h-4 animate-spin">
                                            </div>
                                            Creating Event...
                                        </div>
                                    )
                                    : (
                                        "Create Event"
                                    )}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </Layout>
    );
}
