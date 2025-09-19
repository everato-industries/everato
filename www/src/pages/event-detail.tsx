import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import Layout from "../components/layout";
import api from "../lib/api";
import { isAuthenticated } from "../lib/auth";

interface TicketType {
    id: string;
    name: string;
    event_id: string;
    price: number;
    available_tickets: number;
}

interface Coupon {
    id: string;
    event_id: string;
    code: string;
    discount_percentage: number;
    valid_from: string;
    valid_until: string;
    usage_limit: number;
    usage_count?: number;
    created_at: string;
    updated_at: string;
}

interface Event {
    id: string;
    title: string;
    description: string;
    banner: string;
    icon: string;
    admin_id: string;
    start_time: string;
    end_time: string;
    location: string;
    total_seats: number;
    available_seats: number;
    slug: string;
    // Additional comprehensive fields
    organizer_name?: string;
    organizer_email?: string;
    organizer_phone?: string;
    organization?: string;
    contact_email?: string;
    contact_phone?: string;
    refund_policy?: string;
    terms_and_conditions?: string;
    event_type?: string;
    category?: string;
    max_tickets_per_user?: number;
    booking_start_time?: string;
    booking_end_time?: string;
    tags?: string[];
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
    created_at: string;
    updated_at: string;
    // Deprecated/legacy fields for backward compatibility
    date?: string;
    time?: string;
    endDate?: string;
    endTime?: string;
    venue?: string;
    price?: number;
    image?: string;
    organizer?: {
        name: string;
        email: string;
        avatar?: string;
    };
    capacity?: number;
    availableTickets?: number;
    // New fields for ticket types and coupons
    ticket_types: TicketType[];
    coupons: Coupon[];
}

interface SelectedTicket {
    ticketTypeId: string;
    quantity: number;
    price: number;
}

interface TicketSelection {
    selectedTickets: SelectedTicket[];
    totalPrice: number;
    appliedCoupon?: {
        code: string;
        discountPercentage: number;
    };
    finalPrice: number;
}

export default function EventDetailPage() {
    const { slug } = useParams<{ slug: string }>();
    const [event, setEvent] = useState<Event | null>(null);
    const [loading, setLoading] = useState(true);
    const [ticketSelection, setTicketSelection] = useState<TicketSelection>({
        selectedTickets: [],
        totalPrice: 0,
        finalPrice: 0,
    });
    const [showBookingModal, setShowBookingModal] = useState(false);
    const [couponCode, setCouponCode] = useState("");
    const [couponError, setCouponError] = useState("");

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

    const updateTicketQuantity = (
        ticketTypeId: string,
        quantity: number,
        price: number,
    ) => {
        const existingTickets = ticketSelection.selectedTickets.filter((t) =>
            t.ticketTypeId !== ticketTypeId
        );
        const newTickets = quantity > 0
            ? [...existingTickets, { ticketTypeId, quantity, price }]
            : existingTickets;

        const totalPrice = newTickets.reduce(
            (sum, ticket) => sum + (ticket.quantity * ticket.price),
            0,
        );
        const discount = ticketSelection.appliedCoupon
            ? (totalPrice * ticketSelection.appliedCoupon.discountPercentage /
                100)
            : 0;
        const finalPrice = totalPrice - discount;

        setTicketSelection({
            ...ticketSelection,
            selectedTickets: newTickets,
            totalPrice,
            finalPrice,
        });
    };

    const applyCoupon = (code: string) => {
        if (!event) return;

        const coupon = event.coupons.find((c) =>
            c.code.toLowerCase() === code.toLowerCase()
        );
        if (!coupon) {
            setCouponError("Invalid coupon code");
            return;
        }

        const now = new Date();
        const validFrom = new Date(coupon.valid_from);
        const validUntil = new Date(coupon.valid_until);

        if (now < validFrom || now > validUntil) {
            setCouponError("Coupon has expired or is not yet valid");
            return;
        }

        setCouponError("");
        const totalPrice = ticketSelection.totalPrice;
        const discount = totalPrice * coupon.discount_percentage / 100;
        const finalPrice = totalPrice - discount;

        setTicketSelection({
            ...ticketSelection,
            appliedCoupon: {
                code: coupon.code,
                discountPercentage: coupon.discount_percentage,
            },
            finalPrice,
        });
    };

    const removeCoupon = () => {
        setCouponError("");
        setTicketSelection({
            ...ticketSelection,
            appliedCoupon: undefined,
            finalPrice: ticketSelection.totalPrice,
        });
        setCouponCode("");
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
                                {(event.tags || []).map((tag) => (
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

                            {/* Ticket Selection */}
                            {event.ticket_types &&
                                event.ticket_types.length > 0 && (
                                <div className="mb-8">
                                    <h3 className="mb-4 font-bold text-black text-2xl">
                                        🎫 Select Tickets
                                    </h3>
                                    <div className="space-y-4">
                                        {event.ticket_types.map(
                                            (ticketType) => {
                                                const selectedTicket =
                                                    ticketSelection
                                                        .selectedTickets.find(
                                                            (t) =>
                                                                t.ticketTypeId ===
                                                                    ticketType
                                                                        .id,
                                                        );
                                                const currentQuantity =
                                                    selectedTicket?.quantity ||
                                                    0;
                                                const maxQuantity = Math.min(
                                                    ticketType
                                                        .available_tickets,
                                                    event
                                                        .max_tickets_per_user ||
                                                        10,
                                                );

                                                return (
                                                    <div
                                                        key={ticketType.id}
                                                        className={`border p-4 transition-colors ${
                                                            currentQuantity > 0
                                                                ? "bg-blue-50 border-blue-200"
                                                                : "bg-gray-50 border-gray-200"
                                                        }`}
                                                    >
                                                        <div className="flex justify-between items-start mb-3">
                                                            <div className="flex-1">
                                                                <div className="flex items-center gap-3">
                                                                    <input
                                                                        type="checkbox"
                                                                        id={`ticket-${ticketType.id}`}
                                                                        checked={currentQuantity >
                                                                            0}
                                                                        onChange={(
                                                                            e,
                                                                        ) => {
                                                                            if (
                                                                                e.target
                                                                                    .checked
                                                                            ) {
                                                                                updateTicketQuantity(
                                                                                    ticketType
                                                                                        .id,
                                                                                    1,
                                                                                    ticketType
                                                                                        .price,
                                                                                );
                                                                            } else {
                                                                                updateTicketQuantity(
                                                                                    ticketType
                                                                                        .id,
                                                                                    0,
                                                                                    ticketType
                                                                                        .price,
                                                                                );
                                                                            }
                                                                        }}
                                                                        className="bg-white disabled:opacity-50 border-2 border-gray-300 focus:border-black rounded focus:ring-2 focus:ring-black w-5 h-5 text-black accent-black disabled:cursor-not-allowed"
                                                                        disabled={maxQuantity ===
                                                                            0}
                                                                    />
                                                                    <label
                                                                        htmlFor={`ticket-${ticketType.id}`}
                                                                        className="font-semibold text-black text-lg cursor-pointer"
                                                                    >
                                                                        {ticketType
                                                                            .name}
                                                                    </label>
                                                                </div>
                                                                <p className="mt-1 text-gray-600 text-sm">
                                                                    {maxQuantity ===
                                                                            0
                                                                        ? "Sold Out"
                                                                        : `Available: ${ticketType.available_tickets} tickets`}
                                                                </p>
                                                            </div>
                                                            <div className="text-right">
                                                                <div className="font-bold text-black text-xl">
                                                                    ${ticketType
                                                                        .price
                                                                        .toFixed(
                                                                            2,
                                                                        )}
                                                                </div>
                                                            </div>
                                                        </div>

                                                        {currentQuantity > 0 &&
                                                            (
                                                                <div className="flex justify-between items-center bg-white p-3 border rounded">
                                                                    <span className="text-gray-700">
                                                                        Quantity:
                                                                    </span>
                                                                    <div className="flex items-center gap-2">
                                                                        <button
                                                                            onClick={() =>
                                                                                updateTicketQuantity(
                                                                                    ticketType
                                                                                        .id,
                                                                                    Math.max(
                                                                                        0,
                                                                                        currentQuantity -
                                                                                            1,
                                                                                    ),
                                                                                    ticketType
                                                                                        .price,
                                                                                )}
                                                                            disabled={currentQuantity <=
                                                                                1}
                                                                            className="flex justify-center items-center bg-white hover:bg-gray-50 disabled:opacity-50 border border-gray-300 w-8 h-8 font-semibold text-black transition-colors disabled:cursor-not-allowed"
                                                                        >
                                                                            -
                                                                        </button>
                                                                        <span className="mx-3 min-w-[2rem] font-semibold text-center">
                                                                            {currentQuantity}
                                                                        </span>
                                                                        <button
                                                                            onClick={() =>
                                                                                updateTicketQuantity(
                                                                                    ticketType
                                                                                        .id,
                                                                                    Math.min(
                                                                                        maxQuantity,
                                                                                        currentQuantity +
                                                                                            1,
                                                                                    ),
                                                                                    ticketType
                                                                                        .price,
                                                                                )}
                                                                            disabled={currentQuantity >=
                                                                                maxQuantity}
                                                                            className="flex justify-center items-center bg-white hover:bg-gray-50 disabled:opacity-50 border border-gray-300 w-8 h-8 font-semibold text-black transition-colors disabled:cursor-not-allowed"
                                                                        >
                                                                            +
                                                                        </button>
                                                                    </div>
                                                                </div>
                                                            )}
                                                    </div>
                                                );
                                            },
                                        )}
                                    </div>

                                    {/* Coupon Application */}
                                    {ticketSelection.selectedTickets.length >
                                            0 && (
                                        <div className="bg-white mt-6 p-4 border border-gray-200 rounded">
                                            <h4 className="mb-3 font-semibold text-gray-800">
                                                Have a coupon code?
                                            </h4>
                                            <div className="flex gap-2">
                                                <input
                                                    type="text"
                                                    value={couponCode}
                                                    onChange={(e) =>
                                                        setCouponCode(
                                                            e.target.value
                                                                .toUpperCase(),
                                                        )}
                                                    placeholder="Enter coupon code"
                                                    className="flex-1 input-field"
                                                />
                                                <button
                                                    onClick={() =>
                                                        applyCoupon(couponCode)}
                                                    disabled={!couponCode
                                                        .trim()}
                                                    className="disabled:opacity-50 disabled:cursor-not-allowed btn-primary"
                                                >
                                                    Apply
                                                </button>
                                            </div>
                                            {couponError && (
                                                <p className="mt-2 text-red-600 text-sm">
                                                    {couponError}
                                                </p>
                                            )}
                                            {ticketSelection.appliedCoupon && (
                                                <div className="flex justify-between items-center bg-green-50 mt-3 p-2 border border-green-200 rounded">
                                                    <span className="font-medium text-green-800">
                                                        Coupon "{ticketSelection
                                                            .appliedCoupon
                                                            .code}" applied
                                                        ({ticketSelection
                                                            .appliedCoupon
                                                            .discountPercentage}%
                                                        off)
                                                    </span>
                                                    <button
                                                        onClick={removeCoupon}
                                                        className="text-red-600 hover:text-red-800 text-sm underline"
                                                    >
                                                        Remove
                                                    </button>
                                                </div>
                                            )}
                                        </div>
                                    )}
                                </div>
                            )}

                            {/* Coupons - Admin Only */}
                            {isAuthenticated() && event.coupons &&
                                event.coupons.length > 0 && (
                                <div className="mb-8">
                                    <h3 className="mb-4 font-bold text-black text-2xl">
                                        🎟️ Available Coupons
                                    </h3>
                                    <div className="space-y-4">
                                        {event.coupons.map((coupon) => {
                                            const isValid =
                                                new Date(coupon.valid_until) >
                                                    new Date();
                                            const usagePercent =
                                                coupon.usage_count
                                                    ? (coupon.usage_count /
                                                        coupon.usage_limit) *
                                                        100
                                                    : 0;

                                            return (
                                                <div
                                                    key={coupon.id}
                                                    className={`border p-4 ${
                                                        isValid
                                                            ? "bg-green-50 border-green-200"
                                                            : "bg-gray-50 border-gray-200 opacity-60"
                                                    }`}
                                                >
                                                    <div className="flex justify-between items-start">
                                                        <div>
                                                            <div className="flex items-center gap-2 mb-2">
                                                                <code className="bg-gray-800 px-2 py-1 font-mono text-white text-sm">
                                                                    {coupon
                                                                        .code}
                                                                </code>
                                                                {!isValid && (
                                                                    <span className="bg-red-100 px-2 py-1 text-red-600 text-xs">
                                                                        Expired
                                                                    </span>
                                                                )}
                                                            </div>
                                                            <p className="mb-1 text-gray-600 text-sm">
                                                                Valid until:
                                                                {" "}
                                                                {new Date(
                                                                    coupon
                                                                        .valid_until,
                                                                ).toLocaleDateString()}
                                                            </p>
                                                            <p className="text-gray-600 text-sm">
                                                                Used: {coupon
                                                                    .usage_count ||
                                                                    0} / {coupon
                                                                    .usage_limit}
                                                                {" "}
                                                                times
                                                            </p>
                                                            {usagePercent > 0 &&
                                                                (
                                                                    <div className="bg-gray-200 mt-2 w-full h-2">
                                                                        <div
                                                                            className="bg-blue-500 h-2"
                                                                            style={{
                                                                                width:
                                                                                    `${
                                                                                        Math.min(
                                                                                            usagePercent,
                                                                                            100,
                                                                                        )
                                                                                    }%`,
                                                                            }}
                                                                        >
                                                                        </div>
                                                                    </div>
                                                                )}
                                                        </div>
                                                        <div className="text-right">
                                                            <div className="font-bold text-green-600 text-xl">
                                                                {coupon
                                                                    .discount_percentage}%
                                                                OFF
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            );
                                        })}
                                    </div>
                                </div>
                            )}

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
                                    {ticketSelection.selectedTickets.length > 0
                                        ? (
                                            <div>
                                                <div className="mb-2 font-bold text-black text-2xl">
                                                    Total: ${ticketSelection
                                                        .totalPrice.toFixed(2)}
                                                </div>
                                                {ticketSelection
                                                    .appliedCoupon && (
                                                    <div className="mb-2">
                                                        <div className="text-green-600 text-sm">
                                                            Discount:
                                                            -${(ticketSelection
                                                                .totalPrice -
                                                                ticketSelection
                                                                    .finalPrice)
                                                                .toFixed(2)}
                                                        </div>
                                                        <div className="font-bold text-green-700 text-xl">
                                                            Final:
                                                            ${ticketSelection
                                                                .finalPrice
                                                                .toFixed(2)}
                                                        </div>
                                                    </div>
                                                )}
                                            </div>
                                        )
                                        : (
                                            <div>
                                                <div className="mb-2 font-bold text-black text-3xl">
                                                    Select Tickets
                                                </div>
                                                <div className="text-gray-600">
                                                    Choose from available types
                                                    below
                                                </div>
                                            </div>
                                        )}
                                </div>

                                {ticketSelection.selectedTickets.length > 0 && (
                                    <div className="space-y-3 mb-6">
                                        <h4 className="font-semibold text-black text-sm">
                                            Selected Tickets:
                                        </h4>
                                        {ticketSelection.selectedTickets.map(
                                            (selectedTicket) => {
                                                const ticketType = event
                                                    .ticket_types.find((t) =>
                                                        t.id ===
                                                            selectedTicket
                                                                .ticketTypeId
                                                    );
                                                return (
                                                    <div
                                                        key={selectedTicket
                                                            .ticketTypeId}
                                                        className="flex justify-between items-center bg-gray-50 p-3 rounded"
                                                    >
                                                        <div>
                                                            <div className="font-medium text-black">
                                                                {ticketType
                                                                    ?.name}
                                                            </div>
                                                            <div className="text-gray-600 text-sm">
                                                                {selectedTicket
                                                                    .quantity}
                                                                {" "}
                                                                ×
                                                                ${selectedTicket
                                                                    .price
                                                                    .toFixed(2)}
                                                            </div>
                                                        </div>
                                                        <div className="font-semibold text-black">
                                                            ${(selectedTicket
                                                                .quantity *
                                                                selectedTicket
                                                                    .price)
                                                                .toFixed(2)}
                                                        </div>
                                                    </div>
                                                );
                                            },
                                        )}

                                        {ticketSelection.appliedCoupon && (
                                            <div className="flex justify-between items-center bg-green-50 p-3 border border-green-200 rounded">
                                                <div className="font-medium text-green-800">
                                                    Coupon: {ticketSelection
                                                        .appliedCoupon.code}
                                                </div>
                                                <div className="font-semibold text-green-800">
                                                    -{ticketSelection
                                                        .appliedCoupon
                                                        .discountPercentage}%
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                )}

                                <button
                                    onClick={() => setShowBookingModal(true)}
                                    className="mb-4 w-full btn-primary"
                                    disabled={ticketSelection.selectedTickets
                                        .length === 0}
                                >
                                    {ticketSelection.selectedTickets.length ===
                                            0
                                        ? "Select Tickets First"
                                        : "Proceed to Checkout"}
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
                                    <a
                                        href={`https://www.facebook.com/sharer/sharer.php?u=${
                                            encodeURIComponent(
                                                window.location.href,
                                            )
                                        }`}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="flex-1 text-sm text-center btn-secondary"
                                    >
                                        Facebook
                                    </a>
                                    <a
                                        href={`https://twitter.com/intent/tweet?url=${
                                            encodeURIComponent(
                                                window.location.href,
                                            )
                                        }&text=${
                                            encodeURIComponent(event.title)
                                        }`}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="flex-1 text-sm text-center btn-secondary"
                                    >
                                        Twitter
                                    </a>
                                    <a
                                        href={`https://www.linkedin.com/sharing/share-offsite/?url=${
                                            encodeURIComponent(
                                                window.location.href,
                                            )
                                        }`}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="flex-1 text-sm text-center btn-secondary"
                                    >
                                        LinkedIn
                                    </a>
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
                                {ticketSelection.selectedTickets.map(
                                    (selectedTicket) => {
                                        const ticketType = event.ticket_types
                                            .find((t) =>
                                                t.id ===
                                                    selectedTicket.ticketTypeId
                                            );
                                        return (
                                            <div
                                                key={selectedTicket
                                                    .ticketTypeId}
                                                className="flex justify-between items-center mb-2"
                                            >
                                                <span>
                                                    {ticketType?.name} ×{" "}
                                                    {selectedTicket.quantity}
                                                </span>
                                                <span>
                                                    ${(selectedTicket.quantity *
                                                        selectedTicket.price)
                                                        .toFixed(2)}
                                                </span>
                                            </div>
                                        );
                                    },
                                )}

                                {ticketSelection.appliedCoupon && (
                                    <div className="flex justify-between items-center mb-2 text-green-600">
                                        <span>
                                            Coupon:{" "}
                                            {ticketSelection.appliedCoupon.code}
                                        </span>
                                        <span>
                                            -${(ticketSelection.totalPrice -
                                                ticketSelection.finalPrice)
                                                .toFixed(2)}
                                        </span>
                                    </div>
                                )}

                                <div className="mt-2 pt-2 border-gray-300 border-t">
                                    <div className="flex justify-between items-center font-semibold">
                                        <span>Total</span>
                                        <span>
                                            ${(ticketSelection.appliedCoupon
                                                ? ticketSelection.finalPrice
                                                : ticketSelection.totalPrice)
                                                .toFixed(2)}
                                        </span>
                                    </div>
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
                                to={`/checkout?event=${event.id}&tickets=${
                                    encodeURIComponent(
                                        JSON.stringify(ticketSelection),
                                    )
                                }`}
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
