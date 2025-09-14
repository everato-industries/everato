import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import Layout from "../components/layout";
import api from "../lib/api";

interface Event {
    id: string;
    title: string;
    description: string;
    banner?: string;
    icon?: string;
    admin_id: string;
    start_time: string;
    end_time: string;
    location: string;
    total_seats: number;
    available_seats: number;
    created_at: string;
    updated_at: string;
    slug: string;
    status: string;
    organizer_name?: string;
    organizer_email?: string;
    organizer_phone?: string;
    organization?: string;
    contact_email?: string;
    contact_phone?: string;
    refund_policy?: string;
    terms_and_conditions?: string;
    event_type: string;
    category: string;
    max_tickets_per_user: number;
    booking_start_time?: string;
    booking_end_time?: string;
    tags: string[];
    website_url?: string;
    facebook_url?: string;
    twitter_url?: string;
    instagram_url?: string;
    linkedin_url?: string;
    venue_name?: string;
    address_line1?: string;
    address_line2?: string;
    city?: string;
    state?: string;
    postal_code?: string;
    country?: string;
    latitude?: number;
    longitude?: number;
}

interface TicketSelection {
    quantity: number;
    totalPrice: number;
}

export default function EventDetailPage() {
    const { slug } = useParams<{ slug: string }>();
    const [event, setEvent] = useState<Event | null>(null);
    const [loading, setLoading] = useState(true);
    const [ticketSelection, setTicketSelection] = useState<TicketSelection>({
        quantity: 1,
        totalPrice: 0,
    });
    const [showBookingModal, setShowBookingModal] = useState(false);

    useEffect(() => {
        const fetchEvent = async () => {
            try {
                setLoading(true);
                const response = await api.get(`/events/${slug}`);

                // Handle the API response structure
                if (response.data && response.data.event) {
                    setEvent(response.data.event);
                } else {
                    console.log("No event found in response");
                    setEvent(null);
                }
                setLoading(false);
            } catch (error) {
                console.error("Error fetching event:", error);
                setEvent(null);
                setLoading(false);
            }
        };

        if (slug) {
            fetchEvent();
        }
    }, [slug]);

    const handleQuantityChange = (newQuantity: number) => {
        if (
            event && newQuantity >= 1 &&
            newQuantity <=
                Math.min(
                    event.max_tickets_per_user || 10,
                    event.available_seats,
                )
        ) {
            // Since there's no price in the API response, we'll set it to 0 for now
            // You may need to add a price field to your API or handle pricing differently
            setTicketSelection({
                quantity: newQuantity,
                totalPrice: newQuantity * 0, // No price field in API
            });
        }
    };

    const formatDate = (dateTime: string) => {
        const eventDate = new Date(dateTime);
        if (isNaN(eventDate.getTime())) {
            return "Invalid Date";
        }
        return eventDate.toLocaleDateString("en-US", {
            weekday: "long",
            year: "numeric",
            month: "long",
            day: "numeric",
            hour: "numeric",
            minute: "2-digit",
        });
    };

    if (loading) {
        return (
            <Layout>
                <div className="mx-auto px-4 sm:px-6 lg:px-8 py-12 max-w-7xl">
                    <div className="loading">
                        <div className="bg-gray-200 mb-8 h-96"></div>
                        <div className="bg-gray-200 mb-4 h-8"></div>
                        <div className="bg-gray-200 mb-2 h-4"></div>
                        <div className="bg-gray-200 w-3/4 h-4"></div>
                    </div>
                </div>
            </Layout>
        );
    }

    if (!event) {
        return (
            <Layout>
                <div className="mx-auto px-4 sm:px-6 lg:px-8 py-12 max-w-7xl text-center">
                    <h1 className="mb-4 font-bold text-black text-2xl">
                        Event Not Found
                    </h1>
                    <p className="mb-6 text-gray-600">
                        The event you're looking for doesn't exist or has been
                        removed.
                    </p>
                    <Link to="/events" className="btn-primary">
                        Browse Other Events
                    </Link>
                </div>
            </Layout>
        );
    }

    return (
        <Layout>
            <div className="mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl">
                {/* Breadcrumb */}
                <nav className="mb-8">
                    <ol className="flex items-center space-x-2 text-gray-500 text-sm">
                        <li>
                            <Link to="/" className="hover:text-black">
                                Home
                            </Link>
                        </li>
                        <li>/</li>
                        <li>
                            <Link to="/events" className="hover:text-black">
                                Events
                            </Link>
                        </li>
                        <li>/</li>
                        <li className="text-black">{event.title}</li>
                    </ol>
                </nav>

                <div className="gap-8 grid grid-cols-1 lg:grid-cols-3">
                    {/* Main Content */}
                    <div className="lg:col-span-2">
                        {/* Hero Image */}
                        <div className="flex justify-center items-center bg-gray-100 mb-8 max-h-96 overflow-hidden">
                            <img
                                src={event.banner}
                                alt="banner"
                                className="w-full h-full object-cover aspect-video"
                            />
                        </div>

                        {/* Event Info */}
                        <div className="mb-8">
                            <div className="mb-4">
                                <span className="bg-gray-200 px-3 py-1 text-gray-700 text-sm uppercase tracking-wide">
                                    {event.category}
                                </span>
                            </div>

                            <h1 className="mb-4 font-bold text-black text-4xl">
                                {event.title}
                            </h1>

                            <div className="flex flex-wrap gap-2 mb-6">
                                {event.tags.map((tag) => (
                                    <span
                                        key={tag}
                                        className="bg-black px-2 py-1 text-white text-xs"
                                    >
                                        {tag}
                                    </span>
                                ))}
                            </div>

                            <div className="gap-6 grid grid-cols-1 md:grid-cols-2 mb-8">
                                <div>
                                    <h3 className="mb-2 font-semibold text-black">
                                        📅 Date & Time
                                    </h3>
                                    <p className="text-gray-600">
                                        <span>
                                            From -{" "}
                                            {formatDate(event.start_time)}
                                        </span>
                                        <br />
                                        <span>
                                            To - {formatDate(event.end_time)}
                                        </span>
                                    </p>
                                </div>
                                <div>
                                    <h3 className="mb-2 font-semibold text-black">
                                        📍 Location
                                    </h3>
                                    <p className="text-gray-600">
                                        {event.venue_name || event.location}
                                    </p>
                                    {event.venue_name && event.location && (
                                        <p className="text-gray-500 text-sm">
                                            {event.location}
                                        </p>
                                    )}
                                    {event.city && event.state && (
                                        <p className="text-gray-500 text-sm">
                                            {event.city}, {event.state}{" "}
                                            {event.country}
                                        </p>
                                    )}
                                </div>
                                <div>
                                    <h3 className="mb-2 font-semibold text-black">
                                        👤 Organizer
                                    </h3>
                                    <p className="text-gray-600">
                                        {event.organization + " ~ " +
                                            event.organizer_name}
                                    </p>
                                    {event.organizer_email && (
                                        <p className="text-gray-500 text-sm">
                                            {event.organizer_email}
                                        </p>
                                    )}
                                    {event.contact_email &&
                                        !event.organizer_email && (
                                        <p className="text-gray-500 text-sm">
                                            {event.contact_email}
                                        </p>
                                    )}
                                </div>
                                <div>
                                    <h3 className="mb-2 font-semibold text-black">
                                        🎫 Availability
                                    </h3>
                                    <p className="text-gray-600">
                                        {event.available_seats} of{" "}
                                        {event.total_seats} tickets remaining
                                    </p>
                                </div>
                            </div>

                            <div>
                                <h3 className="mb-4 font-bold text-black text-2xl">
                                    About This Event
                                </h3>
                                <div className="max-w-none prose prose-gray">
                                    <p className="text-gray-600 leading-relaxed">
                                        {event.description}
                                    </p>
                                    <p className="mt-4 font-thin text-gray-400 text-sm italic leading-relaxed">
                                        If you have any questions, feel free to
                                        contact the organizer at{" "}
                                        {event.organizer_email ||
                                            event.contact_email ||
                                            "N/A"}.
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Sidebar - Booking */}
                    <div className="lg:col-span-1">
                        <div className="top-8 sticky">
                            <div className="card">
                                <div className="mb-6 text-center">
                                    <div className="mb-2 font-bold text-black text-3xl">
                                        Free
                                    </div>
                                    <div className="text-gray-600">
                                        per ticket
                                    </div>
                                </div>

                                <div className="space-y-4 mb-6">
                                    <div>
                                        <label className="block mb-2 font-medium text-black text-sm">
                                            Number of Tickets
                                        </label>
                                        <div className="flex items-center border border-gray-300">
                                            <button
                                                onClick={() =>
                                                    handleQuantityChange(
                                                        ticketSelection
                                                            .quantity - 1,
                                                    )}
                                                className="hover:bg-gray-100 px-4 py-2 transition-colors duration-200"
                                                disabled={ticketSelection
                                                    .quantity <= 1}
                                            >
                                                -
                                            </button>
                                            <span className="flex-1 py-2 border-gray-300 border-r border-l text-center">
                                                {ticketSelection.quantity}
                                            </span>
                                            <button
                                                onClick={() =>
                                                    handleQuantityChange(
                                                        ticketSelection
                                                            .quantity + 1,
                                                    )}
                                                className="hover:bg-gray-100 px-4 py-2 transition-colors duration-200"
                                                disabled={ticketSelection
                                                    .quantity >=
                                                    Math.min(
                                                        event
                                                            .max_tickets_per_user ||
                                                            10,
                                                        event.available_seats,
                                                    )}
                                            >
                                                +
                                            </button>
                                        </div>
                                        <p className="mt-1 text-gray-500 text-xs">
                                            Maximum 10 tickets per order
                                        </p>
                                    </div>

                                    <div className="bg-gray-50 p-4">
                                        <div className="flex justify-between items-center">
                                            <span className="font-medium text-black">
                                                Total:
                                            </span>
                                            <span className="font-bold text-black text-xl">
                                                ${ticketSelection.totalPrice}
                                            </span>
                                        </div>
                                    </div>
                                </div>

                                <button
                                    onClick={() => setShowBookingModal(true)}
                                    className="mb-4 w-full btn-primary"
                                    disabled={event.available_seats === 0}
                                >
                                    {event.available_seats === 0
                                        ? "Sold Out"
                                        : "Book Now"}
                                </button>

                                <div className="text-gray-500 text-xs text-center">
                                    <p>
                                        Free cancellation up to 24 hours before
                                        the event
                                    </p>
                                    <p className="mt-1">
                                        Secure payment powered by Stripe
                                    </p>
                                </div>
                            </div>

                            {/* Share Event */}
                            <div className="mt-6 card">
                                <h3 className="mb-4 font-semibold text-black">
                                    Share This Event
                                </h3>
                                <div className="flex space-x-3">
                                    <button className="flex-1 text-sm btn-secondary">
                                        Facebook
                                    </button>
                                    <button className="flex-1 text-sm btn-secondary">
                                        Twitter
                                    </button>
                                    <button className="flex-1 text-sm btn-secondary">
                                        LinkedIn
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Booking Modal */}
            {showBookingModal && (
                <div
                    className="modal-backdrop"
                    onClick={() => setShowBookingModal(false)}
                >
                    <div className="modal" onClick={(e) => e.stopPropagation()}>
                        <h2 className="mb-4 font-bold text-black text-2xl">
                            Complete Your Booking
                        </h2>
                        <div className="mb-6">
                            <h3 className="mb-2 font-semibold text-black">
                                {event.title}
                            </h3>
                            <p className="mb-2 text-gray-600">
                                {formatDate(event.start_time)}
                            </p>
                            <p className="mb-4 text-gray-600">
                                {event.venue_name || event.location}
                            </p>

                            <div className="bg-gray-50 p-4">
                                <div className="flex justify-between items-center mb-2">
                                    <span>
                                        Tickets × {ticketSelection.quantity}
                                    </span>
                                    <span>Free</span>
                                </div>
                                <div className="flex justify-between items-center font-semibold">
                                    <span>Total</span>
                                    <span>Free</span>
                                </div>
                            </div>
                        </div>

                        <div className="flex space-x-4">
                            <button
                                onClick={() => setShowBookingModal(false)}
                                className="flex-1 btn-secondary"
                            >
                                Cancel
                            </button>
                            <Link
                                to={`/checkout?event=${event.id}&quantity=${ticketSelection.quantity}`}
                                className="flex-1 text-center btn-primary"
                            >
                                Proceed to Checkout
                            </Link>
                        </div>
                    </div>
                </div>
            )}
        </Layout>
    );
}
