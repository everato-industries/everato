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

export default api;
