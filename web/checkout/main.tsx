import { useEffect, useMemo, useState } from "react";
import { createRoot } from "react-dom/client";
import {
  AppShell,
  CopyValue,
  Panel,
  StatusPill,
  SurfaceHeader,
} from "../shared/components";
import {
  completeCheckout,
  getCheckoutSessionId,
  loadCheckoutSession,
  type CheckoutCompletion,
  type CheckoutSessionSummary,
  type DataSource,
} from "../shared/api";
import { checkoutSession as fixtureCheckoutSession, outcomeOptions, type StatusTone } from "../shared/data";
import "../shared/styles.css";

const paymentMethods = ["4242 sandbox card", "3DS challenge card", "Bank debit pending"];

function CheckoutApp() {
  const [outcome, setOutcome] = useState(outcomeOptions[0]);
  const [paymentMethod, setPaymentMethod] = useState(paymentMethods[0]);
  const [session, setSession] = useState<CheckoutSessionSummary>(fixtureCheckoutSession);
  const [source, setSource] = useState<DataSource>("fixture");
  const [loadError, setLoadError] = useState<string>();
  const [completion, setCompletion] = useState<CheckoutCompletion>();
  const [isLoading, setIsLoading] = useState(true);
  const [isCompleting, setIsCompleting] = useState(false);
  const sessionId = useMemo(() => getCheckoutSessionId(), []);

  useEffect(() => {
    let active = true;

    loadCheckoutSession(sessionId).then((result) => {
      if (!active) return;
      setSession(result.data);
      setSource(result.source);
      setLoadError(result.error);
      setIsLoading(false);
    });

    return () => {
      active = false;
    };
  }, [sessionId]);

  const confirmation = useMemo(() => {
    if (outcome.id === "success") return "subscription active · invoice paid · payment intent succeeded";
    if (outcome.id === "pending") return "checkout pending · invoice open · async payment awaiting settlement";
    if (outcome.id === "action") return "payment intent requires_action · challenge event queued";
    if (outcome.id === "cancel") return "checkout canceled · return URL selected";
    return "checkout failed · invoice open · retry path available";
  }, [outcome.id]);

  const statusTone = completion ? toneForState(completion.session.paymentIntentStatus) : outcome.tone;
  const apiTone: StatusTone = source === "api" ? "good" : "warning";

  async function handleCompleteCheckout() {
    setIsCompleting(true);
    const result = await completeCheckout(session, outcome.id, paymentMethod);
    setCompletion(result);
    setSession(result.session);
    setSource(result.source);
    setIsCompleting(false);
  }

  return (
    <AppShell active="checkout">
      <SurfaceHeader
        eyebrow="Hosted checkout"
        title="Checkout session"
        meta={`${session.id} · ${session.customer}`}
        actions={
          <a className="button secondary" href="/app/dashboard/">
            View timeline
          </a>
        }
      />

      <div className="checkout-layout">
        <section className="checkout-summary" aria-label="Checkout session summary">
          <div className="summary-head">
            <div>
              <span className="section-kicker">Plan</span>
              <h2>{session.plan}</h2>
            </div>
            <div className="summary-total">
              <span className="section-kicker">Due today</span>
              <strong>{session.price}</strong>
            </div>
          </div>

          <ul className="line-list">
            {session.lineItems.map((line) => (
              <li key={line.label}>
                <span>{line.label}</span>
                <strong>{line.amount}</strong>
              </li>
            ))}
          </ul>

          <div className="confirmation-strip">
            <StatusPill tone={outcome.tone}>{outcome.status}</StatusPill>
            <StatusPill tone="neutral">{paymentMethod}</StatusPill>
            <StatusPill tone={apiTone}>{source === "api" ? "API loaded" : "fixture fallback"}</StatusPill>
          </div>
        </section>

        <div className="surface-grid">
          <Panel title="Customer identity">
            <div className="surface-grid">
              <CopyValue label="Customer" value={session.customer} />
              <CopyValue label="Email" value={session.customerEmail} />
              <CopyValue label="Session source" value={loadError ? `${source}: ${loadError}` : source} />
            </div>
          </Panel>

          <Panel title="Payment method">
            <div className="payment-grid">
              {paymentMethods.map((method) => (
                <button
                  className={`payment-card ${method === paymentMethod ? "active" : ""}`}
                  key={method}
                  onClick={() => setPaymentMethod(method)}
                  type="button"
                >
                  <strong>{method}</strong>
                  <span>Sandbox instrument</span>
                </button>
              ))}
            </div>
          </Panel>
        </div>
      </div>

      <Panel title="Outcome selector">
        <div className="option-grid">
          {outcomeOptions.map((option) => (
            <button
              className={`choice-card ${option.id === outcome.id ? "active" : ""}`}
              key={option.id}
              onClick={() => setOutcome(option)}
              type="button"
            >
              <strong>{option.label}</strong>
              <StatusPill tone={option.tone}>{option.status}</StatusPill>
            </button>
          ))}
        </div>
      </Panel>

      <Panel
        title="Confirmation"
        action={<StatusPill tone={statusTone}>{completion ? completion.session.paymentIntentStatus : outcome.label}</StatusPill>}
      >
        <div className="surface-grid">
          <CopyValue label="Resulting state" value={confirmation} />
          {completion ? (
            <CopyValue
              label="Completion response"
              value={completion.error ? `${completion.message}: ${completion.error}` : completion.message}
            />
          ) : null}
          <div className="confirmation-strip">
            <button className="button" disabled={isLoading || isCompleting} onClick={handleCompleteCheckout} type="button">
              {isCompleting ? "Completing..." : "Complete checkout"}
            </button>
            <a className="button secondary" href={session.returnUrl}>
              Return to app
            </a>
          </div>
        </div>
      </Panel>

      <Panel title="Billing state returned">
        <ul className="object-list">
          <li>
            <span>Checkout session</span>
            <strong>{session.status}</strong>
          </li>
          <li>
            <span>{session.subscriptionId}</span>
            <StatusPill tone={toneForState(session.subscriptionStatus)}>{session.subscriptionStatus}</StatusPill>
          </li>
          <li>
            <span>{session.invoiceId}</span>
            <StatusPill tone={toneForState(session.invoiceStatus)}>{session.invoiceStatus}</StatusPill>
          </li>
          <li>
            <span>{session.paymentIntentId}</span>
            <StatusPill tone={toneForState(session.paymentIntentStatus)}>{session.paymentIntentStatus}</StatusPill>
          </li>
        </ul>
      </Panel>
    </AppShell>
  );
}

function toneForState(state: string): StatusTone {
  const value = state.toLowerCase();
  if (value.includes("succeeded") || value.includes("paid") || value.includes("active") || value.includes("complete")) {
    return "good";
  }
  if (value.includes("fail") || value.includes("declin") || value.includes("canceled") || value.includes("void")) {
    return "danger";
  }
  if (value.includes("pending") || value.includes("open") || value.includes("requires") || value.includes("incomplete")) {
    return "warning";
  }
  return "neutral";
}

createRoot(document.getElementById("root")!).render(<CheckoutApp />);
