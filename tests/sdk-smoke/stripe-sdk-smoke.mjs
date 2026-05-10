import Stripe from "stripe";
import { spawn } from "node:child_process";
import { createServer as createNetServer } from "node:net";
import { createServer as createHTTPServer } from "node:http";
import { mkdir, rm, writeFile } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const testDir = dirname(fileURLToPath(import.meta.url));
const repoRoot = resolve(testDir, "../..");
const smokeDir = resolve(repoRoot, ".billtap/sdk-smoke");
const binaryPath = resolve(smokeDir, "billtap");
const timeoutMs = Number(
  process.env.BILLTAP_STRIPE_SDK_SMOKE_TIMEOUT_MS ?? 30_000,
);
const providedBaseURL = firstEnv(
  "BILLTAP_STRIPE_SDK_SMOKE_BASE_URL",
  "BILLTAP_BASE_URL",
);
const reportJSONPath = resolve(
  repoRoot,
  process.env.BILLTAP_STRIPE_SDK_SMOKE_REPORT_JSON ??
    ".billtap/sdk-smoke/stripe-sdk-smoke-report.json",
);
const reportMDPath = resolve(
  repoRoot,
  process.env.BILLTAP_STRIPE_SDK_SMOKE_REPORT_MD ??
    ".billtap/sdk-smoke/stripe-sdk-smoke-report.md",
);

const report = {
  name: "stripe-sdk-smoke",
  status: "running",
  startedAt: new Date().toISOString(),
  finishedAt: null,
  sdk: {
    name: "stripe-node",
    version: Stripe.PACKAGE_VERSION,
  },
  billtap: {
    baseURL: null,
    startedLocal: !providedBaseURL,
  },
  runId: null,
  checks: [],
  objects: {},
  webhookReceiver: null,
  error: null,
};

async function main() {
  await mkdir(smokeDir, { recursive: true });
  let localServer = null;
  let receiver = null;
  try {
    localServer = providedBaseURL ? null : await startBilltap();
    const baseURL = providedBaseURL
      ? normalizeBaseURL(providedBaseURL)
      : localServer.baseURL;
    const runId = smokeRunId(Boolean(providedBaseURL));
    report.billtap.baseURL = baseURL;
    report.runId = runId;

    receiver = await startWebhookReceiver();
    report.webhookReceiver = {
      url: receiver.url,
      providedURL: Boolean(process.env.BILLTAP_STRIPE_SDK_SMOKE_WEBHOOK_URL),
      received: [],
    };

    const stripe = stripeClient(baseURL);
    console.log(
      providedBaseURL
        ? `stripe SDK smoke using existing Billtap at ${baseURL}`
        : `stripe SDK smoke started Billtap at ${baseURL}`,
    );

    await check("health endpoint ready", async () => {
      await waitForReady(`${baseURL}/healthz`, localServer?.process ?? null);
    });

    await runStripeSDKSmoke(stripe, receiver, runId, Boolean(localServer));
    report.status = "passed";
    console.log(
      `stripe SDK smoke passed: ${report.checks.map((item) => item.name).join(", ")}`,
    );
  } catch (error) {
    report.status = "failed";
    report.error = messageFor(error);
    throw error;
  } finally {
    report.finishedAt = new Date().toISOString();
    if (report.webhookReceiver && receiver) {
      report.webhookReceiver.received = receiver.received.map((event) => ({
        type: event.type,
        id: event.id,
      }));
    }
    await receiver?.close();
    await localServer?.stop();
    await localServer?.removeDatabase();
    await writeReports();
  }
}

async function runStripeSDKSmoke(stripe, receiver, runId, ownsBilltapServer) {
  const customer = await check("customer create/retrieve/list", async () => {
    const created = await stripe.customers.create({
      email: `stripe-sdk-smoke+${runId}@example.test`,
      name: "Stripe SDK Smoke Customer",
      metadata: {
        billtap_smoke: "stripe_sdk",
        run_id: runId,
      },
    });
    assertEqual(created.object, "customer", "customer object");

    const retrieved = await stripe.customers.retrieve(created.id);
    assertEqual(retrieved.id, created.id, "retrieved customer id");

    const listed = await stripe.customers.list({
      email: created.email,
      limit: 10,
    });
    assertListContains(listed, created.id, "customer list");
    return created;
  });
  report.objects.customer = pick(customer, ["id", "object", "email"]);

  const product = await check("product create/retrieve/list", async () => {
    const created = await stripe.products.create({
      name: `Stripe SDK Smoke Product ${runId}`,
      description: "Billtap official Stripe SDK smoke product",
      metadata: {
        billtap_smoke: "stripe_sdk",
        run_id: runId,
      },
    });
    assertEqual(created.object, "product", "product object");

    const retrieved = await stripe.products.retrieve(created.id);
    assertEqual(retrieved.id, created.id, "retrieved product id");

    const listed = await stripe.products.list({ limit: 100 });
    assertListContains(listed, created.id, "product list");
    return created;
  });
  report.objects.product = pick(product, ["id", "object", "name"]);

  const price = await check("price create/retrieve/list", async () => {
    const created = await stripe.prices.create({
      product: product.id,
      currency: "usd",
      unit_amount: 4200,
      lookup_key: `stripe_sdk_smoke_${runId}`,
      recurring: {
        interval: "month",
      },
      metadata: {
        billtap_smoke: "stripe_sdk",
        run_id: runId,
      },
    });
    assertEqual(created.object, "price", "price object");
    assertEqual(created.product, product.id, "price product");

    const retrieved = await stripe.prices.retrieve(created.id);
    assertEqual(retrieved.id, created.id, "retrieved price id");

    const listed = await stripe.prices.list({
      product: product.id,
      active: true,
      type: "recurring",
      limit: 10,
    });
    assertListContains(listed, created.id, "price list");
    return created;
  });
  report.objects.price = pick(price, ["id", "object", "product", "lookup_key"]);

  const webhookEndpoint = await check(
    "webhook endpoint create/retrieve/list",
    async () => {
      const created = await stripe.webhookEndpoints.create({
        url: receiver.url,
        enabled_events: ["checkout.session.completed", "invoice.*"],
      });
      assertEqual(
        created.object,
        "webhook_endpoint",
        "webhook endpoint object",
      );

      const retrieved = await stripe.webhookEndpoints.retrieve(created.id);
      assertEqual(retrieved.id, created.id, "retrieved webhook endpoint id");

      const listed = await stripe.webhookEndpoints.list({ limit: 100 });
      assertListContains(listed, created.id, "webhook endpoint list");

      const updated = await stripe.webhookEndpoints.update(created.id, {
        active: true,
        enabled_events: [
          "checkout.session.completed",
          "invoice.*",
          "payment_intent.*",
        ],
      });
      assertEqual(updated.active, true, "updated webhook endpoint active");
      return updated;
    },
  );
  report.objects.webhookEndpoint = pick(webhookEndpoint, [
    "id",
    "object",
    "url",
    "enabled_events",
    "active",
  ]);

  const checkoutSession = await check(
    "checkout session create/retrieve/list",
    async () => {
      const optionalParamSession = await stripe.checkout.sessions.create({
        customer: customer.id,
        mode: "subscription",
        line_items: [
          {
            price: price.id,
            quantity: 1,
          },
        ],
        allow_promotion_codes: true,
        subscription_data: {
          trial_period_days: 14,
        },
        success_url: "http://127.0.0.1/trial-success",
      });
      assertEqual(
        optionalParamSession.object,
        "checkout.session",
        "checkout session optional params object",
      );

      const created = await stripe.checkout.sessions.create({
        customer: customer.id,
        mode: "subscription",
        line_items: [
          {
            price: price.id,
            quantity: 2,
          },
        ],
        success_url: "http://127.0.0.1/success",
        cancel_url: "http://127.0.0.1/cancel",
      });
      assertEqual(
        created.object,
        "checkout.session",
        "checkout session object",
      );
      assertEqual(created.customer, customer.id, "checkout session customer");

      const retrieved = await stripe.checkout.sessions.retrieve(created.id);
      assertEqual(retrieved.id, created.id, "retrieved checkout session id");

      const listed = await stripe.checkout.sessions.list({ limit: 100 });
      assertListContains(listed, created.id, "checkout session list");
      return created;
    },
  );
  report.objects.checkoutSession = pick(checkoutSession, [
    "id",
    "object",
    "customer",
    "status",
    "payment_status",
  ]);

  const completion = await check(
    "checkout completion via raw SDK request",
    async () => {
      const completed = await stripe.rawRequest(
        "POST",
        `/v1/checkout/sessions/${encodeURIComponent(checkoutSession.id)}/complete`,
        { outcome: "payment_succeeded" },
      );
      assertEqual(
        completed.session.payment_status,
        "paid",
        "completed checkout payment status",
      );
      assert(
        completed.session.subscription,
        "completed checkout should include subscription",
      );
      assert(
        completed.session.invoice,
        "completed checkout should include invoice",
      );
      assert(
        completed.session.payment_intent,
        "completed checkout should include payment intent",
      );
      return completed;
    },
  );
  report.objects.subscription = { id: completion.session.subscription };
  report.objects.invoice = { id: completion.session.invoice };
  report.objects.paymentIntent = { id: completion.session.payment_intent };

  await check("completed checkout retrieval flows", async () => {
    const completedSession = await stripe.checkout.sessions.retrieve(
      checkoutSession.id,
    );
    assertEqual(
      completedSession.payment_status,
      "paid",
      "retrieved completed checkout payment status",
    );

    const subscriptions = await stripe.subscriptions.list({
      customer: customer.id,
      status: "all",
      limit: 10,
    });
    assertListContains(
      subscriptions,
      completedSession.subscription,
      "subscription list",
    );

    const subscription = await stripe.subscriptions.retrieve(
      completedSession.subscription,
    );
    assertEqual(
      subscription.customer,
      customer.id,
      "retrieved subscription customer",
    );

    const invoices = await stripe.invoices.list({
      customer: customer.id,
      subscription: subscription.id,
      limit: 10,
    });
    assertListContains(invoices, completedSession.invoice, "invoice list");

    const invoice = await stripe.invoices.retrieve(completedSession.invoice);
    assertEqual(
      invoice.subscription,
      subscription.id,
      "retrieved invoice subscription",
    );

    const paymentIntents = await stripe.paymentIntents.list({ limit: 100 });
    assertListContains(
      paymentIntents,
      completedSession.payment_intent,
      "payment intent list",
    );

    const paymentIntent = await stripe.paymentIntents.retrieve(
      completedSession.payment_intent,
    );
    assertEqual(
      paymentIntent.status,
      "succeeded",
      "retrieved payment intent status",
    );

    const paymentMethods = await stripe.paymentMethods.list({
      customer: customer.id,
      type: "card",
      limit: 10,
    });
    assert(
      paymentMethods.data.length > 0,
      "payment method list should return a sandbox card",
    );
  });

  const checkoutEvent = await check("event list/retrieve", async () => {
    const events = await stripe.events.list({
      type: "checkout.session.completed",
      limit: 100,
    });
    const event = events.data.find(
      (candidate) => candidate.data?.object?.id === checkoutSession.id,
    );
    assert(
      event,
      "events list should include completed checkout session event",
    );

    const retrieved = await stripe.events.retrieve(event.id);
    assertEqual(retrieved.id, event.id, "retrieved event id");
    assertEqual(
      retrieved.type,
      "checkout.session.completed",
      "retrieved event type",
    );
    return retrieved;
  });
  report.objects.event = pick(checkoutEvent, ["id", "object", "type"]);

  await check("webhook endpoint delete", async () => {
    const deleted = await stripe.webhookEndpoints.del(webhookEndpoint.id);
    assertEqual(deleted.active, false, "deleted webhook endpoint active flag");
  });

  const expectDelivery =
    ownsBilltapServer ||
    process.env.BILLTAP_STRIPE_SDK_SMOKE_EXPECT_DELIVERY === "1";
  if (expectDelivery) {
    await check("webhook receiver observed delivery", async () => {
      await waitFor(
        () => receiver.received.length > 0,
        timeoutMs,
        "webhook delivery",
      );
      assert(
        receiver.received.some(
          (event) => event.type === "checkout.session.completed",
        ),
        "webhook receiver should observe checkout.session.completed",
      );
    });
  }
}

async function startBilltap() {
  await run("go", ["build", "-o", binaryPath, "./cmd/billtap"]);
  const port =
    Number(process.env.BILLTAP_STRIPE_SDK_SMOKE_PORT) || (await findOpenPort());
  const baseURL = `http://127.0.0.1:${port}`;
  const databasePath = resolve(smokeDir, `stripe-sdk-smoke-${process.pid}.db`);
  await removeDatabase(databasePath);

  const child = spawn(binaryPath, [], {
    cwd: repoRoot,
    env: {
      ...process.env,
      BILLTAP_ADDR: `127.0.0.1:${port}`,
      BILLTAP_DATABASE_URL: databasePath,
      BILLTAP_STATIC_DIR: "dist/app",
      BILLTAP_ENV: "test",
    },
    stdio: ["ignore", "pipe", "pipe"],
  });

  child.stdout.on("data", (chunk) =>
    process.stdout.write(`[billtap] ${chunk}`),
  );
  child.stderr.on("data", (chunk) =>
    process.stderr.write(`[billtap] ${chunk}`),
  );

  try {
    await waitForReady(`${baseURL}/healthz`, child);
  } catch (error) {
    if (child.exitCode === null && child.signalCode === null) {
      child.kill("SIGTERM");
      const stopped = await waitForExit(child, 5_000);
      if (!stopped) child.kill("SIGKILL");
    }
    await removeDatabase(databasePath);
    throw error;
  }
  return {
    baseURL,
    process: child,
    stop: async () => {
      if (child.exitCode !== null || child.signalCode !== null) return;
      child.kill("SIGTERM");
      const stopped = await waitForExit(child, 5_000);
      if (!stopped) child.kill("SIGKILL");
    },
    removeDatabase: async () => {
      await removeDatabase(databasePath);
    },
  };
}

async function startWebhookReceiver() {
  if (process.env.BILLTAP_STRIPE_SDK_SMOKE_WEBHOOK_URL) {
    return {
      url: process.env.BILLTAP_STRIPE_SDK_SMOKE_WEBHOOK_URL,
      received: [],
      close: async () => {},
    };
  }

  const received = [];
  const server = createHTTPServer(async (req, res) => {
    const chunks = [];
    req.on("data", (chunk) => chunks.push(chunk));
    req.on("end", () => {
      const body = Buffer.concat(chunks).toString("utf8");
      let parsed = {};
      try {
        parsed = JSON.parse(body);
      } catch {
        parsed = {};
      }
      received.push({
        id: parsed.id,
        type: parsed.type,
        body,
        signature:
          req.headers["billtap-signature"] || req.headers["stripe-signature"],
      });
      res.writeHead(200, { "Content-Type": "application/json" });
      res.end(JSON.stringify({ received: true }));
    });
  });

  const port = await new Promise((resolvePromise, reject) => {
    server.on("error", reject);
    server.listen(0, "127.0.0.1", () => {
      const address = server.address();
      if (address && typeof address === "object") {
        resolvePromise(address.port);
        return;
      }
      reject(new Error("could not allocate a local webhook receiver port"));
    });
  });

  return {
    url: `http://127.0.0.1:${port}/webhooks/stripe`,
    received,
    close: async () => {
      await new Promise((resolvePromise) => server.close(resolvePromise));
    },
  };
}

function stripeClient(baseURL) {
  const url = new URL(baseURL);
  return new Stripe(
    process.env.BILLTAP_STRIPE_SDK_SMOKE_API_KEY ?? "sk_test_billtap_sdk_smoke",
    {
      host: url.hostname,
      port: url.port || (url.protocol === "https:" ? "443" : "80"),
      protocol: url.protocol.replace(":", ""),
      maxNetworkRetries: 0,
      telemetry: false,
      timeout: timeoutMs,
    },
  );
}

async function check(name, fn) {
  const startedAt = new Date().toISOString();
  try {
    const value = await fn();
    report.checks.push({
      name,
      status: "passed",
      startedAt,
      finishedAt: new Date().toISOString(),
    });
    return value;
  } catch (error) {
    report.checks.push({
      name,
      status: "failed",
      startedAt,
      finishedAt: new Date().toISOString(),
      error: messageFor(error),
    });
    throw error;
  }
}

async function run(command, args) {
  await new Promise((resolvePromise, reject) => {
    const child = spawn(command, args, {
      cwd: repoRoot,
      env: process.env,
      stdio: "inherit",
    });
    child.on("error", reject);
    child.on("exit", (code, signal) => {
      if (code === 0) {
        resolvePromise();
        return;
      }
      reject(
        new Error(`${command} ${args.join(" ")} exited with ${signal ?? code}`),
      );
    });
  });
}

async function waitForReady(url, child) {
  await waitFor(
    async () => {
      if (child && (child.exitCode !== null || child.signalCode !== null)) {
        throw new Error(
          `billtap server exited before readiness with ${child.signalCode ?? child.exitCode}`,
        );
      }
      try {
        const response = await fetch(url);
        return response.ok;
      } catch {
        return false;
      }
    },
    timeoutMs,
    url,
  );
}

async function waitFor(predicate, timeout, label) {
  const started = Date.now();
  while (Date.now() - started < timeout) {
    const result = await predicate();
    if (result) return;
    await delay(250);
  }
  throw new Error(`timed out waiting for ${label}`);
}

async function waitForExit(child, timeout) {
  return await new Promise((resolvePromise) => {
    const timer = setTimeout(() => resolvePromise(false), timeout);
    child.once("exit", () => {
      clearTimeout(timer);
      resolvePromise(true);
    });
  });
}

async function findOpenPort() {
  return await new Promise((resolvePromise, reject) => {
    const server = createNetServer();
    server.on("error", reject);
    server.listen(0, "127.0.0.1", () => {
      const address = server.address();
      server.close(() => {
        if (address && typeof address === "object") {
          resolvePromise(address.port);
          return;
        }
        reject(new Error("could not allocate a local port"));
      });
    });
  });
}

async function removeDatabase(databasePath) {
  await Promise.all([
    rm(databasePath, { force: true }),
    rm(`${databasePath}-journal`, { force: true }),
    rm(`${databasePath}-shm`, { force: true }),
    rm(`${databasePath}-wal`, { force: true }),
  ]);
}

async function writeReports() {
  await mkdir(dirname(reportJSONPath), { recursive: true });
  await mkdir(dirname(reportMDPath), { recursive: true });
  await writeFile(reportJSONPath, `${JSON.stringify(report, null, 2)}\n`);
  await writeFile(reportMDPath, markdownReport());
  console.log(`stripe SDK smoke report: ${reportJSONPath}`);
  console.log(`stripe SDK smoke report: ${reportMDPath}`);
}

function markdownReport() {
  const lines = [
    "# Stripe SDK Smoke Report",
    "",
    `- Status: ${report.status}`,
    `- SDK: ${report.sdk.name} ${report.sdk.version}`,
    `- Billtap base URL: ${report.billtap.baseURL}`,
    `- Started local Billtap: ${report.billtap.startedLocal}`,
    `- Run ID: ${report.runId}`,
    "",
    "## Checks",
    "",
  ];
  for (const checkResult of report.checks) {
    lines.push(
      `- ${checkResult.status === "passed" ? "PASS" : "FAIL"} ${checkResult.name}`,
    );
    if (checkResult.error) {
      lines.push(`  - ${checkResult.error}`);
    }
  }
  lines.push(
    "",
    "## Objects",
    "",
    "```json",
    JSON.stringify(report.objects, null, 2),
    "```",
    "",
  );
  if (report.error) {
    lines.push("## Error", "", "```text", report.error, "```", "");
  }
  return `${lines.join("\n")}\n`;
}

function assert(value, message) {
  if (!value) {
    throw new Error(message);
  }
}

function assertEqual(actual, expected, label) {
  if (actual !== expected) {
    throw new Error(
      `${label}: got ${JSON.stringify(actual)}, want ${JSON.stringify(expected)}`,
    );
  }
}

function assertListContains(list, id, label) {
  assert(Array.isArray(list.data), `${label} should be a Stripe list`);
  assert(
    list.data.some((item) => item.id === id),
    `${label} should include ${id}`,
  );
}

function pick(value, keys) {
  const out = {};
  for (const key of keys) {
    if (value?.[key] !== undefined) {
      out[key] = value[key];
    }
  }
  return out;
}

function normalizeBaseURL(value) {
  const url = new URL(value);
  url.pathname = "/";
  url.search = "";
  url.hash = "";
  return url.origin;
}

function firstEnv(...names) {
  for (const name of names) {
    const value = process.env[name]?.trim();
    if (value) return value;
  }
  return "";
}

function smokeRunId(usesProvidedBaseURL) {
  const explicit = process.env.BILLTAP_STRIPE_SDK_SMOKE_RUN_ID?.trim();
  if (explicit) return sanitizeRunId(explicit);
  if (!usesProvidedBaseURL) return "local";
  return sanitizeRunId(`manual_${Date.now()}_${process.pid}`);
}

function sanitizeRunId(value) {
  return value.replace(/[^a-zA-Z0-9_-]/g, "_").slice(0, 48) || "run";
}

function messageFor(error) {
  if (error instanceof Error) return error.stack || error.message;
  return String(error);
}

function delay(ms) {
  return new Promise((resolvePromise) => setTimeout(resolvePromise, ms));
}

try {
  await main();
} catch (error) {
  console.error(messageFor(error));
  process.exitCode = 1;
}
