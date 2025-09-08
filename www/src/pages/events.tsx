import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
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
    category: string;
}

interface FilterOptions {
    category: string;
    location: string;
    priceRange: string;
    dateRange: string;
    sortBy: string;
}

export default function EventsPage() {
    const [events, setEvents] = useState<Event[]>([]);
    const [loading, setLoading] = useState(true);
    const [filters, setFilters] = useState<FilterOptions>({
        category: "",
        location: "",
        priceRange: "",
        dateRange: "",
        sortBy: "date",
    });
    const [searchParams, setSearchParams] = useSearchParams();
    const [searchQuery, setSearchQuery] = useState(searchParams.get("q") || "");

    const categories = [
        "All Categories",
        "Technology",
        "Music",
        "Food & Drink",
        "Business",
        "Sports",
        "Arts & Culture",
        "Health & Wellness",
    ];

    const locations = [
        "All Locations",
        "New York, NY",
        "San Francisco, CA",
        "Los Angeles, CA",
        "Chicago, IL",
        "Austin, TX",
        "Seattle, WA",
    ];

    useEffect(() => {
        const fetchEvents = async () => {
            try {
                setLoading(true);
                // Mock API call - replace with actual API
                setTimeout(() => {
                    const mockEvents: Event[] = [
                        {
                            id: "1",
                            title: "Tech Conference 2025",
                            description:
                                "Annual technology conference featuring the latest in AI and web development.",
                            date: "2025-10-15",
                            location: "San Francisco, CA",
                            price: 299,
                            slug: "tech-conference-2025",
                            category: "Technology",
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
                            category: "Music",
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
                            category: "Food & Drink",
                        },
                        {
                            id: "4",
                            title: "Business Leadership Summit",
                            description:
                                "Learn from industry leaders and network with professionals.",
                            date: "2025-11-05",
                            location: "Chicago, IL",
                            price: 399,
                            slug: "business-leadership-summit",
                            category: "Business",
                        },
                        {
                            id: "5",
                            title: "Art Gallery Opening",
                            description:
                                "Modern art exhibition featuring contemporary artists.",
                            date: "2025-07-22",
                            location: "Los Angeles, CA",
                            price: 75,
                            slug: "art-gallery-opening",
                            category: "Arts & Culture",
                        },
                        {
                            id: "6",
                            title: "Wellness Retreat",
                            description:
                                "Weekend wellness retreat with yoga, meditation, and spa treatments.",
                            date: "2025-09-30",
                            location: "Seattle, WA",
                            price: 249,
                            slug: "wellness-retreat",
                            category: "Health & Wellness",
                        },
                    ];

                    // Apply filters
                    let filteredEvents = mockEvents;

                    if (searchQuery) {
                        filteredEvents = filteredEvents.filter((event) =>
                            event.title.toLowerCase().includes(
                                searchQuery.toLowerCase(),
                            ) ||
                            event.description.toLowerCase().includes(
                                searchQuery.toLowerCase(),
                            )
                        );
                    }

                    if (
                        filters.category &&
                        filters.category !== "All Categories"
                    ) {
                        filteredEvents = filteredEvents.filter((event) =>
                            event.category === filters.category
                        );
                    }

                    if (
                        filters.location && filters.location !== "All Locations"
                    ) {
                        filteredEvents = filteredEvents.filter((event) =>
                            event.location === filters.location
                        );
                    }

                    // Sort events
                    if (filters.sortBy === "price") {
                        filteredEvents.sort((a, b) => a.price - b.price);
                    } else if (filters.sortBy === "title") {
                        filteredEvents.sort((a, b) =>
                            a.title.localeCompare(b.title)
                        );
                    } else {
                        filteredEvents.sort((a, b) =>
                            new Date(a.date).getTime() -
                            new Date(b.date).getTime()
                        );
                    }

                    setEvents(filteredEvents);
                    setLoading(false);
                }, 500);
            } catch (error) {
                console.error("Error fetching events:", error);
                setLoading(false);
            }
        };

        fetchEvents();
    }, [filters, searchQuery]);

    const handleFilterChange = (key: keyof FilterOptions, value: string) => {
        setFilters((prev) => ({ ...prev, [key]: value }));
    };

    const handleSearchSubmit = (e: React.FormEvent) => {
        e.preventDefault();
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
                            className="mx-auto max-w-2xl"
                        >
                            <div className="flex">
                                <input
                                    type="text"
                                    placeholder="Search events..."
                                    className="flex-1 px-6 py-4 focus:outline-none text-black"
                                    value={searchQuery}
                                    onChange={(e) =>
                                        setSearchQuery(e.target.value)}
                                />
                                <button
                                    type="submit"
                                    className="bg-white hover:bg-gray-100 px-8 py-4 font-medium text-black transition-colors duration-200"
                                >
                                    Search
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </section>

            {/* Filters and Events */}
            <section className="py-12">
                <div className="mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl">
                    {/* Filters */}
                    <div className="bg-gray-50 mb-8 p-6">
                        <div className="gap-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5">
                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Category
                                </label>
                                <select
                                    className="px-3 py-2 border border-gray-300 focus:border-black focus:outline-none w-full"
                                    value={filters.category}
                                    onChange={(e) =>
                                        handleFilterChange(
                                            "category",
                                            e.target.value,
                                        )}
                                >
                                    {categories.map((category) => (
                                        <option key={category} value={category}>
                                            {category}
                                        </option>
                                    ))}
                                </select>
                            </div>

                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Location
                                </label>
                                <select
                                    className="px-3 py-2 border border-gray-300 focus:border-black focus:outline-none w-full"
                                    value={filters.location}
                                    onChange={(e) =>
                                        handleFilterChange(
                                            "location",
                                            e.target.value,
                                        )}
                                >
                                    {locations.map((location) => (
                                        <option key={location} value={location}>
                                            {location}
                                        </option>
                                    ))}
                                </select>
                            </div>

                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Price Range
                                </label>
                                <select
                                    className="px-3 py-2 border border-gray-300 focus:border-black focus:outline-none w-full"
                                    value={filters.priceRange}
                                    onChange={(e) =>
                                        handleFilterChange(
                                            "priceRange",
                                            e.target.value,
                                        )}
                                >
                                    <option value="">All Prices</option>
                                    <option value="0-50">$0 - $50</option>
                                    <option value="51-100">$51 - $100</option>
                                    <option value="101-200">$101 - $200</option>
                                    <option value="201+">$201+</option>
                                </select>
                            </div>

                            <div>
                                <label className="block mb-2 font-medium text-black text-sm">
                                    Date Range
                                </label>
                                <select
                                    className="px-3 py-2 border border-gray-300 focus:border-black focus:outline-none w-full"
                                    value={filters.dateRange}
                                    onChange={(e) =>
                                        handleFilterChange(
                                            "dateRange",
                                            e.target.value,
                                        )}
                                >
                                    <option value="">All Dates</option>
                                    <option value="today">Today</option>
                                    <option value="week">This Week</option>
                                    <option value="month">This Month</option>
                                    <option value="quarter">
                                        Next 3 Months
                                    </option>
                                </select>
                            </div>

                            <div>
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
                                    <option value="price">Price</option>
                                    <option value="title">Title</option>
                                </select>
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
                            <div className="gap-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
                                {[1, 2, 3, 4, 5, 6].map((i) => (
                                    <div key={i} className="card loading">
                                        <div className="bg-gray-200 mb-4 h-48">
                                        </div>
                                        <div className="bg-gray-200 mb-2 h-4">
                                        </div>
                                        <div className="bg-gray-200 mb-2 w-3/4 h-4">
                                        </div>
                                        <div className="bg-gray-200 w-1/2 h-4">
                                        </div>
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
                                            <span className="bg-gray-200 px-2 py-1 text-gray-700 text-xs uppercase tracking-wide">
                                                {event.category}
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

                    {/* Pagination - Mock for now */}
                    {events.length > 0 && (
                        <div className="flex justify-center mt-12">
                            <div className="flex space-x-2">
                                <button className="px-4 py-2 btn-secondary">
                                    Previous
                                </button>
                                <button className="px-4 py-2 btn-primary">
                                    1
                                </button>
                                <button className="px-4 py-2 btn-secondary">
                                    2
                                </button>
                                <button className="px-4 py-2 btn-secondary">
                                    3
                                </button>
                                <button className="px-4 py-2 btn-secondary">
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
