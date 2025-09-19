import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { type Event, eventAPI } from "../lib/api";
import Layout from "../components/layout";
import { FaSearch } from "react-icons/fa";
import { FaLocationPin } from "react-icons/fa6";

interface FilterOptions {
    sortBy: string;
}

export default function EventsPage() {
    const [events, setEvents] = useState<Event[]>([]);
    const [loading, setLoading] = useState(true);
    const [currentPage, setCurrentPage] = useState(1);
    const [totalEvents, setTotalEvents] = useState(0);
    const [filters, setFilters] = useState<FilterOptions>({
        sortBy: "date",
    });
    const [searchParams, setSearchParams] = useSearchParams();
    const [searchQuery, setSearchQuery] = useState(searchParams.get("q") || "");

    const eventsPerPage = 6;

    useEffect(() => {
        const fetchEvents = async () => {
            try {
                setLoading(true);

                // Calculate offset for pagination
                const offset = (currentPage - 1) * eventsPerPage;

                // Fetch events from API
                const response = await eventAPI.getAllEvents(
                    eventsPerPage,
                    offset,
                );

                console.log("API Response:", response.data); // Debug log
                console.log(
                    "Current page:",
                    currentPage,
                    "Offset:",
                    offset,
                    "Limit:",
                    eventsPerPage,
                );

                if (response.data && response.data.data) {
                    const fetchedEvents: Event[] =
                        Array.isArray(response.data.data)
                            ? response.data.data
                            : [];

                    // Get total count from API response
                    const totalFromAPI =
                        response.data.pagination?.total_count || 0;

                    // Apply client-side filters for search query
                    let filteredEvents = fetchedEvents;

                    if (searchQuery) {
                        filteredEvents = filteredEvents.filter((event) =>
                            event.title.toLowerCase().includes(
                                searchQuery.toLowerCase(),
                            ) ||
                            event.description.toLowerCase().includes(
                                searchQuery.toLowerCase(),
                            )
                        );
                        // For search, use filtered count since we're searching on client side
                        setTotalEvents(filteredEvents.length);
                    } else {
                        // For normal pagination, use server count
                        setTotalEvents(totalFromAPI);
                        console.log("Using server total count:", totalFromAPI);
                    }

                    // Sort events
                    if (filters.sortBy === "title") {
                        filteredEvents.sort((a, b) =>
                            a.title.localeCompare(b.title)
                        );
                    } else {
                        filteredEvents.sort((a, b) =>
                            new Date(a.start_time).getTime() -
                            new Date(b.start_time).getTime()
                        );
                    }

                    setEvents(filteredEvents);
                } else {
                    setEvents([]);
                    setTotalEvents(0);
                }
            } catch (error) {
                console.error("Error fetching events:", error);
                setEvents([]);
                setTotalEvents(0);
            } finally {
                setLoading(false);
            }
        };

        fetchEvents();
    }, [filters, searchQuery, currentPage, eventsPerPage]);

    const handleFilterChange = (key: keyof FilterOptions, value: string) => {
        // Reset to page 1 when filters change
        setCurrentPage(1);
        setFilters((prev) => ({ ...prev, [key]: value }));
    };

    const handleSearchSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        // Reset to page 1 when searching
        setCurrentPage(1);
        if (searchQuery) {
            setSearchParams({ q: searchQuery });
        } else {
            setSearchParams({});
        }
    };

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
            <section className="bg-black py-16 text-white">
                <div className="mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl">
                    <div className="text-center">
                        <h1 className="mb-6 font-bold text-4xl md:text-6xl">
                            Discover Events
                        </h1>
                        <p className="mx-auto mb-8 max-w-2xl text-gray-300 text-xl">
                            Find amazing events happening around you. From
                            conferences to concerts, workshops to festivals -
                            there's something for everyone.
                        </p>

                        {/* Search Bar */}
                        <form
                            onSubmit={handleSearchSubmit}
                            className="mx-auto mt-8 max-w-3xl"
                        >
                            <div className="flex bg-white/10 shadow-lg backdrop-blur-sm border-2 border-white/20 overflow-hidden">
                                <input
                                    type="text"
                                    placeholder="Search events by title or description..."
                                    className="flex-1 bg-transparent px-6 py-5 focus:outline-none text-white text-lg placeholder-gray-300"
                                    value={searchQuery}
                                    onChange={(e) =>
                                        setSearchQuery(e.target.value)}
                                />
                                <button
                                    type="submit"
                                    className="bg-white hover:bg-gray-100 px-8 py-5 font-semibold text-black text-lg transition-colors duration-200"
                                >
                                    <span className="flex justify-between items-center gap-2">
                                        <FaSearch />
                                        <span>
                                            Search
                                        </span>
                                    </span>
                                </button>
                            </div>
                            {searchQuery && (
                                <div className="mt-3 text-center">
                                    <button
                                        type="button"
                                        onClick={() => {
                                            setSearchQuery("");
                                            setSearchParams({});
                                        }}
                                        className="text-white/80 hover:text-white text-sm underline"
                                    >
                                        Clear search
                                    </button>
                                </div>
                            )}
                        </form>
                    </div>
                </div>
            </section>

            {/* Filters and Events */}
            <section className="py-12">
                <div className="mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl">
                    {/* Filters */}
                    <div className="bg-gray-50 mb-8 card">
                        <div className="flex md:flex-row flex-col md:justify-between md:items-center gap-4">
                            <div className="flex-1">
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Sort By
                                </label>
                                <select
                                    className="px-3 py-2 border border-gray-300 focus:border-black focus:outline-none w-full"
                                    value={filters.sortBy}
                                    onChange={(e) =>
                                        handleFilterChange(
                                            "sortBy",
                                            e.target.value,
                                        )}
                                >
                                    <option value="date">Date</option>
                                    <option value="title">Title</option>
                                </select>
                            </div>
                            <div className="flex-1 self-end">
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Results per page: {eventsPerPage}
                                </label>
                                <div className="text-gray-600 text-sm">
                                    Showing {events.length} of {totalEvents}
                                    {" "}
                                    events
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Results Header */}
                    <div className="flex justify-between items-center mb-8">
                        <h2 className="font-bold text-black text-2xl">
                            {loading
                                ? "Loading..."
                                : `${events.length} Events Found`}
                        </h2>
                        <Link
                            to="/create-event"
                            className="btn-primary"
                        >
                            Create Event
                        </Link>
                    </div>

                    {/* Events Grid */}
                    {loading
                        ? (
                            <div className="events-grid">
                                {[1, 2, 3, 4, 5, 6].map((i) => (
                                    <div key={i} className="event-card loading">
                                        <div className="event-image">
                                            <div className="bg-gray-200 h-full">
                                            </div>
                                        </div>
                                        <div className="event-meta">
                                            <div className="bg-gray-200 mb-2 h-4">
                                            </div>
                                            <div className="bg-gray-200 mb-2 w-3/4 h-4">
                                            </div>
                                        </div>
                                        <div className="event-content">
                                            <div className="bg-gray-200 w-1/2 h-4">
                                            </div>
                                        </div>
                                        <div className="event-location"></div>
                                        <div className="event-footer"></div>
                                    </div>
                                ))}
                            </div>
                        )
                        : events.length === 0
                        ? (
                            <div className="py-12 text-center">
                                <h3 className="mb-4 font-semibold text-black text-xl">
                                    No events found
                                </h3>
                                <p className="mb-6 text-gray-600">
                                    Try adjusting your search criteria or
                                    explore different categories.
                                </p>
                                <Link
                                    to="/create-event"
                                    className="btn-primary"
                                >
                                    Create the First Event
                                </Link>
                            </div>
                        )
                        : (
                            <div className="events-grid">
                                {events.map((event) => (
                                    <article
                                        key={event.id}
                                        className="group cursor-pointer event-card"
                                    >
                                        <div className="event-image">
                                            <div className="flex justify-center items-center bg-gray-100 h-full">
                                                <span className="text-gray-400">
                                                    Event Image
                                                </span>
                                            </div>
                                        </div>

                                        <div className="event-meta">
                                            <span className="bg-gray-200 px-2 py-1 max-w-max text-gray-700 text-xs uppercase tracking-wide">
                                                {event.status}
                                            </span>
                                            <span className="text-gray-500 text-sm">
                                                {formatDate(event.start_time)}
                                            </span>
                                        </div>

                                        <div className="event-content">
                                            <h3 className="font-semibold text-black group-hover:text-gray-700 text-xl transition-colors duration-200">
                                                {event.title}
                                            </h3>
                                            <p className="text-gray-600">
                                                {event.description}
                                            </p>
                                        </div>

                                        <div className="event-location">
                                            <span className="flex items-center gap-1 text-gray-500 text-sm">
                                                <FaLocationPin />
                                                <span>{event.location}</span>
                                            </span>
                                        </div>

                                        <div className="event-footer">
                                            <Link
                                                to={`/events/${event.slug}`}
                                                className="w-full text-center btn-primary"
                                            >
                                                View Details
                                            </Link>
                                            <span className="text-black text-sm text-center italic">
                                                <span className="text-gray-400">
                                                    {event.available_seats} /
                                                </span>
                                                {event.total_seats} seats
                                            </span>
                                        </div>
                                    </article>
                                ))}
                            </div>
                        )}

                    {/* Pagination */}
                    {events.length > 0 && (
                        <div className="flex justify-center mt-12">
                            <div className="flex space-x-2">
                                <button
                                    className={`px-4 py-2 ${
                                        currentPage === 1
                                            ? "btn-secondary opacity-50 cursor-not-allowed"
                                            : "btn-secondary"
                                    }`}
                                    onClick={() =>
                                        currentPage > 1 &&
                                        setCurrentPage(currentPage - 1)}
                                    disabled={currentPage === 1}
                                >
                                    Previous
                                </button>

                                {/* Show page numbers */}
                                {[...Array(
                                    Math.ceil(totalEvents / eventsPerPage) || 1,
                                )].map((_, index) => {
                                    const pageNumber = index + 1;
                                    return (
                                        <button
                                            key={pageNumber}
                                            className={`px-4 py-2 ${
                                                currentPage === pageNumber
                                                    ? "btn-primary"
                                                    : "btn-secondary"
                                            }`}
                                            onClick={() =>
                                                setCurrentPage(pageNumber)}
                                        >
                                            {pageNumber}
                                        </button>
                                    );
                                })}

                                <button
                                    className={`px-4 py-2 ${
                                        currentPage >=
                                                Math.ceil(
                                                    totalEvents / eventsPerPage,
                                                )
                                            ? "btn-secondary opacity-50 cursor-not-allowed"
                                            : "btn-secondary"
                                    }`}
                                    onClick={() =>
                                        currentPage <
                                            Math.ceil(
                                                totalEvents / eventsPerPage,
                                            ) &&
                                        setCurrentPage(currentPage + 1)}
                                    disabled={currentPage >=
                                        Math.ceil(totalEvents / eventsPerPage)}
                                >
                                    Next
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            </section>
        </Layout>
    );
}
