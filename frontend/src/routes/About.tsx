import { useEffect } from "react";
import { Layout } from "../components/layout/Layout";
import { SEO } from "../components/SEO";
import {
  site,
  aboutBio,
  aboutWhatIDo,
  aboutThisSite,
  aboutDetails,
  aboutSkills,
  aboutLinks,
} from "../config";

export function About() {
  useEffect(() => {
    const fills = document.querySelectorAll<HTMLElement>(".skill-fill");
    const timer = setTimeout(() => {
      fills.forEach((el) => {
        el.style.width = el.dataset.width || "0";
      });
    }, 300);
    return () => clearTimeout(timer);
  }, []);

  return (
    <Layout>
      <SEO title="About" path="/about" />

      <div className="mb-8">
        <h1 className="text-[30px] font-bold text-text-primary tracking-tight mb-1">About</h1>
        <p className="text-text-secondary">{site.handle} · {site.fullName} · {site.pronouns}</p>
      </div>

      <div className="text-[17px] leading-[1.75] text-text-secondary">
        {aboutBio.map((p, i) => <p key={i} className="mb-[18px]">{p}</p>)}

        <h2 className="mt-9 mb-3 text-[13px] font-semibold uppercase tracking-[0.06em] text-text-secondary">
          What I do
        </h2>
        {aboutWhatIDo.map((p, i) => <p key={i} className="mb-[18px]">{p}</p>)}

        <h2 className="mt-9 mb-3 text-[13px] font-semibold uppercase tracking-[0.06em] text-text-secondary">
          This site
        </h2>
        {aboutThisSite.map((p, i) =>
          i === 0 ? (
            <p key={i} className="mb-[18px]">
              {p.replace("admin panel", "")}<a href="/admin" className="text-accent">admin panel</a>.
            </p>
          ) : (
            <p key={i} className="mb-[18px]">{p}</p>
          )
        )}
      </div>

      <hr className="border-t border-border my-9" />

      {/* Details */}
      <div>
        <h2 className="font-semibold text-[13px] uppercase tracking-widest text-text-secondary mb-4">
          Details
        </h2>
        <div className="flex flex-col gap-2">
          {aboutDetails.map(({ key, value, highlight }) => (
            <div key={key} className="flex gap-4 text-[16px]">
              <span className="font-semibold text-[14px] text-text-secondary shrink-0 w-[90px] uppercase tracking-[0.03em]">
                {key}
              </span>
              <span className={highlight ? "text-[#287848] font-medium" : "text-text-secondary"}>
                {value}
              </span>
            </div>
          ))}
        </div>
      </div>

      <hr className="border-t border-border my-9" />

      {/* Skills */}
      <div>
        <h2 className="font-semibold text-[13px] uppercase tracking-widest text-text-secondary mb-4">
          Skills
        </h2>
        <div className="flex flex-col gap-2.5">
          {aboutSkills.map(({ name, pct, color }) => (
            <div key={name} className="flex items-center gap-3 text-[15px] text-text-secondary">
              <span className="w-[140px] shrink-0 font-medium">{name}</span>
              <div className="flex-1 h-1.5 bg-surface rounded-full overflow-hidden">
                <div
                  className="skill-fill h-full rounded-full transition-[width] duration-1000 ease-out"
                  style={{ width: "0", background: color }}
                  data-width={`${pct}%`}
                />
              </div>
            </div>
          ))}
        </div>
      </div>

      <hr className="border-t border-border my-9" />

      {/* Links */}
      <div>
        <h2 className="font-semibold text-[13px] uppercase tracking-widest text-text-secondary mb-4">
          Links
        </h2>
        <div className="flex flex-wrap gap-2">
          {aboutLinks.map(({ label, bg, color, href }) => (
            <a
              key={label}
              href={href}
              className="inline-flex items-center text-[15px] font-medium px-3.5 py-1.5 rounded-full no-underline transition-all hover:brightness-95"
              style={{ background: bg, color }}
            >
              {label}
            </a>
          ))}
        </div>
      </div>
    </Layout>
  );
}
