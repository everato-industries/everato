import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Layout from "../components/layout";

interface Event {
    id: string;
    title: string;
    description: string;
    date: string;
    location: string;
    price: number;
    image?: string;
    slug: string;
}

export default function HomePage() {
    const [events, setEvents] = useState<Event[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function fetchData() {
            try {
                // Health check
                const healthResponse = await fetch("/api/v1/health");
                if (healthResponse.ok) {
                    const healthData = await healthResponse.json();
                    console.log("Health check:", healthData);
                }

                // Mock events data for now
                setTimeout(() => {
                    setEvents([
                        {
                            id: "1",
                            title: "Tech Conference 2025",
                            description:
                                "Annual technology conference featuring the latest in AI and web development.",
                            date: "2025-10-15",
                            location: "San Francisco, CA",
                            price: 299,
                            slug: "tech-conference-2025",
                        },
                        {
                            id: "2",
                            title: "Music Festival",
                            description:
                                "Three-day outdoor music festival with top artists.",
                            date: "2025-08-20",
                            location: "Austin, TX",
                            price: 199,
                            slug: "music-festival-2025",
                        },
                        {
                            id: "3",
                            title: "Food & Wine Expo",
                            description:
                                "Culinary experience with renowned chefs and wine tastings.",
                            date: "2025-09-10",
                            location: "New York, NY",
                            price: 149,
                            slug: "food-wine-expo",
                        },
                    ]);
                    setLoading(false);
                }, 1000);
            } catch (error) {
                console.error("Error fetching data:", error);
                setLoading(false);
            }
        }
        fetchData();
    }, []);

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString("en-US", {
            weekday: "long",
            year: "numeric",
            month: "long",
            day: "numeric",
        });
    };

    return (
        <Layout>
            {/* Hero Section */}
            <section className="bg-white py-20">
                <div className="mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl">
                    <div className="mx-auto max-w-3xl text-center">
                        <h1 className="mb-6 font-bold text-black text-5xl md:text-6xl">
                            Discover Amazing Events
                        </h1>
                        <p className="mb-8 text-gray-600 text-xl">
                            Find and book tickets for concerts, conferences,
                            festivals, and more. Create unforgettable
                            experiences with Everato.
                        </p>
                        <div className="flex sm:flex-row flex-col justify-center gap-4">
                            <Link
                                to="/events"
                                className="inline-block text-center btn-primary"
                            >
                                Browse Events
                            </Link>
                            <Link
                                to="/create-event"
                                className="inline-block text-center btn-secondary"
                            >
                                Create Event
                            </Link>
                        </div>
                    </div>
                </div>
            </section>

            {/* Stats Section */}
            <section className="bg-gray-50 py-16">
                <div className="mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl">
                    <div className="gap-8 grid grid-cols-1 md:grid-cols-3 text-center">
                        <div>
                            <div className="mb-2 font-bold text-black text-4xl">
                                10K+
                            </div>
                            <div className="text-gray-600">Events Created</div>
                        </div>
                        <div>
                            <div className="mb-2 font-bold text-black text-4xl">
                                500K+
                            </div>
                            <div className="text-gray-600">Tickets Sold</div>
                        </div>
                        <div>
                            <div className="mb-2 font-bold text-black text-4xl">
                                50+
                            </div>
                            <div className="text-gray-600">
                                Cities Worldwide
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            {/* Featured Events */}
            <section className="bg-white py-20">
                <div className="mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl">
                    <div className="mb-12 text-center">
                        <h2 className="mb-4 font-bold text-black text-4xl">
                            Featured Events
                        </h2>
                        <p className="text-gray-600 text-xl">
                            Don't miss out on these incredible upcoming events
                        </p>
                    </div>

                    {loading
                        ? (
                            <div className="gap-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
                                {[1, 2, 3].map((i) => (
                                    <div key={i} className="card loading">
                                        <div className="bg-gray-200 mb-4 h-48">
                                        </div>
                                        <div className="bg-gray-200 mb-2 h-4">
                                        </div>
                                        <div className="bg-gray-200 w-3/4 h-4">
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )
                        : (
                            <div className="gap-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
                                {events.map((event) => (
                                    <div
                                        key={event.id}
                                        className="group cursor-pointer card"
                                    >
                                        <div className="flex justify-center items-center bg-gray-100 mb-4 h-48">
                                            <span className="text-gray-400">
                                                Event Image
                                            </span>
                                        </div>
                                        <div className="mb-2">
                                            <span className="text-gray-500 text-sm">
                                                {formatDate(event.date)}
                                            </span>
                                        </div>
                                        <h3 className="mb-2 font-semibold text-black group-hover:text-gray-700 text-xl transition-colors duration-200">
                                            {event.title}
                                        </h3>
                                        <p className="mb-4 text-gray-600 line-clamp-2">
                                            {event.description}
                                        </p>
                                        <div className="flex justify-between items-center mb-4">
                                            <span className="text-gray-500 text-sm">
                                                📍 {event.location}
                                            </span>
                                            <span className="font-semibold text-black text-lg">
                                                ${event.price}
                                            </span>
                                        </div>
                                        <Link
                                            to={`/events/${event.slug}`}
                                            className="w-full text-center btn-primary"
                                        >
                                            View Details
                                        </Link>
                                    </div>
                                ))}
                            </div>
                        )}

                    <div className="mt-12 text-center">
                        <Link
                            to="/events"
                            className="btn-secondary"
                        >
                            View All Events
                        </Link>
                    </div>
                </div>
            </section>

            {/* How It Works */}
            <section className="bg-gray-50 py-20">
                <div className="mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl">
                    <div className="mb-12 text-center">
                        <h2 className="mb-4 font-bold text-black text-4xl">
                            How It Works
                        </h2>
                        <p className="text-gray-600 text-xl">
                            Get started in three simple steps
                        </p>
                    </div>

                    <div className="gap-8 grid grid-cols-1 md:grid-cols-3">
                        <div className="text-center">
                            <div className="flex justify-center items-center bg-black mx-auto mb-4 w-16 h-16 font-bold text-white text-2xl">
                                1
                            </div>
                            <h3 className="mb-2 font-semibold text-black text-xl">
                                Discover Events
                            </h3>
                            <p className="text-gray-600">
                                Browse through thousands of events in your area
                                or worldwide
                            </p>
                        </div>

                        <div className="text-center">
                            <div className="flex justify-center items-center bg-black mx-auto mb-4 w-16 h-16 font-bold text-white text-2xl">
                                2
                            </div>
                            <h3 className="mb-2 font-semibold text-black text-xl">
                                Book Tickets
                            </h3>
                            <p className="text-gray-600">
                                Secure your spot with our easy and safe booking
                                process
                            </p>
                        </div>

                        <div className="text-center">
                            <div className="flex justify-center items-center bg-black mx-auto mb-4 w-16 h-16 font-bold text-white text-2xl">
                                3
                            </div>
                            <h3 className="mb-2 font-semibold text-black text-xl">
                                Enjoy Experience
                            </h3>
                            <p className="text-gray-600">
                                Show up with your digital ticket and enjoy the
                                event
                            </p>
                        </div>
                    </div>
                </div>
            </section>

            {/* CTA Section */}
            <section className="bg-black py-20">
                <div className="mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl text-center">
                    <h2 className="mb-4 font-bold text-white text-4xl">
                        Ready to Get Started?
                    </h2>
                    <p className="mb-8 text-gray-300 text-xl">
                        Join thousands of event organizers and attendees on
                        Everato
                    </p>
                    <div className="flex sm:flex-row flex-col justify-center gap-4">
                        <Link
                            to="/auth/register"
                            className="bg-white hover:bg-gray-100 px-8 py-3 font-medium text-black transition-colors duration-200"
                        >
                            Sign Up Free
                        </Link>
                        <Link
                            to="/contact"
                            className="hover:bg-white px-8 py-3 border border-white font-medium text-white hover:text-black transition-colors duration-200"
                        >
                            Contact Sales
                        </Link>
                    </div>
                </div>
            </section>
        </Layout>
    );
}
