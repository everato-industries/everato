import { Link } from "react-router-dom";

export default function Navbar() {
  return (
    <nav className="bg-white shadow-sm">
      <div className="container mx-auto px-4 py-3 flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <Link to="/" className="text-xl font-bold text-gray-800">
            Everato
          </Link>
        </div>

        <div className="hidden md:flex items-center space-x-6">
          <Link to="/" className="text-gray-600 hover:text-blue-600">
            Home
          </Link>
          <Link to="/events" className="text-gray-600 hover:text-blue-600">
            Events
          </Link>
          <Link to="/about" className="text-gray-600 hover:text-blue-600">
            About
          </Link>
        </div>

        <div className="flex items-center space-x-3">
          <Link to="/auth/login" className="px-4 py-2 text-sm text-gray-700 hover:text-blue-600">
            Log in
          </Link>
          <Link
            to="/auth/register"
            className="px-4 py-2 text-sm text-white bg-blue-600 rounded hover:bg-blue-700 transition-colors"
          >
            Sign up
          </Link>
        </div>
      </div>
    </nav>
  );
}
