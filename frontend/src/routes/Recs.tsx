import { Layout } from "../components/layout/Layout";
import { SEO } from "../components/SEO";

interface RecItem {
  name: string;
  href?: string;
  tag: string;
  tagBg: string;
  tagColor: string;
  desc: string;
}

interface RecSection {
  title: string;
  items: RecItem[];
}

const RECS: RecSection[] = [
  {
    title: "Tools",
    items: [
      {
        name: "Ollama",
        href: "https://ollama.ai",
        tag: "CLI",
        tagBg: "#edf5f3",
        tagColor: "#1a7060",
        desc: "Run LLMs locally with one command. Changed my daily workflow completely. Fast, private, free.",
      },
      {
        name: "Zed",
        href: "https://zed.dev",
        tag: "Editor",
        tagBg: "#eaeff8",
        tagColor: "#2e58a0",
        desc: "Fast, minimal editor. Feels like what a text editor should be in 2026. Collaborative built-in.",
      },
      {
        name: "LangGraph",
        href: "https://langchain-ai.github.io/langgraph/",
        tag: "ML",
        tagBg: "#f0eaf5",
        tagColor: "#6a3da0",
        desc: "The least bad way to build agentic workflows. Explicit state machines beat implicit magic chains.",
      },
      {
        name: "Kagi",
        href: "https://kagi.com",
        tag: "Search",
        tagBg: "#f5f0e8",
        tagColor: "#a06020",
        desc: "Paid search that's actually good. No ads, good results, bangs still work. Worth every cent.",
      },
      {
        name: "Fly.io",
        href: "https://fly.io",
        tag: "Infra",
        tagBg: "#eaf3ed",
        tagColor: "#257045",
        desc: "Deploy containers globally, sensible pricing, great DX. My go-to for anything that needs to run 24/7.",
      },
      {
        name: "Obsidian",
        href: "https://obsidian.md",
        tag: "Notes",
        tagBg: "#f0eaf5",
        tagColor: "#6a3da0",
        desc: "Local markdown notes with a graph view. I don't use the graph view. The local part is the point.",
      },
    ],
  },
  {
    title: "Reading",
    items: [
      {
        name: "Attention Is All You Need",
        href: "https://arxiv.org/abs/1706.03762",
        tag: "Paper",
        tagBg: "#eaeff8",
        tagColor: "#2e58a0",
        desc: "Obviously. Read it if you haven't. Then read it again a year later — you'll get more the second time.",
      },
      {
        name: "A Memory Called Empire",
        tag: "Book",
        tagBg: "#f5eaef",
        tagColor: "#a82d5e",
        desc: "Arkady Martine. A diplomat navigating imperial politics with a dead ambassador in her head. Genuinely brilliant.",
      },
      {
        name: "Lilian Weng's Blog",
        href: "https://lilianweng.github.io",
        tag: "Blog",
        tagBg: "#edf5f3",
        tagColor: "#1a7060",
        desc: "The best ML explainers on the internet. Dense, accurate, deep. Bookmark it and read slowly.",
      },
      {
        name: "RAG Survey (Gao et al.)",
        href: "https://arxiv.org/abs/2312.10997",
        tag: "Paper",
        tagBg: "#eaeff8",
        tagColor: "#2e58a0",
        desc: "Comprehensive survey of retrieval-augmented generation approaches. Good map of the space before diving in.",
      },
    ],
  },
  {
    title: "Misc",
    items: [
      {
        name: "Berkeley Mono",
        href: "https://berkeleygraphics.com/typefaces/berkeley-mono/",
        tag: "Font",
        tagBg: "#f5f0e8",
        tagColor: "#a06020",
        desc: "Best monospace font I've used. Worth paying for if you stare at code all day.",
      },
      {
        name: "Practical AI",
        href: "https://changelog.com/practicalai",
        tag: "Podcast",
        tagBg: "#eaf3ed",
        tagColor: "#257045",
        desc: "Stays grounded. Doesn't hype. Actual practitioners talking about actual implementations. Rare.",
      },
      {
        name: "Recurse Center",
        href: "https://recurse.com",
        tag: "Community",
        tagBg: "#ecebf5",
        tagColor: "#4038a0",
        desc: "Three months of self-directed programming retreat. The best learning environment I've been in.",
      },
    ],
  },
];

export function Recs() {
  return (
    <Layout>
      <SEO title="Recommendations" path="/recs" />

      <div className="mb-9">
        <h1 className="text-[30px] font-bold text-text-primary tracking-tight mb-1.5">
          Recommendations
        </h1>
        <p className="text-text-secondary text-[18px]">
          I spy with my crystal eyes. A cave of little treasures of the
          inter-web.
        </p>
      </div>

      {RECS.map((section) => (
        <div key={section.title} className="mb-12">
          <h2 className="text-[22px] font-semibold text-text-primary mb-4">
            {section.title}
          </h2>
          <div>
            {section.items.map((item) => (
              <div
                key={item.name}
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
                    style={{ background: item.tagBg, color: item.tagColor }}
                  >
                    {item.tag}
                  </span>
                </div>
                <p className="text-[16px] text-text-secondary leading-relaxed">
                  {item.desc}
                </p>
              </div>
            ))}
          </div>
        </div>
      ))}
    </Layout>
  );
}
