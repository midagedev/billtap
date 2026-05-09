import { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import {
  AppShell,
  CopyValue,
  Panel,
  StatCard,
  StatusPill,
  SurfaceHeader,
  TimelineTable,
} from "../shared/components";
import {
  createDebugBundle,
  loadDashboardData,
  replayWebhookEvent,
  type DashboardData,
  type DataSource,
  type DebugBundleResult,
} from "../shared/api";
import {
  billingSnapshot,
  dashboardObjects,
  dashboardStats,
  timelineEntries,
  webhookAttempts,
  type DashboardObjectRecord,
  type DashboardObjectType,
  type TimelineEntry,
  type WebhookAttempt,
} from "../shared/data";
import "../shared/styles.css";

const objectTabs: Array<{ id: DashboardObjectType; label: string }> = [
  { id: "customers", label: "Customers" },
  { id: "subscriptions", label: "Subscriptions" },
  { id: "invoices", label: "Invoices" },
  { id: "paymentIntents", label: "Payment intents" },
  { id: "checkoutSessions", label: "Checkout sessions" },
  { id: "webhookEvents", label: "Webhook events" },
];

function DashboardApp() {
  const [dashboard, setDashboard] = useState<DashboardData>({
    billingSnapshot,
    dashboardStats,
    objects: dashboardObjects,
    timelineEntries,
    webhookAttempts,
    appResponses: webhookAttempts.map((attempt) => ({
      id: attempt.id ?? attempt.event,
      eventId: attempt.eventId ?? attempt.event,
      endpoint: attempt.endpoint,
      status: attempt.responseStatus ?? attempt.status,
      body: attempt.responseBody ?? attempt.errorMessage ?? "No response body recorded",
      retryPlan: attempt.retryPlan ?? "retry policy default",
      tone: attempt.tone ?? "neutral",
    })),
    evidence: {
      signatureHeader: "t=1778233312,v1=sandbox_mismatch",
      idempotencyKey: "billtap:evt_test_90:attempt_2",
      failureReason: "handler returned 503 before entitlement update",
    },
    scenario: {
      currentStep: "retry invoice payment",
      clockAdvance: "+3 days",
      assertions: "2 pending",
      tone: "warning",
    },
  });
  const [source, setSource] = useState<DataSource>("fixture");
  const [loadError, setLoadError] = useState<string>();
  const [replayStatus, setReplayStatus] = useState<Record<string, string>>({});
  const [activeType, setActiveType] = useState<DashboardObjectType>("customers");
  const [selectedIds, setSelectedIds] = useState<Partial<Record<DashboardObjectType, string>>>({});
  const [selectedAttemptId, setSelectedAttemptId] = useState<string>();
  const [bundleStatus, setBundleStatus] = useState<string>("No bundle requested");
  const [bundleResult, setBundleResult] = useState<DebugBundleResult>();

  useEffect(() => {
    let active = true;
    loadDashboardData().then((result) => {
      if (!active) return;
      setDashboard(result.data);
      setSource(result.source);
      setLoadError(result.error);
      setSelectedAttemptId((current) => current ?? result.data.webhookAttempts[0]?.id);
    });

    return () => {
      active = false;
    };
  }, []);

  async function handleReplay(attempt: WebhookAttempt) {
    if (!attempt.eventId) return;
    setReplayStatus((current) => ({ ...current, [attempt.eventId!]: "requesting replay" }));
    try {
      const message = await replayWebhookEvent(attempt.eventId);
      setReplayStatus((current) => ({ ...current, [attempt.eventId!]: message }));
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);
      setReplayStatus((current) => ({ ...current, [attempt.eventId!]: message }));
    }
  }

  const activeObjects = dashboard.objects[activeType];
  const selectedObject =
    activeObjects.find((item) => item.id === selectedIds[activeType]) ?? activeObjects[0];
  const selectedAttempt =
    dashboard.webhookAttempts.find((attempt) => attempt.id === selectedAttemptId) ??
    relatedAttempts(dashboard.webhookAttempts, selectedObject)[0] ??
    dashboard.webhookAttempts[0];
  const filteredTimeline = relatedTimeline(dashboard.timelineEntries, selectedObject);
  const visibleTimeline = filteredTimeline.length > 0 ? filteredTimeline : dashboard.timelineEntries.slice(0, 5);
  const visibleResponses = relatedAppResponses(dashboard.appResponses, selectedObject, selectedAttempt);

  async function handleDebugBundle(object: DashboardObjectRecord) {
    setBundleStatus(`Requesting debug bundle for ${object.id}`);
    setBundleResult(undefined);
    try {
      const result = await createDebugBundle({ objectType: object.type, objectId: object.id });
      setBundleResult(result);
      setBundleStatus(result.message);
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);
      setBundleStatus(`Debug bundle export failed: ${message}`);
    }
  }

  return (
    <AppShell active="dashboard">
      <SurfaceHeader
        eyebrow="Developer dashboard"
        title="Debug workspace"
        meta={`${dashboard.billingSnapshot.workspace} · ${source === "api" ? "API timeline" : "fixture timeline"}`}
        actions={
          <>
            <a className="button secondary" href="/app/checkout/">
              Open checkout
            </a>
            <a className="button" href="/app/portal/">
              Open portal
            </a>
          </>
        }
      />

      <div className="status-strip">
        <StatusPill tone={source === "api" ? "good" : "warning"}>
          {source === "api" ? "Dashboard API loaded" : "Fixture fallback"}
        </StatusPill>
        <span>{loadError ?? "Timeline and delivery attempts loaded where available"}</span>
      </div>

      <section className="stat-grid" aria-label="Dashboard metrics">
        {dashboard.dashboardStats.map((stat) => (
          <StatCard key={stat.label} {...stat} />
        ))}
      </section>

      <div className="surface-grid dashboard-grid">
        <div>
          <Panel
            title="Billing objects"
            action={<StatusPill tone={selectedObject.tone}>{selectedObject.status}</StatusPill>}
          >
            <div className="object-tabs" role="tablist" aria-label="Billing object types">
              {objectTabs.map((tab) => (
                <button
                  aria-selected={activeType === tab.id}
                  className={activeType === tab.id ? "active" : ""}
                  key={tab.id}
                  onClick={() => setActiveType(tab.id)}
                  role="tab"
                  type="button"
                >
                  <span>{tab.label}</span>
                  <strong>{dashboard.objects[tab.id].length}</strong>
                </button>
              ))}
            </div>

            <div className="object-browser">
              <ul className="object-record-list" aria-label={`${activeType} list`}>
                {activeObjects.map((object) => (
                  <li key={object.id}>
                    <button
                      className={object.id === selectedObject.id ? "active" : ""}
                      onClick={() =>
                        setSelectedIds((current) => ({
                          ...current,
                          [activeType]: object.id,
                        }))
                      }
                      type="button"
                    >
                      <span>
                        <strong>{object.title}</strong>
                        <small>{object.subtitle}</small>
                      </span>
                      <span>
                        <StatusPill tone={object.tone}>{object.status}</StatusPill>
                        <code>{object.id}</code>
                      </span>
                    </button>
                  </li>
                ))}
              </ul>

              <ObjectDetail object={selectedObject} onExport={() => void handleDebugBundle(selectedObject)} />
            </div>
          </Panel>

          <Panel
            title="Object timeline"
            action={
              <span className="panel-note">
                {filteredTimeline.length > 0 ? `${filteredTimeline.length} linked` : "showing recent"}
              </span>
            }
          >
            <TimelineTable entries={visibleTimeline} />
          </Panel>

          <Panel
            title="Webhook detail evidence"
            action={selectedAttempt ? <StatusPill tone={selectedAttempt.tone ?? "neutral"}>{selectedAttempt.status}</StatusPill> : null}
          >
            {selectedAttempt ? (
              <div className="webhook-evidence-grid">
                <div>
                  <ul className="attempt-picker" aria-label="Webhook attempts">
                    {dashboard.webhookAttempts.map((attempt) => (
                      <li key={attempt.id ?? `${attempt.eventId}-${attempt.attemptNumber}`}>
                        <button
                          className={attempt.id === selectedAttempt.id ? "active" : ""}
                          onClick={() => setSelectedAttemptId(attempt.id)}
                          type="button"
                        >
                          <span>
                            <strong>{attempt.event}</strong>
                            <small>{attempt.endpoint}</small>
                          </span>
                          <StatusPill tone={attempt.tone ?? "neutral"}>{attempt.status}</StatusPill>
                        </button>
                      </li>
                    ))}
                  </ul>
                </div>
                <div className="evidence-stack">
                  <div className="detail-grid compact">
                    <DetailItem label="Event ID" value={selectedAttempt.eventId ?? "-"} code />
                    <DetailItem label="Attempt" value={String(selectedAttempt.attemptNumber ?? selectedAttempt.attempts)} />
                    <DetailItem label="Endpoint" value={selectedAttempt.requestUrl ?? selectedAttempt.endpoint} />
                    <DetailItem label="Response" value={selectedAttempt.responseStatus ?? "-"} />
                    <DetailItem label="Scheduled" value={selectedAttempt.scheduledAt ?? "-"} />
                    <DetailItem label="Delivered" value={selectedAttempt.deliveredAt ?? "-"} />
                    <DetailItem label="Retry plan" value={selectedAttempt.retryPlan ?? "retry policy default"} />
                  </div>
                  <div className="surface-grid">
                    <CopyValue label="Signature header" value={selectedAttempt.signatureHeader ?? dashboard.evidence.signatureHeader} />
                    <CopyValue label="Request headers" value={selectedAttempt.requestHeaders ?? "headers not recorded"} />
                    <CopyValue label="Request body" value={selectedAttempt.requestBody ?? "payload not recorded"} />
                    <CopyValue label="Response body" value={selectedAttempt.responseBody ?? selectedAttempt.errorMessage ?? "body not recorded"} />
                  </div>
                </div>
              </div>
            ) : (
              <p className="empty-state">No webhook attempts recorded yet.</p>
            )}
          </Panel>

          <Panel
            title="Workspace entitlement"
            action={<StatusPill tone="info">{dashboard.billingSnapshot.status}</StatusPill>}
          >
            <div className="snapshot-grid">
              <div>
                <span className="section-kicker">Plan</span>
                <strong>{dashboard.billingSnapshot.plan}</strong>
                <span>{dashboard.billingSnapshot.period}</span>
              </div>
              <div>
                <span className="section-kicker">Seats</span>
                <strong>
                  {dashboard.billingSnapshot.seats.used}/
                  {dashboard.billingSnapshot.seats.basic + dashboard.billingSnapshot.seats.additional} used
                </strong>
                <span>{dashboard.billingSnapshot.seats.pending} pending invitation</span>
              </div>
              <div>
                <span className="section-kicker">Exports</span>
                <strong>{dashboard.billingSnapshot.export.remaining} remaining</strong>
                <span>Renews {dashboard.billingSnapshot.export.renews}</span>
              </div>
              <div>
                <span className="section-kicker">Payment rail</span>
                <strong>{dashboard.billingSnapshot.paymentRail}</strong>
                <span>{dashboard.billingSnapshot.tenantRoute}</span>
              </div>
            </div>
          </Panel>
        </div>

        <aside>
          <Panel
            title="Debug bundle export"
            action={<StatusPill tone={bundleResult?.source === "api" ? "good" : "neutral"}>{bundleResult?.source ?? "ready"}</StatusPill>}
          >
            <div className="debug-bundle">
              <span className="section-kicker">Selected object</span>
              <strong>{selectedObject.id}</strong>
              <span>{selectedObject.label} · {selectedObject.status}</span>
              <button className="button" onClick={() => void handleDebugBundle(selectedObject)} type="button">
                Export debug bundle
              </button>
              <CopyValue label="Bundle status" value={bundleStatus} />
              {bundleResult?.url ? <CopyValue label="Bundle URL" value={bundleResult.url} /> : null}
              {bundleResult ? <CopyValue label="Bundle ID" value={bundleResult.id} /> : null}
            </div>
          </Panel>

          <Panel title="App response area">
            <ul className="app-response-list">
              {visibleResponses.map((response) => (
                <li key={response.id}>
                  <div>
                    <strong>{response.eventId}</strong>
                    <span>{response.endpoint}</span>
                  </div>
                  <StatusPill tone={response.tone}>{response.status}</StatusPill>
                  <CopyValue label="Handler response" value={response.body} />
                  <span>{response.retryPlan}</span>
                </li>
              ))}
            </ul>
          </Panel>

          <Panel title="Webhook delivery attempts">
            <ul className="attempt-list">
              {dashboard.webhookAttempts.map((attempt) => (
                <li key={attempt.id ?? `${attempt.endpoint}-${attempt.event}-${attempt.attemptNumber ?? attempt.attempts}`}>
                  <div>
                    <strong>{attempt.event}</strong>
                    <span>{attempt.endpoint}</span>
                    <span>{attempt.signatureHeader ?? "signature pending"}</span>
                  </div>
                  <div>
                    <StatusPill tone={attempt.tone ?? (attempt.status === "delivered" ? "good" : "warning")}>
                      {attempt.status}
                    </StatusPill>
                    <span>{attempt.attempts} attempt(s)</span>
                  </div>
                  <div className="attempt-detail">
                    <span>Scheduled {attempt.scheduledAt ?? "-"}</span>
                    <span>Delivered {attempt.deliveredAt ?? "-"}</span>
                    <span>Response {attempt.responseStatus ?? "-"}</span>
                    <span>{attempt.retryPlan ?? "retry policy default"}</span>
                  </div>
                  {attempt.eventId ? (
                    <div className="attempt-actions">
                      <button className="button secondary" onClick={() => void handleReplay(attempt)} type="button">
                        Replay
                      </button>
                      <span>{replayStatus[attempt.eventId] ?? "ready"}</span>
                    </div>
                  ) : null}
                </li>
              ))}
            </ul>
          </Panel>

          <Panel title="Copyable failure evidence">
            <div className="surface-grid">
              <CopyValue label="Signature header" value={dashboard.evidence.signatureHeader} />
              <CopyValue label="Idempotency key" value={dashboard.evidence.idempotencyKey} />
              <CopyValue label="Failure reason" value={dashboard.evidence.failureReason} />
            </div>
          </Panel>

          <Panel title="Scenario run">
            <ul className="object-list">
              <li>
                <span>Current step</span>
                <strong>{dashboard.scenario.currentStep}</strong>
              </li>
              <li>
                <span>Clock advance</span>
                <strong>{dashboard.scenario.clockAdvance}</strong>
              </li>
              <li>
                <span>Assertions</span>
                <StatusPill tone={dashboard.scenario.tone}>{dashboard.scenario.assertions}</StatusPill>
              </li>
            </ul>
          </Panel>
        </aside>
      </div>
    </AppShell>
  );
}

function ObjectDetail({ object, onExport }: { object: DashboardObjectRecord; onExport: () => void }) {
  return (
    <section className="object-detail" aria-label={`${object.label} detail`}>
      <div className="object-detail-head">
        <div>
          <span className="section-kicker">{object.label}</span>
          <h2>{object.title}</h2>
          <code>{object.id}</code>
        </div>
        <button className="button secondary" onClick={onExport} type="button">
          Bundle
        </button>
      </div>
      <div className="detail-grid">
        <DetailItem label="Status" value={object.status} />
        <DetailItem label="Subtitle" value={object.subtitle} />
        {object.amount ? <DetailItem label="Amount" value={object.amount} /> : null}
        {object.createdAt ? <DetailItem label="Created" value={object.createdAt} /> : null}
        {object.fields.map((field) => (
          <DetailItem code={field.value.includes("_test") || field.value.startsWith("cus_")} key={field.label} {...field} />
        ))}
      </div>
    </section>
  );
}

function DetailItem({ label, value, code = false }: { label: string; value: string; code?: boolean }) {
  return (
    <div className="detail-item">
      <span>{label}</span>
      {code ? <code>{value}</code> : <strong>{value}</strong>}
    </div>
  );
}

function relatedTimeline(entries: TimelineEntry[], object: DashboardObjectRecord): TimelineEntry[] {
  const needles = new Set([object.id, ...object.fields.map((field) => field.value).filter((value) => value.includes("_"))]);
  return entries.filter((entry) =>
    Array.from(needles).some((needle) =>
      [entry.object, entry.transition, entry.event, entry.delivery, entry.response].some((value) => value.includes(needle)),
    ),
  );
}

function relatedAttempts(attempts: WebhookAttempt[], object: DashboardObjectRecord): WebhookAttempt[] {
  if (object.type === "webhookEvents") {
    return attempts.filter((attempt) => attempt.eventId === object.id || attempt.event === object.title);
  }
  return attempts.filter((attempt) =>
    [attempt.eventId, attempt.event, attempt.requestBody, attempt.responseBody].some((value) => value?.includes(object.id)),
  );
}

function relatedAppResponses(
  responses: DashboardData["appResponses"],
  object: DashboardObjectRecord,
  selectedAttempt?: WebhookAttempt,
) {
  const related = responses.filter(
    (response) =>
      response.eventId === object.id ||
      response.eventId === selectedAttempt?.eventId ||
      response.body.includes(object.id),
  );
  return related.length > 0 ? related : responses.slice(0, 3);
}

createRoot(document.getElementById("root")!).render(<DashboardApp />);
