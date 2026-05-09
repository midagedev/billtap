import { chromium } from "@playwright/test";
import { createServer } from "node:net";
import { mkdir, rm } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { spawn } from "node:child_process";

const testDir = dirname(fileURLToPath(import.meta.url));
const repoRoot = resolve(testDir, "../..");
const smokeDir = resolve(repoRoot, ".billtap/web-smoke");
const binaryPath = resolve(smokeDir, "billtap");
const timeoutMs = Number(process.env.BILLTAP_WEB_SMOKE_TIMEOUT_MS ?? 30_000);

async function main() {
  await mkdir(smokeDir, { recursive: true });
  await run("npm", ["run", "build"]);
  await run("go", ["build", "-o", binaryPath, "./cmd/billtap"]);

  const port = Number(process.env.BILLTAP_WEB_SMOKE_PORT) || await findOpenPort();
  const baseURL = `http://127.0.0.1:${port}`;
  const databasePath = resolve(smokeDir, `web-smoke-${process.pid}.db`);
  await removeDatabase(databasePath);

  const server = spawn(binaryPath, [], {
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

  server.stdout.on("data", (chunk) => process.stdout.write(`[billtap] ${chunk}`));
  server.stderr.on("data", (chunk) => process.stderr.write(`[billtap] ${chunk}`));

  const stopServer = async () => {
    if (server.exitCode !== null || server.signalCode !== null) return;
    server.kill("SIGTERM");
    const stopped = await waitForExit(server, 5_000);
    if (!stopped) server.kill("SIGKILL");
  };

  try {
    await waitForReady(`${baseURL}/healthz`, server);
    const seed = await seedBillingData(baseURL);
    const checks = smokeChecks(seed);
    await runBrowserSmoke(baseURL, checks);
    console.log(`web smoke passed: ${checks.map((check) => check.path).join(", ")}`);
  } finally {
    await stopServer();
    await removeDatabase(databasePath);
  }
}

function smokeChecks(seed) {
  return [
    {
      name: "dashboard",
      path: "/app/dashboard/",
      texts: ["Debug workspace", "Billing objects", "Object timeline"],
    },
    {
      name: "checkout",
      path: `/app/checkout/?session_id=${encodeURIComponent(seed.checkoutSession.id)}`,
      texts: ["Checkout session", "Outcome selector", "Complete checkout"],
    },
    {
      name: "portal",
      path: `/app/portal/?customer_id=${encodeURIComponent(seed.customer.id)}`,
      texts: ["Subscription management", "Current subscription", "Cancellation"],
    },
  ];
}

async function seedBillingData(baseURL) {
  const customer = await postForm(`${baseURL}/v1/customers`, {
    email: "web-smoke@example.test",
    name: "Web Smoke User",
  });
  const product = await postForm(`${baseURL}/v1/products`, {
    name: "Web Smoke Team",
  });
  const price = await postForm(`${baseURL}/v1/prices`, {
    product: product.id,
    currency: "usd",
    unit_amount: "9900",
    "recurring[interval]": "month",
  });
  const checkoutSession = await postForm(`${baseURL}/v1/checkout/sessions`, {
    customer: customer.id,
    mode: "subscription",
    "line_items[0][price]": price.id,
    "line_items[0][quantity]": "1",
    success_url: "http://127.0.0.1/success",
    cancel_url: "http://127.0.0.1/cancel",
  });

  await postJSON(`${baseURL}/api/checkout/sessions/${encodeURIComponent(checkoutSession.id)}/complete`, {
    outcome: "payment_succeeded",
  });

  return { customer, product, price, checkoutSession };
}

async function postForm(url, values) {
  const body = new URLSearchParams();
  for (const [key, value] of Object.entries(values)) {
    body.set(key, value);
  }

  const response = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body,
  });
  return await readJSONResponse(response, url);
}

async function postJSON(url, value) {
  const response = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(value),
  });
  return await readJSONResponse(response, url);
}

async function readJSONResponse(response, url) {
  const body = await response.text();
  if (!response.ok) {
    throw new Error(`${url} returned ${response.status}: ${body}`);
  }
  return JSON.parse(body);
}

async function runBrowserSmoke(baseURL, checks) {
  const browser = await chromium.launch({
    headless: process.env.PLAYWRIGHT_HEADLESS !== "0",
  });
  const page = await browser.newPage();
  const consoleErrors = [];
  let currentRoute = "";

  page.on("console", (message) => {
    if (message.type() === "error") {
      consoleErrors.push(`${currentRoute || "unknown route"}: ${message.text()}`);
    }
  });
  page.on("pageerror", (error) => {
    consoleErrors.push(`${currentRoute || "unknown route"}: ${error.message}`);
  });

  try {
    for (const check of checks) {
      currentRoute = check.path;
      const response = await page.goto(`${baseURL}${check.path}`, {
        waitUntil: "domcontentloaded",
        timeout: timeoutMs,
      });
      await page.waitForLoadState("load", { timeout: timeoutMs });

      if (!response || !response.ok()) {
        throw new Error(`${check.name} returned ${response?.status() ?? "no response"}`);
      }

      for (const text of check.texts) {
        await page.getByText(text, { exact: true }).first().waitFor({
          state: "visible",
          timeout: timeoutMs,
        });
      }
    }

    if (consoleErrors.length > 0) {
      throw new Error(`console errors during web smoke:\n${consoleErrors.join("\n")}`);
    }
  } finally {
    await browser.close();
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
      reject(new Error(`${command} ${args.join(" ")} exited with ${signal ?? code}`));
    });
  });
}

async function findOpenPort() {
  return await new Promise((resolvePromise, reject) => {
    const server = createServer();
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

async function waitForReady(url, server) {
  const started = Date.now();
  while (Date.now() - started < timeoutMs) {
    if (server.exitCode !== null || server.signalCode !== null) {
      throw new Error(`billtap server exited before readiness with ${server.signalCode ?? server.exitCode}`);
    }

    try {
      const response = await fetch(url);
      if (response.ok) return;
    } catch {
      // Server is still starting.
    }

    await delay(250);
  }
  throw new Error(`timed out waiting for ${url}`);
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

async function removeDatabase(databasePath) {
  await Promise.all([
    rm(databasePath, { force: true }),
    rm(`${databasePath}-journal`, { force: true }),
    rm(`${databasePath}-shm`, { force: true }),
    rm(`${databasePath}-wal`, { force: true }),
  ]);
}

function delay(ms) {
  return new Promise((resolvePromise) => setTimeout(resolvePromise, ms));
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
