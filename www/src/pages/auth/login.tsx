import { useState } from "react";
import { Link } from "react-router-dom";

export default function LoginPage() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError("");

        try {
            const response = await fetch("/api/v1/auth/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ email, password }),
            });

            if (response.ok) {
                window.location.href = "/";
            } else {
                const data = await response.json();
                setError(data.message || "Login failed. Please try again.");
            }
        } catch (err) {
            setError("An error occurred. Please try again.");
            console.error(err);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="bg-primary min-h-screen flex items-center justify-center">
            <main className="w-full max-w-sm mx-auto">
                <div className="bg-secondary rounded-lg shadow-lg p-8">
                    <form onSubmit={handleSubmit} className="space-y-6">
                        <div className="mb-4 text-center">
                            <h1 className="text-3xl font-bold fg-secondary mb-1">
                                Login
                            </h1>
                            <p className="fg-primary/60 text-sm">
                                Please enter your email and password to login.
                            </p>
                        </div>

                        {error && (
                            <div className="p-3 bg-red-100 border border-red-400 text-red-700 rounded">
                                {error}
                            </div>
                        )}

                        <div>
                            <label
                                htmlFor="email"
                                className="block mb-1 font-medium"
                            >
                                Email
                            </label>
                            <input
                                type="email"
                                id="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none"
                            />
                        </div>
                        <div>
                            <label
                                htmlFor="password"
                                className="block mb-1 font-medium"
                            >
                                Password
                            </label>
                            <input
                                type="password"
                                id="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none"
                            />
                        </div>
                        <button
                            type="submit"
                            className="btn w-full py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
                            disabled={isLoading}
                        >
                            {isLoading ? "Logging in..." : "Login"}
                        </button>
                    </form>
                    <p className="mt-6 text-center text-sm fg-primary/70">
                        Don't have an account?
                        <Link
                            to="/auth/register"
                            className="text-accent hover:underline ml-1"
                        >
                            Register here
                        </Link>
                    </p>
                </div>
            </main>
        </div>
    );
}
