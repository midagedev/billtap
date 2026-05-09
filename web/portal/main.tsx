import { useEffect, useState } from "react";
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
import "../shared/styles.css";

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

  async function runAction(label: string, operation: () => Promise<{ data: PortalData; source: DataSource; message: string; error?: string }>) {
    setBusyAction(label);
    setActionStatus(`${label} pending`);
    try {
      const result = await operation();
      setPortal(result.data);
      setDraftSeats(result.data.subscription.seats);
      setSource(result.source);
      setActionStatus(result.error ? `${result.message}: ${result.error}` : result.message);
    } catch (error) {
      setActionStatus(error instanceof Error ? error.message : String(error));
    } finally {
      setBusyAction(undefined);
    }
  }

  const subscription = portal.subscription;
  const cancelMode = subscription.status === "canceled" ? "immediate" : subscription.cancelAtPeriodEnd ? "period" : "none";

  return (
    <AppShell active="portal">
      <SurfaceHeader
        eyebrow="Hosted billing portal"
        title="Subscription management"
        meta={`${subscription.workspace} · ${subscription.owner}`}
        actions={
          <a className="button secondary" href="/app/dashboard/">
            View evidence
          </a>
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
    </AppShell>
  );
}

createRoot(document.getElementById("root")!).render(<PortalApp />);
