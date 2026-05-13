const appBase = normalizeBase(import.meta.env.BASE_URL || "/app/");
const publicBasePath = appBase.endsWith("/app/") ? appBase.slice(0, -"/app/".length) : appBase.replace(/\/$/, "");

export function appHref(path = ""): string {
  return joinBase(appBase, path);
}

export function apiHref(path: string): string {
  return joinBase(publicBasePath, path);
}

function normalizeBase(value: string): string {
  const trimmed = value.trim();
  if (!trimmed || trimmed === "/") return "/app/";
  const withLeading = trimmed.startsWith("/") ? trimmed : `/${trimmed}`;
  return withLeading.endsWith("/") ? withLeading : `${withLeading}/`;
}

function joinBase(base: string, path: string): string {
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;
  const normalizedBase = base.endsWith("/") ? base.slice(0, -1) : base;
  if (!normalizedBase) return normalizedPath;
  return `${normalizedBase}${normalizedPath}`;
}
