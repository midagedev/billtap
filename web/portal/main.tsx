import { useEffect, useState, type FormEvent } from "react";
import { createRoot } from "react-dom/client";
import {
  AppShell,
  CopyValue,
  Panel,
  StatusPill,
  SurfaceHeader,
} from "../shared/components";
import {
  cancelPortalSubscription,
  changePortalPlan,
  changePortalSeats,
  getPortalCustomerId,
  loadPortalData,
  resumePortalSubscription,
  updatePortalPaymentMethod,
  type DataSource,
  type PortalCancelMode,
  type PortalData,
  type PortalPaymentOutcome,
} from "../shared/api";
import { billingSnapshot, invoices, portalPlans } from "../shared/data";
import { appHref } from "../shared/basePath";
import "../shared/styles.css";

const securePaymentFrame = `<!doctype html>
<html>
  <head>
    <style>
      body { margin: 0; font: 14px system-ui, sans-serif; color: #172033; }
      .grid { display: grid; gap: 8px; }
      .row { display: grid; grid-template-columns: 1fr 0.6fr; gap: 8px; }
      input { width: 100%; min-height: 38px; box-sizing: border-box; border: 1px solid #cbd5e1; border-radius: 8px; padding: 0 10px; font: inherit; }
    </style>
  </head>
  <body>
    <div class="grid">
      <input name="number" autocomplete="cc-number" placeholder="4242 4242 4242 4242" />
      <div class="row">
        <input name="expiry" autocomplete="cc-exp" placeholder="12/34" />
        <input name="cvc" autocomplete="cc-csc" placeholder="123" />
      </div>
    </div>
    <script>
      const post = () => parent.postMessage({
        type: 'billtap-card-input',
        number: document.querySelector('input[name="number"]').value
      }, '*');
      document.addEventListener('input', post);
    </script>
  </body>
</html>`;

const fixturePortalData: PortalData = {
  subscription: {
    id: "sub_test_29",
    customerId: "cus_test_workspace",
    customerEmail: billingSnapshot.owner,
    workspace: billingSnapshot.workspace,
    owner: billingSnapshot.owner,
    planId: "team",
    planName: billingSnapshot.plan,
    price: "$99/mo",
    status: billingSnapshot.status,
    period: billingSnapshot.period,
    seats: billingSnapshot.seats.basic + billingSnapshot.seats.additional,
    usedSeats: billingSnapshot.seats.used,
    pendingSeats: billingSnapshot.seats.pending,
    cancelAtPeriodEnd: false,
    exportRemaining: billingSnapshot.export.remaining,
    exportManualLimit: billingSnapshot.export.manualLimit,
    exportRenewal: billingSnapshot.export.renews,
  },
  plans: portalPlans,
  invoices,
  paymentMethod: {
    id: "pm_fixture_card",
    label: "Visa sandbox ending 4242",
    status: "saved",
  },
};

function PortalApp() {
  const [portal, setPortal] = useState<PortalData>(fixturePortalData);
  const [source, setSource] = useState<DataSource>("fixture");
  const [loadError, setLoadError] = useState<string>();
  const [actionStatus, setActionStatus] = useState("Ready");
  const [busyAction, setBusyAction] = useState<string>();
  const [draftSeats, setDraftSeats] = useState(fixturePortalData.subscription.seats);
  const [isEditingPaymentMethod, setIsEditingPaymentMethod] = useState(false);
  const [portalCardNumber, setPortalCardNumber] = useState("4242424242424242");

  useEffect(() => {
    let active = true;
    loadPortalData(getPortalCustomerId()).then((result) => {
      if (!active) return;
      setPortal(result.data);
      setDraftSeats(result.data.subscription.seats);
      setSource(result.source);
      setLoadError(result.error);
    });
    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    const listener = (event: MessageEvent) => {
      const data = event.data as { type?: string; number?: string };
      if (data?.type === "billtap-card-input" && typeof data.number === "string") {
        setPortalCardNumber(data.number);
      }
    };
    window.addEventListener("message", listener);
    return () => window.removeEventListener("message", listener);
  }, []);

  async function runAction(label: string, operation: () => Promise<{ data: PortalData; source: DataSource; message: string; error?: string }>) {
    setBusyAction(label);
    setActionStatus(`${label} pending`);
    try {
      const result = await operation();
      setPortal(result.data);
      setDraftSeats(result.data.subscription.seats);
      setSource(result.source);
      setActionStatus(result.error ? `${result.message}: ${result.error}` : result.message);
      if (!result.error && shouldRedirectAfterPortalAction(label)) {
        window.setTimeout(() => window.location.assign(getPortalReturnUrl()), 0);
      }
    } catch (error) {
      setActionStatus(error instanceof Error ? error.message : String(error));
    } finally {
      setBusyAction(undefined);
    }
  }

  const subscription = portal.subscription;
  const cancelMode = subscription.status === "canceled" ? "immediate" : subscription.cancelAtPeriodEnd ? "period" : "none";
  const returnUrl = getPortalReturnUrl();

  async function handlePortalPaymentSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const outcome = portalPaymentOutcomeForCard(portalCardNumber);
    await runAction("payment method", () => updatePortalPaymentMethod(portal, outcome));
    setIsEditingPaymentMethod(false);
  }

  return (
    <AppShell active="portal">
      <div data-testid="page-container-main">
        <SurfaceHeader
          eyebrow="Hosted billing portal"
          title="Subscription management"
          meta={`${subscription.workspace} · ${subscription.owner}`}
          actions={
            <div className="confirmation-strip">
              <a className="button secondary" data-testid="return-to-business-link" href={returnUrl}>
                Return to business
              </a>
              <a className="button secondary" href={appHref("dashboard/")}>
                View evidence
              </a>
            </div>
          }
        />

      <div className="status-strip">
        <StatusPill tone={source === "api" ? "good" : "warning"}>
          {source === "api" ? "Portal API loaded" : "Fixture portal"}
        </StatusPill>
        <span>{loadError ?? actionStatus}</span>
      </div>

      <div className="portal-layout">
        <div className="subscription-stack">
          <section className="portal-card">
            <div className="panel-header">
              <div>
                <span className="section-kicker">Current subscription</span>
                <h2>{subscription.planName}</h2>
              </div>
              <StatusPill tone={cancelMode === "none" ? "good" : "warning"}>
                {cancelMode === "none" ? subscription.status : `cancel ${cancelMode}`}
              </StatusPill>
            </div>
            <div className="snapshot-grid">
              <div>
                <span className="section-kicker">Period</span>
                <strong>{subscription.period}</strong>
                <span>Renews unless canceled</span>
              </div>
              <div>
                <span className="section-kicker">Exports</span>
                <strong>{subscription.exportRemaining} remaining</strong>
                <span>{subscription.exportManualLimit} manual allowance</span>
              </div>
            </div>
          </section>

          <Panel title="Plan change">
            <div className="plan-grid">
              {portal.plans.map((item) => (
                <button
                  className={`plan-card ${item.id === subscription.planId ? "active" : ""}`}
                  key={item.id}
                  disabled={busyAction !== undefined}
                  onClick={() => void runAction("plan change", () => changePortalPlan(portal, item.id))}
                  type="button"
                >
                  <strong>
                    {item.name} · {item.price}
                  </strong>
                  <span>{item.detail}</span>
                </button>
              ))}
            </div>
          </Panel>

          <Panel title="Seat quantity">
            <div className="seat-control">
              <div>
                <span className="section-kicker">Purchased seats</span>
                <strong>{subscription.seats} total</strong>
                <span>{subscription.usedSeats} used</span>
              </div>
              <div className="stepper" aria-label="Seat quantity control">
                <button onClick={() => setDraftSeats((value) => Math.max(1, value - 1))} type="button">
                  -
                </button>
                <strong>{draftSeats}</strong>
                <button onClick={() => setDraftSeats((value) => value + 1)} type="button">
                  +
                </button>
                <button
                  className="button secondary"
                  disabled={busyAction !== undefined}
                  onClick={() => void runAction("seat change", () => changePortalSeats(portal, draftSeats))}
                  type="button"
                >
                  Apply
                </button>
              </div>
            </div>
          </Panel>
        </div>

        <aside className="subscription-stack">
          <Panel title="Payment method update">
            <div className="stripe-portal-payment">
              <div className="stripe-portal-card">
                <div>
                  <span className="section-kicker">Current payment method</span>
                  <strong>{portal.paymentMethod.status === "failed" ? "Payment method failed" : "Visa ending 4242"}</strong>
                  <span>{portal.paymentMethod.label}</span>
                </div>
                <div className="stripe-portal-actions">
                  <button type="button" onClick={() => setIsEditingPaymentMethod(true)}>
                    Edit
                  </button>
                  <a role="button" aria-label="Change saved payment method" onClick={() => setIsEditingPaymentMethod(true)}>
                    <svg viewBox="0 0 20 20" aria-hidden="true">
                      <path d="M4 13.5V16h2.5l7.4-7.4-2.5-2.5L4 13.5Zm11.7-6.7a.9.9 0 0 0 0-1.3l-1.2-1.2a.9.9 0 0 0-1.3 0l-.9.9 2.5 2.5.9-.9Z" />
                    </svg>
                  </a>
                </div>
              </div>
              {isEditingPaymentMethod ? (
                <form className="stripe-portal-form" onSubmit={handlePortalPaymentSubmit}>
                  <label className="stripe-radio-row" htmlFor="radio-add">
                    <input id="radio-add" name="payment-method-choice" type="radio" defaultChecked />
                    <span>Add payment method</span>
                  </label>
                  <iframe
                    name="__privateStripeFrameBilltap"
                    title="Secure payment input frame"
                    srcDoc={securePaymentFrame}
                  />
                  <div className="confirmation-strip">
                    <button className="button" data-testid="confirm" disabled={busyAction !== undefined} type="submit">
                      Update
                    </button>
                    <button
                      className="button secondary"
                      data-test="cancel"
                      onClick={() => setIsEditingPaymentMethod(false)}
                      type="button"
                    >
                      Go back
                    </button>
                  </div>
                </form>
              ) : null}
            </div>
            <div className="option-grid">
              <button
                className={`choice-card ${portal.paymentMethod.status !== "failed" ? "active" : ""}`}
                disabled={busyAction !== undefined}
                onClick={() => void runAction("payment method", () => updatePortalPaymentMethod(portal, "succeeds"))}
                type="button"
              >
                <strong>Update succeeds</strong>
                <StatusPill tone="good">saved</StatusPill>
              </button>
              <button
                className={`choice-card ${portal.paymentMethod.status === "failed" ? "active" : ""}`}
                disabled={busyAction !== undefined}
                onClick={() => void runAction("payment method", () => updatePortalPaymentMethod(portal, "fails"))}
                type="button"
              >
                <strong>Update fails</strong>
                <StatusPill tone="danger">declined</StatusPill>
              </button>
            </div>
          </Panel>

          <Panel title="Cancellation">
            <div className="confirmation-strip">
              <button
                className="button warning"
                disabled={busyAction !== undefined}
                onClick={() => void runAction("cancel", () => cancelPortalSubscription(portal, "period"))}
                type="button"
              >
                At period end
              </button>
              <button
                className="button danger"
                disabled={busyAction !== undefined}
                onClick={() => void runAction("cancel", () => cancelPortalSubscription(portal, "immediate"))}
                type="button"
              >
                Immediate
              </button>
              <button
                className="button secondary"
                disabled={busyAction !== undefined}
                onClick={() => void runAction("resume", () => resumePortalSubscription(portal))}
                type="button"
              >
                Resume
              </button>
            </div>
          </Panel>

          <Panel title="Invoice history">
            <ul className="invoice-list">
              {portal.invoices.map((invoice) => (
                <li key={invoice.id}>
                  <div>
                    <strong>{invoice.period}</strong>
                    <span>{invoice.id}</span>
                  </div>
                  <div>
                    <StatusPill tone={invoice.status === "paid" ? "good" : "neutral"}>
                      {invoice.status}
                    </StatusPill>
                    <span>{invoice.amount}</span>
                  </div>
                </li>
              ))}
            </ul>
          </Panel>

          <Panel title="Pending portal result">
            <div className="surface-grid">
              <CopyValue label="Plan request" value={subscription.planId} />
              <CopyValue label="Seat request" value={`${subscription.seats} seats`} />
              <CopyValue label="Payment method" value={`${portal.paymentMethod.label} · ${portal.paymentMethod.status}`} />
              <CopyValue label="Action status" value={actionStatus} />
            </div>
          </Panel>
        </aside>
      </div>
      </div>
    </AppShell>
  );
}

function getPortalReturnUrl(location: Location = window.location): string {
  const params = new URLSearchParams(location.search);
  return params.get("return_url") ?? params.get("returnUrl") ?? appHref("dashboard/");
}

function shouldRedirectAfterPortalAction(action: string, location: Location = window.location): boolean {
  const params = new URLSearchParams(location.search);
  const redirectOnAction = params.get("redirect_on_action") ?? params.get("redirectOnAction");
  if (redirectOnAction !== "true" && redirectOnAction !== "1") {
    return false;
  }
  const normalizedAction = action.toLowerCase();
  const flow = (params.get("flow") ?? params.get("flow_data[type]") ?? "").toLowerCase();
  if (flow.includes("payment_method")) {
    return normalizedAction === "payment method";
  }
  if (flow.includes("subscription_cancel")) {
    return normalizedAction === "cancel";
  }
  return normalizedAction === "payment method" || normalizedAction === "cancel";
}

function portalPaymentOutcomeForCard(cardNumber: string): PortalPaymentOutcome {
  const normalized = cardNumber.replace(/\D/g, "");
  if (normalized === "4000000000000002" || normalized === "4000000000000341") {
    return "fails";
  }
  return "succeeds";
}

createRoot(document.getElementById("root")!).render(<PortalApp />);
