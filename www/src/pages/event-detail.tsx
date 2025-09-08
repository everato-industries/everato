import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import Layout from "../components/layout";

interface Event {
    id: string;
    title: string;
    description: string;
    date: string;
    time: string;
    endDate?: string;
    endTime?: string;
    location: string;
    venue: string;
    price: number;
    image?: string;
    slug: string;
    category: string;
    organizer: {
        name: string;
        email: string;
        avatar?: string;
    };
    capacity: number;
    availableTickets: number;
    tags: string[];
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
                // Mock API call - replace with actual API
                setTimeout(() => {
                    // Mock event data based on slug
                    const mockEvent: Event = {
                        id: "1",
                        title: slug === "tech-conference-2025"
                            ? "Tech Conference 2025"
                            : "Sample Event",
                        description:
                            "Join us for an incredible experience that you'll never forget. This event brings together amazing people, great content, and unforgettable moments. Whether you're here to learn, network, or simply enjoy yourself, we've created something special just for you.",
                        date: "2025-10-15",
                        time: "09:00",
                        endDate: "2025-10-17",
                        endTime: "18:00",
                        location: "San Francisco, CA",
                        venue:
                            "Moscone Center, 747 Howard St, San Francisco, CA 94103",
                        price: 299,
                        slug: slug || "",
                        category: "Technology",
                        organizer: {
                            name: "Tech Events Inc",
                            email: "contact@techevents.com",
                        },
                        capacity: 1000,
                        availableTickets: 850,
                        tags: [
                            "Technology",
                            "AI",
                            "Web Development",
                            "Networking",
                        ],
                    };

                    setEvent(mockEvent);
                    setTicketSelection({
                        quantity: 1,
                        totalPrice: mockEvent.price,
                    });
                    setLoading(false);
                }, 500);
            } catch (error) {
                console.error("Error fetching event:", error);
                setLoading(false);
            }
        };

        fetchEvent();
    }, [slug]);

    const handleQuantityChange = (newQuantity: number) => {
        if (
            event && newQuantity >= 1 &&
            newQuantity <= Math.min(10, event.availableTickets)
        ) {
            setTicketSelection({
                quantity: newQuantity,
                totalPrice: newQuantity * event.price,
            });
        }
    };

    const formatDate = (date: string, time?: string) => {
        const eventDate = new Date(`${date}T${time || "00:00"}`);
        return eventDate.toLocaleDateString("en-US", {
            weekday: "long",
            year: "numeric",
            month: "long",
            day: "numeric",
            ...(time && { hour: "numeric", minute: "2-digit" }),
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
                        <div className="flex justify-center items-center bg-gray-100 mb-8 h-96">
                            <span className="text-gray-400 text-xl">
                                Event Image
                            </span>
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
                                        {formatDate(event.date, event.time)}
                                        {event.endDate && (
                                            <span>
                                                - {formatDate(
                                                    event.endDate,
                                                    event.endTime,
                                                )}
                                            </span>
                                        )}
                                    </p>
                                </div>
                                <div>
                                    <h3 className="mb-2 font-semibold text-black">
                                        📍 Location
                                    </h3>
                                    <p className="text-gray-600">
                                        {event.venue}
                                    </p>
                                    <p className="text-gray-500 text-sm">
                                        {event.location}
                                    </p>
                                </div>
                                <div>
                                    <h3 className="mb-2 font-semibold text-black">
                                        👤 Organizer
                                    </h3>
                                    <p className="text-gray-600">
                                        {event.organizer.name}
                                    </p>
                                    <p className="text-gray-500 text-sm">
                                        {event.organizer.email}
                                    </p>
                                </div>
                                <div>
                                    <h3 className="mb-2 font-semibold text-black">
                                        🎫 Availability
                                    </h3>
                                    <p className="text-gray-600">
                                        {event.availableTickets} of{" "}
                                        {event.capacity} tickets remaining
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
                                    <p className="mt-4 text-gray-600 leading-relaxed">
                                        This is a fantastic opportunity to
                                        connect with like-minded individuals,
                                        learn from industry experts, and
                                        discover new perspectives. We've
                                        carefully curated an experience that
                                        combines valuable content with
                                        meaningful networking opportunities.
                                    </p>
                                    <p className="mt-4 text-gray-600 leading-relaxed">
                                        Don't miss out on this unique event.
                                        Spaces are limited and filling up fast.
                                        Secure your spot today and be part of
                                        something extraordinary.
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
                                        ${event.price}
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
                                                        10,
                                                        event.availableTickets,
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
                                    disabled={event.availableTickets === 0}
                                >
                                    {event.availableTickets === 0
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
                                {formatDate(event.date, event.time)}
                            </p>
                            <p className="mb-4 text-gray-600">{event.venue}</p>

                            <div className="bg-gray-50 p-4">
                                <div className="flex justify-between items-center mb-2">
                                    <span>
                                        Tickets × {ticketSelection.quantity}
                                    </span>
                                    <span>${ticketSelection.totalPrice}</span>
                                </div>
                                <div className="flex justify-between items-center font-semibold">
                                    <span>Total</span>
                                    <span>${ticketSelection.totalPrice}</span>
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
