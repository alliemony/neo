import { Link } from "react-router-dom";
import { usePages } from "../../hooks/usePages";

export function Header() {
  const { pages } = usePages();

  return (
    <header className="border-b-2 border-border">
      <nav className="max-w-6xl mx-auto px-4 py-4 flex justify-between items-center">
        <Link
          to="/"
          className="font-heading text-2xl font-bold text-text-primary no-underline"
        >
          neo
        </Link>
        <div className="flex gap-6 font-heading text-sm">
          <Link
            to="/"
            className="text-text-secondary hover:text-accent no-underline"
          >
            blog
          </Link>
          {pages.map((page) => (
            <Link
              key={page.slug}
              to={`/page/${page.slug}`}
              className="text-text-secondary hover:text-accent no-underline"
            >
              {page.title.toLowerCase()}
            </Link>
          ))}
        </div>
      </nav>
    </header>
  );
}
