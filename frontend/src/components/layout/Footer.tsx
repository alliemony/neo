export function Footer() {
  const year = new Date().getFullYear();

  return (
    <footer className="border-t-2 border-border mt-12">
      <div className="max-w-6xl mx-auto px-4 py-6 text-center text-text-secondary text-sm font-body">
        <p>neo &middot; built with care &middot; {year}</p>
      </div>
    </footer>
  );
}
