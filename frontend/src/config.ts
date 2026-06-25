// ─── Site ────────────────────────────────────────────────────────────────────

export const site = {
  name: "allieg.dev",
  handle: "@allieg",
  fullName: "Allie Goh",
  pronouns: "she/her",
  tagline:
    "welcome to my digital garden — a place for notes, experiments, and thoughts that are still growing.",
  footer: "allieg.dev © 2026 — made with too much coffee",
};

// ─── Home — "Now" section ────────────────────────────────────────────────────

export const nowItems = [
  "Building this digital garden",
  "Running local LLM experiments",
  "Reading: A Memory Called Empire",
  "Location: somewhere with good coffee",
];

// ─── About ───────────────────────────────────────────────────────────────────

export const aboutBio = [
  "Hi — I'm Allie. I build things with language models, tinker with systems, and occasionally write about what I find.",
  "This site is my digital garden — a place for notes that aren't finished, experiments that might be wrong, and thoughts that are still growing. Nothing here is meant to be authoritative. It's just what I'm thinking about.",
];

export const aboutWhatIDo = [
  "My work sits at the intersection of NLP, agentic systems, and tooling. I'm interested in how language models can be integrated into real workflows — not demos, but things that actually change how people work.",
  "I also make games occasionally, contribute to open-source when I find gaps worth filling, and spend too much time thinking about RAG pipeline architecture.",
];

export const aboutThisSite = [
  "Built with React and a Go backend. Content managed via a lightweight admin panel.",
  "If something's broken, I probably know. If you want to say hi — the links are below.",
];

export const aboutDetails: { key: string; value: string; highlight?: boolean }[] = [
  { key: "Location", value: "somewhere with good coffee" },
  { key: "Focus",    value: "NLP, LLMs, agentic systems, Python" },
  { key: "Status",   value: "online", highlight: true },
];

export const aboutSkills: { name: string; pct: number; color: string }[] = [
  { name: "NLP / LLMs",      pct: 92, color: "#1a7060" },
  { name: "Python",           pct: 88, color: "#2e58a0" },
  { name: "Agentic Systems",  pct: 80, color: "#6a3da0" },
  { name: "Web / HTML / CSS", pct: 74, color: "#a06020" },
  { name: "Game Dev",         pct: 55, color: "#a82d5e" },
];

export const aboutLinks: { label: string; bg: string; color: string; href: string }[] = [
  { label: "GitHub",     bg: "#edf5f3", color: "#1a7060", href: "#" },
  { label: "HuggingFace",bg: "#f5f0e8", color: "#a06020", href: "#" },
  { label: "itch.io",    bg: "#f5eaef", color: "#a82d5e", href: "#" },
  { label: "Mastodon",   bg: "#ecebf5", color: "#4038a0", href: "#" },
  { label: "Email",      bg: "#eaeff8", color: "#2e58a0", href: "mailto:allie.sudo@proton.me" },
  { label: "Resume / CV",bg: "#f0eaf5", color: "#6a3da0", href: "#" },
];

// ─── Tag colour palette ───────────────────────────────────────────────────────
// Each tag string is hashed to a stable index into this array.
// Add, remove, or reorder entries to adjust which colours appear.

export const tagPalette: { bg: string; text: string }[] = [
  { bg: "#edf5f3", text: "#1a7060" }, // teal
  { bg: "#f0eaf5", text: "#6a3da0" }, // purple
  { bg: "#eaeff8", text: "#2e58a0" }, // blue
  { bg: "#f5eaef", text: "#a82d5e" }, // rose
  { bg: "#eaf3ed", text: "#257045" }, // green
  { bg: "#f5f0e8", text: "#a06020" }, // amber
  { bg: "#f5ebe8", text: "#b04030" }, // coral
  { bg: "#ecebf5", text: "#4038a0" }, // indigo
];
