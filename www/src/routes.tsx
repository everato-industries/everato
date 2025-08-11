import { Route, Routes } from "react-router-dom";
import LoginPage from "./pages/auth/login";
import HomePage from "./pages/home";

export default function AppRoutes() {
    return (
        <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/auth/login" element={<LoginPage />} />

            <Route
                path="*"
                element={
                    <div className="text-center text-red-500">
                        404 Not Found
                    </div>
                }
            />
        </Routes>
    );
}
