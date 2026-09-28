import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Layout from "../components/layout";
import { bookingAPI, type CreateBookingRequest, eventAPI } from "../lib/api";
import { useAuth } from "../hooks/useAuth";

// ─── Types ────────────────────────────────────────────────────────────────────

interface TicketSelection {
    ticket_type_id: string;
    name: string;
    price: number;
    quantity: number;
    subtotal: number;
}

interface AvailabilityItem {
    id: string;
    name: string;
    price: number;
    total_available: number;
    remaining: number;
}

// ─── Step 1: Ticket Selector ──────────────────────────────────────────────────

function TicketSelector({
    availability,
    selections,
    onSelectionChange,
}: {
    availability: AvailabilityItem[];
    selections: Record<string, number>;
    onSelectionChange: (id: string, qty: number) => void;
}) {
    return (
        <div className="space-y-4">
            <h2 className="font-semibold text-xl">Select Tickets</h2>
            {availability.map((ticket) => (
                <div
                    key={ticket.id}
                    className="flex justify-between items-center p-5 border border-gray-200 rounded-lg"
                >
                    <div>
                        <p className="font-semibold text-gray-900">{ticket.name}</p>
                        <p className="text-gray-500 text-sm">
                            ₹{ticket.price.toFixed(2)} per ticket
                        </p>
                        <p className="text-gray-400 text-xs">
                            {ticket.remaining} remaining
                        </p>
                    </div>
                    <div className="flex items-center space-x-3">
                        <button
                            type="button"
                            disabled={(selections[ticket.id] ?? 0) === 0}
                            onClick={() =>
                                onSelectionChange(
                                    ticket.id,
                                    Math.max(0, (selections[ticket.id] ?? 0) - 1),
                                )
                            }
                            className="flex justify-center items-center bg-gray-100 hover:bg-gray-200 disabled:opacity-40 rounded-full w-8 h-8 font-bold text-xl transition"
                        >
                            −
                        </button>
                        <span className="w-6 font-medium text-center">
                            {selections[ticket.id] ?? 0}
                        </span>
                        <button
                            type="button"
                            disabled={(selections[ticket.id] ?? 0) >= ticket.remaining}
                            onClick={() =>
                                onSelectionChange(
                                    ticket.id,
                                    Math.min(ticket.remaining, (selections[ticket.id] ?? 0) + 1),
                                )
                            }
                            className="flex justify-center items-center bg-gray-100 hover:bg-gray-200 disabled:opacity-40 rounded-full w-8 h-8 font-bold text-xl transition"
                        >
                            +
                        </button>
                    </div>
                </div>
            ))}
        </div>
    );
}

// ─── Step 2: Booking Summary ──────────────────────────────────────────────────

function BookingSummary({
    eventTitle,
    selections,
    totalAmount,
    couponCode,
    onCouponChange,
    onConfirm,
    onBack,
    loading,
}: {
    eventTitle: string;
    selections: TicketSelection[];
    totalAmount: number;
    couponCode: string;
    onCouponChange: (code: string) => void;
    onConfirm: () => void;
    onBack: () => void;
    loading: boolean;
}) {
    return (
        <div className="space-y-6">
            <h2 className="font-semibold text-xl">Review Booking</h2>
            <div className="bg-gray-50 p-5 rounded-lg">
                <p className="font-semibold text-lg mb-4">{eventTitle}</p>
                <div className="space-y-2 divide-y divide-gray-200">
                    {selections.map((s) => (
                        <div key={s.ticket_type_id} className="flex justify-between py-2">
                            <span className="text-gray-700">
                                {s.name} × {s.quantity}
                            </span>
                            <span className="font-medium">₹{s.subtotal.toFixed(2)}</span>
                        </div>
                    ))}
                </div>
                <div className="flex justify-between items-center pt-4 border-gray-300 border-t font-bold text-lg">
                    <span>Total</span>
                    <span>₹{totalAmount.toFixed(2)}</span>
                </div>
            </div>

            {/* Coupon field */}
            <div>
                <label className="block mb-1 text-gray-700 text-sm font-medium">
                    Coupon Code (optional)
                </label>
                <input
                    type="text"
                    value={couponCode}
                    onChange={(e) => onCouponChange(e.target.value.toUpperCase())}
                    placeholder="Enter coupon code"
                    className="input-field"
                />
            </div>

            <div className="flex gap-3">
                <button type="button" onClick={onBack} className="btn-secondary flex-1">
                    ← Back
                </button>
                <button
                    type="button"
                    onClick={onConfirm}
                    disabled={loading}
                    className="btn-primary flex-1 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {loading ? "Creating booking…" : "Confirm Booking"}
                </button>
            </div>
        </div>
    );
}

// ─── Step 3: Confirmation ─────────────────────────────────────────────────────

function BookingConfirmation({
    bookingId,
    eventTitle,
    finalAmount,
    discount,
}: {
    bookingId: string;
    eventTitle: string;
    finalAmount: number;
    discount: number;
}) {
    const navigate = useNavigate();

    return (
        <div className="py-10 text-center space-y-6">
            <div className="flex justify-center items-center bg-green-100 mx-auto rounded-full w-20 h-20">
                <svg className="w-10 h-10 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
            </div>
            <div>
                <h2 className="font-bold text-2xl text-gray-900">Booking Confirmed!</h2>
                <p className="mt-2 text-gray-600">Your booking for <strong>{eventTitle}</strong> is confirmed.</p>
            </div>
            <div className="bg-gray-50 p-5 rounded-lg text-left space-y-2">
                <div className="flex justify-between">
                    <span className="text-gray-600">Booking ID</span>
                    <span className="font-mono text-sm font-medium">{bookingId}</span>
                </div>
                {discount > 0 && (
                    <div className="flex justify-between text-green-600">
                        <span>Discount Applied</span>
                        <span>−₹{discount.toFixed(2)}</span>
                    </div>
                )}
                <div className="flex justify-between font-bold">
                    <span>Amount Paid</span>
                    <span>₹{finalAmount.toFixed(2)}</span>
                </div>
            </div>
            <p className="text-gray-500 text-sm">
                A confirmation email will be sent to you shortly.
            </p>
            <div className="flex gap-3 justify-center">
                <button
                    onClick={() => navigate("/my-tickets")}
                    className="btn-primary"
                >
                    View My Tickets
                </button>
                <button
                    onClick={() => navigate("/events")}
                    className="btn-secondary"
                >
                    Browse More Events
                </button>
            </div>
        </div>
    );
}

// ─── Main Booking Page ────────────────────────────────────────────────────────

export default function BookingPage() {
    const { slug } = useParams<{ slug: string }>();
    const { user } = useAuth();
    const navigate = useNavigate();

    const [step, setStep] = useState<1 | 2 | 3>(1);
    const [eventTitle, setEventTitle] = useState("");
    const [eventId, setEventId] = useState("");
    const [availability, setAvailability] = useState<AvailabilityItem[]>([]);
    const [selections, setSelections] = useState<Record<string, number>>({});
    const [couponCode, setCouponCode] = useState("");
    const [loading, setLoading] = useState(true);
    const [bookingLoading, setBookingLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [confirmedBooking, setConfirmedBooking] = useState<{
        id: string;
        finalAmount: number;
        discount: number;
    } | null>(null);

    // Fetch event details + availability
    useEffect(() => {
        if (!slug) return;

        const fetchEventAndAvailability = async () => {
            try {
                setLoading(true);
                const eventRes = await eventAPI.getEvent(slug);
                const ev = eventRes.data.data?.event ?? eventRes.data;
                setEventTitle(ev.title);
                setEventId(ev.id);

                const availRes = await bookingAPI.checkAvailability(ev.id);
                const items: AvailabilityItem[] =
                    availRes.data.data?.availability ??
                    availRes.data.data ??
                    availRes.data ??
                    [];
                setAvailability(Array.isArray(items) ? items : []);
            } catch (err) {
                setError("Failed to load event information. Please try again.");
            } finally {
                setLoading(false);
            }
        };

        fetchEventAndAvailability();
    }, [slug]);

    // Compute selected tickets as a list
    const selectedTickets: TicketSelection[] = availability
        .filter((t) => (selections[t.id] ?? 0) > 0)
        .map((t) => ({
            ticket_type_id: t.id,
            name: t.name,
            price: t.price,
            quantity: selections[t.id],
            subtotal: t.price * selections[t.id],
        }));

    const totalAmount = selectedTickets.reduce((sum, t) => sum + t.subtotal, 0);

    const handleConfirmBooking = async () => {
        if (!user) {
            navigate("/auth/login", { state: { from: { pathname: `/events/${slug}/book` } } });
            return;
        }

        setBookingLoading(true);
        setError(null);

        try {
            const payload: CreateBookingRequest = {
                event_id: eventId,
                tickets: selectedTickets.map((t) => ({
                    ticket_type_id: t.ticket_type_id,
                    quantity: t.quantity,
                })),
                coupon_code: couponCode || undefined,
            };

            const res = await bookingAPI.createBooking(payload);
            const data = res.data.data;

            setConfirmedBooking({
                id: data.booking?.id ?? data.id ?? "N/A",
                finalAmount: data.final_amount ?? totalAmount,
                discount: data.discount ?? 0,
            });
            setStep(3);
        } catch (err: unknown) {
            const axErr = err as { response?: { data?: { message?: string } } };
            setError(
                axErr.response?.data?.message ??
                "Failed to create booking. Please try again."
            );
        } finally {
            setBookingLoading(false);
        }
    };

    if (loading) {
        return (
            <Layout>
                <div className="flex justify-center items-center min-h-[60vh]">
                    <div className="animate-spin rounded-full border-b-2 border-black w-10 h-10" />
                </div>
            </Layout>
        );
    }

    return (
        <Layout>
            <div className="mx-auto px-4 py-10 max-w-xl">
                {/* Step indicator */}
                {step < 3 && (
                    <div className="flex items-center mb-8 space-x-2">
                        {[
                            { n: 1, label: "Select" },
                            { n: 2, label: "Review" },
                        ].map(({ n, label }, idx) => (
                            <div key={n} className="flex items-center">
                                <div
                                    className={`flex items-center justify-center w-8 h-8 rounded-full text-sm font-bold ${
                                        step >= n
                                            ? "bg-black text-white"
                                            : "bg-gray-200 text-gray-500"
                                    }`}
                                >
                                    {n}
                                </div>
                                <span
                                    className={`ml-1 text-sm ${step >= n ? "text-black font-medium" : "text-gray-400"}`}
                                >
                                    {label}
                                </span>
                                {idx < 1 && <div className="mx-3 flex-1 h-px bg-gray-200 w-8" />}
                            </div>
                        ))}
                    </div>
                )}

                {error && (
                    <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
                        {error}
                    </div>
                )}

                {step === 1 && (
                    <>
                        <TicketSelector
                            availability={availability}
                            selections={selections}
                            onSelectionChange={(id, qty) =>
                                setSelections((prev) => ({ ...prev, [id]: qty }))
                            }
                        />
                        <div className="mt-6 flex justify-between items-center">
                            <div className="text-gray-700">
                                <span className="font-medium">Total: </span>
                                <span className="font-bold text-xl">₹{totalAmount.toFixed(2)}</span>
                            </div>
                            <button
                                type="button"
                                disabled={selectedTickets.length === 0}
                                onClick={() => setStep(2)}
                                className="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                Continue →
                            </button>
                        </div>
                    </>
                )}

                {step === 2 && (
                    <BookingSummary
                        eventTitle={eventTitle}
                        selections={selectedTickets}
                        totalAmount={totalAmount}
                        couponCode={couponCode}
                        onCouponChange={setCouponCode}
                        onConfirm={handleConfirmBooking}
                        onBack={() => setStep(1)}
                        loading={bookingLoading}
                    />
                )}

                {step === 3 && confirmedBooking && (
                    <BookingConfirmation
                        bookingId={confirmedBooking.id}
                        eventTitle={eventTitle}
                        finalAmount={confirmedBooking.finalAmount}
                        discount={confirmedBooking.discount}
                    />
                )}
            </div>
        </Layout>
    );
}
