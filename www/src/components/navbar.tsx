import { useEffect, useState } from "react";
import { Link, useLocation } from "react-router-dom";
import api from "../lib/api";

export default function Navbar() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [serverInfo, setServerInfo] = useState<Record<string, string> | null>(
    null,
  );
  const location = useLocation();

  const navLinks = [
    { href: "/", label: "Home" },
    { href: "/events", label: "Events" },
  ];

  useEffect(() => {
    api.get("/dashboard/info")
      .then((response) => {
        setServerInfo(response.data);
      })
      .catch((error) => {
        console.error("Failed to fetch server info:", error);
      });
  }, []);

  return (
    <nav className="top-0 z-40 sticky bg-white border-gray-200 border-b">
      <div className="mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link to="/" className="h-[3rem] font-bold text-black text-2xl">
            <div className="flex items-center space-x-2 h-[3rem]">
              <span>{serverInfo ? serverInfo.org_name : "Everato"}</span>
              <span className="self-end font-thin text-gray-400 text-sm italic">
                {serverInfo ? "~ Powered by, Everato" : ""}
              </span>
            </div>
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-8">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                to={link.href}
                className={`font-medium transition-colors duration-200 ${
                  location.pathname === link.href
                    ? "text-black border-b-2 border-black pb-1"
                    : "text-gray-600 hover:text-black"
                }`}
              >
                {link.label}
              </Link>
            ))}
          </div>

          {/* Mobile menu button */}
          <div className="md:hidden">
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="text-gray-600 hover:text-black transition-colors duration-200"
            >
              <svg
                className="w-6 h-6"
                fill="none"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                {isMenuOpen
                  ? <path d="M6 18L18 6M6 6l12 12" />
                  : <path d="M4 6h16M4 12h16M4 18h16" />}
              </svg>
            </button>
          </div>
        </div>

        {/* Mobile Navigation */}
        {isMenuOpen && (
          <div className="md:hidden py-4 border-gray-200 border-t">
            <div className="flex flex-col space-y-4">
              {navLinks.map((link) => (
                <Link
                  key={link.href}
                  to={link.href}
                  className={`font-medium transition-colors duration-200 ${
                    location.pathname === link.href
                      ? "text-black"
                      : "text-gray-600 hover:text-black"
                  }`}
                  onClick={() => setIsMenuOpen(false)}
                >
                  {link.label}
                </Link>
              ))}
            </div>
          </div>
        )}
      </div>
    </nav>
  );
}
