import { Link } from "react-router-dom";

export default function Footer() {
  const currentYear = new Date().getFullYear();

  const footerLinks = {
    events: [
      { href: "/events", label: "Browse Events" },
      { href: "/create-event", label: "Create Event" },
      { href: "/organizer", label: "For Organizers" },
      { href: "/pricing", label: "Pricing" },
    ],
  };

  return (
    <footer className="bg-white border-gray-200 border-t w-full">
      <div className="mx-auto px-4 sm:px-6 lg:px-8 py-12 w-full max-w-7xl">
        <div className="gap-8 grid grid-cols-1 md:grid-cols-2">
          {/* Brand Section */}
          <div className="col-span-1">
            <Link to="/" className="font-bold text-black text-2xl">
              Everato
            </Link>
            <p className="mt-4 max-w-xs text-gray-600">
              Modern event management platform for creating, managing, and
              attending memorable events.
            </p>
          </div>

          <div className="col-span-1 text-start md:text-end">
            <h3 className="font-semibold text-black text-sm uppercase tracking-wider">
              Events
            </h3>
            <ul className="space-y-3 mt-4">
              {footerLinks.events.map((link) => (
                <li key={link.href}>
                  <Link
                    to={link.href}
                    className="text-gray-600 hover:text-black transition-colors duration-200"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* Bottom Section */}
        <div className="mt-12 pt-8 border-gray-200 border-t">
          <div className="flex md:flex-row flex-col justify-between items-center">
            <p className="text-gray-600">
              © {currentYear} Everato. All rights reserved.
            </p>
            <div className="flex space-x-6 mt-4 md:mt-0">
              <Link
                to="/privacy"
                className="text-gray-600 hover:text-black transition-colors duration-200"
              >
                Privacy
              </Link>
              <Link
                to="/terms"
                className="text-gray-600 hover:text-black transition-colors duration-200"
              >
                Terms
              </Link>
              <Link
                to="/cookies"
                className="text-gray-600 hover:text-black transition-colors duration-200"
              >
                Cookies
              </Link>
              <Link
                to="/admin"
                className="text-gray-600 hover:text-black transition-colors duration-200"
              >
                Admin
              </Link>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
