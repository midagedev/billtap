import {
  billingSnapshot,
  checkoutSession as fixtureCheckoutSession,
  dashboardStats,
  dashboardObjects,
  invoices as fixtureInvoices,
  portalPlans,
  timelineEntries,
  webhookAttempts,
  type CheckoutOutcomeId,
  type AppResponseRecord,
  type DashboardObjectRecord,
  type DashboardObjectType,
  type Invoice,
  type PortalPlan,
  type StatusTone,
  type TimelineEntry,
  type WebhookAttempt,
} from "./data";

export type DataSource = "api" | "fixture";

export type LoadResult<T> = {
  data: T;
  source: DataSource;
  error?: string;
};

export type CheckoutLineItem = {
  label: string;
  amount: string;
};

export type CheckoutSessionSummary = {
  id: string;
  customer: string;
  customerEmail: string;
  plan: string;
  price: string;
  status: string;
  paymentStatus: string;
  subscriptionId: string;
  subscriptionStatus: string;
  invoiceId: string;
  invoiceStatus: string;
  paymentIntentId: string;
  paymentIntentStatus: string;
  returnUrl: string;
  lineItems: CheckoutLineItem[];
};

export type CheckoutCompletion = {
  session: CheckoutSessionSummary;
  source: DataSource;
  message: string;
  error?: string;
};

export type DashboardData = {
  billingSnapshot: typeof billingSnapshot;
  dashboardStats: typeof dashboardStats;
  objects: Record<DashboardObjectType, DashboardObjectRecord[]>;
  timelineEntries: TimelineEntry[];
  webhookAttempts: WebhookAttempt[];
  appResponses: AppResponseRecord[];
  evidence: {
    signatureHeader: string;
    idempotencyKey: string;
    failureReason: string;
  };
  scenario: {
    currentStep: string;
    clockAdvance: string;
    assertions: string;
    tone: StatusTone;
  };
};

export type DebugBundleRequest = {
  objectType: DashboardObjectType;
  objectId: string;
};

export type DebugBundleResult = {
  id: string;
  message: string;
  source: DataSource;
  url?: string;
  createdAt?: string;
};

export type PortalPaymentMethod = {
  id: string;
  label: string;
  status: string;
};

export type PortalSubscriptionSummary = {
  id: string;
  customerId: string;
  customerEmail: string;
  workspace: string;
  owner: string;
  planId: string;
  planName: string;
  price: string;
  status: string;
  period: string;
  seats: number;
  usedSeats: number;
  pendingSeats: number;
  cancelAtPeriodEnd: boolean;
  canceledAt?: string;
  exportRemaining: number;
  exportManualLimit: number;
  exportRenewal: string;
};

export type PortalData = {
  subscription: PortalSubscriptionSummary;
  plans: PortalPlan[];
  invoices: Invoice[];
  paymentMethod: PortalPaymentMethod;
};

export type PortalAction =
  | "plan-change"
  | "seat-change"
  | "cancel"
  | "resume"
  | "payment-method";

export type PortalActionResult = {
  data: PortalData;
  source: DataSource;
  action: PortalAction;
  message: string;
  error?: string;
};

export type PortalCancelMode = "period" | "immediate";

export type PortalPaymentOutcome = "succeeds" | "fails";

export async function replayWebhookEvent(eventId: string): Promise<string> {
  const body = await postJSON(`/api/events/${encodeURIComponent(eventId)}/replay`, {});
  const record = asRecord(body);
  return record
    ? readString(record, ["message", "status", "result"], "Replay requested")
    : "Replay requested";
}

export async function createDebugBundle(request: DebugBundleRequest): Promise<DebugBundleResult> {
  const payload = {
    object_type: request.objectType,
    object_id: request.objectId,
    include: ["timeline", "webhook_events", "delivery_attempts", "app_responses", "assertions"],
  };

  try {
    const body = await postJSON("/api/debug-bundles", payload);
    const record = asRecord(body) ?? {};
    const bundle = firstRecord(record, ["bundle", "debug_bundle", "debugBundle", "data"]) ?? record;
    return {
      id:
        readString(bundle, ["id", "bundle_id", "bundleId"], "") ||
        `dbg_${request.objectId}`,
      message: readString(record, ["message", "status"], "Debug bundle export created"),
      source: "api",
      url: readString(bundle, ["url", "download_url", "downloadUrl"], "") || undefined,
      createdAt: readString(bundle, ["created_at", "createdAt"], "") || undefined,
    };
  } catch (error) {
    if (!isEndpointUnavailable(error)) throw error;
    return {
      id: `fixture_bundle_${request.objectId}`,
      message: "Fixture debug bundle prepared because /api/debug-bundles is not available",
      source: "fixture",
      createdAt: new Date().toISOString(),
    };
  }
}

const fixtureSession: CheckoutSessionSummary = {
  ...fixtureCheckoutSession,
  lineItems: [...fixtureCheckoutSession.lineItems],
};

export function getCheckoutSessionId(location: Location = window.location): string | undefined {
  const params = new URLSearchParams(location.search);
  const queryValue =
    params.get("session_id") ??
    params.get("checkout_session_id") ??
    params.get("sessionId") ??
    params.get("cs") ??
    params.get("id");

  if (queryValue) return queryValue;

  const parts = location.pathname.split("/").filter(Boolean);
  const checkoutIndex = parts.lastIndexOf("checkout");
  if (checkoutIndex >= 0 && parts[checkoutIndex + 1]) return parts[checkoutIndex + 1];

  const last = parts.at(-1);
  return last?.startsWith("cs_") ? last : undefined;
}

export async function loadCheckoutSession(sessionId: string | undefined): Promise<LoadResult<CheckoutSessionSummary>> {
  if (!sessionId) {
    return { data: fixtureSession, source: "fixture", error: "No checkout session id in URL" };
  }

  try {
    const body = await getJSON(`/v1/checkout/sessions/${encodeURIComponent(sessionId)}`);
    return { data: normalizeCheckoutSession(body, sessionId), source: "api" };
  } catch (error) {
    return {
      data: { ...fixtureSession, id: sessionId },
      source: "fixture",
      error: errorMessage(error),
    };
  }
}

export async function completeCheckout(
  session: CheckoutSessionSummary,
  outcome: CheckoutOutcomeId,
  paymentMethod: string,
): Promise<CheckoutCompletion> {
  const payload = {
    outcome: apiOutcome(outcome),
    payment_method: paymentMethod,
    payment_method_label: paymentMethod,
  };
  const paths = [
    `/api/checkout/sessions/${encodeURIComponent(session.id)}/complete`,
    `/api/checkout/${encodeURIComponent(session.id)}/complete`,
    `/v1/checkout/sessions/${encodeURIComponent(session.id)}/complete`,
  ];

  for (const path of paths) {
    try {
      const body = await postJSON(path, payload);
      return {
        session: normalizeCheckoutSession(body, session.id),
        source: "api",
        message: "Checkout completion accepted by API",
      };
    } catch (error) {
      if (!isEndpointUnavailable(error)) {
        return {
          session,
          source: "api",
          message: "Checkout completion failed",
          error: errorMessage(error),
        };
      }
    }
  }

  return {
    session: simulateCheckoutCompletion(session, outcome),
    source: "fixture",
    message: "Fixture completion used because no checkout completion endpoint responded",
  };
}

export function getPortalCustomerId(location: Location = window.location): string | undefined {
  const params = new URLSearchParams(location.search);
  const queryValue =
    params.get("customer_id") ??
    params.get("customer") ??
    params.get("customerId") ??
    params.get("cus") ??
    params.get("id");

  if (queryValue) return queryValue;

  const parts = location.pathname.split("/").filter(Boolean);
  const portalIndex = parts.lastIndexOf("portal");
  if (portalIndex >= 0 && parts[portalIndex + 1]) return parts[portalIndex + 1];

  const last = parts.at(-1);
  return last?.startsWith("cus_") ? last : undefined;
}

export async function loadPortalData(customerId: string | undefined): Promise<LoadResult<PortalData>> {
  const fallback = fixturePortalData(customerId);
  const portalPaths = customerId
    ? [
        `/api/portal/customers/${encodeURIComponent(customerId)}`,
        `/api/portal?customer_id=${encodeURIComponent(customerId)}`,
      ]
    : ["/api/portal"];

  for (const path of portalPaths) {
    try {
      const body = await getJSON(path);
      return { data: normalizePortalData(body, fallback), source: "api" };
    } catch (error) {
      if (!isEndpointUnavailable(error)) {
        return { data: fallback, source: "fixture", error: errorMessage(error) };
      }
    }
  }

  try {
    const data = await loadPortalDataFromStripeLikeApi(customerId, fallback);
    if (data) return { data, source: "api" };
  } catch (error) {
    return { data: fallback, source: "fixture", error: errorMessage(error) };
  }

  return {
    data: fallback,
    source: "fixture",
    error: "Portal API not available; fixture portal state loaded",
  };
}

export async function changePortalPlan(data: PortalData, planId: string): Promise<PortalActionResult> {
  const payload = {
    customer_id: data.subscription.customerId,
    plan_id: planId,
    plan: planId,
  };
  const fallback = simulatePortalPlanChange(data, planId);
  return postPortalAction(
    `/api/portal/subscriptions/${encodeURIComponent(data.subscription.id)}/plan-change`,
    payload,
    fallback,
    "plan-change",
  );
}

export async function changePortalSeats(data: PortalData, seats: number): Promise<PortalActionResult> {
  const nextSeats = Math.max(1, Math.trunc(seats));
  const payload = {
    customer_id: data.subscription.customerId,
    seats: nextSeats,
    quantity: nextSeats,
  };
  const fallback = simulatePortalSeatChange(data, nextSeats);
  return postPortalAction(
    `/api/portal/subscriptions/${encodeURIComponent(data.subscription.id)}/seat-change`,
    payload,
    fallback,
    "seat-change",
  );
}

export async function cancelPortalSubscription(
  data: PortalData,
  mode: PortalCancelMode,
): Promise<PortalActionResult> {
  const payload = {
    customer_id: data.subscription.customerId,
    mode,
    cancel_at_period_end: mode === "period",
  };
  const fallback = simulatePortalCancellation(data, mode);
  return postPortalAction(
    `/api/portal/subscriptions/${encodeURIComponent(data.subscription.id)}/cancel`,
    payload,
    fallback,
    "cancel",
  );
}

export async function resumePortalSubscription(data: PortalData): Promise<PortalActionResult> {
  const payload = { customer_id: data.subscription.customerId };
  const fallback = simulatePortalResume(data);
  return postPortalAction(
    `/api/portal/subscriptions/${encodeURIComponent(data.subscription.id)}/resume`,
    payload,
    fallback,
    "resume",
  );
}

export async function updatePortalPaymentMethod(
  data: PortalData,
  outcome: PortalPaymentOutcome,
): Promise<PortalActionResult> {
  const payload = {
    customer_id: data.subscription.customerId,
    outcome,
    simulate: outcome === "succeeds" ? "success" : "failure",
  };
  const fallback = simulatePortalPaymentMethod(data, outcome);
  return postPortalAction(
    `/api/portal/customers/${encodeURIComponent(data.subscription.customerId)}/payment-method`,
    payload,
    fallback,
    "payment-method",
  );
}

function fixturePortalData(customerId?: string): PortalData {
  const totalSeats = billingSnapshot.seats.basic + billingSnapshot.seats.additional;
  return {
    subscription: {
      id: dashboardObjects.subscriptions[0]?.id ?? "sub_test_29",
      customerId: customerId ?? fixtureCheckoutSession.customer,
      customerEmail: billingSnapshot.owner,
      workspace: billingSnapshot.workspace,
      owner: billingSnapshot.owner,
      planId: "team",
      planName: billingSnapshot.plan,
      price: "$99/mo",
      status: billingSnapshot.status,
      period: billingSnapshot.period,
      seats: totalSeats,
      usedSeats: billingSnapshot.seats.used,
      pendingSeats: billingSnapshot.seats.pending,
      cancelAtPeriodEnd: false,
      exportRemaining: billingSnapshot.export.remaining,
      exportManualLimit: billingSnapshot.export.manualLimit,
      exportRenewal: billingSnapshot.export.renews,
    },
    plans: portalPlans,
    invoices: [...fixtureInvoices],
    paymentMethod: {
      id: "pm_fixture_card",
      label: "Visa sandbox ending 4242",
      status: "saved",
    },
  };
}

async function loadPortalDataFromStripeLikeApi(
  customerId: string | undefined,
  fallback: PortalData,
): Promise<PortalData | undefined> {
  const [customersResult, subscriptionsResult, invoicesResult, paymentIntentsResult] = await Promise.allSettled([
    customerId ? getJSON(`/v1/customers/${encodeURIComponent(customerId)}`) : getJSON("/v1/customers"),
    getJSON("/v1/subscriptions"),
    getJSON("/v1/invoices"),
    getJSON("/v1/payment_intents"),
  ]);
  if (customersResult.status !== "fulfilled" || subscriptionsResult.status !== "fulfilled") return undefined;

  const customerRecord = customerId
    ? asRecord(customersResult.value)
    : asRecord(readCollection(customersResult.value, ["data"])[0]);
  if (!customerRecord) return undefined;

  const resolvedCustomerId = readString(customerRecord, ["id"], fallback.subscription.customerId);
  const subscriptions: Record<string, unknown>[] = readCollection(subscriptionsResult.value, ["data"])
    .flatMap((item) => {
      const record = asRecord(item);
      return record ? [record] : [];
    })
    .filter((item) => readObjectId(item.customer) === resolvedCustomerId || readString(item, ["customer"], "") === resolvedCustomerId);
  const invoices: Record<string, unknown>[] = invoicesResult.status === "fulfilled"
    ? readCollection(invoicesResult.value, ["data"]).flatMap((item) => {
        const record = asRecord(item);
        return record ? [record] : [];
      })
    : [];
  const paymentIntents: Record<string, unknown>[] = paymentIntentsResult.status === "fulfilled"
    ? readCollection(paymentIntentsResult.value, ["data"]).flatMap((item) => {
        const record = asRecord(item);
        return record ? [record] : [];
      })
    : [];

  return normalizePortalData(
    {
      customer: customerRecord,
      subscription: subscriptions[0],
      invoices: invoices.filter((invoice) => readString(invoice, ["customer"], "") === resolvedCustomerId),
      payment_intents: paymentIntents.filter((intent) => readString(intent, ["customer"], "") === resolvedCustomerId),
    },
    fallback,
  );
}

async function postPortalAction(
  path: string,
  payload: unknown,
  fallback: PortalData,
  action: PortalAction,
): Promise<PortalActionResult> {
  try {
    const body = await postJSON(path, payload);
    return {
      data: normalizePortalData(body, fallback),
      source: "api",
      action,
      message: portalActionMessage(action, "api"),
    };
  } catch (error) {
    if (!isEndpointUnavailable(error)) {
      return {
        data: fallback,
        source: "fixture",
        action,
        message: portalActionMessage(action, "fixture"),
        error: errorMessage(error),
      };
    }
    return {
      data: fallback,
      source: "fixture",
      action,
      message: portalActionMessage(action, "fixture"),
    };
  }
}

function normalizePortalData(value: unknown, fallback: PortalData): PortalData {
  const root = asRecord(value);
  const state = firstRecord(root, ["state", "portal_state", "data"]) ?? root ?? {};
  const customer = firstRecord(state, ["customer"]);
  const subscription = firstRecord(state, ["subscription"]);
  const summary = firstRecord(state, ["summary"]);
  const invoicesSource = readCollection(state, ["invoices", "invoice_history"]);
  const paymentMethods = firstRecord(root, ["payment_method"]);

  const items = readArray(subscription, ["items", "line_items"]) ?? [];
  const firstItem = asRecord(items[0]);
  const planId =
    readString(subscription, ["metadata_plan", "plan"], "") ||
    readString(firstRecord(subscription, ["metadata"]), ["plan"], "") ||
    planIdForPrice(readString(firstItem, ["price"], ""));
  const plan = portalPlans.find((item) => item.id === planId) ?? portalPlans.find((item) => item.id === fallback.subscription.planId) ?? portalPlans[0];
  const seats = readNumber(firstItem, ["quantity"], fallback.subscription.seats);
  const periodStart = readString(subscription, ["current_period_start", "currentPeriodStart"], "");
  const periodEnd = readString(subscription, ["current_period_end", "currentPeriodEnd"], "");
  const cancelAtPeriodEnd = readBool(subscription?.cancel_at_period_end) || readBool(subscription?.cancelAtPeriodEnd) || readBool(summary?.cancel_at_period_end);

  return {
    subscription: {
      ...fallback.subscription,
      id: readString(subscription, ["id"], fallback.subscription.id),
      customerId:
        readString(subscription, ["customer", "customer_id", "customerId"], "") ||
        readString(customer, ["id"], fallback.subscription.customerId),
      customerEmail: readString(customer, ["email"], fallback.subscription.customerEmail),
      owner: readString(customer, ["email"], fallback.subscription.owner),
      planId: plan.id,
      planName: plan.name,
      price: plan.price,
      status:
        readString(subscription, ["status"], "") ||
        readString(summary, ["subscription_status", "subscriptionStatus"], fallback.subscription.status),
      period: formatPeriod(periodStart, periodEnd, fallback.subscription.period),
      seats,
      cancelAtPeriodEnd,
      canceledAt: readString(subscription, ["canceled_at", "canceledAt"], "") || undefined,
    },
    plans: portalPlans,
    invoices: invoicesSource.length > 0 ? invoicesSource.flatMap((item) => normalizePortalInvoice(item)) : fallback.invoices,
    paymentMethod: {
      id: readString(paymentMethods, ["payment_method", "paymentMethod", "id"], fallback.paymentMethod.id),
      label: paymentMethodLabel(paymentMethods, fallback.paymentMethod.label),
      status: readString(paymentMethods, ["status", "outcome"], fallback.paymentMethod.status),
    },
  };
}

function normalizePortalInvoice(value: unknown): Invoice[] {
  const record = asRecord(value);
  if (!record) return [];
  const id = readString(record, ["id"], "");
  if (!id) return [];
  return [{
    id,
    period: formatDateField(record, ["created_at", "createdAt", "created"]) ?? "Invoice",
    amount:
      readMoney(record, ["total", "amount_paid", "amount_due", "subtotal"], ["currency"]) ??
      "$0.00",
    status: readString(record, ["status"], "draft"),
  }];
}

function simulatePortalPlanChange(data: PortalData, planId: string): PortalData {
  const plan = portalPlans.find((item) => item.id === planId) ?? data.plans[0];
  return {
    ...data,
    subscription: {
      ...data.subscription,
      planId: plan.id,
      planName: plan.name,
      price: plan.price,
      status: data.subscription.status === "canceled" ? "active" : data.subscription.status,
    },
  };
}

function simulatePortalSeatChange(data: PortalData, seats: number): PortalData {
  return {
    ...data,
    subscription: {
      ...data.subscription,
      seats,
      usedSeats: Math.min(data.subscription.usedSeats, seats),
    },
  };
}

function simulatePortalCancellation(data: PortalData, mode: PortalCancelMode): PortalData {
  return {
    ...data,
    subscription: {
      ...data.subscription,
      status: mode === "immediate" ? "canceled" : data.subscription.status,
      cancelAtPeriodEnd: mode === "period",
      canceledAt: mode === "immediate" ? new Date().toISOString() : undefined,
    },
  };
}

function simulatePortalResume(data: PortalData): PortalData {
  return {
    ...data,
    subscription: {
      ...data.subscription,
      status: data.subscription.status === "canceled" ? "active" : data.subscription.status,
      cancelAtPeriodEnd: false,
      canceledAt: undefined,
    },
  };
}

function simulatePortalPaymentMethod(data: PortalData, outcome: PortalPaymentOutcome): PortalData {
  return {
    ...data,
    paymentMethod: {
      id: outcome === "succeeds" ? "pm_fixture_updated" : data.paymentMethod.id,
      label: outcome === "succeeds" ? "Sandbox card saved" : "Sandbox card declined",
      status: outcome === "succeeds" ? "succeeded" : "failed",
    },
  };
}

function portalActionMessage(action: PortalAction, source: DataSource): string {
  const prefix = source === "api" ? "Portal API" : "Fixture portal";
  switch (action) {
    case "plan-change":
      return `${prefix} plan change applied`;
    case "seat-change":
      return `${prefix} seat quantity applied`;
    case "cancel":
      return `${prefix} cancellation applied`;
    case "resume":
      return `${prefix} resume applied`;
    case "payment-method":
      return `${prefix} payment method simulation applied`;
  }
}

function planIdForPrice(priceId: string): string {
  const lower = priceId.toLowerCase();
  if (lower.includes("starter")) return "starter";
  if (lower.includes("scale")) return "scale";
  return "team";
}

function formatPeriod(start: string, end: string, fallback: string): string {
  if (!start && !end) return fallback;
  return [formatDateOnly(start), formatDateOnly(end)].filter(Boolean).join(" - ") || fallback;
}

function formatDateOnly(value: string): string {
  if (!value) return "";
  const date = new Date(value);
  if (Number.isNaN(date.valueOf())) return value;
  return date.toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

function paymentMethodLabel(record: Record<string, unknown> | undefined, fallback: string): string {
  const id = readString(record, ["payment_method", "paymentMethod", "id"], "");
  if (!id) return fallback;
  const status = readString(record, ["status"], "");
  return status === "failed" ? "Sandbox card declined" : `${id} saved`;
}

function readBool(value: unknown): boolean {
  if (typeof value === "boolean") return value;
  if (typeof value === "number") return value !== 0;
  if (typeof value === "string") return value === "true" || value === "1";
  return false;
}

export async function loadDashboardData(search: string = window.location.search): Promise<LoadResult<DashboardData>> {
  const query = filterDashboardQuery(search);
  const [timeline, attempts, customers, subscriptions, invoicesResult, paymentIntents, checkoutSessions, events] =
    await Promise.allSettled([
    getJSON(`/api/timeline${query}`),
    getJSON(`/api/delivery-attempts${query}`),
      getJSON("/v1/customers"),
      getJSON("/v1/subscriptions"),
      getJSON("/v1/invoices"),
      getJSON("/v1/payment_intents"),
      getJSON("/v1/checkout/sessions"),
      getJSON("/v1/events"),
    ]);

  const apiTimeline = timeline.status === "fulfilled" ? normalizeTimelineEntries(timeline.value) : undefined;
  const apiAttempts = attempts.status === "fulfilled" ? normalizeWebhookAttempts(attempts.value) : undefined;
  const apiObjects: Partial<Record<DashboardObjectType, DashboardObjectRecord[]>> = {
    customers: customers.status === "fulfilled" ? normalizeObjectCollection(customers.value, "customers") : undefined,
    subscriptions:
      subscriptions.status === "fulfilled" ? normalizeObjectCollection(subscriptions.value, "subscriptions") : undefined,
    invoices: invoicesResult.status === "fulfilled" ? normalizeObjectCollection(invoicesResult.value, "invoices") : undefined,
    paymentIntents:
      paymentIntents.status === "fulfilled" ? normalizeObjectCollection(paymentIntents.value, "paymentIntents") : undefined,
    checkoutSessions:
      checkoutSessions.status === "fulfilled"
        ? normalizeObjectCollection(checkoutSessions.value, "checkoutSessions")
        : undefined,
    webhookEvents: events.status === "fulfilled" ? normalizeObjectCollection(events.value, "webhookEvents") : undefined,
  };
  const objects = mergeDashboardObjects(apiObjects);
  const hasApiObjects = Object.values(apiObjects).some((items) => items && items.length > 0);
  const source: DataSource = apiTimeline || apiAttempts || hasApiObjects ? "api" : "fixture";
  const data: DashboardData = {
    billingSnapshot,
    dashboardStats: source === "api" ? liveDashboardStats(apiTimeline, apiAttempts) : dashboardStats,
    objects,
    timelineEntries: apiTimeline ?? timelineEntries,
    webhookAttempts: apiAttempts ?? webhookAttempts,
    appResponses: appResponsesFromAttempts(apiAttempts ?? webhookAttempts),
    evidence: dashboardEvidence(apiAttempts),
    scenario: {
      currentStep: "retry invoice payment",
      clockAdvance: "+3 days",
      assertions: "2 pending",
      tone: "warning",
    },
  };

  if (source === "api") return { data, source };

  const error = [timeline, attempts, customers, subscriptions, invoicesResult, paymentIntents, checkoutSessions, events]
    .map((result) => (result.status === "rejected" ? errorMessage(result.reason) : ""))
    .filter(Boolean)
    .join("; ");

  return { data, source, error: error || "Dashboard API not available" };
}

function filterDashboardQuery(search: string): string {
  const input = new URLSearchParams(search);
  const output = new URLSearchParams();
  for (const key of [
    "customerId",
    "subscriptionId",
    "checkoutSessionId",
    "invoiceId",
    "paymentIntentId",
    "eventId",
    "scenarioRunId",
  ]) {
    const value = input.get(key);
    if (value) output.set(key, value);
  }

  const rendered = output.toString();
  return rendered ? `?${rendered}` : "";
}

async function getJSON(path: string): Promise<unknown> {
  return requestJSON(path, { method: "GET" });
}

async function postJSON(path: string, body: unknown): Promise<unknown> {
  return requestJSON(path, {
    method: "POST",
    body: JSON.stringify(body),
    headers: { "Content-Type": "application/json" },
  });
}

async function requestJSON(path: string, init: RequestInit): Promise<unknown> {
  const response = await fetch(path, {
    ...init,
    headers: {
      Accept: "application/json",
      ...init.headers,
    },
  });

  if (!response.ok) {
    throw new ApiError(`${response.status} ${response.statusText}`, response.status);
  }

  if (response.status === 204) {
    return {};
  }

  const contentType = response.headers.get("content-type") ?? "";
  if (!contentType.includes("application/json")) {
    const text = await response.text();
    if (!text.trim()) return {};
    return { message: text.trim() };
  }

  return response.json();
}

class ApiError extends Error {
  constructor(message: string, readonly status: number) {
    super(message);
  }
}

function isEndpointUnavailable(error: unknown): boolean {
  return (
    error instanceof ApiError &&
    ([404, 405, 501].includes(error.status) || error.message.startsWith("Expected JSON"))
  );
}

function normalizeCheckoutSession(value: unknown, fallbackId: string): CheckoutSessionSummary {
  const root = asRecord(value);
  const session =
    firstRecord(root, ["session", "checkout_session", "checkoutSession", "data"]) ?? root ?? {};
  const customer = firstRecord(session, ["customer"]);
  const price = firstRecord(session, ["price"]);
  const subscription = firstRecord(root, ["subscription"]) ?? firstRecord(session, ["subscription"]);
  const invoice = firstRecord(root, ["invoice"]) ?? firstRecord(session, ["invoice", "latest_invoice"]);
  const paymentIntent =
    firstRecord(root, ["payment_intent", "paymentIntent"]) ??
    firstRecord(session, ["payment_intent", "paymentIntent"]);

  return {
    id: readString(session, ["id"], fallbackId),
    customer: readObjectId(session.customer) ?? readString(session, ["customer_id", "customer"], fixtureSession.customer),
    customerEmail:
      readString(customer, ["email"], "") ||
      readString(session, ["customer_email", "customerEmail"], fixtureSession.customerEmail),
    plan:
      readString(price, ["nickname", "lookup_key"], "") ||
      readString(session, ["plan", "plan_name", "planName"], fixtureSession.plan),
    price:
      readMoney(session, ["amount_total", "amountSubtotal", "amount_total"], ["currency"]) ||
      readString(session, ["price", "amount"], fixtureSession.price),
    status: readString(session, ["status"], fixtureSession.status),
    paymentStatus: readString(session, ["payment_status", "paymentStatus"], fixtureSession.paymentStatus),
    subscriptionId: readObjectId(session.subscription) ?? readString(subscription, ["id"], fixtureSession.subscriptionId),
    subscriptionStatus: readString(subscription, ["status"], fixtureSession.subscriptionStatus),
    invoiceId:
      readObjectId(session.invoice) ??
      readObjectId(session.latest_invoice) ??
      readString(invoice, ["id"], fixtureSession.invoiceId),
    invoiceStatus: readString(invoice, ["status"], fixtureSession.invoiceStatus),
    paymentIntentId:
      readObjectId(session.payment_intent) ??
      readString(paymentIntent, ["id"], fixtureSession.paymentIntentId),
    paymentIntentStatus: readString(paymentIntent, ["status"], fixtureSession.paymentIntentStatus),
    returnUrl: readString(session, ["return_url", "returnUrl", "success_url", "successUrl"], fixtureSession.returnUrl),
    lineItems: normalizeLineItems(root, session),
  };
}

function normalizeLineItems(root: Record<string, unknown> | undefined, session: Record<string, unknown>): CheckoutLineItem[] {
  const source = readArray(root, ["line_items", "lineItems"]) ?? readArray(session, ["line_items", "lineItems"]);
  if (!source) return [...fixtureSession.lineItems];

  const items = source.flatMap((item) => {
    const record = asRecord(item);
    if (!record) return [];
    const price = firstRecord(record, ["price"]);
    return {
      label:
        readString(record, ["label", "description", "name"], "") ||
        readString(price, ["nickname", "lookup_key"], "Line item"),
      amount:
        readMoney(record, ["amount_total", "amount", "unit_amount"], ["currency"]) ||
        readString(record, ["display_amount"], fixtureSession.price),
    };
  });

  return items.length > 0 ? items : [...fixtureSession.lineItems];
}

function normalizeTimelineEntries(value: unknown): TimelineEntry[] {
  const source = readCollection(value, ["entries", "timeline", "data", "items"]);
  return source.flatMap((item) => {
    const record = asRecord(item);
    if (!record) return [];
    const delivery = readString(record, ["delivery", "delivery_status", "deliveryStatus", "status"], "-");
    const response = readString(record, ["response", "app_response", "appResponse", "response_status"], "-");
    return {
      time: formatTime(readString(record, ["time", "timestamp", "created_at", "created"], "")),
      object: readString(record, ["object_id", "objectId", "subject", "object"], "-"),
      transition: readString(record, ["transition", "message", "action"], "-"),
      event: readString(record, ["event", "event_type", "eventType", "type", "action"], "-"),
      delivery,
      response,
      tone: toneForStatus(`${delivery} ${response}`),
    };
  });
}

function normalizeWebhookAttempts(value: unknown): WebhookAttempt[] {
  const source = readCollection(value, ["attempts", "delivery_attempts", "deliveryAttempts", "data", "items"]);
  return source.flatMap((item, index) => {
    const record = asRecord(item);
    if (!record) return [];
    const request = firstRecord(record, ["request"]);
    const response = firstRecord(record, ["response"]);
    const event = firstRecord(record, ["event", "webhook_event", "webhookEvent"]);
    const endpoint = firstRecord(record, ["endpoint", "webhook_endpoint", "webhookEndpoint"]);
    const requestHeaders =
      firstRecord(request, ["headers"]) ??
      firstRecord(record, ["request_headers", "requestHeaders", "headers"]);
    const responseHeaders = firstRecord(response, ["headers"]);
    const attemptNumber = readNumber(record, ["attempt_number", "attemptNumber", "attempt"], Number.NaN);
    const responseStatus =
      readString(record, ["response_status", "responseStatus", "status_code", "statusCode"], "") ||
      readString(response, ["status", "status_code", "statusCode"], "");
    const errorMessage =
      readString(record, ["error", "error_message", "errorMessage"], "") ||
      readString(response, ["error", "error_message", "errorMessage"], "");
    const status = readString(record, ["status", "delivery_status", "deliveryStatus"], "-");
    const signatureHeader =
      readString(record, ["signature", "signature_header", "signatureHeader"], "") ||
      readHeader(requestHeaders, ["stripe-signature", "billtap-signature"]) ||
      readHeader(responseHeaders, ["stripe-signature", "billtap-signature"]);

    return {
      id: readString(record, ["id"], "") || undefined,
      eventId:
        readString(record, ["event_id", "eventId"], "") ||
        readEventId(record.event) ||
        readString(event, ["id"], "") ||
        undefined,
      endpointId:
        readString(record, ["endpoint_id", "endpointId"], "") ||
        readObjectId(record.endpoint) ||
        readString(endpoint, ["id"], "") ||
        undefined,
      endpoint:
        redactEvidenceUrl(readString(record, ["endpoint_url", "endpointUrl", "request_url", "requestUrl", "url"], "")) ||
        redactEvidenceUrl(readString(request, ["url"], "")) ||
        readString(endpoint, ["url"], "-"),
      event:
        readString(record, ["event", "event_type", "eventType", "type"], "") ||
        readString(event, ["type", "event_type", "eventType"], "-"),
      attemptNumber: Number.isNaN(attemptNumber) ? undefined : attemptNumber,
      attempts: readNumber(record, ["attempts", "attempt_count", "attemptCount"], Number.isNaN(attemptNumber) ? 1 : attemptNumber),
      next: formatTime(readString(record, ["next", "next_retry_at", "nextRetryAt"], "")) || "-",
      status,
      tone: toneForStatus(`${status} ${responseStatus} ${errorMessage}`),
      scheduledAt: optionalTime(record, ["scheduled_at", "scheduledAt", "scheduled"]),
      deliveredAt: optionalTime(record, ["delivered_at", "deliveredAt", "delivered"]),
      requestUrl:
        redactEvidenceUrl(readString(record, ["request_url", "requestUrl"], "")) ||
        redactEvidenceUrl(readString(request, ["url"], "")) ||
        undefined,
      requestHeaders: formatHeaders(requestHeaders),
      requestBody:
        redactEvidenceText(readString(record, ["request_body", "requestBody"], "")) ||
        redactEvidenceText(readString(request, ["body", "raw_body", "rawBody"], "")) ||
        undefined,
      signatureHeader: signatureHeader ? redactSignature(signatureHeader) : undefined,
      responseStatus: responseStatus || undefined,
      responseBody:
        redactEvidenceText(readString(record, ["response_body", "responseBody"], "")) ||
        redactEvidenceText(readString(response, ["body", "raw_body", "rawBody"], "")) ||
        undefined,
      errorMessage: errorMessage || undefined,
      retryPlan:
        readString(record, ["retry_plan", "retryPlan", "retry_policy", "retryPolicy"], "") ||
        retrySummary(record) ||
        undefined,
    };
  });
}

function normalizeObjectCollection(value: unknown, type: DashboardObjectType): DashboardObjectRecord[] {
  const source = readCollection(value, ["data", "items", type]);
  return source.flatMap((item) => {
    const record = asRecord(item);
    if (!record) return [];
    return normalizeObjectRecord(record, type);
  });
}

function normalizeObjectRecord(record: Record<string, unknown>, type: DashboardObjectType): DashboardObjectRecord[] {
  const id = readString(record, ["id"], "");
  if (!id) return [];

  const status = objectStatus(record, type);
  const amount = objectAmount(record, type);
  const createdAt = optionalTime(record, ["created_at", "createdAt", "created", "completed_at", "completedAt"]);
  const customer = readObjectId(record.customer) ?? readString(record, ["customer", "customer_id", "customerId"], "");
  const subscription =
    readObjectId(record.subscription) ??
    readString(record, ["subscription", "subscription_id", "subscriptionId"], "");
  const invoice =
    readObjectId(record.invoice) ??
    readString(record, ["invoice", "invoice_id", "invoiceId", "latest_invoice", "latestInvoice"], "");
  const paymentIntent =
    readObjectId(record.payment_intent) ??
    readString(record, ["payment_intent", "paymentIntent", "payment_intent_id", "paymentIntentId"], "");

  if (type === "customers") {
    const email = readString(record, ["email"], "unknown customer");
    return [
      {
        id,
        type,
        label: "Customer",
        title: email,
        subtitle: readString(record, ["name"], id),
        status: status || "created",
        tone: toneForStatus(status || "created"),
        createdAt,
        fields: compactFields([
          ["Customer ID", id],
          ["Name", readString(record, ["name"], "-")],
          ["Email", email],
          ["Created", createdAt],
        ]),
      },
    ];
  }

  if (type === "subscriptions") {
    return [
      {
        id,
        type,
        label: "Subscription",
        title: readString(record, ["plan", "description"], "") || id,
        subtitle: customer || "customer pending",
        status,
        tone: toneForStatus(status),
        createdAt,
        amount,
        fields: compactFields([
          ["Subscription ID", id],
          ["Customer", customer],
          ["Latest invoice", invoice],
          ["Current period start", formatDateField(record, ["current_period_start", "currentPeriodStart"])],
          ["Current period end", formatDateField(record, ["current_period_end", "currentPeriodEnd"])],
          ["Cancel at period end", readString(record, ["cancel_at_period_end", "cancelAtPeriodEnd"], "")],
        ]),
      },
    ];
  }

  if (type === "invoices") {
    return [
      {
        id,
        type,
        label: "Invoice",
        title: amount || id,
        subtitle: subscription || customer || "billing object",
        status,
        tone: toneForStatus(status),
        createdAt,
        amount,
        fields: compactFields([
          ["Invoice ID", id],
          ["Customer", customer],
          ["Subscription", subscription],
          ["Payment intent", paymentIntent],
          ["Amount due", readMoney(record, ["amount_due", "amountDue"], ["currency"])],
          ["Amount paid", readMoney(record, ["amount_paid", "amountPaid"], ["currency"])],
          ["Attempt count", readString(record, ["attempt_count", "attemptCount"], "")],
          ["Next payment attempt", formatDateField(record, ["next_payment_attempt", "nextPaymentAttempt"])],
        ]),
      },
    ];
  }

  if (type === "paymentIntents") {
    return [
      {
        id,
        type,
        label: "Payment intent",
        title: amount || id,
        subtitle: invoice || customer || "payment object",
        status,
        tone: toneForStatus(status),
        createdAt,
        amount,
        fields: compactFields([
          ["Payment intent ID", id],
          ["Customer", customer],
          ["Invoice", invoice],
          ["Amount", amount],
          ["Failure code", readString(record, ["failure_code", "failureCode"], "none")],
          ["Failure message", readString(record, ["failure_message", "failureMessage"], "none")],
        ]),
      },
    ];
  }

  if (type === "checkoutSessions") {
    return [
      {
        id,
        type,
        label: "Checkout session",
        title: readString(record, ["mode"], "subscription"),
        subtitle: customer || "customer pending",
        status,
        tone: toneForStatus(status),
        createdAt,
        amount,
        fields: compactFields([
          ["Checkout session ID", id],
          ["Customer", customer],
          ["Subscription", subscription],
          ["Invoice", invoice],
          ["Payment intent", paymentIntent],
          ["Payment status", readString(record, ["payment_status", "paymentStatus"], "")],
          ["Hosted URL", readString(record, ["url"], "")],
        ]),
      },
    ];
  }

  const eventType = readString(record, ["type", "event_type", "eventType"], "event");
  const request = firstRecord(record, ["request"]);
  const billtap = firstRecord(record, ["billtap"]);
  const data = firstRecord(record, ["data"]);
  const objectPayload = asRecord(data?.object);
  const relatedObject = readObjectId(data?.object) ?? readString(objectPayload, ["id"], "");
  return [
    {
      id,
      type,
      label: "Webhook event",
      title: eventType,
      subtitle: relatedObject || readString(billtap, ["source"], "event payload"),
      status: readString(record, ["status"], "") || `${readString(record, ["pending_webhooks", "pendingWebhooks"], "0")} pending`,
      tone: toneForStatus(`${eventType} ${readString(record, ["pending_webhooks", "pendingWebhooks"], "")}`),
      createdAt: optionalTime(record, ["created_at", "createdAt", "created"]),
      fields: compactFields([
        ["Event ID", id],
        ["Event type", eventType],
        ["Related object", relatedObject],
        ["Source", readString(billtap, ["source"], "")],
        ["Sequence", readString(billtap, ["sequence"], "")],
        ["Request ID", readString(request, ["id"], "")],
        ["Idempotency key", readString(request, ["idempotency_key", "idempotencyKey"], "")],
      ]),
    },
  ];
}

function mergeDashboardObjects(
  apiObjects: Partial<Record<DashboardObjectType, DashboardObjectRecord[]>>,
): Record<DashboardObjectType, DashboardObjectRecord[]> {
  return {
    customers: nonEmpty(apiObjects.customers, dashboardObjects.customers),
    subscriptions: nonEmpty(apiObjects.subscriptions, dashboardObjects.subscriptions),
    invoices: nonEmpty(apiObjects.invoices, dashboardObjects.invoices),
    paymentIntents: nonEmpty(apiObjects.paymentIntents, dashboardObjects.paymentIntents),
    checkoutSessions: nonEmpty(apiObjects.checkoutSessions, dashboardObjects.checkoutSessions),
    webhookEvents: nonEmpty(apiObjects.webhookEvents, dashboardObjects.webhookEvents),
  };
}

function nonEmpty<T>(value: T[] | undefined, fallback: T[]): T[] {
  return value && value.length > 0 ? value : fallback;
}

function objectStatus(record: Record<string, unknown>, type: DashboardObjectType): string {
  if (type === "checkoutSessions") {
    const status = readString(record, ["status"], "");
    const paymentStatus = readString(record, ["payment_status", "paymentStatus"], "");
    return [status, paymentStatus].filter(Boolean).join(" / ") || "created";
  }
  return readString(record, ["status"], "") || (readString(record, ["active"], "") === "true" ? "active" : "created");
}

function objectAmount(record: Record<string, unknown>, type: DashboardObjectType): string | undefined {
  if (type === "invoices") {
    return (
      readMoney(record, ["total", "amount_due", "amountDue", "subtotal"], ["currency"]) ||
      undefined
    );
  }
  if (type === "paymentIntents") {
    return readMoney(record, ["amount"], ["currency"]) || undefined;
  }
  return readMoney(record, ["amount_total", "amountTotal", "total"], ["currency"]) || undefined;
}

function compactFields(items: Array<[string, string | undefined]>): { label: string; value: string }[] {
  return items
    .filter(([, value]) => value !== undefined && value !== "")
    .map(([label, value]) => ({ label, value: value ?? "-" }));
}

function formatDateField(record: Record<string, unknown>, keys: string[]): string | undefined {
  const value = readString(record, keys, "");
  if (!value) return undefined;
  return formatTime(value);
}

function appResponsesFromAttempts(attempts: WebhookAttempt[]): AppResponseRecord[] {
  return attempts
    .filter((attempt) => attempt.responseStatus || attempt.responseBody || attempt.errorMessage)
    .map((attempt) => ({
      id: attempt.id ?? `${attempt.eventId ?? attempt.event}-${attempt.attemptNumber ?? attempt.attempts}`,
      eventId: attempt.eventId ?? attempt.event,
      endpoint: attempt.endpoint,
      status: attempt.responseStatus ?? attempt.status,
      body: attempt.responseBody ?? attempt.errorMessage ?? "No response body recorded",
      retryPlan: attempt.retryPlan ?? "retry policy default",
      tone: attempt.tone ?? toneForStatus(`${attempt.status} ${attempt.responseStatus ?? ""}`),
    }));
}

function dashboardEvidence(attempts?: WebhookAttempt[]): DashboardData["evidence"] {
  const attempt = attempts?.find((item) => item.signatureHeader || item.errorMessage || item.responseBody);
  return {
    signatureHeader: attempt?.signatureHeader ?? "t=1778233312,v1=sandbox_mismatch",
    idempotencyKey: attempt?.eventId ? `billtap:${attempt.eventId}:attempt_${attempt.attemptNumber ?? attempt.attempts}` : "billtap:evt_test_90:attempt_2",
    failureReason:
      attempt?.errorMessage ??
      attempt?.responseBody ??
      "handler returned 503 before entitlement update",
  };
}

function retrySummary(record: Record<string, unknown>): string {
  const maxAttempts = readNumber(record, ["max_attempts", "maxAttempts"], Number.NaN);
  const backoff = readArray(record, ["backoff", "retry_backoff", "retryBackoff"]);
  const parts = [];
  if (!Number.isNaN(maxAttempts)) parts.push(`max ${maxAttempts}`);
  if (backoff?.length) parts.push(`backoff ${backoff.map(String).join(", ")}`);
  return parts.join(" · ");
}

function optionalTime(record: Record<string, unknown>, keys: string[]): string | undefined {
  const value = readString(record, keys, "");
  return value ? formatTime(value) : undefined;
}

function formatHeaders(headers: Record<string, unknown> | undefined): string | undefined {
  if (!headers) return undefined;
  const lines = Object.entries(headers).flatMap(([key, value]) => {
    if (typeof value === "string" || typeof value === "number") return `${key}: ${redactHeaderValue(key, String(value))}`;
    if (Array.isArray(value)) return `${key}: ${value.map((part) => redactHeaderValue(key, String(part))).join(", ")}`;
    return [];
  });
  return lines.length > 0 ? lines.join("\n") : undefined;
}

function readHeader(headers: Record<string, unknown> | undefined, keys: string[]): string | undefined {
  if (!headers) return undefined;
  const match = Object.entries(headers).find(([key]) => keys.includes(key.toLowerCase()));
  const value = match?.[1];
  if (typeof value === "string") return redactHeaderValue(match?.[0] ?? "", value);
  if (Array.isArray(value)) return value.map((part) => redactHeaderValue(match?.[0] ?? "", String(part))).join(", ");
  return undefined;
}

function redactHeaderValue(key: string, value: string): string {
  const normalized = key.toLowerCase();
  if (normalized === "billtap-signature" || normalized === "stripe-signature") return redactSignature(value);
  if (
    normalized === "authorization" ||
    normalized === "cookie" ||
    normalized === "set-cookie" ||
    normalized.includes("api-key") ||
    normalized.includes("secret") ||
    normalized.includes("token")
  ) {
    return "****";
  }
  return value;
}

function redactSignature(value: string): string {
  return value.replace(/(v\d=)[^,\s]+/g, "$1****");
}

function redactEvidenceUrl(value: string): string {
  if (!value) return value;
  try {
    const parsed = new URL(value);
    for (const key of [...parsed.searchParams.keys()]) {
      const normalized = key.toLowerCase();
      if (normalized.includes("key") || normalized.includes("secret") || normalized.includes("token")) {
        parsed.searchParams.set(key, "****");
      }
    }
    return parsed.toString();
  } catch {
    return value;
  }
}

function redactEvidenceText(value: string): string {
  if (!value) return value;
  return value
    .replace(/("(?:client_)?secret"\s*:\s*)"[^"]*"/gi, '$1"****"')
    .replace(/("authorization"\s*:\s*)"[^"]*"/gi, '$1"****"')
    .replace(/("api_?key"\s*:\s*)"[^"]*"/gi, '$1"****"')
    .replace(/("card"\s*:\s*\{[^}]*"number"\s*:\s*)"[^"]*"/gi, '$1"****"');
}

function readCollection(value: unknown, keys: string[]): unknown[] {
  if (Array.isArray(value)) return value;
  const record = asRecord(value);
  const nested = record ? readArray(record, keys) : undefined;
  return nested ?? [];
}

function liveDashboardStats(entries?: TimelineEntry[], attempts?: WebhookAttempt[]): typeof dashboardStats {
  const retryCount = attempts?.filter((attempt) => /retry/i.test(attempt.status)).length ?? 0;
  return [
    {
      label: "Timeline events",
      value: String(entries?.length ?? 0),
      meta: "from /api/timeline",
      tone: "info",
    },
    {
      label: "Webhook attempts",
      value: String(attempts?.length ?? 0),
      meta: retryCount ? `${retryCount} retries scheduled` : "delivery API loaded",
      tone: retryCount ? "warning" : "good",
    },
    { label: "Scenario runs", value: "-", meta: "waiting for scenario API", tone: "neutral" },
    { label: "Debug bundles", value: "-", meta: "not requested", tone: "neutral" },
  ];
}

function simulateCheckoutCompletion(session: CheckoutSessionSummary, outcome: CheckoutOutcomeId): CheckoutSessionSummary {
  const next = { ...session };

  if (outcome === "success") {
    return {
      ...next,
      status: "complete",
      paymentStatus: "paid",
      subscriptionStatus: "active",
      invoiceStatus: "paid",
      paymentIntentStatus: "succeeded",
    };
  }

  if (outcome === "pending") {
    return {
      ...next,
      status: "complete",
      paymentStatus: "pending",
      subscriptionStatus: "incomplete",
      invoiceStatus: "open",
      paymentIntentStatus: "processing",
    };
  }

  if (outcome === "action") {
    return {
      ...next,
      status: "complete",
      paymentStatus: "unpaid",
      subscriptionStatus: "incomplete",
      invoiceStatus: "open",
      paymentIntentStatus: "requires_action",
    };
  }

  if (outcome === "cancel") {
    return {
      ...next,
      status: "expired",
      paymentStatus: "unpaid",
      subscriptionStatus: "not created",
      invoiceStatus: "void",
      paymentIntentStatus: "canceled",
    };
  }

  return {
    ...next,
    status: "complete",
    paymentStatus: "unpaid",
    subscriptionStatus: "incomplete",
    invoiceStatus: "open",
    paymentIntentStatus: "requires_payment_method",
  };
}

function apiOutcome(outcome: CheckoutOutcomeId): string {
  if (outcome === "declined") return "payment_failed";
  if (outcome === "funds") return "insufficient_funds";
  if (outcome === "expired") return "expired_card";
  if (outcome === "action") return "requires_action";
  if (outcome === "cancel") return "canceled";
  if (outcome === "pending") return "payment_pending";
  return "payment_succeeded";
}

function toneForStatus(value: string): StatusTone {
  const lower = value.toLowerCase();
  if (lower.includes("fail") || lower.includes("declin") || lower.includes("mismatch") || lower.includes("blocked")) {
    return "danger";
  }
  if (lower.includes("retry") || lower.includes("pending") || lower.includes("timeout")) return "warning";
  if (lower.includes("duplicate") || lower.includes("held") || lower.includes("action")) return "info";
  if (lower.includes("deliver") || lower.includes("success") || lower.includes("200")) return "good";
  return "neutral";
}

function formatTime(value: string): string {
  if (!value) return "-";
  const date = new Date(value);
  if (!Number.isNaN(date.valueOf())) {
    return date.toLocaleTimeString("en-US", {
      hour12: false,
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  }
  return value;
}

function readMoney(record: Record<string, unknown> | undefined, amountKeys: string[], currencyKeys: string[]): string | undefined {
  const amount = readNumber(record, amountKeys, Number.NaN);
  if (Number.isNaN(amount)) return undefined;
  const currency = readString(record, currencyKeys, "usd").toUpperCase();
  const normalized = amount > 999 ? amount / 100 : amount;
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency,
  }).format(normalized);
}

function readString(record: Record<string, unknown> | undefined, keys: string[], fallback: string): string {
  if (!record) return fallback;
  for (const key of keys) {
    const value = record[key];
    if (typeof value === "string" && value.trim()) return value;
    if (typeof value === "number") return String(value);
  }
  return fallback;
}

function readNumber(record: Record<string, unknown> | undefined, keys: string[], fallback: number): number {
  if (!record) return fallback;
  for (const key of keys) {
    const value = record[key];
    if (typeof value === "number") return value;
    if (typeof value === "string") {
      const parsed = Number(value);
      if (!Number.isNaN(parsed)) return parsed;
    }
  }
  return fallback;
}

function readArray(record: Record<string, unknown> | undefined, keys: string[]): unknown[] | undefined {
  if (!record) return undefined;
  for (const key of keys) {
    const value = record[key];
    if (Array.isArray(value)) return value;
    const nested = asRecord(value);
    const nestedData = readArray(nested, ["data"]);
    if (nestedData) return nestedData;
  }
  return undefined;
}

function firstRecord(record: Record<string, unknown> | undefined, keys: string[]): Record<string, unknown> | undefined {
  if (!record) return undefined;
  for (const key of keys) {
    const value = asRecord(record[key]);
    if (value) return value;
  }
  return undefined;
}

function readObjectId(value: unknown): string | undefined {
  if (typeof value === "string") return value;
  const record = asRecord(value);
  return record ? readString(record, ["id"], "") || undefined : undefined;
}

function readEventId(value: unknown): string | undefined {
  if (typeof value === "string") return value.startsWith("evt_") ? value : undefined;
  return readObjectId(value);
}

function asRecord(value: unknown): Record<string, unknown> | undefined {
  return value && typeof value === "object" && !Array.isArray(value) ? (value as Record<string, unknown>) : undefined;
}

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}
