import { resolve } from "node:path";
import { fileURLToPath } from "node:url";
import react from "@vitejs/plugin-react";
import { defineConfig, type Plugin } from "vite";

const appRoutes = new Set(["checkout", "dashboard", "portal"]);
const projectRoot = fileURLToPath(new URL(".", import.meta.url));
const publicBasePath = normalizePublicBasePath(firstNonEmpty(process.env.BILLTAP_PUBLIC_BASE_PATH, process.env.PUBLIC_BASE_PATH));
const appBasePath = joinBasePath(publicBasePath, "/app/");

function appPathDevFallback(): Plugin {
  return {
    name: "billtap-app-path-dev-fallback",
    apply: "serve",
    configureServer(server) {
      server.middlewares.use((req, _res, next) => {
        const [pathname, query] = req.url?.split("?") ?? ["", ""];
        const appPath = stripBasePath(pathname, publicBasePath);
        const match = appPath.match(/^\/app\/(checkout|dashboard|portal)\/?$/);

        if (match && appRoutes.has(match[1])) {
          req.url = `${appBasePath}${match[1]}/index.html${query ? `?${query}` : ""}`;
        } else if (appPath === "/app/" || appPath === "/app") {
          req.url = `${appBasePath}dashboard/index.html${query ? `?${query}` : ""}`;
        }

        next();
      });
    },
  };
}

function firstNonEmpty(...values: Array<string | undefined>): string {
  for (const value of values) {
    if (value?.trim()) return value;
  }
  return "";
}

function normalizePublicBasePath(value: string): string {
  const trimmed = value.trim();
  if (!trimmed || trimmed === "/") return "";
  const withLeading = trimmed.startsWith("/") ? trimmed : `/${trimmed}`;
  return withLeading.replace(/\/+$/, "");
}

function joinBasePath(basePath: string, path: string): string {
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;
  if (!basePath) return normalizedPath;
  return `${basePath}${normalizedPath}`;
}

function stripBasePath(pathname: string, basePath: string): string {
  if (!basePath) return pathname;
  if (pathname === basePath) return "/";
  if (pathname.startsWith(`${basePath}/`)) return pathname.slice(basePath.length);
  return pathname;
}

export default defineConfig({
  root: "web",
  base: appBasePath,
  plugins: [react(), appPathDevFallback()],
  server: {
    host: "127.0.0.1",
    port: 5173,
  },
  preview: {
    host: "127.0.0.1",
    port: 4173,
  },
  build: {
    outDir: "../dist/app",
    emptyOutDir: true,
    rollupOptions: {
      input: {
        checkout: resolve(projectRoot, "web/checkout/index.html"),
        dashboard: resolve(projectRoot, "web/dashboard/index.html"),
        portal: resolve(projectRoot, "web/portal/index.html"),
      },
    },
  },
});
