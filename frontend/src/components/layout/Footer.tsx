import { site } from "../../config";

export function Footer() {
  return (
    <footer className="max-w-[700px] mx-auto px-6 py-8 text-center text-text-secondary text-sm border-t border-border">
      {site.footer}
    </footer>
  );
}
