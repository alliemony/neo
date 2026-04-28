import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { HelmetProvider } from "react-helmet-async";
import "@fontsource/jetbrains-mono/400.css";
import "@fontsource/jetbrains-mono/700.css";
import "@fontsource/inter/400.css";
import "@fontsource/inter/600.css";
import "highlight.js/styles/github.css";
import "./index.css";
import App from "./App";
import { ThemeProvider } from "./themes/ThemeProvider";

async function enableApiMocks() {
  if (import.meta.env.DEV && import.meta.env.VITE_ENABLE_API_MOCKS === "true") {
    const { worker } = await import("./mocks/browser");
    await worker.start({ onUnhandledRequest: "bypass" });
  }
}

void enableApiMocks().then(() => {
  createRoot(document.getElementById("root")!).render(
    <StrictMode>
      <HelmetProvider>
        <ThemeProvider>
          <App />
        </ThemeProvider>
      </HelmetProvider>
    </StrictMode>,
  );
});
