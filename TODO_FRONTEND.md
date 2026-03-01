# Frontend TODO - MVP Focused

**Last Updated:** March 1, 2026  
**Based On:** REPORT.md MVP Gap Analysis  
**Reference:** React + TypeScript + TailwindCSS + Vite Stack

---

## 🚨 **CRITICAL PATH - MVP BLOCKERS (P0)**

These features MUST be completed before MVP can ship to clients.

### **CP-1: Authentication Integration** ⏱️ 1-2 days

**Status:** Partially implemented - needs real backend integration  
**Priority:** P0 - CRITICAL  
**Files:** `www/src/pages/auth/`, `www/src/lib/api.ts`, `www/src/hooks/`

**Tasks:**

#### 1. Connect Login Form
- [ ] Update `www/src/pages/auth/login.tsx` to call `POST /api/v1/auth/login`
- [ ] Handle JWT token storage in localStorage
- [ ] Redirect to dashboard on successful login
- [ ] Display error messages for invalid credentials
- [ ] Add loading state during authentication

#### 2. Connect Registration Form
- [ ] Update `www/src/pages/auth/register.tsx` to call `POST /api/v1/auth/register`
- [ ] Add form validation (email format, password strength)
- [ ] Handle registration success → auto-login
- [ ] Display email verification message
- [ ] Show validation errors inline

#### 3. Create Authentication Context
**File:** `www/src/contexts/AuthContext.tsx`

```typescript
interface AuthContextType {
    user: User | null;
    login: (email: string, password: string) => Promise<void>;
    register: (data: RegisterData) => Promise<void>;
    logout: () => void;
    isAuthenticated: boolean;
    loading: boolean;
}

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
    // Implementation
};
```

#### 4. Create useAuth Hook
**File:** `www/src/hooks/useAuth.ts`

```typescript
export function useAuth() {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within AuthProvider");
    }
    return context;
}
```

#### 5. Protected Route Component
**File:** `www/src/components/ProtectedRoute.tsx`

```typescript
export const ProtectedRoute: React.FC<{children: React.ReactNode}> = ({ children }) => {
    const { isAuthenticated, loading } = useAuth();
    
    if (loading) return <LoadingSpinner />;
    if (!isAuthenticated) return <Navigate to="/login" />;
    
    return <>{children}</>;
};
```

#### 6. JWT Token Management
- [ ] Store access token in localStorage
- [ ] Add token to all API requests via axios interceptor
- [ ] Implement token refresh logic
- [ ] Handle token expiry (401 responses)
- [ ] Auto-logout on token expiration

**API Integration:**
```typescript
// www/src/lib/api.ts
api.interceptors.request.use((config) => {
    const token = localStorage.getItem("access_token");
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

api.interceptors.response.use(
    (response) => response,
    async (error) => {
        if (error.response?.status === 401) {
            // Handle token refresh or logout
        }
        return Promise.reject(error);
    }
);
```

**Testing Checklist:**
- [ ] Login with valid credentials → Dashboard
- [ ] Login with invalid credentials → Error message
- [ ] Register new user → Auto-login
- [ ] Access protected route without auth → Redirect to login
- [ ] Logout → Clear tokens and redirect to home
- [ ] Token expiry → Auto-logout

---

### **CP-2: Ticket Booking Flow** ⏱️ 5-7 days

**Status:** 0% Complete  
**Priority:** P0 - CRITICAL BLOCKER  
**Complexity:** Highest - Multi-step flow with payment

**New Files to Create:**

#### 1. Ticket Selection Component
**File:** `www/src/components/booking/TicketSelector.tsx`

```typescript
interface TicketSelectorProps {
    eventId: string;
    ticketTypes: TicketType[];
    onSelectionChange: (selection: TicketSelection[]) => void;
}

export default function TicketSelector({ eventId, ticketTypes, onSelectionChange }: TicketSelectorProps) {
    const [quantities, setQuantities] = useState<Record<string, number>>({});
    
    // Render ticket type cards with quantity selectors
    // Calculate total price
    // Show availability
}
```

#### 2. Booking Summary Component
**File:** `www/src/components/booking/BookingSummary.tsx`

```typescript
interface BookingSummaryProps {
    eventTitle: string;
    selectedTickets: TicketSelection[];
    totalAmount: number;
    onConfirm: () => void;
    onEdit: () => void;
}
```

#### 3. Booking Page
**File:** `www/src/pages/booking.tsx`

**Multi-step flow:**
1. **Step 1:** Select tickets and quantity
2. **Step 2:** Review booking summary
3. **Step 3:** Payment (integrate Stripe)
4. **Step 4:** Confirmation with ticket details

```typescript
export default function BookingPage() {
    const [step, setStep] = useState(1);
    const [selectedTickets, setSelectedTickets] = useState([]);
    const [booking, setBooking] = useState(null);
    const { eventId } = useParams();
    
    // Step 1: Ticket Selection
    // Step 2: Review Summary
    // Step 3: Payment
    // Step 4: Confirmation
}
```

#### 4. Payment Integration Component
**File:** `www/src/components/payment/StripeCheckout.tsx`

**Dependencies:**
```bash
cd www
pnpm install @stripe/stripe-js @stripe/react-stripe-js
```

```typescript
import { Elements } from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY);

export function StripeCheckout({ bookingId, amount }: StripeCheckoutProps) {
    const [clientSecret, setClientSecret] = useState("");
    
    useEffect(() => {
        // Create payment intent
        api.post("/payments/create-intent", { bookingId, amount })
            .then(res => setClientSecret(res.data.clientSecret));
    }, []);
    
    return (
        <Elements stripe={stripePromise} options={{ clientSecret }}>
            <CheckoutForm />
        </Elements>
    );
}
```

#### 5. My Tickets Page
**File:** `www/src/pages/my-tickets.tsx`

```typescript
export default function MyTicketsPage() {
    const [tickets, setTickets] = useState([]);
    const { user } = useAuth();
    
    useEffect(() => {
        api.get(`/bookings/user/${user.id}`)
            .then(res => setTickets(res.data.data.bookings));
    }, [user.id]);
    
    return (
        <Layout>
            <h1>My Tickets</h1>
            <div className="grid gap-4">
                {tickets.map(ticket => (
                    <TicketCard key={ticket.id} ticket={ticket} />
                ))}
            </div>
        </Layout>
    );
}
```

#### 6. Ticket Display Component
**File:** `www/src/components/ticket/TicketCard.tsx`

**Features:**
- Show QR code
- Event details
- Ticket type and number
- Download button
- Check-in status

**Dependencies:**
```bash
pnpm install react-qr-code jspdf html2canvas
```

```typescript
import QRCode from "react-qr-code";

export function TicketCard({ ticket }: TicketCardProps) {
    const downloadTicket = async () => {
        // Generate PDF with jsPDF
    };
    
    return (
        <div className="ticket-card">
            <QRCode value={ticket.qrCodeData} size={200} />
            <div className="ticket-info">
                <h3>{ticket.eventTitle}</h3>
                <p>Ticket #{ticket.ticketNumber}</p>
                <p>Type: {ticket.ticketType}</p>
                {ticket.isCheckedIn && <Badge>Checked In</Badge>}
            </div>
            <button onClick={downloadTicket}>Download PDF</button>
        </div>
    );
}
```

**API Endpoints to Integrate:**
```typescript
// www/src/lib/api.ts

export const bookingAPI = {
    createBooking: (data: CreateBookingRequest) => 
        api.post("/bookings/create", data),
    
    getUserBookings: (userId: string) => 
        api.get(`/bookings/user/${userId}`),
    
    getBookingDetails: (bookingId: string) => 
        api.get(`/bookings/${bookingId}`),
    
    checkAvailability: (eventId: string) => 
        api.get(`/events/${eventId}/availability`),
};

export const paymentAPI = {
    createIntent: (bookingId: string, amount: number) => 
        api.post("/payments/create-intent", { bookingId, amount }),
    
    verifyPayment: (paymentId: string) => 
        api.post("/payments/verify", { paymentId }),
    
    getPaymentStatus: (paymentId: string) => 
        api.get(`/payments/${paymentId}/status`),
};
```

**Environment Variables:**
```env
# www/.env
VITE_STRIPE_PUBLISHABLE_KEY=pk_test_...
VITE_API_URL=http://localhost:8989/api/v1
```

**Routing Updates:**
**File:** `www/src/main.tsx`

```typescript
{
    path: "/events/:eventId/book",
    element: <ProtectedRoute><BookingPage /></ProtectedRoute>
},
{
    path: "/my-tickets",
    element: <ProtectedRoute><MyTicketsPage /></ProtectedRoute>
},
{
    path: "/booking/:bookingId/confirmation",
    element: <ProtectedRoute><BookingConfirmation /></ProtectedRoute>
}
```

**Testing Checklist:**
- [ ] Select tickets and see price calculation
- [ ] Cannot exceed available quantity
- [ ] Create booking successfully
- [ ] Complete payment flow (test mode)
- [ ] See confirmation after payment
- [ ] View tickets in "My Tickets"
- [ ] QR code displays correctly
- [ ] Download ticket as PDF

---

### **CP-3: QR Scanner Page (Admin)** ⏱️ 2-3 days

**Status:** 0% Complete  
**Priority:** P0 - CRITICAL for event management  
**User:** Super Admin only

**Dependencies:**
```bash
cd www
pnpm install @zxing/browser
```

**New Files to Create:**

#### 1. QR Scanner Component
**File:** `www/src/components/scanner/QRScanner.tsx`

```typescript
import { BrowserMultiFormatReader } from "@zxing/browser";

export function QRScanner({ onScan }: { onScan: (data: string) => void }) {
    const [scanning, setScanning] = useState(false);
    const videoRef = useRef<HTMLVideoElement>(null);
    
    const startScan = async () => {
        const codeReader = new BrowserMultiFormatReader();
        try {
            const result = await codeReader.decodeFromVideoDevice(
                undefined, 
                videoRef.current!,
                (result, error) => {
                    if (result) {
                        onScan(result.getText());
                    }
                }
            );
        } catch (err) {
            console.error("Scanner error:", err);
        }
    };
    
    return (
        <div className="scanner-container">
            <video ref={videoRef} className="scanner-video" />
            <button onClick={startScan}>Start Scanning</button>
        </div>
    );
}
```

#### 2. Scanner Page
**File:** `www/src/pages/scanner.tsx`

```typescript
export default function ScannerPage() {
    const [result, setResult] = useState<ScanResult | null>(null);
    const [loading, setLoading] = useState(false);
    const { user } = useAuth();
    
    const handleScan = async (qrData: string) => {
        setLoading(true);
        try {
            const response = await api.post("/checkin/validate", {
                qr_code_data: qrData,
                admin_id: user.id,
            });
            
            if (response.data.data.valid) {
                // Mark attendance
                await api.post("/checkin/mark", {
                    ticket_id: response.data.data.ticket_id,
                    admin_id: user.id,
                });
                setResult({ success: true, data: response.data.data });
            } else {
                setResult({ success: false, message: response.data.message });
            }
        } catch (error) {
            setResult({ success: false, message: "Invalid QR code" });
        } finally {
            setLoading(false);
        }
    };
    
    return (
        <Layout>
            <h1>Ticket Scanner</h1>
            <QRScanner onScan={handleScan} />
            {loading && <LoadingSpinner />}
            {result && <ScanResult result={result} />}
        </Layout>
    );
}
```

#### 3. Scan Result Display
**File:** `www/src/components/scanner/ScanResult.tsx`

```typescript
export function ScanResult({ result }: { result: ScanResult }) {
    if (result.success) {
        return (
            <div className="scan-success">
                <CheckCircleIcon className="text-green-500" />
                <h2>Valid Ticket ✓</h2>
                <p>Ticket: {result.data.ticketNumber}</p>
                <p>User: {result.data.userName}</p>
                <p>Type: {result.data.ticketType}</p>
            </div>
        );
    }
    
    return (
        <div className="scan-error">
            <XCircleIcon className="text-red-500" />
            <h2>Invalid Ticket ✗</h2>
            <p>{result.message}</p>
        </div>
    );
}
```

#### 4. Manual Entry Fallback
**File:** `www/src/components/scanner/ManualEntry.tsx`

```typescript
export function ManualEntry({ onSubmit }: { onSubmit: (ticketNumber: string) => void }) {
    const [ticketNumber, setTicketNumber] = useState("");
    
    return (
        <div className="manual-entry">
            <h3>Manual Ticket Entry</h3>
            <input 
                type="text" 
                placeholder="Enter ticket number"
                value={ticketNumber}
                onChange={(e) => setTicketNumber(e.target.value)}
            />
            <button onClick={() => onSubmit(ticketNumber)}>Validate</button>
        </div>
    );
}
```

**Routing:**
```typescript
{
    path: "/scanner",
    element: <ProtectedRoute roles={["SUPER_ADMIN"]}><ScannerPage /></ProtectedRoute>
}
```

**Testing:**
- [ ] Camera permissions granted
- [ ] QR code scanned successfully
- [ ] Valid ticket shows success
- [ ] Invalid ticket shows error
- [ ] Duplicate scan prevented
- [ ] Manual entry works

---

### **CP-4: Replace Mock Data with Real APIs** ⏱️ 2-3 days

**Status:** CRITICAL - All pages use mock data  
**Priority:** P0  

**Files to Update:**

#### 1. Home Page
**File:** `www/src/pages/home.tsx`

```typescript
// BEFORE (mock)
const events = mockEvents;

// AFTER (real API)
const [events, setEvents] = useState([]);
const [loading, setLoading] = useState(true);

useEffect(() => {
    api.get("/events/featured")
        .then(res => setEvents(res.data.data.events))
        .catch(err => console.error(err))
        .finally(() => setLoading(false));
}, []);

if (loading) return <LoadingSpinner />;
```

#### 2. Events Page
**File:** `www/src/pages/events.tsx`

```typescript
// Add filtering and pagination
const [events, setEvents] = useState([]);
const [filters, setFilters] = useState({
    category: "",
    location: "",
    page: 1,
    limit: 12,
});

const fetchEvents = async () => {
    const params = new URLSearchParams(filters);
    const response = await api.get(`/events?${params}`);
    setEvents(response.data.data.events);
};

useEffect(() => {
    fetchEvents();
}, [filters]);
```

#### 3. Event Detail Page
**File:** `www/src/pages/event-detail.tsx`

```typescript
// Get event by slug
const { slug } = useParams();
const [event, setEvent] = useState(null);
const [loading, setLoading] = useState(true);

useEffect(() => {
    api.get(`/events/${slug}`)
        .then(res => setEvent(res.data.data.event))
        .catch(err => {
            if (err.response?.status === 404) {
                navigate("/404");
            }
        })
        .finally(() => setLoading(false));
}, [slug]);
```

#### 4. Dashboard Page
**File:** `www/src/pages/dashboard.tsx`

```typescript
// Replace all mock stats
const [stats, setStats] = useState(null);

useEffect(() => {
    api.get("/dashboard/stats")
        .then(res => setStats(res.data.data))
        .catch(err => console.error(err));
}, []);

// Use real data
<StatsCard title="Total Events" value={stats?.totalEvents} />
```

#### 5. Admin Page
**File:** `www/src/pages/admin.tsx`

```typescript
// Admin-specific stats
useEffect(() => {
    api.get("/admin/dashboard/stats")
        .then(res => setStats(res.data.data))
        .catch(err => console.error(err));
}, []);
```

**API Integration Checklist:**
- [ ] Home page shows featured events
- [ ] Events page with filters working
- [ ] Event detail fetches by slug
- [ ] Dashboard shows real stats
- [ ] Admin dashboard shows real data
- [ ] Loading states for all API calls
- [ ] Error handling for failed requests
- [ ] Empty states when no data

---

### **CP-5: Error Handling & Loading States** ⏱️ 1-2 days

**Status:** Basic implementation  
**Priority:** P0 - Better UX  

**Components to Create:**

#### 1. Error Boundary
**File:** `www/src/components/ErrorBoundary.tsx`

```typescript
export class ErrorBoundary extends React.Component<
    {children: React.ReactNode},
    {hasError: boolean}
> {
    state = { hasError: false };
    
    static getDerivedStateFromError() {
        return { hasError: true };
    }
    
    componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
        console.error("Error caught:", error, errorInfo);
    }
    
    render() {
        if (this.state.hasError) {
            return <ErrorFallback />;
        }
        return this.props.children;
    }
}
```

#### 2. Loading Spinner
**File:** `www/src/components/ui/LoadingSpinner.tsx`

```typescript
export function LoadingSpinner({ fullPage = false }: { fullPage?: boolean }) {
    const className = fullPage 
        ? "fixed inset-0 flex items-center justify-center bg-white/80"
        : "flex justify-center p-8";
    
    return (
        <div className={className}>
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
        </div>
    );
}
```

#### 3. Loading Skeleton
**File:** `www/src/components/ui/Skeleton.tsx`

```typescript
export function Skeleton({ className }: { className?: string }) {
    return (
        <div className={`animate-pulse bg-gray-200 rounded ${className}`} />
    );
}

export function EventCardSkeleton() {
    return (
        <div className="card">
            <Skeleton className="h-48 w-full" />
            <Skeleton className="h-6 w-3/4 mt-4" />
            <Skeleton className="h-4 w-1/2 mt-2" />
        </div>
    );
}
```

#### 4. Toast Notifications
**File:** `www/src/components/ui/Toast.tsx`

**Dependencies:**
```bash
pnpm install react-hot-toast
```

```typescript
import toast, { Toaster } from "react-hot-toast";

// In Layout component
<Toaster position="top-right" />

// Usage
toast.success("Booking confirmed!");
toast.error("Payment failed");
toast.loading("Processing...");
```

#### 5. Error Display Component
**File:** `www/src/components/ui/ErrorMessage.tsx`

```typescript
export function ErrorMessage({ 
    message, 
    retry 
}: { 
    message: string; 
    retry?: () => void 
}) {
    return (
        <div className="error-container">
            <XCircleIcon className="text-red-500" />
            <p>{message}</p>
            {retry && <button onClick={retry}>Retry</button>}
        </div>
    );
}
```

**Implementation Checklist:**
- [ ] Wrap app in ErrorBoundary
- [ ] Add LoadingSpinner to all data-fetching pages
- [ ] Replace immediate renders with Skeleton components
- [ ] Add Toast for success/error messages
- [ ] Retry mechanism for failed API calls
- [ ] Offline state detection

---

## 🔥 **HIGH PRIORITY - MVP ENHANCEMENT (P1)**

### **P1-1: User Profile & Settings** ⏱️ 2-3 days

**Files to Create:**
- `www/src/pages/profile.tsx`
- `www/src/pages/settings.tsx`
- `www/src/components/profile/EditProfile.tsx`

**Features:**
- View/edit user information
- Change password
- Email preferences
- Account deletion

---

### **P1-2: Admin Event Management** ⏱️ 2-3 days

**Files:**
- `www/src/pages/admin/create-event.tsx`
- `www/src/pages/admin/edit-event.tsx`
- `www/src/components/admin/EventForm.tsx`

**Features:**
- Create new events
- Edit existing events
- Manage ticket types
- Upload event images

---

### **P1-3: Enhanced UI Components** ⏱️ 2-3 days

**Components to Build:**
- Modal component
- Dropdown/Select
- Date picker
- Image upload
- Pagination
- Search/Filter components

---

## 📊 **NICE TO HAVE - POST-MVP (P2)**

### **P2-1: Social Features**
- Event sharing
- Reviews and ratings
- Social login

### **P2-2: Advanced Search**
- Full-text search
- Auto-complete
- Filter presets
- Search history

### **P2-3: Performance Optimization**
- React.lazy code splitting
- Image lazy loading
- Virtual scrolling
- Bundle size optimization

---

## 📅 **IMPLEMENTATION ROADMAP**

### **Week 1-2: Core Features (Days 1-14)**
- ✅ Day 1-2: Authentication integration (CP-1)
- ✅ Day 3-9: Booking flow (CP-2)
- ✅ Day 10-12: QR scanner (CP-3)
- ✅ Day 13-14: Replace mock data (CP-4)

### **Week 3-4: Polish & Enhancement (Days 15-28)**
- ✅ Day 15-16: Error handling (CP-5)
- ✅ Day 17-19: User profile (P1-1)
- ✅ Day 20-22: Admin event management (P1-2)
- ✅ Day 23-25: Enhanced components (P1-3)
- ✅ Day 26-28: Testing and bug fixes

---

## ✅ **MVP COMPLETION CHECKLIST**

### Core Features
- [ ] Login/register working
- [ ] Protected routes implemented
- [ ] Booking flow complete
- [ ] Payment integration working
- [ ] QR codes displaying
- [ ] Scanner functional
- [ ] All mock data replaced

### UI/UX
- [ ] Loading states on all pages
- [ ] Error handling comprehensive
- [ ] Mobile responsive
- [ ] Toast notifications
- [ ] Empty states

### Testing
- [ ] Authentication flow tested
- [ ] Booking flow tested
- [ ] Payment flow tested (test mode)
- [ ] QR scanner tested
- [ ] Cross-browser tested

---

**Total Estimated Effort:** 15-20 development days (3-4 weeks for 1 developer)  
**Dependencies:** Backend APIs must be ready
