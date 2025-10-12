import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getAdminUser, isAuthenticated } from "../lib/auth";
import { adminAPI, type Event as BackendEvent, eventAPI } from "../lib/api";

interface TicketType {
    id: string;
    name: string;
    price: number;
    quantity: number;
    description?: string;
}

interface EventFormData {
    name: string;
    description: string;
    start_time: string;
    end_time: string;
    venue: {
        name: string;
        address: string;
        city: string;
        state: string;
        country: string;
    };
    capacity: number;
    is_public: boolean;
    status: "draft" | "published" | "cancelled" | "ended";
    tags: string;
    ticket_types: TicketType[];
}

// Use the backend event interface
type Event = BackendEvent;

const EditEvent: React.FC = () => {
    const { slug } = useParams<{ slug: string }>();
    const navigate = useNavigate();

    const [loading, setLoading] = useState(true);
    const [submitLoading, setSubmitLoading] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");
    const [event, setEvent] = useState<Event | null>(null);
    const [hasManagePermission, setHasManagePermission] = useState(false);

    const [formData, setFormData] = useState<EventFormData>({
        name: "",
        description: "",
        start_time: "",
        end_time: "",
        venue: {
            name: "",
            address: "",
            city: "",
            state: "",
            country: "",
        },
        capacity: 0,
        is_public: true,
        status: "draft",
        tags: "",
        ticket_types: [],
    });

    // Check authentication and permissions
    useEffect(() => {
        const checkAuth = async () => {
            const authenticated = await isAuthenticated();
            if (!authenticated) {
                navigate("/admin/login");
                return;
            }

            // Get admin user data from localStorage
            const adminUser = getAdminUser();
            if (!adminUser) {
                setError("Admin user data not found");
                setLoading(false);
                return;
            }

            // Check permissions from stored admin data
            try {
                // Get current admin's detailed info to verify permissions
                const response = await adminAPI.getAdminById(adminUser.id);

                if (response.status === 200) {
                    const adminData = response.data.admin;
                    const hasManage = adminData.role === "SUPER_ADMIN" ||
                        adminData.permissions.includes("MANAGE_EVENTS") ||
                        adminData.permissions.includes("EDIT_EVENT");

                    setHasManagePermission(hasManage);

                    if (!hasManage) {
                        setError("You don't have permission to edit events");
                        setLoading(false);
                        return;
                    }
                }
            } catch (err) {
                console.error("Permission check failed:", err);
                setError("Failed to verify permissions");
                setLoading(false);
                return;
            }
        };

        checkAuth();
    }, [navigate]);

    // Load event data
    useEffect(() => {
        if (!hasManagePermission || !slug) return;

        const loadEventData = async () => {
            try {
                if (!slug) {
                    throw new Error("Event slug is missing");
                }
                const response = await eventAPI.getEvent(slug);

                if (response.status === 200) {
                    // Backend returns { message: "...", event: {...} }
                    const eventData: Event = response.data.event;
                    setEvent(eventData);

                    // Transform event data to form data
                    setFormData({
                        name: eventData.title, // Backend uses 'title'
                        description: eventData.description,
                        start_time: new Date(eventData.start_time).toISOString()
                            .slice(0, 16),
                        end_time: new Date(eventData.end_time).toISOString()
                            .slice(0, 16),
                        venue: {
                            name: eventData.venue_name || "",
                            address: eventData.address_line1 || "",
                            city: eventData.city || "",
                            state: eventData.state || "",
                            country: eventData.country || "",
                        },
                        capacity: eventData.total_seats,
                        is_public: true, // Default since backend doesn't have this field
                        status: eventData.status as
                            | "draft"
                            | "published"
                            | "cancelled"
                            | "ended",
                        tags: Array.isArray(eventData.tags)
                            ? eventData.tags.join(", ")
                            : String(eventData.tags || ""),
                        ticket_types: eventData.ticket_types || [],
                    });
                }
            } catch (err) {
                console.error("Failed to load event:", err);
                const message = err instanceof Error
                    ? err.message
                    : "Failed to load event";
                setError(message);
            } finally {
                setLoading(false);
            }
        };

        loadEventData();
    }, [slug, hasManagePermission]);

    const handleInputChange = (
        e: React.ChangeEvent<
            HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
        >,
    ) => {
        const { name, value, type } = e.target;

        if (name.includes(".")) {
            const [section, field] = name.split(".");
            setFormData((prev) => ({
                ...prev,
                [section]: {
                    ...(prev[section as keyof EventFormData] as Record<
                        string,
                        unknown
                    >),
                    [field]: type === "number" ? Number(value) : value,
                },
            }));
        } else {
            setFormData((prev) => ({
                ...prev,
                [name]: type === "number" ? Number(value) : value,
            }));
        }
    };

    const handleUpdateEvent = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!hasManagePermission) {
            setError("You don't have permission to edit events");
            return;
        }

        setSubmitLoading(true);
        setError("");
        setSuccess("");

        try {
            const payload = {
                ...formData,
                tags: formData.tags.split(",").map((tag: string) => tag.trim())
                    .filter((tag: string) => tag),
            };

            // Update event using axios API
            if (!slug) {
                throw new Error("Event slug is missing");
            }
            const response = await eventAPI.updateEvent(slug, payload);

            if (response.status === 200) {
                setSuccess("Event updated successfully!");
                setTimeout(() => {
                    navigate("/admin");
                }, 2000);
            }
        } catch (err) {
            console.error("Failed to update event:", err);
            const message = err instanceof Error
                ? err.message
                : "Failed to update event";
            setError(message);
        } finally {
            setSubmitLoading(false);
        }
    };

    const handleStartEvent = async () => {
        if (!event || !hasManagePermission) return;

        try {
            if (!slug) {
                throw new Error("Event slug is missing");
            }
            const response = await eventAPI.startEvent(slug);
            if (response.status === 200) {
                setSuccess("Event started successfully!");
                setEvent({ ...event, status: "published" });
            }
        } catch (err) {
            console.error("Failed to start event:", err);
            setError("Failed to start event");
        }
    };

    const handleEndEvent = async () => {
        if (!event || !hasManagePermission) return;

        try {
            if (!slug) {
                throw new Error("Event slug is missing");
            }
            const response = await eventAPI.endEvent(slug);
            if (response.status === 200) {
                setSuccess("Event ended successfully!");
                setEvent({ ...event, status: "ended" });
            }
        } catch (err) {
            console.error("Failed to end event:", err);
            setError("Failed to end event");
        }
    };
    const handleDeleteEvent = async () => {
        if (!event || !hasManagePermission) return;

        const confirmed = window.confirm(
            "Are you sure you want to delete this event? This action cannot be undone.",
        );
        if (!confirmed) return;

        try {
            if (!slug) {
                throw new Error("Event slug is missing");
            }
            const response = await eventAPI.deleteEvent(slug);
            if (response.status === 200) {
                setSuccess("Event deleted successfully!");
                setTimeout(() => {
                    navigate("/admin");
                }, 2000);
            }
        } catch (err) {
            console.error("Failed to delete event:", err);
            setError("Failed to delete event");
        }
    };

    const addTicketType = () => {
        setFormData((prev) => ({
            ...prev,
            ticket_types: [
                ...prev.ticket_types,
                {
                    id: `temp_${Date.now()}`,
                    name: "",
                    price: 0,
                    quantity: 0,
                    description: "",
                },
            ],
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

    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="text-xl">Loading event...</div>
            </div>
        );
    }

    if (error && !event) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="text-center">
                    <div className="mb-4 text-red-500 text-xl">{error}</div>
                    <button
                        onClick={() => navigate("/admin")}
                        className="btn-secondary"
                    >
                        Back to Admin
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="bg-white p-8 min-h-screen">
            <div className="mx-auto max-w-4xl">
                {/* Header */}
                <div className="flex justify-between items-center mb-8">
                    <div>
                        <h1 className="font-bold text-black text-3xl">
                            Edit Event
                        </h1>
                        {event && (
                            <p className="mt-2 text-gray-600">
                                Event: {event.title} | Status:
                                <span
                                    className={`ml-1 px-2 py-1 rounded text-xs font-medium ${
                                        event.status === "published"
                                            ? "bg-green-100 text-green-800"
                                            : event.status === "draft"
                                            ? "bg-yellow-100 text-yellow-800"
                                            : event.status === "ended"
                                            ? "bg-gray-100 text-gray-800"
                                            : "bg-red-100 text-red-800"
                                    }`}
                                >
                                    {event.status.toUpperCase()}
                                </span>
                            </p>
                        )}
                    </div>
                    <button
                        onClick={() => navigate("/admin")}
                        className="btn-secondary"
                    >
                        Back to Admin
                    </button>
                </div>

                {/* Action Buttons */}
                {event && hasManagePermission && (
                    <div className="flex gap-4 bg-gray-50 mb-8 p-4 rounded">
                        <h3 className="mr-4 font-semibold text-black text-lg">
                            Quick Actions:
                        </h3>

                        {event.status === "draft" && (
                            <button
                                onClick={handleStartEvent}
                                className="btn-primary"
                            >
                                Start Event
                            </button>
                        )}

                        {(event.status === "published" ||
                            event.status === "draft") && (
                            <button
                                onClick={handleEndEvent}
                                className="btn-secondary"
                            >
                                End Event
                            </button>
                        )}

                        <button
                            onClick={handleDeleteEvent}
                            className="bg-red-600 hover:bg-red-700 px-4 py-2 rounded text-white transition-colors"
                        >
                            Delete Event
                        </button>
                    </div>
                )}

                {/* Messages */}
                {error && (
                    <div className="bg-red-100 mb-6 p-4 border border-red-400 rounded text-red-700">
                        {error}
                    </div>
                )}

                {success && (
                    <div className="bg-green-100 mb-6 p-4 border border-green-400 rounded text-green-700">
                        {success}
                    </div>
                )}

                {/* Edit Form */}
                <form onSubmit={handleUpdateEvent} className="space-y-8">
                    {/* Basic Information */}
                    <div className="card">
                        <h2 className="mb-6 font-semibold text-black text-xl">
                            Basic Information
                        </h2>

                        <div className="gap-6 grid grid-cols-1 md:grid-cols-2">
                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Event Name *
                                </label>
                                <input
                                    type="text"
                                    name="name"
                                    value={formData.name}
                                    onChange={handleInputChange}
                                    required
                                    className="input-field"
                                    placeholder="Enter event name"
                                />
                            </div>

                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Capacity *
                                </label>
                                <input
                                    type="number"
                                    name="capacity"
                                    value={formData.capacity}
                                    onChange={handleInputChange}
                                    required
                                    min="1"
                                    className="input-field"
                                    placeholder="Maximum attendees"
                                />
                            </div>
                        </div>

                        <div className="mt-6">
                            <label className="block mb-2 font-medium text-black text-sm">
                                Description *
                            </label>
                            <textarea
                                name="description"
                                value={formData.description}
                                onChange={handleInputChange}
                                required
                                rows={4}
                                className="resize-none input-field"
                                placeholder="Describe your event..."
                            />
                        </div>

                        <div className="gap-6 grid grid-cols-1 md:grid-cols-2 mt-6">
                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Start Date & Time *
                                </label>
                                <input
                                    type="datetime-local"
                                    name="start_time"
                                    value={formData.start_time}
                                    onChange={handleInputChange}
                                    required
                                    className="input-field"
                                />
                            </div>

                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    End Date & Time *
                                </label>
                                <input
                                    type="datetime-local"
                                    name="end_time"
                                    value={formData.end_time}
                                    onChange={handleInputChange}
                                    required
                                    className="input-field"
                                />
                            </div>
                        </div>

                        <div className="gap-6 grid grid-cols-1 md:grid-cols-2 mt-6">
                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Status
                                </label>
                                <select
                                    name="status"
                                    value={formData.status}
                                    onChange={handleInputChange}
                                    className="input-field"
                                >
                                    <option value="draft">Draft</option>
                                    <option value="published">Published</option>
                                    <option value="cancelled">Cancelled</option>
                                    <option value="ended">Ended</option>
                                </select>
                            </div>

                            <div className="flex items-center pt-8">
                                <label className="flex items-center">
                                    <input
                                        type="checkbox"
                                        name="is_public"
                                        checked={formData.is_public}
                                        onChange={(e) =>
                                            setFormData((prev) => ({
                                                ...prev,
                                                is_public: e.target.checked,
                                            }))}
                                        className="mr-2"
                                    />
                                    <span className="font-medium text-black text-sm">
                                        Public Event
                                    </span>
                                </label>
                            </div>
                        </div>

                        <div className="mt-6">
                            <label className="block mb-2 font-medium text-black text-sm">
                                Tags (comma separated)
                            </label>
                            <input
                                type="text"
                                name="tags"
                                value={formData.tags}
                                onChange={handleInputChange}
                                className="input-field"
                                placeholder="conference, tech, networking"
                            />
                        </div>
                    </div>

                    {/* Venue Information */}
                    <div className="card">
                        <h2 className="mb-6 font-semibold text-black text-xl">
                            Venue Information
                        </h2>

                        <div className="gap-6 grid grid-cols-1 md:grid-cols-2">
                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Venue Name *
                                </label>
                                <input
                                    type="text"
                                    name="venue.name"
                                    value={formData.venue.name}
                                    onChange={handleInputChange}
                                    required
                                    className="input-field"
                                    placeholder="Enter venue name"
                                />
                            </div>

                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Address *
                                </label>
                                <input
                                    type="text"
                                    name="venue.address"
                                    value={formData.venue.address}
                                    onChange={handleInputChange}
                                    required
                                    className="input-field"
                                    placeholder="Street address"
                                />
                            </div>
                        </div>

                        <div className="gap-6 grid grid-cols-1 md:grid-cols-3 mt-6">
                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    City *
                                </label>
                                <input
                                    type="text"
                                    name="venue.city"
                                    value={formData.venue.city}
                                    onChange={handleInputChange}
                                    required
                                    className="input-field"
                                    placeholder="City"
                                />
                            </div>

                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    State *
                                </label>
                                <input
                                    type="text"
                                    name="venue.state"
                                    value={formData.venue.state}
                                    onChange={handleInputChange}
                                    required
                                    className="input-field"
                                    placeholder="State/Province"
                                />
                            </div>

                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Country *
                                </label>
                                <input
                                    type="text"
                                    name="venue.country"
                                    value={formData.venue.country}
                                    onChange={handleInputChange}
                                    required
                                    className="input-field"
                                    placeholder="Country"
                                />
                            </div>
                        </div>
                    </div>

                    {/* Ticket Types */}
                    <div className="card">
                        <div className="flex justify-between items-center mb-6">
                            <h2 className="font-semibold text-black text-xl">
                                Ticket Types
                            </h2>
                            <button
                                type="button"
                                onClick={addTicketType}
                                className="btn-primary"
                            >
                                Add Ticket Type
                            </button>
                        </div>

                        {formData.ticket_types.length === 0
                            ? (
                                <p className="py-8 text-gray-500 text-center">
                                    No ticket types added yet. Click "Add Ticket
                                    Type" to get started.
                                </p>
                            )
                            : (
                                <div className="space-y-4">
                                    {formData.ticket_types.map((
                                        ticket,
                                        index,
                                    ) => (
                                        <div
                                            key={ticket.id}
                                            className="p-4 border border-gray-200 rounded"
                                        >
                                            <div className="gap-4 grid grid-cols-1 md:grid-cols-4">
                                                <div>
                                                    <label className="block mb-2 font-medium text-black text-sm">
                                                        Ticket Name *
                                                    </label>
                                                    <input
                                                        type="text"
                                                        value={ticket.name}
                                                        onChange={(e) =>
                                                            updateTicketType(
                                                                index,
                                                                "name",
                                                                e.target.value,
                                                            )}
                                                        required
                                                        className="input-field"
                                                        placeholder="e.g., General Admission"
                                                    />
                                                </div>

                                                <div>
                                                    <label className="block mb-2 font-medium text-black text-sm">
                                                        Price ($) *
                                                    </label>
                                                    <input
                                                        type="number"
                                                        step="0.01"
                                                        min="0"
                                                        value={ticket.price}
                                                        onChange={(e) =>
                                                            updateTicketType(
                                                                index,
                                                                "price",
                                                                Number(
                                                                    e.target
                                                                        .value,
                                                                ),
                                                            )}
                                                        required
                                                        className="input-field"
                                                        placeholder="0.00"
                                                    />
                                                </div>

                                                <div>
                                                    <label className="block mb-2 font-medium text-black text-sm">
                                                        Quantity *
                                                    </label>
                                                    <input
                                                        type="number"
                                                        min="1"
                                                        value={ticket.quantity}
                                                        onChange={(e) =>
                                                            updateTicketType(
                                                                index,
                                                                "quantity",
                                                                Number(
                                                                    e.target
                                                                        .value,
                                                                ),
                                                            )}
                                                        required
                                                        className="input-field"
                                                        placeholder="100"
                                                    />
                                                </div>

                                                <div className="flex items-end">
                                                    <button
                                                        type="button"
                                                        onClick={() =>
                                                            removeTicketType(
                                                                index,
                                                            )}
                                                        className="bg-red-600 hover:bg-red-700 px-4 py-2 rounded w-full text-white transition-colors"
                                                    >
                                                        Remove
                                                    </button>
                                                </div>
                                            </div>

                                            <div className="mt-4">
                                                <label className="block mb-2 font-medium text-black text-sm">
                                                    Description (Optional)
                                                </label>
                                                <textarea
                                                    value={ticket.description ||
                                                        ""}
                                                    onChange={(e) =>
                                                        updateTicketType(
                                                            index,
                                                            "description",
                                                            e.target.value,
                                                        )}
                                                    rows={2}
                                                    className="resize-none input-field"
                                                    placeholder="Additional details about this ticket type..."
                                                />
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                    </div>

                    {/* Submit Button */}
                    <div className="flex justify-end space-x-4">
                        <button
                            type="button"
                            onClick={() => navigate("/admin")}
                            className="btn-secondary"
                            disabled={submitLoading}
                        >
                            Cancel
                        </button>

                        <button
                            type="submit"
                            disabled={submitLoading || !hasManagePermission}
                            className="disabled:opacity-50 disabled:cursor-not-allowed btn-primary"
                        >
                            {submitLoading ? "Updating..." : "Update Event"}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default EditEvent;
