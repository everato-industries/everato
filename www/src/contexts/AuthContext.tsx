import React, {
    createContext,
    useCallback,
    useContext,
    useEffect,
    useState,
} from "react";
import api from "../lib/api";

// ─── Types ────────────────────────────────────────────────────────────────────

export interface User {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
    verified: boolean;
}

export interface RegisterData {
    firstName: string;
    lastName: string;
    email: string;
    password: string;
}

interface AuthContextType {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
    loading: boolean;
    login: (email: string, password: string) => Promise<void>;
    register: (data: RegisterData) => Promise<void>;
    logout: () => void;
}

// ─── Storage helpers ──────────────────────────────────────────────────────────

const TOKEN_KEY = "everato_user_token";
const USER_KEY = "everato_user_data";

function persistAuth(token: string, user: User) {
    localStorage.setItem(TOKEN_KEY, token);
    localStorage.setItem(USER_KEY, JSON.stringify(user));
}

function clearPersistedAuth() {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
}

function loadPersistedToken(): string | null {
    return localStorage.getItem(TOKEN_KEY);
}

function loadPersistedUser(): User | null {
    try {
        const raw = localStorage.getItem(USER_KEY);
        return raw ? (JSON.parse(raw) as User) : null;
    } catch {
        return null;
    }
}

// ─── Context ──────────────────────────────────────────────────────────────────

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(loadPersistedUser);
    const [token, setToken] = useState<string | null>(loadPersistedToken);
    const [loading, setLoading] = useState(false);

    // Keep axios default header in sync with token
    useEffect(() => {
        if (token) {
            api.defaults.headers.common["Authorization"] = `Bearer ${token}`;
        } else {
            delete api.defaults.headers.common["Authorization"];
        }
    }, [token]);

    const login = useCallback(async (email: string, password: string) => {
        setLoading(true);
        try {
            const response = await api.post("/auth/login", { email, password });
            const { token: newToken, user: rawUser } = response.data as {
                token: string;
                user: User;
            };

            persistAuth(newToken, rawUser);
            setToken(newToken);
            setUser(rawUser);
        } finally {
            setLoading(false);
        }
    }, []);

    const register = useCallback(async (data: RegisterData) => {
        setLoading(true);
        try {
            await api.post("/auth/register", {
                firstName: data.firstName,
                lastName: data.lastName,
                email: data.email,
                password: data.password,
            });
            // Registration succeeds — user must verify email then login
        } finally {
            setLoading(false);
        }
    }, []);

    const logout = useCallback(() => {
        clearPersistedAuth();
        setToken(null);
        setUser(null);
        delete api.defaults.headers.common["Authorization"];
    }, []);

    const value: AuthContextType = {
        user,
        token,
        isAuthenticated: !!(token && user),
        loading,
        login,
        register,
        logout,
    };

    return (
        <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
    );
}

// ─── Hook ─────────────────────────────────────────────────────────────────────

/**
 * useAuth — access the user authentication context.
 * Must be used inside <AuthProvider>.
 */
export function useAuth(): AuthContextType {
    const ctx = useContext(AuthContext);
    if (!ctx) {
        throw new Error("useAuth must be used within an <AuthProvider>");
    }
    return ctx;
}

export default AuthContext;
