export function appHref(path = ""): string {
  return joinBase(joinBase(publicBasePath(), "/app/"), path);
}

export function apiHref(path: string): string {
  return joinBase(publicBasePath(), path);
}

function publicBasePath(pathname = globalThis.location?.pathname ?? "/"): string {
  const match = pathname.match(/^(.*?)(?:\/app(?:\/|$)|\/checkout(?:\/|$)|\/portal(?:\/|$))/);
  if (!match) return "";
  return normalizeBasePath(match[1]);
}

function normalizeBasePath(value: string): string {
  const trimmed = value.trim();
  if (!trimmed || trimmed === "/") return "";
  const withLeading = trimmed.startsWith("/") ? trimmed : `/${trimmed}`;
  return withLeading.replace(/\/+$/, "");
}

function joinBase(base: string, path: string): string {
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;
  const normalizedBase = base.endsWith("/") ? base.slice(0, -1) : base;
  if (!normalizedBase) return normalizedPath;
  return `${normalizedBase}${normalizedPath}`;
}
