import axios from "axios";
import { clearAuthData, getAuthToken, updateLastActivity } from "./auth";

// Supporting interfaces
export interface TicketType {
    id: string;
    name: string;
    price: number;
    quantity: number;
    description?: string;
}

export interface Coupon {
    id: string;
    code: string;
    discount_percentage?: number;
    discount_amount?: number;
    description?: string;
}

// Event interface for API responses (matches backend response structure)
export interface Event {
    id: string;
    title: string;
    description: string;
    banner?: string;
    icon?: string;
    start_time: string;
    end_time: string;
    location: string;
    status: string;
    slug: string;
    total_seats: number;
    available_seats: number;
    created_at: string;
    updated_at: string;
    admin_id?: string;
    tags: string[];
    // Venue fields
    venue_name?: string;
    address_line1?: string;
    address_line2?: string;
    city?: string;
    state?: string;
    postal_code?: string;
    country?: string;
    // Additional fields from backend
    organizer_name?: string;
    organizer_email?: string;
    event_type?: string;
    category?: string;
    ticket_types?: TicketType[];
    coupons?: Coupon[];
}

// Create axios instance with default configuration
const api = axios.create({
    baseURL: "http://localhost:8080/api/v1",
    headers: {
        "Content-Type": "application/json",
    },
    withCredentials: true, // Include cookies in requests
});

// Request interceptor to add auth token and update activity
api.interceptors.request.use(
    (config) => {
        // Get token from auth utility (checks both localStorage and cookies)
        const token = getAuthToken();
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }

        // Update last activity timestamp
        updateLastActivity();

        return config;
    },
    (error) => {
        return Promise.reject(error);
    },
);

// Response interceptor for error handling
api.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        // Handle common errors
        if (error.response?.status === 401) {
            // Unauthorized - clear auth data and redirect to admin login
            clearAuthData();
            // Use React Router navigation instead of direct location change
            // The component using this should handle the redirect properly
            console.warn("Unauthorized access - authentication required");
        } else if (error.response?.status === 403) {
            // Forbidden - user doesn't have permission
            console.error("Access denied: Insufficient permissions");
        }
        return Promise.reject(error);
    },
);

// Event API functions
export const eventAPI = {
    // Get recent events (for home page)
    getRecentEvents: (limit: number = 6) => {
        return api.get(`/events/recent?limit=${limit}`);
    },

    // Get all events with pagination (for events page)
    getAllEvents: (limit: number = 6, offset: number = 0) => {
        return api.get(`/events/all?limit=${limit}&offset=${offset}`);
    },

    // Get all events with advanced filtering, pagination, and sorting (for admin)
    getAllEventsWithFilters: (params: {
        limit?: number;
        offset?: number;
        search?: string;
        sortBy?: "title" | "created_at" | "start_time";
        sortOrder?: "asc" | "desc";
    }) => {
        const queryParams = new URLSearchParams();

        if (params.limit) queryParams.append("limit", params.limit.toString());
        if (params.offset) {
            queryParams.append("offset", params.offset.toString());
        }
        if (params.search) queryParams.append("search", params.search);
        if (params.sortBy) queryParams.append("sortBy", params.sortBy);
        if (params.sortOrder) queryParams.append("sortOrder", params.sortOrder);

        return api.get(`/events/all?${queryParams.toString()}`);
    },

    // Create a new event
    createEvent: (
        eventData: Omit<Event, "id" | "created_at" | "updated_at" | "admin_id">,
    ) => {
        return api.post("/events/create", eventData);
    },

    // Get event by slug
    getEvent: (slug: string) => {
        return api.get(`/events/${slug}`);
    },

    // Update an event
    updateEvent: (slug: string, eventData: Partial<Event>) => {
        return api.put(`/events/${slug}`, eventData);
    },

    // Delete an event
    deleteEvent: (slug: string) => {
        return api.delete(`/events/${slug}`);
    },

    // Start an event (change status to published)
    startEvent: (slug: string) => {
        return api.post(`/events/${slug}/start`);
    },

    // End an event (change status to ended)
    endEvent: (slug: string) => {
        return api.post(`/events/${slug}/end`);
    },
};

// Admin API functions
export const adminAPI = {
    // Get current admin info by ID
    getAdminById: (adminId: string) => {
        return api.get(`/admin/${adminId}`);
    },

    // Get admin by username
    getAdminByUsername: (username: string) => {
        return api.get(`/admin/u/${username}`);
    },

    // Get all permissions
    getAllPermissions: () => {
        return api.get("/admin/permissions");
    },

    // Get all roles
    getAllRoles: () => {
        return api.get("/admin/roles");
    },
};

// ─── Booking API ──────────────────────────────────────────────────────────────

export interface CreateBookingRequest {
    event_id: string;
    tickets: Array<{
        ticket_type_id: string;
        quantity: number;
    }>;
    coupon_code?: string;
}

export const bookingAPI = {
    /** Create a new booking (auth required) */
    createBooking: (data: CreateBookingRequest) =>
        api.post("/bookings/create", data),

    /** Get all bookings for a specific user */
    getUserBookings: (userId: string) =>
        api.get(`/bookings/user/${userId}`),

    /** Get a single booking's details */
    getBookingDetails: (bookingId: string) =>
        api.get(`/bookings/${bookingId}`),

    /** Cancel (soft-delete) a booking */
    cancelBooking: (bookingId: string) =>
        api.delete(`/bookings/${bookingId}`),

    /** Check real-time ticket availability for an event */
    checkAvailability: (eventId: string) =>
        api.get(`/events/${eventId}/availability`),
};

// ─── User API ─────────────────────────────────────────────────────────────────

export const userAPI = {
    /** Register a new user account */
    register: (data: {
        firstName: string;
        lastName: string;
        email: string;
        password: string;
    }) => api.post("/auth/register", data),

    /** Login and receive a JWT token */
    login: (email: string, password: string) =>
        api.post("/auth/login", { email, password }),

    /** Refresh the JWT access token */
    refresh: () => api.post("/auth/refresh"),
};

// ─── Ticket API ───────────────────────────────────────────────────────────────

export interface Ticket {
    id: string;
    booking_id: string;
    ticket_type_id?: string;
    ticket_type_name?: string;
    ticket_number: string;
    qr_code_data: string;
    qr_code_image?: string; // base64 PNG
    is_checked_in: boolean;
    checked_in_at?: string;
    created_at: string;
}

export const ticketAPI = {
    /** Get all tickets and QR codes for a booking */
    getTicketsByBooking: (bookingId: string) =>
        api.get(`/tickets/booking/${bookingId}`),

    /** Get a single ticket by its UUID */
    getTicketById: (ticketId: string) =>
        api.get(`/tickets/${ticketId}`),
};

// ─── Check-in / Attendance API ────────────────────────────────────────────────

export interface CheckinScanRequest {
    qr_data: string;
    location?: string;
    device_info?: string;
}

export interface AttendanceStats {
    total_tickets: number;
    checked_in_count: number;
    not_checked_in_count: number;
}

export interface EventAttendee {
    ticket_id: string;
    ticket_number: string;
    is_checked_in: boolean;
    checked_in_at?: string;
    user_email: string;
    first_name: string;
    last_name: string;
    ticket_type_name: string;
}

export const checkinAPI = {
    /** Scan and validate QR code for check-in */
    scan: (data: CheckinScanRequest) =>
        api.post("/checkin/scan", data),

    /** Get check-in statistics for an event */
    getStats: (eventId: string) =>
        api.get(`/checkin/events/${eventId}/stats`),

    /** Get full attendee list for an event */
    getAttendees: (eventId: string) =>
        api.get(`/checkin/events/${eventId}/attendees`),
};

// ─── Payment API ──────────────────────────────────────────────────────────────

export interface PaymentClientConfig {
    provider: "razorpay" | "stripe" | "cashfree" | "none";
    enabled: boolean | string;
    key_id?: string;
    publishable_key?: string;
    app_id?: string;
    env?: string;
    currency?: string;
    message?: string;
}

export interface PaymentOrderResponse {
    order_id: string;
    provider: string;
    amount: number;
    currency: string;
    key_id?: string;
    client_secret?: string;
    session_id?: string;
    metadata?: Record<string, string>;
}

export interface VerifyPaymentRequest {
    booking_id: string;
    order_id: string;
    payment_id: string;
    signature?: string;
}

export const paymentAPI = {
    /** Get active payment provider configuration and client public keys */
    getConfig: () => api.get("/payments/config"),

    /** Create payment order with active provider for a booking */
    createOrder: (bookingId: string) =>
        api.post("/payments/create-order", { booking_id: bookingId }),

    /** Verify payment signature or completion */
    verifyPayment: (data: VerifyPaymentRequest) =>
        api.post("/payments/verify", data),
};

export default api;
