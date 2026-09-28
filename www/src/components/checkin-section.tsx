import { useCallback, useEffect, useState } from "react";
import {
    type AttendanceStats,
    checkinAPI,
    type Event,
    type EventAttendee,
    eventAPI,
} from "../lib/api";

export default function CheckinSection() {
    const [events, setEvents] = useState<Event[]>([]);
    const [selectedEventId, setSelectedEventId] = useState<string>("");
    const [eventsLoading, setEventsLoading] = useState(true);

    // Stats & Attendees
    const [stats, setStats] = useState<AttendanceStats | null>(null);
    const [attendees, setAttendees] = useState<EventAttendee[]>([]);
    const [loadingData, setLoadingData] = useState(false);
    const [searchFilter, setSearchFilter] = useState("");
    const [statusFilter, setStatusFilter] = useState<"ALL" | "CHECKED_IN" | "PENDING">("ALL");

    // Scan input state
    const [qrInput, setQrInput] = useState("");
    const [location, setLocation] = useState("Main Entrance");
    const [deviceInfo, setDeviceInfo] = useState("Admin Dashboard Scanner");
    const [scanning, setScanning] = useState(false);
    const [scanResult, setScanResult] = useState<{
        type: "success" | "warning" | "error";
        title: string;
        message: string;
        ticketNumber?: string;
        checkedInAt?: string;
    } | null>(null);

    // Load events on mount
    useEffect(() => {
        const fetchEvents = async () => {
            try {
                setEventsLoading(true);
                const res = await eventAPI.getAllEvents(100);
                const list = res.data.data?.events ?? res.data.events ?? [];
                setEvents(list);
                if (list.length > 0) {
                    setSelectedEventId(list[0].id);
                }
            } catch (err) {
                console.error("Failed to load events for check-in:", err);
            } finally {
                setEventsLoading(false);
            }
        };

        fetchEvents();
    }, []);

    // Load stats and attendees for selected event
    const refreshEventData = useCallback(async (eventId: string) => {
        if (!eventId) return;
        try {
            setLoadingData(true);
            const [statsRes, attendeesRes] = await Promise.allSettled([
                checkinAPI.getStats(eventId),
                checkinAPI.getAttendees(eventId),
            ]);

            if (statsRes.status === "fulfilled") {
                const s = statsRes.value.data.data ?? statsRes.value.data;
                setStats(s);
            }
            if (attendeesRes.status === "fulfilled") {
                const a = attendeesRes.value.data.data?.attendees ?? attendeesRes.value.data.attendees ?? [];
                setAttendees(Array.isArray(a) ? a : []);
            }
        } catch (err) {
            console.error("Error refreshing check-in data:", err);
        } finally {
            setLoadingData(false);
        }
    }, []);

    useEffect(() => {
        if (selectedEventId) {
            refreshEventData(selectedEventId);
            setScanResult(null);
        }
    }, [selectedEventId, refreshEventData]);

    // Handle Scan / Submit QR
    const handleCheckinSubmit = async (e?: React.FormEvent) => {
        if (e) e.preventDefault();
        const trimmed = qrInput.trim();
        if (!trimmed) return;

        try {
            setScanning(true);
            setScanResult(null);

            const res = await checkinAPI.scan({
                qr_data: trimmed,
                location: location.trim() || undefined,
                device_info: deviceInfo.trim() || undefined,
            });

            const data = res.data.data ?? res.data;
            setScanResult({
                type: "success",
                title: "Check-in Successful! 🎉",
                message: `Ticket verified and attendee admitted.`,
                ticketNumber: data.ticket_number,
                checkedInAt: data.checked_in_at,
            });

            // Clear input for next scan
            setQrInput("");
            // Refresh stats and attendee list
            refreshEventData(selectedEventId);
        } catch (err: any) {
            const errData = err.response?.data;
            const status = err.response?.status;

            if (status === 409) {
                // Already checked in
                setScanResult({
                    type: "warning",
                    title: "Ticket Already Checked In ⚠️",
                    message: "This ticket has already been used for entry.",
                    ticketNumber: errData?.ticket_number,
                    checkedInAt: errData?.checked_in_at,
                });
            } else if (status === 422 || status === 400) {
                // Invalid signature / QR
                setScanResult({
                    type: "error",
                    title: "Invalid or Tampered Ticket ❌",
                    message: errData?.message ?? "QR code signature verification failed.",
                });
            } else if (status === 404) {
                setScanResult({
                    type: "error",
                    title: "Ticket Not Found ❌",
                    message: "No ticket record found matching this QR code.",
                });
            } else {
                setScanResult({
                    type: "error",
                    title: "Check-in Failed",
                    message: errData?.error ?? "An unexpected server error occurred.",
                });
            }
        } finally {
            setScanning(false);
        }
    };

    // Filter attendees
    const filteredAttendees = attendees.filter((a) => {
        const matchesSearch =
            searchFilter === "" ||
            a.ticket_number.toLowerCase().includes(searchFilter.toLowerCase()) ||
            a.user_email.toLowerCase().includes(searchFilter.toLowerCase()) ||
            `${a.first_name} ${a.last_name}`.toLowerCase().includes(searchFilter.toLowerCase()) ||
            a.ticket_type_name.toLowerCase().includes(searchFilter.toLowerCase());

        const matchesStatus =
            statusFilter === "ALL" ||
            (statusFilter === "CHECKED_IN" && a.is_checked_in) ||
            (statusFilter === "PENDING" && !a.is_checked_in);

        return matchesSearch && matchesStatus;
    });

    const totalTickets = Number(stats?.total_tickets ?? 0);
    const checkedInCount = Number(stats?.checked_in_count ?? 0);
    const notCheckedInCount = Number(stats?.not_checked_in_count ?? 0);
    const checkinPercentage =
        totalTickets > 0 ? Math.round((checkedInCount / totalTickets) * 100) : 0;

    return (
        <div className="space-y-8">
            {/* Header & Event Selector */}
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-6 rounded-2xl border border-gray-200 shadow-sm">
                <div>
                    <h2 className="text-2xl font-bold text-gray-900">
                        Event Check-in & Scanner
                    </h2>
                    <p className="text-sm text-gray-500 mt-1">
                        Scan attendee QR codes and track gate entries in real time.
                    </p>
                </div>

                <div className="flex items-center gap-3">
                    <label htmlFor="event-select" className="text-sm font-medium text-gray-700 whitespace-nowrap">
                        Select Event:
                    </label>
                    {eventsLoading ? (
                        <div className="animate-pulse bg-gray-200 h-10 w-48 rounded-lg" />
                    ) : (
                        <select
                            id="event-select"
                            value={selectedEventId}
                            onChange={(e) => setSelectedEventId(e.target.value)}
                            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-black focus:border-black block p-2.5 min-w-[200px]"
                        >
                            {events.map((evt) => (
                                <option key={evt.id} value={evt.id}>
                                    {evt.title}
                                </option>
                            ))}
                        </select>
                    )}
                    <button
                        onClick={() => refreshEventData(selectedEventId)}
                        disabled={loadingData}
                        className="p-2.5 text-gray-600 hover:text-black border border-gray-300 rounded-lg bg-gray-50 hover:bg-gray-100 transition-colors"
                        title="Refresh Data"
                    >
                        🔄
                    </button>
                </div>
            </div>

            {/* Attendance Metrics */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
                <div className="bg-white border border-gray-200 rounded-xl p-5 shadow-sm">
                    <div className="text-sm font-medium text-gray-500">Total Registered</div>
                    <div className="text-3xl font-extrabold text-gray-900 mt-2">
                        {loadingData ? "…" : totalTickets}
                    </div>
                    <div className="text-xs text-gray-400 mt-1">Confirmed tickets issued</div>
                </div>

                <div className="bg-white border border-emerald-200 rounded-xl p-5 shadow-sm">
                    <div className="flex justify-between items-center">
                        <span className="text-sm font-medium text-emerald-800">Checked In</span>
                        <span className="text-xs font-bold bg-emerald-100 text-emerald-800 px-2 py-0.5 rounded-full">
                            {checkinPercentage}%
                        </span>
                    </div>
                    <div className="text-3xl font-extrabold text-emerald-600 mt-2">
                        {loadingData ? "…" : checkedInCount}
                    </div>
                    <div className="w-full bg-gray-100 h-2 rounded-full mt-3 overflow-hidden">
                        <div
                            className="bg-emerald-500 h-full rounded-full transition-all duration-500"
                            style={{ width: `${checkinPercentage}%` }}
                        />
                    </div>
                </div>

                <div className="bg-white border border-gray-200 rounded-xl p-5 shadow-sm">
                    <div className="text-sm font-medium text-gray-500">Remaining to Check In</div>
                    <div className="text-3xl font-extrabold text-amber-600 mt-2">
                        {loadingData ? "…" : notCheckedInCount}
                    </div>
                    <div className="text-xs text-gray-400 mt-1">Expected at entrance</div>
                </div>
            </div>

            {/* QR Scan Station */}
            <div className="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm space-y-6">
                <div className="flex items-center justify-between border-b border-gray-100 pb-4">
                    <div>
                        <h3 className="text-lg font-bold text-gray-900 flex items-center gap-2">
                            <span>📷</span>
                            <span>QR Code Scanner Terminal</span>
                        </h3>
                        <p className="text-xs text-gray-500 mt-0.5">
                            Scan with a barcode/QR reader device or paste the signed QR data string.
                        </p>
                    </div>
                    <span className="text-xs bg-black text-white px-2.5 py-1 rounded-full font-mono">
                        HMAC SHA-256 Verified
                    </span>
                </div>

                {/* Scan Result Alert */}
                {scanResult && (
                    <div
                        className={`p-4 rounded-xl border flex items-start justify-between gap-4 animate-in fade-in slide-in-from-top-2 duration-200 ${
                            scanResult.type === "success"
                                ? "bg-emerald-50 border-emerald-200 text-emerald-900"
                                : scanResult.type === "warning"
                                ? "bg-amber-50 border-amber-200 text-amber-900"
                                : "bg-red-50 border-red-200 text-red-900"
                        }`}
                    >
                        <div>
                            <h4 className="font-bold text-base">{scanResult.title}</h4>
                            <p className="text-sm mt-1">{scanResult.message}</p>
                            {scanResult.ticketNumber && (
                                <p className="text-xs font-mono mt-1 font-semibold">
                                    Ticket #{scanResult.ticketNumber}
                                    {scanResult.checkedInAt && (
                                        <span className="ml-2 font-normal opacity-80">
                                            (Admitted at {new Date(scanResult.checkedInAt).toLocaleTimeString("en-IN")})
                                        </span>
                                    )}
                                </p>
                            )}
                        </div>
                        <button
                            onClick={() => setScanResult(null)}
                            className="text-gray-400 hover:text-gray-700 p-1"
                        >
                            ✕
                        </button>
                    </div>
                )}

                {/* Scanner Form */}
                <form onSubmit={handleCheckinSubmit} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            QR Payload / Scan Data <span className="text-red-500">*</span>
                        </label>
                        <div className="relative">
                            <input
                                type="text"
                                value={qrInput}
                                onChange={(e) => setQrInput(e.target.value)}
                                placeholder="Paste or scan QR code string here (e.g. {&quot;tid&quot;:&quot;...&quot;,&quot;sig&quot;:&quot;...&quot;})"
                                className="w-full bg-gray-50 border border-gray-300 rounded-xl px-4 py-3 text-sm font-mono focus:bg-white focus:ring-2 focus:ring-black focus:border-black outline-none transition-all pr-24"
                                autoFocus
                            />
                            {qrInput && (
                                <button
                                    type="button"
                                    onClick={() => setQrInput("")}
                                    className="absolute right-3 top-3 text-xs text-gray-400 hover:text-gray-600 bg-gray-200 px-2 py-1 rounded"
                                >
                                    Clear
                                </button>
                            )}
                        </div>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-xs font-medium text-gray-600 mb-1">
                                Gate / Check-in Location
                            </label>
                            <input
                                type="text"
                                value={location}
                                onChange={(e) => setLocation(e.target.value)}
                                placeholder="e.g. Main Entrance Gate 2"
                                className="w-full bg-gray-50 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-1 focus:ring-black outline-none"
                            />
                        </div>
                        <div>
                            <label className="block text-xs font-medium text-gray-600 mb-1">
                                Scanner Device Info
                            </label>
                            <input
                                type="text"
                                value={deviceInfo}
                                onChange={(e) => setDeviceInfo(e.target.value)}
                                placeholder="e.g. Desk-Tablet-01"
                                className="w-full bg-gray-50 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-1 focus:ring-black outline-none"
                            />
                        </div>
                    </div>

                    <button
                        type="submit"
                        disabled={scanning || !qrInput.trim()}
                        className="w-full sm:w-auto px-6 py-2.5 bg-black text-white text-sm font-semibold rounded-xl hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center gap-2"
                    >
                        {scanning ? (
                            <>
                                <div className="animate-spin rounded-full border-2 border-white border-t-transparent w-4 h-4" />
                                <span>Verifying & Admitting…</span>
                            </>
                        ) : (
                            <>
                                <span>✓</span>
                                <span>Admit & Validate Ticket</span>
                            </>
                        )}
                    </button>
                </form>
            </div>

            {/* Attendee Roster Table */}
            <div className="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm space-y-4">
                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                    <div>
                        <h3 className="text-lg font-bold text-gray-900">
                            Attendee Roster ({filteredAttendees.length})
                        </h3>
                        <p className="text-xs text-gray-500 mt-0.5">
                            Real-time list of all booked ticket holders and their check-in status.
                        </p>
                    </div>

                    {/* Filter controls */}
                    <div className="flex flex-wrap items-center gap-3">
                        <input
                            type="text"
                            value={searchFilter}
                            onChange={(e) => setSearchFilter(e.target.value)}
                            placeholder="Search name, email, ticket…"
                            className="bg-gray-50 border border-gray-300 text-xs rounded-lg px-3 py-2 w-48 focus:ring-1 focus:ring-black outline-none"
                        />
                        <div className="flex border border-gray-200 rounded-lg overflow-hidden text-xs">
                            <button
                                onClick={() => setStatusFilter("ALL")}
                                className={`px-3 py-2 ${
                                    statusFilter === "ALL"
                                        ? "bg-black text-white font-medium"
                                        : "bg-white text-gray-600 hover:bg-gray-50"
                                }`}
                            >
                                All
                            </button>
                            <button
                                onClick={() => setStatusFilter("CHECKED_IN")}
                                className={`px-3 py-2 ${
                                    statusFilter === "CHECKED_IN"
                                        ? "bg-black text-white font-medium"
                                        : "bg-white text-gray-600 hover:bg-gray-50"
                                }`}
                            >
                                Checked In
                            </button>
                            <button
                                onClick={() => setStatusFilter("PENDING")}
                                className={`px-3 py-2 ${
                                    statusFilter === "PENDING"
                                        ? "bg-black text-white font-medium"
                                        : "bg-white text-gray-600 hover:bg-gray-50"
                                }`}
                            >
                                Pending
                            </button>
                        </div>
                    </div>
                </div>

                {/* Table */}
                <div className="overflow-x-auto border border-gray-200 rounded-xl">
                    <table className="w-full text-left text-sm text-gray-600">
                        <thead className="bg-gray-50 text-xs text-gray-500 uppercase tracking-wider border-b border-gray-200">
                            <tr>
                                <th className="px-4 py-3">Attendee</th>
                                <th className="px-4 py-3">Ticket Type</th>
                                <th className="px-4 py-3">Ticket #</th>
                                <th className="px-4 py-3">Status</th>
                                <th className="px-4 py-3">Checked In At</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {loadingData ? (
                                <tr>
                                    <td colSpan={5} className="py-8 text-center text-gray-400">
                                        Loading attendee list…
                                    </td>
                                </tr>
                            ) : filteredAttendees.length === 0 ? (
                                <tr>
                                    <td colSpan={5} className="py-8 text-center text-gray-400">
                                        No attendees found matching filter.
                                    </td>
                                </tr>
                            ) : (
                                filteredAttendees.map((att) => (
                                    <tr key={att.ticket_id} className="hover:bg-gray-50/50">
                                        <td className="px-4 py-3">
                                            <div className="font-medium text-gray-900">
                                                {att.first_name} {att.last_name}
                                            </div>
                                            <div className="text-xs text-gray-400">{att.user_email}</div>
                                        </td>
                                        <td className="px-4 py-3">
                                            <span className="text-xs bg-gray-100 text-gray-700 px-2 py-1 rounded">
                                                {att.ticket_type_name}
                                            </span>
                                        </td>
                                        <td className="px-4 py-3 font-mono text-xs text-gray-700">
                                            {att.ticket_number}
                                        </td>
                                        <td className="px-4 py-3">
                                            <span
                                                className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-semibold ${
                                                    att.is_checked_in
                                                        ? "bg-emerald-100 text-emerald-800"
                                                        : "bg-gray-100 text-gray-500"
                                                }`}
                                            >
                                                {att.is_checked_in ? "✓ Checked In" : "Pending"}
                                            </span>
                                        </td>
                                        <td className="px-4 py-3 text-xs text-gray-500">
                                            {att.checked_in_at
                                                ? new Date(att.checked_in_at).toLocaleTimeString("en-IN", {
                                                      hour: "2-digit",
                                                      minute: "2-digit",
                                                      second: "2-digit",
                                                  })
                                                : "—"}
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}
