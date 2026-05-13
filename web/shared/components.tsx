import type { ReactNode } from "react";
import { appHref } from "./basePath";
import type { StatusTone, TimelineEntry } from "./data";

type AppShellProps = {
  active: "dashboard" | "checkout" | "portal";
  children: ReactNode;
};

const navItems = [
  { id: "dashboard", label: "Dashboard", href: appHref("dashboard/") },
  { id: "checkout", label: "Checkout", href: appHref("checkout/") },
  { id: "portal", label: "Portal", href: appHref("portal/") },
] as const;

export function AppShell({ active, children }: AppShellProps) {
  const logoUrl = `${import.meta.env.BASE_URL}billtap-mark.svg`;

  return (
    <div className="app-shell">
      <aside className="sidebar" aria-label="Billtap surfaces">
        <a className="brand" href={appHref("dashboard/")} aria-label="Billtap dashboard">
          <img src={logoUrl} alt="" />
          <span>
            <strong>Billtap</strong>
            <small>Billing lab</small>
          </span>
        </a>
        <nav className="surface-nav">
          {navItems.map((item) => (
            <a key={item.id} className={item.id === active ? "active" : ""} href={item.href}>
              {item.label}
            </a>
          ))}
        </nav>
        <div className="boundary-note">
          <span>Mode</span>
          <strong>Sandbox only</strong>
        </div>
      </aside>
      <main className="surface-main">{children}</main>
    </div>
  );
}

type SurfaceHeaderProps = {
  eyebrow: string;
  title: string;
  meta: string;
  actions?: ReactNode;
};

export function SurfaceHeader({ eyebrow, title, meta, actions }: SurfaceHeaderProps) {
  return (
    <header className="surface-header">
      <div>
        <p>{eyebrow}</p>
        <h1>{title}</h1>
        <span>{meta}</span>
      </div>
      {actions ? <div className="header-actions">{actions}</div> : null}
    </header>
  );
}

type StatCardProps = {
  label: string;
  value: string;
  meta: string;
  tone: StatusTone;
};

export function StatCard({ label, value, meta, tone }: StatCardProps) {
  return (
    <section className="stat-card">
      <span>{label}</span>
      <strong>{value}</strong>
      <small className={`tone-${tone}`}>{meta}</small>
    </section>
  );
}

export function StatusPill({ tone, children }: { tone: StatusTone; children: ReactNode }) {
  return <span className={`status-pill tone-${tone}`}>{children}</span>;
}

export function Panel({ title, action, children }: { title: string; action?: ReactNode; children: ReactNode }) {
  return (
    <section className="panel">
      <div className="panel-header">
        <h2>{title}</h2>
        {action}
      </div>
      {children}
    </section>
  );
}

export function TimelineTable({ entries }: { entries: TimelineEntry[] }) {
  return (
    <div className="timeline-table" role="table" aria-label="Billing timeline">
      <div className="timeline-row timeline-head" role="row">
        <span>Time</span>
        <span>Object</span>
        <span>Transition</span>
        <span>Webhook</span>
        <span>Delivery</span>
        <span>App response</span>
      </div>
      {entries.map((entry) => (
        <div className="timeline-row" role="row" key={`${entry.time}-${entry.object}`}>
          <span>{entry.time}</span>
          <code>{entry.object}</code>
          <strong>{entry.transition}</strong>
          <span>{entry.event}</span>
          <StatusPill tone={entry.tone}>{entry.delivery}</StatusPill>
          <span>{entry.response}</span>
        </div>
      ))}
    </div>
  );
}

export function CopyValue({ label, value }: { label: string; value: string }) {
  return (
    <label className="copy-value">
      <span>{label}</span>
      <input readOnly value={value} aria-label={label} />
    </label>
  );
}
