import { AuthProvider } from "./contexts/AuthContext";
import AppRoutes from "./routes";

function App() {
    return (
        <AuthProvider>
            <div className="min-h-screen bg-white">
                <AppRoutes />
            </div>
        </AuthProvider>
    );
}

export default App;
