// Authentication utility functions for admin panel

export interface AdminUser {
    id: string;
    name: string;
    email: string;
    username: string;
    role: string;
    permissions: string[];
}

export interface AuthTokens {
    token: string;
    expiresAt?: string;
}

export interface AuthResponse {
    token: string;
    user: AdminUser;
    message: string;
}

// Storage keys
const STORAGE_KEYS = {
    TOKEN: "everato_admin_token",
    USER: "everato_admin_user",
    SESSION: "everato_admin_session",
} as const;

// Cookie settings
const COOKIE_OPTIONS = {
    secure: import.meta.env.PROD,
    sameSite: "strict" as const,
    maxAge: 7 * 24 * 60 * 60 * 1000, // 7 days
} as const;

/**
 * Set cookie with proper options
 */
export function setCookie(
    name: string,
    value: string,
    options: Record<string, unknown> = {},
): void {
    const opts = { ...COOKIE_OPTIONS, ...options };
    let cookieString = `${name}=${encodeURIComponent(value)}`;

    if (opts.maxAge) {
        const expires = new Date(Date.now() + opts.maxAge);
        cookieString += `; expires=${expires.toUTCString()}`;
    }

    cookieString += `; path=/`;

    if (opts.secure) {
        cookieString += `; secure`;
    }

    if (opts.sameSite) {
        cookieString += `; samesite=${opts.sameSite}`;
    }

    document.cookie = cookieString;
}

/**
 * Get cookie value by name
 */
export function getCookie(name: string): string | null {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) {
        const cookieValue = parts.pop()?.split(";").shift();
        return cookieValue ? decodeURIComponent(cookieValue) : null;
    }
    return null;
}

/**
 * Remove cookie
 */
export function removeCookie(name: string): void {
    document.cookie =
        `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
}

/**
 * Save admin authentication data to localStorage and cookies
 */
export function saveAuthData(authResponse: AuthResponse): void {
    try {
        // Save to localStorage
        localStorage.setItem(STORAGE_KEYS.TOKEN, authResponse.token);
        localStorage.setItem(
            STORAGE_KEYS.USER,
            JSON.stringify(authResponse.user),
        );

        // Save to cookies (for server-side access)
        setCookie(STORAGE_KEYS.TOKEN, authResponse.token);
        setCookie(STORAGE_KEYS.USER, JSON.stringify(authResponse.user));

        // Save session timestamp
        const sessionData = {
            loginTime: new Date().toISOString(),
            lastActivity: new Date().toISOString(),
        };
        localStorage.setItem(STORAGE_KEYS.SESSION, JSON.stringify(sessionData));
        setCookie(STORAGE_KEYS.SESSION, JSON.stringify(sessionData));

        console.log("✅ Admin authentication data saved successfully");
    } catch (error) {
        console.error("❌ Error saving auth data:", error);
        throw new Error("Failed to save authentication data");
    }
}

/**
 * Retrieve admin authentication token
 */
export function getAuthToken(): string | null {
    // Try localStorage first, then cookies
    return localStorage.getItem(STORAGE_KEYS.TOKEN) ||
        getCookie(STORAGE_KEYS.TOKEN);
}

/**
 * Retrieve admin user data
 */
export function getAdminUser(): AdminUser | null {
    try {
        // Try localStorage first, then cookies
        const userStr = localStorage.getItem(STORAGE_KEYS.USER) ||
            getCookie(STORAGE_KEYS.USER);
        return userStr ? JSON.parse(userStr) : null;
    } catch (error) {
        console.error("❌ Error parsing admin user data:", error);
        return null;
    }
}

/**
 * Check if admin is authenticated
 */
export function isAuthenticated(): boolean {
    const token = getAuthToken();
    const user = getAdminUser();
    return !!(token && user);
}

/**
 * Clear all authentication data
 */
export function clearAuthData(): void {
    try {
        // Clear localStorage
        localStorage.removeItem(STORAGE_KEYS.TOKEN);
        localStorage.removeItem(STORAGE_KEYS.USER);
        localStorage.removeItem(STORAGE_KEYS.SESSION);

        // Clear cookies
        removeCookie(STORAGE_KEYS.TOKEN);
        removeCookie(STORAGE_KEYS.USER);
        removeCookie(STORAGE_KEYS.SESSION);

        console.log("✅ Admin authentication data cleared");
    } catch (error) {
        console.error("❌ Error clearing auth data:", error);
    }
}

/**
 * Update last activity timestamp
 */
export function updateLastActivity(): void {
    try {
        const existingSession = localStorage.getItem(STORAGE_KEYS.SESSION);
        if (existingSession) {
            const sessionData = JSON.parse(existingSession);
            sessionData.lastActivity = new Date().toISOString();

            localStorage.setItem(
                STORAGE_KEYS.SESSION,
                JSON.stringify(sessionData),
            );
            setCookie(STORAGE_KEYS.SESSION, JSON.stringify(sessionData));
        }
    } catch (error) {
        console.error("❌ Error updating last activity:", error);
    }
}

/**
 * Check if user has specific permission
 */
export function hasPermission(permission: string): boolean {
    const user = getAdminUser();
    if (!user) return false;

    // Super admin has all permissions
    if (user.role === "SUPER_ADMIN") return true;

    return user.permissions?.includes(permission) || false;
}

/**
 * Get all admin permissions
 */
export function getPermissions(): string[] {
    const user = getAdminUser();
    return user?.permissions || [];
}

/**
 * Format admin display name
 */
export function getDisplayName(): string {
    const user = getAdminUser();
    if (!user) return "Unknown Admin";

    return user.name || user.username || user.email || "Admin User";
}
