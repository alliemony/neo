import { Link, useLocation } from "react-router-dom";
import { site } from "../../config";

const NAV_LINKS = [
  { to: "/", label: "home" },
  { to: "/blog", label: "blog" },
  { to: "/widgets", label: "widgets" },
  { to: "/about", label: "about" },
  { to: "/recs", label: "recs" },
  { to: "/music", label: "music" },
];

export function Header() {
  const location = useLocation();

  function isActive(to: string) {
    if (to === "/") return location.pathname === "/";
    return location.pathname.startsWith(to);
  }

  return (
    <nav className="max-w-[700px] mx-auto px-6 py-5 flex items-center gap-4 flex-wrap">
      <Link
        to="/"
        className="font-bold text-lg text-accent no-underline mr-auto tracking-tight hover:text-accent-alt"
        style={{ fontFamily: "inherit" }}
      >
        {site.name}
      </Link>
      <ul className="flex list-none gap-0.5 flex-wrap p-0 m-0">
        {NAV_LINKS.map(({ to, label }) => (
          <li key={to} className="flex">
            <Link
              to={to}
              className={`text-[15px] no-underline px-2.5 py-1 rounded-md transition-all duration-150 ${
                isActive(to)
                  ? "text-accent bg-[oklch(0.95_0.03_25)]"
                  : "text-text-secondary hover:text-text-primary hover:bg-surface"
              }`}
            >
              {label}
            </Link>
          </li>
        ))}
      </ul>
    </nav>
  );
}
