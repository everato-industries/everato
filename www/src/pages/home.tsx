import { useEffect } from "react";
import Layout from "../components/Layout";

export default function HomePage() {
    // This ensures hydration issues are avoided
    useEffect(() => {
        async function fetch_health() {
            try {
                const response = await fetch("/api/v1/health");
                if (!response.ok) {
                    console.error("Health check failed");
                }
                const data = await response.json();
                console.log("Health check from API: ", data);
            } catch (error) {
                console.error("Error fetching health:", error);
            }
        }
        fetch_health(); // Call the async function
    }, []);

    return (
        <Layout>
            <div className="container mx-auto p-6">
                <h2 className="text-blue-400 text-3xl font-bold mb-4">
                    Welcome to Everato!
                </h2>
                <p className="text-gray-700 mb-8">
                    Your modern event management platform.
                    <br />
                    Explore events, manage tickets, and more.
                </p>

                {/* Featured events could go here */}
                <div className="bg-gray-50 p-6 rounded-lg shadow-sm">
                    <h3 className="text-xl font-semibold mb-3">
                        Featured Events
                    </h3>
                    <p className="text-gray-600">
                        Discover upcoming events in your area.
                    </p>

                    {/* Event cards would be mapped here */}
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-4">
                        {/* This would be populated from an API in a real implementation */}
                        <div className="bg-white p-4 rounded shadow">
                            Coming soon!
                        </div>
                    </div>
                </div>
            </div>
        </Layout>
    );
}
