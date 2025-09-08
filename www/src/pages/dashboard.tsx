import { useEffect, useState } from "react";
import Layout from "../components/layout";

interface DashboardStats {
    totalEvents: number;
    totalTicketsSold: number;
    totalRevenue: number;
    upcomingEvents: number;
}

interface RecentEvent {
    id: string;
    title: string;
    date: string;
    ticketsSold: number;
    revenue: number;
    status: "upcoming" | "ongoing" | "completed";
}

export default function DashboardPage() {
    const [stats, setStats] = useState<DashboardStats>({
        totalEvents: 0,
        totalTicketsSold: 0,
        totalRevenue: 0,
        upcomingEvents: 0,
    });
    const [recentEvents, setRecentEvents] = useState<RecentEvent[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchDashboardData();
    }, []);

    const fetchDashboardData = async () => {
        try {
            setLoading(true);
            // Mock API call - replace with actual API
            setTimeout(() => {
                setStats({
                    totalEvents: 12,
                    totalTicketsSold: 3456,
                    totalRevenue: 52340,
                    upcomingEvents: 5,
                });

                setRecentEvents([
                    {
                        id: "1",
                        title: "Tech Conference 2025",
                        date: "2025-10-15",
                        ticketsSold: 245,
                        revenue: 12500,
                        status: "upcoming",
                    },
                    {
                        id: "2",
                        title: "Music Festival",
                        date: "2025-08-20",
                        ticketsSold: 156,
                        revenue: 8900,
                        status: "upcoming",
                    },
                    {
                        id: "3",
                        title: "Food & Wine Expo",
                        date: "2025-07-10",
                        ticketsSold: 89,
                        revenue: 4560,
                        status: "completed",
                    },
                ]);
                setLoading(false);
            }, 500);
        } catch (error) {
            console.error("Error fetching dashboard data:", error);
            setLoading(false);
        }
    };

    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: "USD",
        }).format(amount);
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
        });
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case "upcoming":
                return "bg-blue-100 text-blue-800";
            case "ongoing":
                return "bg-green-100 text-green-800";
            case "completed":
                return "bg-gray-100 text-gray-800";
            default:
                return "bg-gray-100 text-gray-800";
        }
    };

    return (
        <Layout>
            <div className="mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl">
                {/* Header */}
                <div className="mb-8">
                    <h1 className="font-bold text-black text-3xl">Dashboard</h1>
                    <p className="mt-2 text-gray-600">
                        Welcome back! Here's an overview of your events.
                    </p>
                </div>

                {/* Stats Cards */}
                {loading
                    ? (
                        <div className="gap-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 mb-8">
                            {[1, 2, 3, 4].map((i) => (
                                <div key={i} className="card loading">
                                    <div className="bg-gray-200 mb-2 h-4"></div>
                                    <div className="bg-gray-200 h-6"></div>
                                </div>
                            ))}
                        </div>
                    )
                    : (
                        <div className="gap-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 mb-8">
                            <div className="card">
                                <div className="mb-2 font-medium text-gray-600 text-sm">
                                    Total Events
                                </div>
                                <div className="font-bold text-black text-3xl">
                                    {stats.totalEvents}
                                </div>
                            </div>
                            <div className="card">
                                <div className="mb-2 font-medium text-gray-600 text-sm">
                                    Tickets Sold
                                </div>
                                <div className="font-bold text-black text-3xl">
                                    {stats.totalTicketsSold.toLocaleString()}
                                </div>
                            </div>
                            <div className="card">
                                <div className="mb-2 font-medium text-gray-600 text-sm">
                                    Total Revenue
                                </div>
                                <div className="font-bold text-black text-3xl">
                                    {formatCurrency(stats.totalRevenue)}
                                </div>
                            </div>
                            <div className="card">
                                <div className="mb-2 font-medium text-gray-600 text-sm">
                                    Upcoming Events
                                </div>
                                <div className="font-bold text-black text-3xl">
                                    {stats.upcomingEvents}
                                </div>
                            </div>
                        </div>
                    )}

                <div className="gap-8 grid grid-cols-1 lg:grid-cols-3">
                    {/* Recent Events */}
                    <div className="lg:col-span-2">
                        <div className="card">
                            <div className="flex justify-between items-center mb-6">
                                <h2 className="font-bold text-black text-xl">
                                    Recent Events
                                </h2>
                                <button className="text-sm btn-secondary">
                                    View All
                                </button>
                            </div>

                            {loading
                                ? (
                                    <div className="space-y-4">
                                        {[1, 2, 3].map((i) => (
                                            <div key={i} className="loading">
                                                <div className="bg-gray-200 mb-2 h-4">
                                                </div>
                                                <div className="bg-gray-200 w-3/4 h-4">
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                )
                                : (
                                    <div className="space-y-4">
                                        {recentEvents.map((event) => (
                                            <div
                                                key={event.id}
                                                className="flex justify-between items-center hover:shadow-sm p-4 border border-gray-200 transition-shadow duration-200"
                                            >
                                                <div className="flex-1">
                                                    <h3 className="mb-1 font-semibold text-black">
                                                        {event.title}
                                                    </h3>
                                                    <p className="text-gray-600 text-sm">
                                                        {formatDate(event.date)}
                                                    </p>
                                                </div>
                                                <div className="mr-4 text-right">
                                                    <div className="font-medium text-black text-sm">
                                                        {event.ticketsSold}{" "}
                                                        tickets
                                                    </div>
                                                    <div className="text-gray-600 text-sm">
                                                        {formatCurrency(
                                                            event.revenue,
                                                        )}
                                                    </div>
                                                </div>
                                                <div>
                                                    <span
                                                        className={`px-2 py-1 text-xs font-medium uppercase tracking-wide ${
                                                            getStatusColor(
                                                                event.status,
                                                            )
                                                        }`}
                                                    >
                                                        {event.status}
                                                    </span>
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                )}
                        </div>
                    </div>

                    {/* Quick Actions */}
                    <div className="lg:col-span-1">
                        <div className="card">
                            <h2 className="mb-6 font-bold text-black text-xl">
                                Quick Actions
                            </h2>
                            <div className="space-y-3">
                                <button className="w-full btn-primary">
                                    Create New Event
                                </button>
                                <button className="w-full btn-secondary">
                                    View All Events
                                </button>
                                <button className="w-full btn-secondary">
                                    Export Data
                                </button>
                                <button className="w-full btn-secondary">
                                    Settings
                                </button>
                            </div>
                        </div>

                        {/* Recent Activity */}
                        <div className="mt-6 card">
                            <h2 className="mb-6 font-bold text-black text-xl">
                                Recent Activity
                            </h2>
                            {loading
                                ? (
                                    <div className="space-y-3">
                                        {[1, 2, 3].map((i) => (
                                            <div key={i} className="loading">
                                                <div className="bg-gray-200 h-3">
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                )
                                : (
                                    <div className="space-y-3">
                                        <div className="text-gray-600 text-sm">
                                            <span className="font-medium">
                                                John Doe
                                            </span>{" "}
                                            purchased 2 tickets for Tech
                                            Conference 2025
                                        </div>
                                        <div className="text-gray-600 text-sm">
                                            <span className="font-medium">
                                                Jane Smith
                                            </span>{" "}
                                            registered for Music Festival
                                        </div>
                                        <div className="text-gray-600 text-sm">
                                            <span className="font-medium">
                                                Bob Johnson
                                            </span>{" "}
                                            left a 5-star review
                                        </div>
                                        <div className="text-gray-600 text-sm">
                                            Event{" "}
                                            <span className="font-medium">
                                                Food & Wine Expo
                                            </span>{" "}
                                            sold out
                                        </div>
                                    </div>
                                )}
                        </div>
                    </div>
                </div>
            </div>
        </Layout>
    );
}
