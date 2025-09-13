import axios from "axios";
import { clearAuthData, getAuthToken, updateLastActivity } from "./auth";

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
            window.location.href = "/admin";
        } else if (error.response?.status === 403) {
            // Forbidden - user doesn't have permission
            console.error("Access denied: Insufficient permissions");
        }
        return Promise.reject(error);
    },
);

export default api;
