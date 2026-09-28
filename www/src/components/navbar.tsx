import { useEffect, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import api from "../lib/api";
import { useAuth } from "../hooks/useAuth";

export default function Navbar() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [serverInfo, setServerInfo] = useState<Record<string, string> | null>(
    null,
  );
  const location = useLocation();
  const navigate = useNavigate();
  const { isAuthenticated, user, logout } = useAuth();

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

  const handleLogout = () => {
    logout();
    navigate("/");
  };

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

            {/* Auth buttons */}
            {isAuthenticated ? (
              <div className="flex items-center space-x-4">
                <Link
                  to="/my-tickets"
                  className="text-gray-600 hover:text-black font-medium transition-colors duration-200"
                >
                  My Tickets
                </Link>
                <span className="text-gray-400 text-sm">
                  {user?.firstName} {user?.lastName}
                </span>
                <button
                  onClick={handleLogout}
                  className="btn-secondary text-sm"
                >
                  Logout
                </button>
              </div>
            ) : (
              <div className="flex items-center space-x-3">
                <Link to="/auth/login" className="text-gray-600 hover:text-black font-medium transition-colors duration-200">
                  Sign in
                </Link>
                <Link to="/auth/register" className="btn-primary text-sm">
                  Sign up
                </Link>
              </div>
            )}
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
              {isAuthenticated ? (
                <>
                  <Link
                    to="/my-tickets"
                    className="text-gray-600 hover:text-black font-medium"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    My Tickets
                  </Link>
                  <button
                    onClick={() => { handleLogout(); setIsMenuOpen(false); }}
                    className="text-left text-gray-600 hover:text-black font-medium"
                  >
                    Logout ({user?.firstName})
                  </button>
                </>
              ) : (
                <>
                  <Link
                    to="/auth/login"
                    className="text-gray-600 hover:text-black font-medium"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    Sign in
                  </Link>
                  <Link
                    to="/auth/register"
                    className="text-gray-600 hover:text-black font-medium"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    Sign up
                  </Link>
                </>
              )}
            </div>
          </div>
        )}
      </div>
    </nav>
  );
}
