import { useState, useEffect } from 'react';
import { Layout } from "../components/layout/Layout";
import { SEO } from "../components/SEO";
import { getRecs } from "../services/api";
import type { RecSection } from "../types/post";

export function Recs() {
  const [sections, setSections] = useState<RecSection[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    getRecs()
      .then((data) => setSections(data.sections ?? []))
      .catch(() => setError(true))
      .finally(() => setLoading(false));
  }, []);

  return (
    <Layout>
      <SEO title="Recommendations" path="/recs" />

      <div className="mb-9">
        <h1 className="text-[30px] font-bold text-text-primary tracking-tight mb-1.5">
          Recommendations
        </h1>
        <p className="text-text-secondary text-[18px]">
          I spy with my crystal eyes. A cave of little treasures of the inter-web.
        </p>
      </div>

      {loading && (
        <p className="text-text-secondary">Loading…</p>
      )}

      {error && (
        <p className="text-text-secondary">Could not load recommendations. Try again later.</p>
      )}

      {!loading && !error && sections.length === 0 && (
        <p className="text-text-secondary">No recommendations yet.</p>
      )}

      {sections.map((section) => (
        <div key={section.title} className="mb-12">
          <h2 className="text-[22px] font-semibold text-text-primary mb-4">
            {section.title}
          </h2>
          <div>
            {section.items.map((item) => (
              <div
                key={item.id}
                className="py-3.5 border-t border-border first:border-t-0"
              >
                <div className="flex items-center gap-2.5 mb-1 flex-wrap">
                  <span className="font-semibold text-[17px] text-text-primary">
                    {item.href ? (
                      <a
                        href={item.href}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-text-primary no-underline hover:text-accent transition-colors"
                      >
                        {item.name}
                      </a>
                    ) : (
                      item.name
                    )}
                  </span>
                  <span
                    className="text-[11px] font-semibold px-2 py-0.5 rounded-full uppercase tracking-[0.03em]"
                    style={{ background: item.tag_bg, color: item.tag_color }}
                  >
                    {item.tag}
                  </span>
                </div>
                <p className="text-[16px] text-text-secondary leading-relaxed">
                  {item.description}
                </p>
              </div>
            ))}
          </div>
        </div>
      ))}
    </Layout>
  );
}
