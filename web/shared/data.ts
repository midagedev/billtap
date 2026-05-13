import { appHref } from "./basePath";

export type StatusTone = "good" | "warning" | "danger" | "neutral" | "info";

export type TimelineEntry = {
  time: string;
  object: string;
  transition: string;
  event: string;
  delivery: string;
  response: string;
  tone: StatusTone;
};

export type WebhookAttempt = {
  id?: string;
  eventId?: string;
  endpointId?: string;
  endpoint: string;
  event: string;
  attemptNumber?: number;
  attempts: number;
  next: string;
  status: string;
  tone?: StatusTone;
  scheduledAt?: string;
  deliveredAt?: string;
  requestUrl?: string;
  requestHeaders?: string;
  requestBody?: string;
  signatureHeader?: string;
  responseStatus?: string;
  responseBody?: string;
  errorMessage?: string;
  retryPlan?: string;
};

export type DashboardObjectType =
  | "customers"
  | "subscriptions"
  | "invoices"
  | "paymentIntents"
  | "checkoutSessions"
  | "webhookEvents";

export type ObjectField = {
  label: string;
  value: string;
};

export type DashboardObjectRecord = {
  id: string;
  type: DashboardObjectType;
  label: string;
  title: string;
  subtitle: string;
  status: string;
  tone: StatusTone;
  createdAt?: string;
  amount?: string;
  fields: ObjectField[];
};

export type AppResponseRecord = {
  id: string;
  eventId: string;
  endpoint: string;
  status: string;
  body: string;
  retryPlan: string;
  tone: StatusTone;
};

export type Invoice = {
  id: string;
  period: string;
  amount: string;
  status: string;
};

export type PortalPlan = {
  id: string;
  name: string;
  price: string;
  detail: string;
};

export type CheckoutOutcomeId =
  | "success"
  | "declined"
  | "funds"
  | "expired"
  | "incorrect_cvc"
  | "processing_error"
  | "action"
  | "cancel"
  | "pending";

export type CheckoutOutcome = {
  id: CheckoutOutcomeId;
  label: string;
  status: string;
  tone: StatusTone;
};

export const billingSnapshot = {
  workspace: "ws_demo",
  owner: "billing-owner@demo.local",
  tenantRoute: "platform-card",
  paymentRail: "sandbox card",
  plan: "Team Annual",
  status: "trialing",
  period: "May 1 - May 31",
  seats: {
    basic: 5,
    additional: 3,
    used: 6,
    pending: 1,
  },
  export: {
    defaultLimit: 1000,
    manualLimit: 200,
    remaining: 836,
    renews: "Jun 1",
  },
};

export const dashboardStats = [
  { label: "Open workspaces", value: "12", meta: "3 with pending assertions", tone: "info" as const },
  { label: "Webhook attempts", value: "84", meta: "7 retries scheduled", tone: "warning" as const },
  { label: "Scenario runs", value: "18", meta: "15 passing", tone: "good" as const },
  { label: "Debug bundles", value: "5", meta: "latest 14 min ago", tone: "neutral" as const },
];

export const timelineEntries: TimelineEntry[] = [
  {
    time: "10:42:08",
    object: "cs_test_42",
    transition: "checkout completed",
    event: "checkout.session.completed",
    delivery: "delivered",
    response: "200 OK",
    tone: "good",
  },
  {
    time: "10:42:10",
    object: "in_test_61",
    transition: "invoice paid",
    event: "invoice.payment_succeeded",
    delivery: "retry queued",
    response: "503 timeout",
    tone: "warning",
  },
  {
    time: "10:43:12",
    object: "sub_test_29",
    transition: "entitlement asserted",
    event: "customer.subscription.updated",
    delivery: "duplicate held",
    response: "idempotent",
    tone: "info",
  },
  {
    time: "10:47:31",
    object: "evt_test_90",
    transition: "replay requested",
    event: "invoice.payment_failed",
    delivery: "blocked",
    response: "signature mismatch",
    tone: "danger",
  },
];

export const webhookAttempts: WebhookAttempt[] = [
  {
    id: "delatt_test_02",
    eventId: "evt_test_61",
    endpoint: "http://app.local/webhooks/stripe",
    event: "invoice.payment_succeeded",
    attemptNumber: 2,
    attempts: 2,
    next: "+2m",
    status: "retrying",
    tone: "warning",
    scheduledAt: "10:44:10",
    requestUrl: "http://app.local/webhooks/stripe",
    signatureHeader: "t=1778233312,v1=sandbox_retry",
    requestHeaders: "Stripe-Signature: t=1778233312,v1=sandbox_retry",
    requestBody: "invoice.payment_succeeded payload",
    responseStatus: "503",
    responseBody: "timeout before entitlement update",
    errorMessage: "handler returned 503",
    retryPlan: "max 5 · backoff 10s, 30s, 2m, 10m",
  },
  {
    id: "delatt_test_03",
    eventId: "evt_test_refund",
    endpoint: "http://app.local/webhooks/connect",
    event: "charge.refunded",
    attemptNumber: 1,
    attempts: 1,
    next: "-",
    status: "delivered",
    tone: "good",
    deliveredAt: "10:42:13",
    requestUrl: "http://app.local/webhooks/connect",
    signatureHeader: "t=1778233333,v1=sandbox_connect",
    responseStatus: "200",
    responseBody: "ok",
    retryPlan: "none",
  },
  {
    id: "delatt_test_04",
    eventId: "evt_test_29",
    endpoint: "http://app.local/webhooks/stripe",
    event: "customer.subscription.updated",
    attemptNumber: 1,
    attempts: 1,
    next: "-",
    status: "held duplicate",
    tone: "info",
    scheduledAt: "10:43:12",
    requestUrl: "http://app.local/webhooks/stripe",
    signatureHeader: "t=1778233392,v1=sandbox_duplicate",
    requestBody: "customer.subscription.updated payload",
    responseStatus: "idempotent",
    responseBody: "duplicate event already applied",
    retryPlan: "duplicate held by reliability rule",
  },
];

export const dashboardObjects: Record<DashboardObjectType, DashboardObjectRecord[]> = {
  customers: [
    {
      id: "cus_test_workspace",
      type: "customers",
      label: "Customer",
      title: "owner@demo.local",
      subtitle: "Example workspace owner",
      status: "active",
      tone: "good",
      createdAt: "10:40:12",
      fields: [
        { label: "Customer ID", value: "cus_test_workspace" },
        { label: "Name", value: "Example Owner" },
        { label: "Email", value: "owner@demo.local" },
        { label: "Workspace", value: "ws_demo" },
      ],
    },
  ],
  subscriptions: [
    {
      id: "sub_test_29",
      type: "subscriptions",
      label: "Subscription",
      title: "Team Annual",
      subtitle: "cus_test_workspace",
      status: "trialing",
      tone: "info",
      createdAt: "10:42:09",
      amount: "$1,188.00",
      fields: [
        { label: "Subscription ID", value: "sub_test_29" },
        { label: "Customer", value: "cus_test_workspace" },
        { label: "Latest invoice", value: "in_test_61" },
        { label: "Current period", value: "May 1 - May 31" },
        { label: "Cancel at period end", value: "false" },
      ],
    },
  ],
  invoices: [
    {
      id: "in_test_61",
      type: "invoices",
      label: "Invoice",
      title: "May 2026 renewal",
      subtitle: "sub_test_29",
      status: "paid after retry",
      tone: "warning",
      createdAt: "10:42:10",
      amount: "$1,188.00",
      fields: [
        { label: "Invoice ID", value: "in_test_61" },
        { label: "Customer", value: "cus_test_workspace" },
        { label: "Subscription", value: "sub_test_29" },
        { label: "Payment intent", value: "pi_test_61" },
        { label: "Attempt count", value: "2" },
      ],
    },
  ],
  paymentIntents: [
    {
      id: "pi_test_61",
      type: "paymentIntents",
      label: "Payment intent",
      title: "$1,188.00 USD",
      subtitle: "in_test_61",
      status: "succeeded",
      tone: "good",
      createdAt: "10:42:11",
      amount: "$1,188.00",
      fields: [
        { label: "Payment intent ID", value: "pi_test_61" },
        { label: "Invoice", value: "in_test_61" },
        { label: "Customer", value: "cus_test_workspace" },
        { label: "Failure code", value: "none" },
      ],
    },
  ],
  checkoutSessions: [
    {
      id: "cs_test_42",
      type: "checkoutSessions",
      label: "Checkout session",
      title: "Team Annual checkout",
      subtitle: "cus_test_workspace",
      status: "complete",
      tone: "good",
      createdAt: "10:42:08",
      amount: "$1,188.00",
      fields: [
        { label: "Checkout session ID", value: "cs_test_42" },
        { label: "Customer", value: "cus_test_workspace" },
        { label: "Subscription", value: "sub_test_29" },
        { label: "Invoice", value: "in_test_61" },
        { label: "Payment status", value: "paid" },
      ],
    },
  ],
  webhookEvents: [
    {
      id: "evt_test_61",
      type: "webhookEvents",
      label: "Webhook event",
      title: "invoice.payment_succeeded",
      subtitle: "in_test_61",
      status: "retrying",
      tone: "warning",
      createdAt: "10:42:10",
      fields: [
        { label: "Event ID", value: "evt_test_61" },
        { label: "Event type", value: "invoice.payment_succeeded" },
        { label: "Source", value: "checkout" },
        { label: "Idempotency key", value: "billtap:invoice.payment_succeeded:in_test_61" },
      ],
    },
    {
      id: "evt_test_refund",
      type: "webhookEvents",
      label: "Webhook event",
      title: "charge.refunded",
      subtitle: "connect rail",
      status: "delivered",
      tone: "good",
      createdAt: "10:42:13",
      fields: [
        { label: "Event ID", value: "evt_test_refund" },
        { label: "Event type", value: "charge.refunded" },
        { label: "Source", value: "connect" },
        { label: "Idempotency key", value: "billtap:charge.refunded:ch_test_refund" },
      ],
    },
  ],
};

export const checkoutSession = {
  id: "cs_test_checkout_demo",
  customer: "cus_test_workspace",
  customerEmail: "owner@demo.local",
  plan: "Team Annual",
  price: "$1,188 / year",
  status: "open",
  paymentStatus: "unpaid",
  subscriptionId: "sub_test_pending",
  subscriptionStatus: "not created",
  invoiceId: "in_test_pending",
  invoiceStatus: "draft",
  paymentIntentId: "pi_test_pending",
  paymentIntentStatus: "requires_payment_method",
  returnUrl: appHref("dashboard/"),
  lineItems: [
    { label: "Team Annual base", amount: "$948.00" },
    { label: "3 additional seats", amount: "$240.00" },
    { label: "Sandbox tax", amount: "$0.00" },
  ],
};

export const outcomeOptions: CheckoutOutcome[] = [
  { id: "success", label: "Payment succeeds", status: "Ready", tone: "good" },
  { id: "declined", label: "Card declined", status: "Failure", tone: "danger" },
  { id: "funds", label: "Insufficient funds", status: "Failure", tone: "danger" },
  { id: "expired", label: "Expired card", status: "Failure", tone: "warning" },
  { id: "incorrect_cvc", label: "Incorrect CVC", status: "Failure", tone: "danger" },
  { id: "processing_error", label: "Processing error", status: "Failure", tone: "danger" },
  { id: "action", label: "Requires action", status: "Challenge", tone: "info" },
  { id: "cancel", label: "User cancels", status: "Return", tone: "neutral" },
  { id: "pending", label: "Async payment pending", status: "Pending", tone: "warning" },
];

export const portalPlans: PortalPlan[] = [
  { id: "starter", name: "Starter", price: "$49/mo", detail: "2 seats, 200 exports" },
  { id: "team", name: "Team", price: "$99/mo", detail: "5 seats, 1,000 exports" },
  { id: "scale", name: "Scale", price: "$249/mo", detail: "15 seats, 5,000 exports" },
];

export const invoices: Invoice[] = [
  { id: "in_test_061", period: "May 2026", amount: "$1,188.00", status: "paid" },
  { id: "in_test_044", period: "Apr 2026", amount: "$99.00", status: "paid" },
  { id: "in_test_031", period: "Mar 2026", amount: "$99.00", status: "void" },
];
