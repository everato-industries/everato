import axios from "axios";
import { clearAuthData, getAuthToken, updateLastActivity } from "./auth";

// Event interface for API responses
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

    // Create a new event
    createEvent: (
        eventData: Omit<Event, "id" | "created_at" | "updated_at" | "admin_id">,
    ) => {
        return api.post("/events/create", eventData);
    },
};

export default api;
