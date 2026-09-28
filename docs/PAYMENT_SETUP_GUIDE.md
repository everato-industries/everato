# Payment Gateway Setup & Self-Hosting Guide

This guide walks developers and organizers through setting up and configuring payment collection on a self-hosted **Everato** instance.

Everato includes native, pluggable support for three leading payment gateways:
- **Razorpay** (Recommended for India, UPI, Cards, Netbanking & International)
- **Stripe** (Recommended for Global/International card payments in USD, EUR, GBP, etc.)
- **Cashfree PG** (Alternative for India with low transaction fees and instant settlements)

---

## Table of Contents

1. [How Payments Work in Everato](#how-payments-work-in-everato)
2. [Quick Selection Matrix](#quick-selection-matrix)
3. [Provider Setup Guides](#provider-setup-guides)
   - [Option 1: Setting up Razorpay](#option-1-setting-up-razorpay)
   - [Option 2: Setting up Stripe](#option-2-setting-up-stripe)
   - [Option 3: Setting up Cashfree](#option-3-setting-up-cashfree)
4. [Testing Locally with Webhooks (ngrok)](#testing-locally-with-webhooks-ngrok)
5. [Going Live & Collecting Real Money](#going-live--collecting-real-money)
6. [Troubleshooting & FAQs](#troubleshooting--faqs)

---

## How Payments Work in Everato

Everato follows a secure, server-verified payment lifecycle:

```
[Attendee]                             [Everato Server]                     [Payment Gateway]
     │                                        │                                     │
     │ 1. Reserve Tickets (Create Booking)    │                                     │
     ├───────────────────────────────────────>│ (Booking status: PENDING)           │
     │                                        │                                     │
     │ 2. Request Payment Checkout            │                                     │
     ├───────────────────────────────────────>│ 3. Initialize Order/Intent          │
     │                                        ├────────────────────────────────────>│
     │                                        │<────────────────────────────────────┤
     │                                        │ (Returns order_id, client keys)     │
     │<───────────────────────────────────────┤                                     │
     │                                        │                                     │
     │ 4. Open Payment Modal & Authorize      │                                     │
     ├─────────────────────────────────────────────────────────────────────────────>│
     │<─────────────────────────────────────────────────────────────────────────────┤
     │                                        │                                     │
     │ 5. Submit Payment Signature / Token    │                                     │
     ├───────────────────────────────────────>│6. Cryptographically Verify Signature│
     │                                        │    - Updates booking to CONFIRMED   │
     │                                        │    - Generates HMAC-signed QR passes│
     │                                        │    - Sends confirmation email       │
     │<───────────────────────────────────────┤                                     │
     │                                        │                                     │
     │                                        │ 7. Asynchronous Webhook (Fallback)  │
     │                                        │<────────────────────────────────────┤
```

> **Zero-Risk Free Events**: If an event ticket price is `0.00` (or 100% coupon applied), Everato automatically bypasses the payment gateway and confirms the booking immediately.

---

## Quick Selection Matrix

| Feature | Razorpay | Stripe | Cashfree |
|---|---|---|---|
| **Primary Region** | India & Global | Worldwide (40+ countries) | India |
| **Payment Modes** | UPI, Cards, Netbanking, Wallets | Credit/Debit Cards, Apple Pay, Google Pay | UPI, Cards, Netbanking, Wallets |
| **Default Currency** | `INR` | `USD` (or EUR, GBP, AUD, etc.) | `INR` |
| **Sandbox/Test Mode** | Instant test keys | Instant test keys | Instant sandbox app ID |
| **Best For** | Indian events, tech conferences, hackathons | Global international conferences | Indian businesses prioritizing low fees |

---

## Provider Setup Guides

### Option 1: Setting up Razorpay

#### 1. Create a Razorpay Account
1. Sign up at [https://dashboard.razorpay.com/](https://dashboard.razorpay.com/).
2. You can immediately access **Test Mode** without submitting business documents.

#### 2. Obtain API Keys
1. In the Razorpay Dashboard, toggle the switch to **Test Mode** (top-right or sidebar).
2. Go to **Account & Settings** → **API Keys** (under *Developer Controls*).
3. Click **Generate Key**.
4. Copy the **Key ID** (`rzp_test_...`) and **Key Secret**.

#### 3. Configure Everato
In your `.env` file:
```bash
# Payment Provider
PAYMENT_PROVIDER=razorpay
PAYMENT_CURRENCY=INR

# Razorpay Keys
RAZORPAY_KEY_ID=rzp_test_your_key_id
RAZORPAY_KEY_SECRET=your_razorpay_secret_key
RAZORPAY_WEBHOOK_SECRET=your_custom_webhook_secret_string
```
Or in `config.yaml`:
```yaml
payment:
    provider: razorpay
    currency: INR
```

#### 4. Configure Webhooks
1. In Razorpay Dashboard, go to **Account & Settings** → **Webhooks**.
2. Click **Add New Webhook**.
3. **Webhook URL**: `https://your-domain.com/api/v1/payments/webhook`
4. **Secret**: Enter the same string you put in `RAZORPAY_WEBHOOK_SECRET`.
5. **Active Events**: Check `order.paid` and `payment.captured`.
6. Click **Save**.

#### 5. Test Credentials
- **UPI**: Use any VPA like `success@razorpay`.
- **Cards**: Use card number `4000 0000 0000 0002` (any future expiry, any 3-digit CVV, OTP `123456`).

---

### Option 2: Setting up Stripe

#### 1. Create a Stripe Account
1. Sign up at [https://dashboard.stripe.com/register](https://dashboard.stripe.com/register).
2. Stripe provides immediate sandbox access in **Test Mode**.

#### 2. Obtain API Keys
1. In the Stripe Dashboard, make sure **Test mode** is toggled ON.
2. Go to **Developers** → **API keys**.
3. Copy:
   - **Publishable key** (`pk_test_...`)
   - **Secret key** (`sk_test_...`)

#### 3. Configure Everato
In your `.env` file:
```bash
# Payment Provider
PAYMENT_PROVIDER=stripe
PAYMENT_CURRENCY=USD # Or EUR, GBP, INR, etc.

# Stripe Keys
STRIPE_SECRET_KEY=sk_test_your_secret_key
STRIPE_PUBLISHABLE_KEY=pk_test_your_publishable_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_signing_secret
```
Or in `config.yaml`:
```yaml
payment:
    provider: stripe
    currency: USD
```

#### 4. Configure Webhooks
1. In Stripe Dashboard, go to **Developers** → **Webhooks**.
2. Click **Add endpoint**.
3. **Endpoint URL**: `https://your-domain.com/api/v1/payments/webhook`
4. **Select events**: Select `payment_intent.succeeded`.
5. Click **Add endpoint**.
6. Reveal the **Signing secret** (`whsec_...`) and paste it into `STRIPE_WEBHOOK_SECRET` in `.env`.

#### 5. Test Credentials
- **Cards**: Card number `4242 4242 4242 4242` (any future expiry, any 3-digit CVC).

---

### Option 3: Setting up Cashfree

#### 1. Create a Cashfree Account
1. Sign up at [https://merchant.cashfree.com/](https://merchant.cashfree.com/).
2. Switch to the **Payment Gateway** product and select the **Sandbox / Test Environment**.

#### 2. Obtain API Keys
1. Go to **Developers** → **API Keys**.
2. Click **Generate API Keys**.
3. Copy the **App ID** and **Secret Key**.

#### 3. Configure Everato
In your `.env` file:
```bash
# Payment Provider
PAYMENT_PROVIDER=cashfree
PAYMENT_CURRENCY=INR

# Cashfree Keys
CASHFREE_APP_ID=your_cashfree_app_id
CASHFREE_SECRET_KEY=your_cashfree_secret_key
CASHFREE_ENV=TEST # Change to PRODUCTION when going live
```
Or in `config.yaml`:
```yaml
payment:
    provider: cashfree
    currency: INR
```

#### 4. Configure Webhooks
1. In Cashfree Merchant Dashboard, go to **Developers** → **Webhooks**.
2. **Endpoint URL**: `https://your-domain.com/api/v1/payments/webhook`
3. Check the event: `PAYMENT_SUCCESS_WEBHOOK`.

---

## Testing Locally with Webhooks (ngrok)

Payment gateways cannot send webhook callbacks to `localhost`. When testing locally, use a secure tunnel such as **ngrok**:

1. Install ngrok (`brew install ngrok` or `npm install -g ngrok`).
2. Start your Everato server on port 8080:
   ```bash
   make dev # or ./bin/everato
   ```
3. In another terminal, start the tunnel:
   ```bash
   ngrok http 8080
   ```
4. Copy your forwarding HTTPS URL (e.g. `https://a1b2-c3d4.ngrok-free.app`).
5. Set your gateway webhook URL to:
   ```
   https://a1b2-c3d4.ngrok-free.app/api/v1/payments/webhook
   ```
6. Complete a test booking on `http://localhost:5173` (or `8080`).
7. Watch the terminal logs to verify webhook delivery and ticket QR pass generation!

---

## Going Live & Collecting Real Money

When you are ready to collect real ticket payments:

### 1. Complete Business KYC
- Log in to your provider dashboard (Razorpay, Stripe, or Cashfree).
- Complete account activation and submit business documentation (PAN, GSTIN, Company registration, or Individual identity proofs depending on jurisdiction).
- Add your bank account details for automated payouts.

### 2. Switch to Live Credentials
1. Toggle from **Test Mode** to **Live Mode** in your provider dashboard.
2. Generate **Live API Keys**.
3. Update your production `.env` file on your server:

```bash
# Example: Live Razorpay Configuration
PAYMENT_PROVIDER=razorpay
PAYMENT_CURRENCY=INR
RAZORPAY_KEY_ID=rzp_live_...
RAZORPAY_KEY_SECRET=...
RAZORPAY_WEBHOOK_SECRET=...

# Example: Live Stripe Configuration
PAYMENT_PROVIDER=stripe
PAYMENT_CURRENCY=USD
STRIPE_SECRET_KEY=sk_live_...
STRIPE_PUBLISHABLE_KEY=pk_live_...
STRIPE_WEBHOOK_SECRET=whsec_...

# Example: Live Cashfree Configuration
PAYMENT_PROVIDER=cashfree
PAYMENT_CURRENCY=INR
CASHFREE_APP_ID=...
CASHFREE_SECRET_KEY=...
CASHFREE_ENV=PRODUCTION
```

### 3. Ensure HTTPS is Enabled
Payment gateways require your site to be served over **HTTPS** in production. Ensure you have SSL/TLS configured on your server or reverse proxy (Nginx / Caddy / Cloudflare).

### 4. Restart the Everato Service
```bash
# If using systemd
sudo systemctl restart everato

# If running directly
./bin/everato
```

---

## Troubleshooting & FAQs

### Q1: What happens if a customer pays, but closes the browser before redirection?
**A:** Everato's webhook handler (`POST /api/v1/payments/webhook`) acts as a safeguard. The payment gateway notifies Everato asynchronously. Upon receiving the valid webhook, Everato automatically confirms the booking, generates the HMAC-signed tickets, and sends the ticket email.

### Q2: What if I don't want to accept payments and only host free events?
**A:** You do not need to configure any payment credentials! If `PAYMENT_PROVIDER` is unset or credentials are empty, free bookings are automatically confirmed without calling any payment gateway.

### Q3: Why am I getting "signature verification failed"?
**A:**
- Ensure `RAZORPAY_KEY_SECRET` or `STRIPE_WEBHOOK_SECRET` matches the exact secret from your dashboard.
- Ensure your server system clock is synced via NTP (out-of-sync clocks can cause webhook signature rejections).

### Q4: Can I switch payment providers later?
**A:** Yes! Because Everato uses a unified database schema and common `PaymentGateway` interface, you can switch providers simply by changing `PAYMENT_PROVIDER` and the corresponding keys in your `.env` file and restarting the application. Past tickets and booking records remain unaffected.
