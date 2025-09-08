import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import Layout from "../../components/layout";

interface RegisterFormData {
    firstName: string;
    lastName: string;
    email: string;
    password: string;
    confirmPassword: string;
}

interface ApiError {
    message: string;
}

export default function RegisterPage() {
    const [formData, setFormData] = useState<RegisterFormData>({
        firstName: "",
        lastName: "",
        email: "",
        password: "",
        confirmPassword: "",
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState(false);
    const navigate = useNavigate();

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
        // Clear error when user starts typing
        if (error) setError(null);
    };

    const validateForm = (): boolean => {
        if (formData.password !== formData.confirmPassword) {
            setError("Passwords do not match");
            return false;
        }
        if (formData.password.length < 6) {
            setError("Password must be at least 6 characters long");
            return false;
        }
        return true;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        if (!validateForm()) {
            setLoading(false);
            return;
        }

        try {
            const response = await fetch("/api/v1/auth/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    firstName: formData.firstName,
                    lastName: formData.lastName,
                    email: formData.email,
                    password: formData.password,
                }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || "Registration failed");
            }

            setSuccess(true);
            setTimeout(() => {
                navigate("/auth/login");
            }, 2000);
        } catch (err) {
            const error = err as ApiError;
            setError(error.message || "An error occurred during registration");
        } finally {
            setLoading(false);
        }
    };

    if (success) {
        return (
            <Layout showNavbar={false} showFooter={false}>
                <div className="flex justify-center items-center bg-gray-50 min-h-screen">
                    <div className="w-full max-w-md text-center">
                        <div className="bg-green-50 px-6 py-4 border border-green-200 rounded text-green-700">
                            <h2 className="mb-2 font-semibold text-xl">
                                Registration Successful!
                            </h2>
                            <p>
                                Please check your email to verify your account.
                            </p>
                            <p className="mt-2 text-sm">
                                Redirecting to login page...
                            </p>
                        </div>
                    </div>
                </div>
            </Layout>
        );
    }

    return (
        <Layout showNavbar={false} showFooter={false}>
            <div className="flex justify-center items-center bg-gray-50 px-4 sm:px-6 lg:px-8 py-12 min-h-screen">
                <div className="space-y-8 w-full max-w-md">
                    {/* Header */}
                    <div className="text-center">
                        <Link to="/" className="font-bold text-black text-3xl">
                            Everato
                        </Link>
                        <h2 className="mt-6 font-bold text-black text-3xl">
                            Create your account
                        </h2>
                        <p className="mt-2 text-gray-600">
                            Or{" "}
                            <Link
                                to="/auth/login"
                                className="font-medium text-black hover:underline"
                            >
                                sign in to your existing account
                            </Link>
                        </p>
                    </div>

                    {/* Form */}
                    <form className="space-y-6 mt-8" onSubmit={handleSubmit}>
                        {error && (
                            <div className="bg-red-50 px-4 py-3 border border-red-200 text-red-700">
                                {error}
                            </div>
                        )}

                        <div className="space-y-4">
                            <div className="gap-4 grid grid-cols-2">
                                <div>
                                    <label
                                        htmlFor="firstName"
                                        className="sr-only"
                                    >
                                        First Name
                                    </label>
                                    <input
                                        id="firstName"
                                        name="firstName"
                                        type="text"
                                        required
                                        className="input-field"
                                        placeholder="First name"
                                        value={formData.firstName}
                                        onChange={handleInputChange}
                                    />
                                </div>
                                <div>
                                    <label
                                        htmlFor="lastName"
                                        className="sr-only"
                                    >
                                        Last Name
                                    </label>
                                    <input
                                        id="lastName"
                                        name="lastName"
                                        type="text"
                                        required
                                        className="input-field"
                                        placeholder="Last name"
                                        value={formData.lastName}
                                        onChange={handleInputChange}
                                    />
                                </div>
                            </div>

                            <div>
                                <label htmlFor="email" className="sr-only">
                                    Email address
                                </label>
                                <input
                                    id="email"
                                    name="email"
                                    type="email"
                                    autoComplete="email"
                                    required
                                    className="input-field"
                                    placeholder="Email address"
                                    value={formData.email}
                                    onChange={handleInputChange}
                                />
                            </div>

                            <div>
                                <label htmlFor="password" className="sr-only">
                                    Password
                                </label>
                                <input
                                    id="password"
                                    name="password"
                                    type="password"
                                    autoComplete="new-password"
                                    required
                                    className="input-field"
                                    placeholder="Password (min. 6 characters)"
                                    value={formData.password}
                                    onChange={handleInputChange}
                                />
                            </div>

                            <div>
                                <label
                                    htmlFor="confirmPassword"
                                    className="sr-only"
                                >
                                    Confirm Password
                                </label>
                                <input
                                    id="confirmPassword"
                                    name="confirmPassword"
                                    type="password"
                                    autoComplete="new-password"
                                    required
                                    className="input-field"
                                    placeholder="Confirm password"
                                    value={formData.confirmPassword}
                                    onChange={handleInputChange}
                                />
                            </div>
                        </div>

                        <div className="flex items-center">
                            <input
                                id="agree-terms"
                                name="agree-terms"
                                type="checkbox"
                                required
                                className="border-gray-300 focus:ring-0 w-4 h-4 text-black"
                            />
                            <label
                                htmlFor="agree-terms"
                                className="block ml-2 text-gray-600 text-sm"
                            >
                                I agree to the{" "}
                                <Link
                                    to="/terms"
                                    className="font-medium text-black hover:underline"
                                >
                                    Terms of Service
                                </Link>{" "}
                                and{" "}
                                <Link
                                    to="/privacy"
                                    className="font-medium text-black hover:underline"
                                >
                                    Privacy Policy
                                </Link>
                            </label>
                        </div>

                        <div>
                            <button
                                type="submit"
                                disabled={loading}
                                className={`btn-primary w-full ${
                                    loading
                                        ? "opacity-50 cursor-not-allowed"
                                        : ""
                                }`}
                            >
                                {loading
                                    ? "Creating account..."
                                    : "Create account"}
                            </button>
                        </div>
                    </form>

                    {/* Social Register Options */}
                    <div className="mt-6">
                        <div className="relative">
                            <div className="absolute inset-0 flex items-center">
                                <div className="border-gray-300 border-t w-full" />
                            </div>
                            <div className="relative flex justify-center text-sm">
                                <span className="bg-gray-50 px-2 text-gray-500">
                                    Or continue with
                                </span>
                            </div>
                        </div>

                        <div className="gap-3 grid grid-cols-2 mt-6">
                            <button
                                type="button"
                                className="flex justify-center items-center btn-secondary"
                            >
                                <svg
                                    className="mr-2 w-5 h-5"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        fill="currentColor"
                                        d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                                    />
                                    <path
                                        fill="currentColor"
                                        d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                                    />
                                    <path
                                        fill="currentColor"
                                        d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                                    />
                                    <path
                                        fill="currentColor"
                                        d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                                    />
                                </svg>
                                Google
                            </button>

                            <button
                                type="button"
                                className="flex justify-center items-center btn-secondary"
                            >
                                <svg
                                    className="mr-2 w-5 h-5"
                                    fill="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        fillRule="evenodd"
                                        d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
                                        clipRule="evenodd"
                                    />
                                </svg>
                                GitHub
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </Layout>
    );
}
