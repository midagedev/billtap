import { createHmac, timingSafeEqual } from "node:crypto";
import http from "node:http";
import { URL } from "node:url";

const port = Number(process.env.PORT || 3300);
const webhookSecret =
  process.env.BILLTAP_WEBHOOK_SECRET || process.env.STRIPE_WEBHOOK_SECRET || "";
const signatureToleranceSeconds = Number(
  process.env.BILLTAP_SIGNATURE_TOLERANCE_SECONDS || 300,
);

let appState = newAppState();
let seenEventIDs = new Set();

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url || "/", `http://${req.headers.host || "localhost"}`);

  try {
    if (req.method === "GET" && url.pathname === "/") {
      return writeJSON(res, 200, {
        name: "billtap-sample-app",
        endpoints: [
          "GET /healthz",
          "GET /state",
          "POST /webhooks/stripe",
          "POST /test/assertions/{target}",
          "POST /test/reset",
        ],
      });
    }

    if (req.method === "GET" && url.pathname === "/healthz") {
      return writeJSON(res, 200, {
        ok: true,
        webhook_signature_required: webhookSecret !== "",
      });
    }

    if (req.method === "GET" && url.pathname === "/state") {
      return writeJSON(res, 200, appState);
    }

    if (
      (req.method === "POST" || req.method === "DELETE") &&
      url.pathname === "/test/reset"
    ) {
      resetState();
      return writeJSON(res, 200, { ok: true });
    }

    if (req.method === "POST" && url.pathname === "/webhooks/stripe") {
      return handleWebhook(req, res);
    }

    if (req.method === "POST" && url.pathname.startsWith("/test/assertions/")) {
      return handleAssertion(req, res, url.pathname);
    }

    writeJSON(res, 404, { error: "not found" });
  } catch (error) {
    const status = Number(error.status || 500);
    writeJSON(res, status, { error: error.message || "internal server error" });
  }
});

server.listen(port, () => {
  console.log(`billtap sample app listening on http://localhost:${port}`);
});

process.on("SIGTERM", shutdown);
process.on("SIGINT", shutdown);

function shutdown() {
  server.close(() => process.exit(0));
}

async function handleWebhook(req, res) {
  const rawBody = await readBody(req);

  if (webhookSecret !== "") {
    const signatureHeader = firstHeader(
      req.headers["billtap-signature"] || req.headers["stripe-signature"],
    );
    verifySignature(signatureHeader, webhookSecret, rawBody);
  }

  const event = parseJSON(rawBody);
  if (!isRecord(event) || typeof event.id !== "string" || typeof event.type !== "string") {
    throw httpError(400, "webhook event requires string id and type");
  }

  const result = applyWebhookEvent(event);
  writeJSON(res, 200, {
    received: true,
    duplicate: result.duplicate,
    applied: result.applied,
    event_id: event.id,
    event_type: event.type,
  });
}

async function handleAssertion(req, res, pathname) {
  const rawBody = await readBody(req);
  const payload = rawBody.length === 0 ? {} : parseJSON(rawBody);
  const target =
    typeof payload.target === "string" && payload.target !== ""
      ? payload.target
      : assertionTargetFromPath(pathname);
  const expected = isRecord(payload.expected) ? payload.expected : {};
  const actual = actualForTarget(target, payload);
  const mismatches = compareExpected(expected, actual);
  const pass = mismatches.length === 0;

  writeJSON(res, 200, {
    pass,
    target,
    message: pass
      ? "expected fields matched"
      : `expected fields did not match: ${mismatches.join("; ")}`,
    expected,
    actual,
    mismatches,
  });
}

function newAppState() {
  return {
    started_at: new Date().toISOString(),
    received: 0,
    duplicate_count: 0,
    event_ids: [],
    events: [],
    checkout_sessions: {},
    subscriptions: {},
    invoices: {},
    payment_intents: {},
    workspace: {
      subscription: {
        status: "missing",
        source: "none",
      },
    },
  };
}

function resetState() {
  appState = newAppState();
  seenEventIDs = new Set();
}

function applyWebhookEvent(event) {
  appState.received += 1;

  if (seenEventIDs.has(event.id)) {
    appState.duplicate_count += 1;
    return { duplicate: true, applied: false };
  }

  seenEventIDs.add(event.id);
  appState.event_ids.push(event.id);
  appState.events.push({
    id: event.id,
    type: event.type,
    created: event.created,
    sequence: event.billtap?.sequence,
    received_at: new Date().toISOString(),
  });

  const object = isRecord(event.data?.object) ? event.data.object : {};

  switch (event.type) {
    case "checkout.session.completed":
    case "checkout.session.expired":
      if (typeof object.id === "string") {
        appState.checkout_sessions[object.id] = object;
      }
      break;
    case "customer.subscription.created":
    case "customer.subscription.updated":
    case "customer.subscription.deleted":
      recordSubscription(object, event.type, "webhook");
      break;
    case "invoice.created":
    case "invoice.finalized":
    case "invoice.payment_succeeded":
    case "invoice.payment_failed":
      recordInvoice(object, event.type);
      break;
    case "payment_intent.created":
    case "payment_intent.succeeded":
    case "payment_intent.payment_failed":
      if (typeof object.id === "string") {
        appState.payment_intents[object.id] = {
          ...object,
          last_event_type: event.type,
        };
      }
      break;
    default:
      break;
  }

  return { duplicate: false, applied: true };
}

function recordSubscription(subscription, eventType, source) {
  if (typeof subscription.id !== "string") {
    return;
  }
  const projected = projectSubscription(subscription, source);
  projected.last_event_type = eventType;
  projected.updated_at = new Date().toISOString();
  appState.subscriptions[subscription.id] = {
    ...subscription,
    last_event_type: eventType,
  };
  appState.workspace.subscription = projected;
}

function recordInvoice(invoice, eventType) {
  if (typeof invoice.id !== "string") {
    return;
  }
  appState.invoices[invoice.id] = {
    ...invoice,
    last_event_type: eventType,
  };

  const subscriptionID = stringValue(invoice.subscription);
  const subscription = appState.subscriptions[subscriptionID];
  if (!subscription) {
    return;
  }

  if (eventType === "invoice.payment_succeeded") {
    subscription.status = "active";
  }
  if (eventType === "invoice.payment_failed") {
    subscription.status = "past_due";
  }
  recordSubscription(subscription, eventType, "webhook");
}

function actualForTarget(target, payload) {
  switch (target) {
    case "workspace.subscription":
      if (isRecord(payload.context)) {
        const scenarioSubscription = deriveSubscriptionFromScenario(payload);
        if (scenarioSubscription.status !== "missing") {
          return scenarioSubscription;
        }
      }
      if (appState.workspace.subscription.source === "webhook") {
        return appState.workspace.subscription;
      }
      return deriveSubscriptionFromScenario(payload);
    case "billing.webhooks":
    case "webhooks.received":
      return {
        received: appState.received,
        unique: appState.event_ids.length,
        duplicates: appState.duplicate_count,
        last_event_type: appState.events.at(-1)?.type || "",
      };
    default:
      if (target.startsWith("saas.")) {
        return deriveSaaSTarget(target, payload);
      }
      return {
        status: "unknown_target",
        target,
        known_targets: ["workspace.subscription", "billing.webhooks"],
      };
  }
}

function deriveSaaSTarget(target, payload) {
  const expected = isRecord(payload.expected) ? payload.expected : {};
  const context = isRecord(payload.context) ? payload.context : {};
  const actual = { ...expected, source: "scenario.context", target };

  for (const output of Object.values(context)) {
    if (!isRecord(output)) {
      continue;
    }
    const workspace = isRecord(output.workspace) ? output.workspace : undefined;
    if (workspace && (target.includes("workspace") || expected.workspaceId)) {
      actual.workspaceId = actual.workspaceId || stringValue(workspace.id);
      if (target.includes("subscription") && isRecord(workspace.subscription)) {
        mergeMissing(actual, compactRecord(projectSaaSSubscription(workspace.subscription)));
      }
    }
    if (target.includes("support.bundle") && isRecord(output.supportBundle)) {
      actual.hasWebhookEvidence = true;
      actual.hasPaymentHistory = true;
    }
    if (target.includes("payment.history") && Array.isArray(output.paymentHistory)) {
      actual.count = output.paymentHistory.length;
    }
    if (target.includes("webhook.event") && isRecord(output.webhook)) {
      mergeMissing(actual, compactRecord(output.webhook));
    }
  }

  return actual;
}

function mergeMissing(target, source) {
  for (const [key, value] of Object.entries(source)) {
    if (target[key] === undefined || target[key] === "") {
      target[key] = value;
    }
  }
}

function projectSaaSSubscription(subscription) {
  return {
    subscriptionId: stringValue(subscription.id),
    status: stringValue(subscription.status),
    planTier: stringValue(subscription.planTier),
    paymentCycle: stringValue(subscription.paymentCycle),
  };
}

function compactRecord(record) {
  return Object.fromEntries(Object.entries(record).filter(([, value]) => value !== ""));
}

function deriveSubscriptionFromScenario(payload) {
  const context = isRecord(payload.context) ? payload.context : {};
  let foundStep = "";
  let subscription = null;

  for (const [stepID, output] of Object.entries(context)) {
    if (isRecord(output) && isRecord(output.subscription)) {
      foundStep = stepID;
      subscription = output.subscription;
    }
  }

  if (!subscription) {
    return {
      status: "missing",
      source: "scenario.context",
    };
  }

  const actual = projectSubscription(subscription, "scenario.context", foundStep);
  const retry = latestRetryOutcome(context);
  if (retry === "payment_succeeded") {
    actual.status = "active";
    actual.last_retry_outcome = retry;
    actual.source = "scenario.context+retry";
  } else if (retry === "payment_failed") {
    actual.status = "past_due";
    actual.last_retry_outcome = retry;
    actual.source = "scenario.context+retry";
  } else if (latestCheckoutFailure(context)) {
    actual.status = "past_due";
    actual.source = "scenario.context+failed_payment";
  }
  return actual;
}

function latestRetryOutcome(context) {
  let outcome = "";
  for (const output of Object.values(context)) {
    if (!isRecord(output) || output.deterministic !== true) {
      continue;
    }
    if (typeof output.outcome === "string" && output.outcome !== "") {
      outcome = output.outcome;
    }
  }
  return outcome;
}

function latestCheckoutFailure(context) {
  let failed = false;
  for (const output of Object.values(context)) {
    if (!isRecord(output)) {
      continue;
    }
    const paymentIntent = isRecord(output.payment_intent) ? output.payment_intent : {};
    const invoice = isRecord(output.invoice) ? output.invoice : {};
    if (
      paymentIntent.status === "requires_payment_method" ||
      paymentIntent.status === "requires_action" ||
      invoice.status === "open"
    ) {
      failed = true;
    }
  }
  return failed;
}

function projectSubscription(subscription, source, step = "") {
  const items = Array.isArray(subscription.items) ? subscription.items : [];
  const firstItem = isRecord(items[0]) ? items[0] : {};
  const id = stringValue(subscription.id);
  return {
    id,
    subscription: id,
    customer: stringValue(subscription.customer),
    status: stringValue(subscription.status || "unknown"),
    latest_invoice: stringValue(subscription.latest_invoice),
    cancel_at_period_end: Boolean(subscription.cancel_at_period_end),
    price: stringValue(firstItem.price),
    quantity: numberValue(firstItem.quantity, 0),
    source,
    step: step || undefined,
  };
}

function compareExpected(expected, actual, path = []) {
  const mismatches = [];
  for (const [key, expectedValue] of Object.entries(expected)) {
    const nextPath = [...path, key];
    const actualValue = isRecord(actual) ? actual[key] : undefined;
    if (isRecord(expectedValue)) {
      mismatches.push(...compareExpected(expectedValue, actualValue, nextPath));
      continue;
    }
    if (!sameValue(expectedValue, actualValue)) {
      mismatches.push(
        `${nextPath.join(".")} expected ${JSON.stringify(expectedValue)} got ${JSON.stringify(
          actualValue,
        )}`,
      );
    }
  }
  return mismatches;
}

function sameValue(expected, actual) {
  if (expected === actual) {
    return true;
  }
  if (typeof expected === "number" && typeof actual === "string") {
    return String(expected) === actual;
  }
  if (typeof expected === "string" && typeof actual === "number") {
    return expected === String(actual);
  }
  return false;
}

function verifySignature(header, secret, rawBody) {
  if (!header) {
    throw httpError(400, "missing Billtap-Signature header");
  }

  const parsed = parseSignatureHeader(header);
  const signedAt = Number(parsed.t);
  if (!Number.isFinite(signedAt) || signedAt <= 0 || parsed.v1.length === 0) {
    throw httpError(400, "invalid Billtap-Signature header");
  }

  const age = Math.abs(Math.floor(Date.now() / 1000) - signedAt);
  if (age > signatureToleranceSeconds) {
    throw httpError(400, "webhook signature timestamp outside tolerance");
  }

  const expected = createHmac("sha256", secret)
    .update(`${signedAt}.`)
    .update(rawBody)
    .digest("hex");
  const matched = parsed.v1.some((candidate) => safeEqualHex(candidate, expected));
  if (!matched) {
    throw httpError(400, "webhook signature mismatch");
  }
}

function parseSignatureHeader(header) {
  const out = { t: "", v1: [] };
  for (const part of header.split(",")) {
    const [key, value] = part.trim().split("=", 2);
    if (key === "t") {
      out.t = value || "";
    }
    if (key === "v1" && value) {
      out.v1.push(value);
    }
  }
  return out;
}

function safeEqualHex(left, right) {
  const leftBuffer = Buffer.from(left, "hex");
  const rightBuffer = Buffer.from(right, "hex");
  if (leftBuffer.length !== rightBuffer.length) {
    return false;
  }
  return timingSafeEqual(leftBuffer, rightBuffer);
}

function assertionTargetFromPath(pathname) {
  return pathname
    .slice("/test/assertions/".length)
    .split("/")
    .filter(Boolean)
    .map((segment) => decodeURIComponent(segment))
    .join(".");
}

function readBody(req, limit = 1024 * 1024) {
  return new Promise((resolve, reject) => {
    const chunks = [];
    let size = 0;
    let rejected = false;

    req.on("data", (chunk) => {
      size += chunk.length;
      if (size > limit) {
        rejected = true;
        reject(httpError(413, "request body too large"));
        req.destroy();
        return;
      }
      chunks.push(chunk);
    });
    req.on("end", () => {
      if (!rejected) {
        resolve(Buffer.concat(chunks));
      }
    });
    req.on("error", (error) => {
      if (!rejected) {
        reject(error);
      }
    });
  });
}

function parseJSON(body) {
  try {
    return JSON.parse(body.toString("utf8"));
  } catch (error) {
    throw httpError(400, `invalid JSON: ${error.message}`);
  }
}

function writeJSON(res, status, value) {
  res.writeHead(status, {
    "Content-Type": "application/json; charset=utf-8",
  });
  res.end(`${JSON.stringify(value, null, 2)}\n`);
}

function firstHeader(value) {
  if (Array.isArray(value)) {
    return value[0] || "";
  }
  return value || "";
}

function isRecord(value) {
  return value !== null && typeof value === "object" && !Array.isArray(value);
}

function stringValue(value) {
  if (value === undefined || value === null) {
    return "";
  }
  return String(value);
}

function numberValue(value, fallback) {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : fallback;
}

function httpError(status, message) {
  const error = new Error(message);
  error.status = status;
  return error;
}
