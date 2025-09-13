import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { type Event, eventAPI } from "../lib/api";
import Layout from "../components/layout";

interface PublicStats {
    totalEvents: number;
    upcomingEvents: number;
    totalActiveEvents: number;
}

export default function DashboardPage() {
    const [stats, setStats] = useState<PublicStats>({
        totalEvents: 0,
        upcomingEvents: 0,
        totalActiveEvents: 0,
    });
    const [recentEvents, setRecentEvents] = useState<Event[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchPublicData();
    }, []);

    const fetchPublicData = async () => {
        try {
            setLoading(true);

            const recentResponse = await eventAPI.getRecentEvents(6);
            if (
                recentResponse.data && recentResponse.data.data &&
                recentResponse.data.data.events
            ) {
                const events = recentResponse.data.data.events;
                setRecentEvents(events);

                const now = new Date();
                const upcomingCount = events.filter((event: Event) =>
                    new Date(event.start_time) > now
                ).length;

                const activeCount = events.filter((event: Event) =>
                    event.status === "active" || event.status === "published"
                ).length;

                setStats({
                    totalEvents: events.length,
                    upcomingEvents: upcomingCount,
                    totalActiveEvents: activeCount,
                });
            }
        } catch (error) {
            console.error("Error fetching public data:", error);
        } finally {
            setLoading(false);
        }
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    const getStatusColor = (status: string) => {
        switch (status.toLowerCase()) {
            case "active":
            case "published":
                return "bg-green-100 text-green-800";
            case "draft":
                return "bg-yellow-100 text-yellow-800";
            case "cancelled":
                return "bg-red-100 text-red-800";
            default:
                return "bg-gray-100 text-gray-800";
        }
    };

    const isUpcoming = (startTime: string) => {
        return new Date(startTime) > new Date();
    };

    return (
        <Layout>
            <div className="mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl">
                <div className="mb-8">
                    <h1 className="mb-2 font-bold text-black text-3xl">
                        Events Dashboard
                    </h1>
                    <p className="text-gray-600">
                        Discover what's happening in your area and beyond
                    </p>
                </div>

                <div className="gap-6 grid grid-cols-1 md:grid-cols-3 mb-8">
                    <div className="card">
                        <h3 className="mb-2 font-medium text-gray-500 text-sm uppercase tracking-wide">
                            Recent Events
                        </h3>
                        <p className="font-bold text-black text-3xl">
                            {loading ? "..." : stats.totalEvents}
                        </p>
                    </div>
                    <div className="card">
                        <h3 className="mb-2 font-medium text-gray-500 text-sm uppercase tracking-wide">
                            Upcoming Events
                        </h3>
                        <p className="font-bold text-black text-3xl">
                            {loading ? "..." : stats.upcomingEvents}
                        </p>
                    </div>
                    <div className="card">
                        <h3 className="mb-2 font-medium text-gray-500 text-sm uppercase tracking-wide">
                            Active Events
                        </h3>
                        <p className="font-bold text-black text-3xl">
                            {loading ? "..." : stats.totalActiveEvents}
                        </p>
                    </div>
                </div>

                <div className="card">
                    <div className="flex justify-between items-center px-6 py-4 border-b">
                        <h2 className="font-semibold text-black text-xl">
                            Recent Events
                        </h2>
                        <Link
                            to="/events"
                            className="text-black text-sm hover:underline"
                        >
                            View all events →
                        </Link>
                    </div>

                    {loading
                        ? (
                            <div className="p-6">
                                <div className="space-y-4">
                                    {[1, 2, 3].map((i) => (
                                        <div
                                            key={i}
                                            className="bg-gray-200 rounded h-16 animate-pulse"
                                        />
                                    ))}
                                </div>
                            </div>
                        )
                        : recentEvents.length === 0
                        ? (
                            <div className="p-6 text-center">
                                <p className="text-gray-500">No events found</p>
                            </div>
                        )
                        : (
                            <div className="divide-y divide-gray-300">
                                {recentEvents.map((event) => (
                                    <div
                                        key={event.id}
                                        className="hover:bg-gray-50 px-6 py-4 transition-colors duration-200"
                                    >
                                        <div className="flex justify-between items-center">
                                            <div className="flex-1">
                                                <div className="flex items-center space-x-3">
                                                    <Link
                                                        to={`/events/${event.slug}`}
                                                        className="font-semibold text-black hover:text-gray-700"
                                                    >
                                                        {event.title}
                                                    </Link>
                                                    <span
                                                        className={`px-2 py-1 text-xs font-medium uppercase tracking-wide ${
                                                            getStatusColor(
                                                                event.status,
                                                            )
                                                        }`}
                                                    >
                                                        {event.status}
                                                    </span>
                                                    {isUpcoming(
                                                        event.start_time,
                                                    ) && (
                                                        <span className="bg-blue-100 px-2 py-1 font-medium text-blue-800 text-xs uppercase tracking-wide">
                                                            Upcoming
                                                        </span>
                                                    )}
                                                </div>
                                                <div className="flex items-center space-x-4 mt-1 text-gray-500 text-sm">
                                                    <span>
                                                        📅 {formatDate(
                                                            event.start_time,
                                                        )}
                                                    </span>
                                                    <span>
                                                        📍 {event.location}
                                                    </span>
                                                    <span>
                                                        🎫{" "}
                                                        {event.available_seats}
                                                        {" "}
                                                        / {event.total_seats}
                                                        {" "}
                                                        seats
                                                    </span>
                                                </div>
                                            </div>
                                            <div className="ml-4">
                                                <Link
                                                    to={`/events/${event.slug}`}
                                                    className="text-sm btn-primary"
                                                >
                                                    View Details
                                                </Link>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                </div>

                <div className="gap-6 grid grid-cols-1 md:grid-cols-2 mt-8">
                    <div className="card">
                        <h3 className="mb-3 font-semibold text-black text-lg">
                            Explore Events
                        </h3>
                        <p className="mb-4 text-gray-600">
                            Browse through all available events and find
                            something you'll love
                        </p>
                        <Link to="/events" className="btn-primary">
                            Browse All Events
                        </Link>
                    </div>

                    <div className="bg-gray-50 card">
                        <h3 className="mb-3 font-semibold text-black text-lg">
                            Event Organizer?
                        </h3>
                        <p className="mb-4 text-gray-600">
                            Create and manage your own events with our
                            easy-to-use platform
                        </p>
                        <Link to="/auth/login" className="btn-secondary">
                            Get Started
                        </Link>
                    </div>
                </div>
            </div>
        </Layout>
    );
}
