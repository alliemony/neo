import { Layout } from "../components/layout/Layout";
import { SEO } from "../components/SEO";
import { WidgetEmbed } from "../components/blog/WidgetEmbed";

interface WidgetDef {
  id: string;
  title: string;
  description: string;
  type: string;
  typeBg: string;
  typeColor: string;
  featured?: boolean;
}

const WIDGETS: WidgetDef[] = [
  {
    id: "sentiment",
    title: "Sentiment Classifier",
    description: "Fine-tuned DistilBERT for sentiment analysis. Part of my ongoing NLP experiments series.",
    type: "Gradio",
    typeBg: "#f5ebe8",
    typeColor: "#b04030",
    featured: true,
  },
  {
    id: "gpt2",
    title: "GPT-2 Text Completion",
    description: "Classic GPT-2 running via HuggingFace Inference. Type a prompt, see what it continues.",
    type: "HuggingFace",
    typeBg: "#f5f0e8",
    typeColor: "#a06020",
    featured: true,
  },
];

export function Widgets() {
  return (
    <Layout>
      <SEO title="Widgets" path="/widgets" />

      <div className="mb-8">
        <h1 className="text-[30px] font-bold text-text-primary tracking-tight mb-1.5">Widgets</h1>
        <p className="text-text-secondary text-[18px]">experiments, demos, and interactive prototypes</p>
      </div>

      <div className="flex flex-col gap-8">
        {WIDGETS.map((w) => (
          <div key={w.id} className="border-t border-border pt-8 first:border-t-0 first:pt-0">
            <div className="flex items-center gap-2 mb-2">
              <span
                className="text-[11px] font-semibold px-2 py-0.5 rounded-full uppercase tracking-[0.03em]"
                style={{ background: w.typeBg, color: w.typeColor }}
              >
                {w.type}
              </span>
              {w.featured && <span className="text-[#a06020] text-[13px]">★</span>}
            </div>
            <h2 className="text-[20px] font-semibold text-text-primary mb-1.5">{w.title}</h2>
            <p className="text-[16px] text-text-secondary leading-relaxed mb-4">{w.description}</p>
            <div className="rounded-lg overflow-hidden border border-border">
              <WidgetEmbed widgetId={w.id} />
            </div>
          </div>
        ))}

        {WIDGETS.length === 0 && (
          <div className="text-center py-16 text-text-secondary border border-dashed border-border rounded-lg">
            <p className="font-medium">No widgets yet</p>
            <p className="text-sm mt-1">check back soon</p>
          </div>
        )}
      </div>
    </Layout>
  );
}
