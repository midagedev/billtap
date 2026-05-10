import { useEffect, useMemo, useState, type FormEvent } from "react";
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
import {
  checkoutSession as fixtureCheckoutSession,
  outcomeOptions,
  type CheckoutOutcomeId,
  type StatusTone,
} from "../shared/data";
import "../shared/styles.css";

const paymentMethods = [
  "4242 sandbox card",
  "3DS challenge card",
  "Bank debit pending",
];
const checkoutSuccessMessage = "Your payment has been processed successfully.";

function CheckoutApp() {
  const [outcome, setOutcome] = useState(outcomeOptions[0]);
  const [paymentMethod, setPaymentMethod] = useState(paymentMethods[0]);
  const [session, setSession] = useState<CheckoutSessionSummary>(
    fixtureCheckoutSession,
  );
  const [source, setSource] = useState<DataSource>("fixture");
  const [loadError, setLoadError] = useState<string>();
  const [completion, setCompletion] = useState<CheckoutCompletion>();
  const [isLoading, setIsLoading] = useState(true);
  const [isCompleting, setIsCompleting] = useState(false);
  const [cardNumber, setCardNumber] = useState("4242424242424242");
  const [cardExpiry, setCardExpiry] = useState("12/34");
  const [cardCvc, setCardCvc] = useState("123");
  const [billingName, setBillingName] = useState("Jane Doe");
  const [billingCountry, setBillingCountry] = useState("US");
  const [formMessage, setFormMessage] = useState<string>();
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
    if (outcome.id === "success")
      return "subscription active · invoice paid · payment intent succeeded";
    if (outcome.id === "pending")
      return "checkout pending · invoice open · async payment awaiting settlement";
    if (outcome.id === "action")
      return "payment intent requires_action · challenge event queued";
    if (outcome.id === "cancel")
      return "checkout canceled · return URL selected";
    return "checkout failed · invoice open · retry path available";
  }, [outcome.id]);

  const statusTone = completion
    ? toneForState(completion.session.paymentIntentStatus)
    : outcome.tone;
  const apiTone: StatusTone = source === "api" ? "good" : "warning";
  const paymentSucceeded = completion
    ? isSuccessfulCompletion(completion)
    : false;

  useEffect(() => {
    if (!paymentSucceeded || !window.opener) return;
    window.opener.postMessage({ success: true, data: "결제성공" }, "*");
    window.opener.focus?.();
  }, [paymentSucceeded]);

  async function handleCompleteCheckout() {
    setIsCompleting(true);
    setFormMessage(undefined);
    const result = await completeCheckout(session, outcome.id, paymentMethod);
    setCompletion(result);
    setSession(result.session);
    setSource(result.source);
    setIsCompleting(false);
  }

  async function handleStripeCompatibleSubmit(
    event: FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault();
    setIsCompleting(true);
    setFormMessage(undefined);
    const cardOutcome = checkoutOutcomeForCard(cardNumber);
    const result = await completeCheckout(
      session,
      cardOutcome,
      paymentMethodForCard(cardNumber),
    );
    setCompletion(result);
    setSession(result.session);
    setSource(result.source);
    setIsCompleting(false);
    if (!isSuccessfulCompletion(result)) {
      setFormMessage(failureMessageForOutcome(cardOutcome));
    }
  }

  if (paymentSucceeded) {
    return (
      <AppShell active="checkout">
        <section className="stripe-success-shell" aria-live="polite">
          <span className="stripe-success-icon" aria-hidden="true">
            OK
          </span>
          <h1>Payment successful</h1>
          <h2>{checkoutSuccessMessage}</h2>
          <a className="button secondary" href={session.returnUrl}>
            Return to app
          </a>
        </section>
      </AppShell>
    );
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
        <section
          className="checkout-summary"
          aria-label="Checkout session summary"
        >
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
            <StatusPill tone={apiTone}>
              {source === "api" ? "API loaded" : "fixture fallback"}
            </StatusPill>
          </div>
        </section>

        <div className="surface-grid">
          <Panel title="Customer identity">
            <div className="surface-grid">
              <CopyValue label="Customer" value={session.customer} />
              <CopyValue label="Email" value={session.customerEmail} />
              <CopyValue
                label="Session source"
                value={loadError ? `${source}: ${loadError}` : source}
              />
            </div>
          </Panel>

          <Panel title="Payment method">
            <form
              className="stripe-compatible-form"
              onSubmit={handleStripeCompatibleSubmit}
            >
              <button
                className="payment-method-tab"
                type="button"
                onClick={() => setPaymentMethod(paymentMethods[0])}
              >
                Card
              </button>
              <label>
                <span>Email</span>
                <input
                  id="email"
                  value={session.customerEmail}
                  onChange={() => undefined}
                  autoComplete="email"
                />
              </label>
              <label>
                <span>Card number</span>
                <input
                  id="cardNumber"
                  value={cardNumber}
                  onChange={(event) => setCardNumber(event.target.value)}
                  autoComplete="cc-number"
                  inputMode="numeric"
                />
              </label>
              <div className="stripe-form-row">
                <label>
                  <span>Expiry</span>
                  <input
                    id="cardExpiry"
                    value={cardExpiry}
                    onChange={(event) => setCardExpiry(event.target.value)}
                    autoComplete="cc-exp"
                  />
                </label>
                <label>
                  <span>CVC</span>
                  <input
                    id="cardCvc"
                    value={cardCvc}
                    onChange={(event) => setCardCvc(event.target.value)}
                    autoComplete="cc-csc"
                  />
                </label>
              </div>
              <label>
                <span>Name on card</span>
                <input
                  name="billingName"
                  value={billingName}
                  onChange={(event) => setBillingName(event.target.value)}
                  autoComplete="cc-name"
                />
              </label>
              <label>
                <span>Country</span>
                <select
                  name="billingCountry"
                  value={billingCountry}
                  onChange={(event) => setBillingCountry(event.target.value)}
                >
                  <option value="US">United States</option>
                  <option value="KR">South Korea</option>
                  <option value="JP">Japan</option>
                  <option value="GB">United Kingdom</option>
                </select>
              </label>
              {formMessage ? (
                <div className="stripe-form-error" role="alert">
                  {formMessage}
                </div>
              ) : null}
              <button
                className="button"
                disabled={isLoading || isCompleting}
                type="submit"
              >
                {isCompleting ? "Processing..." : "Pay"}
              </button>
            </form>
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
        action={
          <StatusPill tone={statusTone}>
            {completion
              ? completion.session.paymentIntentStatus
              : outcome.label}
          </StatusPill>
        }
      >
        <div className="surface-grid">
          <CopyValue label="Resulting state" value={confirmation} />
          {completion ? (
            <CopyValue
              label="Completion response"
              value={
                completion.error
                  ? `${completion.message}: ${completion.error}`
                  : completion.message
              }
            />
          ) : null}
          <div className="confirmation-strip">
            <button
              className="button"
              disabled={isLoading || isCompleting}
              onClick={handleCompleteCheckout}
              type="button"
            >
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
            <StatusPill tone={toneForState(session.subscriptionStatus)}>
              {session.subscriptionStatus}
            </StatusPill>
          </li>
          <li>
            <span>{session.invoiceId}</span>
            <StatusPill tone={toneForState(session.invoiceStatus)}>
              {session.invoiceStatus}
            </StatusPill>
          </li>
          <li>
            <span>{session.paymentIntentId}</span>
            <StatusPill tone={toneForState(session.paymentIntentStatus)}>
              {session.paymentIntentStatus}
            </StatusPill>
          </li>
        </ul>
      </Panel>
    </AppShell>
  );
}

function toneForState(state: string): StatusTone {
  const value = state.toLowerCase();
  if (
    value.includes("succeeded") ||
    value.includes("paid") ||
    value.includes("active") ||
    value.includes("complete")
  ) {
    return "good";
  }
  if (
    value.includes("fail") ||
    value.includes("declin") ||
    value.includes("canceled") ||
    value.includes("void")
  ) {
    return "danger";
  }
  if (
    value.includes("pending") ||
    value.includes("open") ||
    value.includes("requires") ||
    value.includes("incomplete")
  ) {
    return "warning";
  }
  return "neutral";
}

function checkoutOutcomeForCard(cardNumber: string): CheckoutOutcomeId {
  const normalized = cardNumber.replace(/\D/g, "");
  switch (normalized) {
    case "4000000000000002":
      return "declined";
    case "4000000000009995":
      return "funds";
    case "4000000000000069":
      return "expired";
    case "4000000000000127":
      return "incorrect_cvc";
    case "4000000000000119":
      return "processing_error";
    case "4000000000003220":
      return "action";
    default:
      return "success";
  }
}

function paymentMethodForCard(cardNumber: string): string {
  const normalized = cardNumber.replace(/\D/g, "");
  switch (normalized) {
    case "4000000000000002":
      return "pm_card_visa_chargeDeclined";
    case "4000000000009995":
      return "pm_card_visa_chargeDeclinedInsufficientFunds";
    case "4000000000003220":
      return "pm_card_threeDSecure2Required";
    default:
      return "pm_card_visa";
  }
}

function failureMessageForOutcome(outcome: CheckoutOutcomeId): string {
  if (outcome === "funds") return "Your card has insufficient funds.";
  if (outcome === "expired") return "Your card has expired.";
  if (outcome === "incorrect_cvc")
    return "Your card's security code is incorrect.";
  if (outcome === "processing_error")
    return "An error occurred while processing your card. Try again later.";
  if (outcome === "action")
    return "Payment failed because authentication is required.";
  return "Your card was declined.";
}

function isSuccessfulCompletion(completion: CheckoutCompletion): boolean {
  const paymentStatus = completion.session.paymentStatus.toLowerCase();
  const paymentIntentStatus =
    completion.session.paymentIntentStatus.toLowerCase();
  return (
    paymentStatus === "paid" ||
    paymentStatus === "no_payment_required" ||
    paymentIntentStatus === "succeeded"
  );
}

createRoot(document.getElementById("root")!).render(<CheckoutApp />);
