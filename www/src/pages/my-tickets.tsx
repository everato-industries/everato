import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Layout from "../components/layout";
import { bookingAPI, ticketAPI, type Ticket } from "../lib/api";
import { useAuth } from "../hooks/useAuth";

// ─── Types ────────────────────────────────────────────────────────────────────

interface BookingTicketItem {
    id: string;
    ticket_type_name?: string;
    quantity: number;
    price_per_unit: number;
    subtotal: number;
}

interface Booking {
    id: string;
    event_id: string;
    event_title?: string;
    event_start_time?: string;
    event_end_time?: string;
    status: string;
    total_amount: number;
    discount_amount?: number;
    final_amount: number;
    coupon_code?: string;
    created_at: string;
    tickets?: BookingTicketItem[];
}

// ─── Status Badge ─────────────────────────────────────────────────────────────

function StatusBadge({ status }: { status: string }) {
    const colors: Record<string, string> = {
        CONFIRMED: "bg-green-100 text-green-700",
        PENDING: "bg-yellow-100 text-yellow-700",
        CANCELLED: "bg-red-100 text-red-700",
        COMPLETED: "bg-blue-100 text-blue-700",
    };

    return (
        <span
            className={`px-3 py-1 rounded-full text-xs font-medium ${
                colors[status] ?? "bg-gray-100 text-gray-600"
            }`}
        >
            {status}
        </span>
    );
}

// ─── QR Ticket Modal ──────────────────────────────────────────────────────────

interface TicketModalProps {
    booking: Booking;
    onClose: () => void;
}

function TicketModal({ booking, onClose }: TicketModalProps) {
    const [tickets, setTickets] = useState<Ticket[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [selectedIndex, setSelectedIndex] = useState(0);

    useEffect(() => {
        const fetchTickets = async () => {
            try {
                setLoading(true);
                const res = await ticketAPI.getTicketsByBooking(booking.id);
                const list = res.data.data?.tickets ?? res.data.tickets ?? [];
                setTickets(list);
            } catch (err: unknown) {
                console.error("Error loading tickets:", err);
                setError("Could not load QR ticket. Please try again.");
            } finally {
                setLoading(false);
            }
        };

        fetchTickets();
    }, [booking.id]);

    const activeTicket = tickets[selectedIndex];

    const eventDate = booking.event_start_time
        ? new Date(booking.event_start_time).toLocaleDateString("en-IN", {
              year: "numeric",
              month: "long",
              day: "numeric",
              hour: "2-digit",
              minute: "2-digit",
          })
        : "Date TBD";

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
            <div className="bg-white rounded-2xl max-w-md w-full overflow-hidden shadow-2xl animate-in fade-in zoom-in-95 duration-200">
                {/* Header */}
                <div className="bg-gray-900 text-white p-6 relative">
                    <button
                        onClick={onClose}
                        className="absolute top-4 right-4 text-gray-400 hover:text-white p-1 rounded-lg transition-colors"
                        aria-label="Close"
                    >
                        ✕
                    </button>
                    <span className="text-xs uppercase tracking-wider text-gray-400 font-semibold">
                        Admission Pass
                    </span>
                    <h3 className="text-xl font-bold mt-1 text-white truncate">
                        {booking.event_title ?? "Event Ticket"}
                    </h3>
                    <p className="text-sm text-gray-300 mt-1">{eventDate}</p>
                </div>

                {/* Content */}
                <div className="p-6">
                    {loading && (
                        <div className="py-12 text-center">
                            <div className="animate-spin rounded-full border-b-2 border-black w-8 h-8 mx-auto" />
                            <p className="text-sm text-gray-500 mt-3">Loading ticket pass…</p>
                        </div>
                    )}

                    {error && (
                        <div className="bg-red-50 border border-red-200 text-red-700 p-4 rounded-xl text-sm text-center">
                            {error}
                        </div>
                    )}

                    {!loading && !error && tickets.length === 0 && (
                        <div className="py-8 text-center text-gray-500 text-sm">
                            No tickets generated for this booking yet.
                        </div>
                    )}

                    {!loading && !error && activeTicket && (
                        <div className="text-center space-y-4">
                            {/* Ticket Selector if multiple */}
                            {tickets.length > 1 && (
                                <div className="flex justify-center gap-2 mb-2 overflow-x-auto pb-1">
                                    {tickets.map((t, idx) => (
                                        <button
                                            key={t.id}
                                            onClick={() => setSelectedIndex(idx)}
                                            className={`px-3 py-1 text-xs rounded-full font-medium transition-all ${
                                                selectedIndex === idx
                                                    ? "bg-black text-white"
                                                    : "bg-gray-100 text-gray-600 hover:bg-gray-200"
                                            }`}
                                        >
                                            Ticket #{idx + 1}
                                        </button>
                                    ))}
                                </div>
                            )}

                            {/* Ticket Badge & Number */}
                            <div>
                                <span
                                    className={`inline-block px-3 py-1 rounded-full text-xs font-semibold ${
                                        activeTicket.is_checked_in
                                            ? "bg-amber-100 text-amber-800"
                                            : "bg-emerald-100 text-emerald-800"
                                    }`}
                                >
                                    {activeTicket.is_checked_in ? "✓ Already Checked In" : "● Valid for Entry"}
                                </span>
                                <p className="font-mono text-xs text-gray-500 mt-1 select-all">
                                    {activeTicket.ticket_number}
                                </p>
                            </div>

                            {/* QR Code Container */}
                            <div className="p-4 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200 flex flex-col items-center justify-center">
                                {activeTicket.qr_code_image ? (
                                    <img
                                        src={
                                            activeTicket.qr_code_image.startsWith("data:")
                                                ? activeTicket.qr_code_image
                                                : `data:image/png;base64,${activeTicket.qr_code_image}`
                                        }
                                        alt={`QR Code for ${activeTicket.ticket_number}`}
                                        className="w-52 h-52 object-contain bg-white p-2 rounded-xl shadow-sm"
                                    />
                                ) : (
                                    <div className="w-52 h-52 flex items-center justify-center bg-gray-100 text-gray-400 text-xs">
                                        QR image not available
                                    </div>
                                )}
                                <p className="text-xs text-gray-400 mt-2">
                                    Scan at venue gate for entry
                                </p>
                            </div>

                            {/* Meta info */}
                            <div className="text-xs text-gray-500 bg-gray-50 rounded-lg p-3 space-y-1 text-left">
                                {activeTicket.ticket_type_name && (
                                    <div className="flex justify-between">
                                        <span>Tier:</span>
                                        <span className="font-medium text-gray-900">{activeTicket.ticket_type_name}</span>
                                    </div>
                                )}
                                {activeTicket.checked_in_at && (
                                    <div className="flex justify-between text-amber-700">
                                        <span>Checked in at:</span>
                                        <span>{new Date(activeTicket.checked_in_at).toLocaleString("en-IN")}</span>
                                    </div>
                                )}
                            </div>
                        </div>
                    )}
                </div>

                {/* Footer action */}
                <div className="bg-gray-50 p-4 border-t border-gray-100 flex gap-3 justify-end">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                    >
                        Close
                    </button>
                    {activeTicket?.qr_code_image && (
                        <a
                            href={
                                activeTicket.qr_code_image.startsWith("data:")
                                    ? activeTicket.qr_code_image
                                    : `data:image/png;base64,${activeTicket.qr_code_image}`
                            }
                            download={`ticket-${activeTicket.ticket_number}.png`}
                            className="px-4 py-2 text-sm font-medium text-white bg-black rounded-lg hover:bg-gray-800 transition-colors"
                        >
                            Save QR Code
                        </a>
                    )}
                </div>
            </div>
        </div>
    );
}

// ─── Booking Card ─────────────────────────────────────────────────────────────

function BookingCard({
    booking,
    onShowTickets,
}: {
    booking: Booking;
    onShowTickets: (b: Booking) => void;
}) {
    const eventDate = booking.event_start_time
        ? new Date(booking.event_start_time).toLocaleDateString("en-IN", {
              year: "numeric",
              month: "long",
              day: "numeric",
              hour: "2-digit",
              minute: "2-digit",
          })
        : "Date TBD";

    return (
        <div className="bg-white border border-gray-200 rounded-xl p-6 shadow-sm hover:shadow-md transition-shadow">
            {/* Header */}
            <div className="flex justify-between items-start mb-4">
                <div className="flex-1 mr-4">
                    <h3 className="font-semibold text-lg text-gray-900">
                        {booking.event_title ?? "Event"}
                    </h3>
                    <p className="text-gray-500 text-sm mt-1">{eventDate}</p>
                </div>
                <StatusBadge status={booking.status} />
            </div>

            {/* Booking details */}
            <div className="bg-gray-50 rounded-lg p-4 space-y-2 text-sm">
                <div className="flex justify-between">
                    <span className="text-gray-500">Booking ID</span>
                    <span className="font-mono text-xs text-gray-700">{booking.id.slice(0, 8)}…</span>
                </div>
                <div className="flex justify-between">
                    <span className="text-gray-500">Booked on</span>
                    <span className="text-gray-700">
                        {new Date(booking.created_at).toLocaleDateString("en-IN")}
                    </span>
                </div>
                {booking.coupon_code && (
                    <div className="flex justify-between">
                        <span className="text-gray-500">Coupon</span>
                        <span className="text-green-600 font-mono">{booking.coupon_code}</span>
                    </div>
                )}
                {booking.discount_amount != null && booking.discount_amount > 0 && (
                    <div className="flex justify-between text-green-600">
                        <span>Discount</span>
                        <span>−₹{Number(booking.discount_amount).toFixed(2)}</span>
                    </div>
                )}
                <div className="flex justify-between font-semibold border-t border-gray-200 pt-2">
                    <span>Amount Paid</span>
                    <span>₹{Number(booking.final_amount).toFixed(2)}</span>
                </div>
            </div>

            {/* Ticket items */}
            {booking.tickets && booking.tickets.length > 0 && (
                <div className="mt-4">
                    <p className="text-xs text-gray-400 mb-2 uppercase tracking-wide">Tickets</p>
                    <div className="space-y-1">
                        {booking.tickets.map((t) => (
                            <div key={t.id} className="flex justify-between text-sm">
                                <span className="text-gray-600">
                                    {t.ticket_type_name ?? "Ticket"} × {t.quantity}
                                </span>
                                <span className="text-gray-700">
                                    ₹{Number(t.subtotal).toFixed(2)}
                                </span>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {/* View QR Code Action */}
            <div className="mt-5 pt-4 border-t border-gray-100 flex gap-3">
                <button
                    onClick={() => onShowTickets(booking)}
                    className="w-full flex items-center justify-center gap-2 py-2.5 px-4 bg-black text-white text-sm font-medium rounded-lg hover:bg-gray-800 transition-colors shadow-sm"
                >
                    <span>📱</span>
                    <span>View QR Ticket / Pass</span>
                </button>
            </div>
        </div>
    );
}

// ─── My Tickets Page ──────────────────────────────────────────────────────────

export default function MyTicketsPage() {
    const { user } = useAuth();
    const [bookings, setBookings] = useState<Booking[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [selectedBooking, setSelectedBooking] = useState<Booking | null>(null);

    useEffect(() => {
        if (!user?.id) return;

        const fetchBookings = async () => {
            try {
                setLoading(true);
                const res = await bookingAPI.getUserBookings(user.id);
                // Handle different possible response shapes
                const data =
                    res.data.data?.bookings ??
                    res.data.data ??
                    res.data ??
                    [];
                setBookings(Array.isArray(data) ? data : []);
            } catch {
                setError("Failed to load your bookings. Please try again.");
            } finally {
                setLoading(false);
            }
        };

        fetchBookings();
    }, [user?.id]);

    return (
        <Layout>
            <div className="mx-auto px-4 sm:px-6 lg:px-8 py-10 max-w-3xl">
                <div className="flex justify-between items-center mb-8">
                    <h1 className="font-bold text-3xl text-gray-900">My Tickets</h1>
                    <Link to="/events" className="btn-secondary text-sm">
                        Browse Events
                    </Link>
                </div>

                {loading && (
                    <div className="flex justify-center py-20">
                        <div className="animate-spin rounded-full border-b-2 border-black w-10 h-10" />
                    </div>
                )}

                {error && (
                    <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                        {error}
                    </div>
                )}

                {!loading && !error && bookings.length === 0 && (
                    <div className="text-center py-20 space-y-4">
                        <div className="text-6xl">🎟️</div>
                        <h2 className="font-semibold text-xl text-gray-700">
                            No bookings yet
                        </h2>
                        <p className="text-gray-500">
                            When you book tickets for events, they'll appear here.
                        </p>
                        <Link to="/events" className="btn-primary inline-block mt-4">
                            Browse Events
                        </Link>
                    </div>
                )}

                {!loading && bookings.length > 0 && (
                    <div className="space-y-4">
                        {bookings.map((booking) => (
                            <BookingCard
                                key={booking.id}
                                booking={booking}
                                onShowTickets={(b) => setSelectedBooking(b)}
                            />
                        ))}
                    </div>
                )}

                {/* QR Code Pass Modal */}
                {selectedBooking && (
                    <TicketModal
                        booking={selectedBooking}
                        onClose={() => setSelectedBooking(null)}
                    />
                )}
            </div>
        </Layout>
    );
}
