import { resolve } from "node:path";
import { fileURLToPath } from "node:url";
import react from "@vitejs/plugin-react";
import { defineConfig, type Plugin } from "vite";

const appRoutes = new Set(["checkout", "dashboard", "portal"]);
const projectRoot = fileURLToPath(new URL(".", import.meta.url));

function appPathDevFallback(): Plugin {
  return {
    name: "billtap-app-path-dev-fallback",
    apply: "serve",
    configureServer(server) {
      server.middlewares.use((req, _res, next) => {
        const [pathname, query] = req.url?.split("?") ?? ["", ""];
        const match = pathname.match(/^\/app\/(checkout|dashboard|portal)\/?$/);

        if (match && appRoutes.has(match[1])) {
          req.url = `/app/${match[1]}/index.html${query ? `?${query}` : ""}`;
        } else if (pathname === "/app/" || pathname === "/app") {
          req.url = `/app/dashboard/index.html${query ? `?${query}` : ""}`;
        }

        next();
      });
    },
  };
}

export default defineConfig({
  root: "web",
  base: "/app/",
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
