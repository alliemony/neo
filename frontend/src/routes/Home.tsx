import { Link } from "react-router-dom";
import { Layout } from "../components/layout/Layout";
import { PostCard } from "../components/blog/PostCard";
import { TagPill } from "../components/blog/TagPill";
import { SEO } from "../components/SEO";
import { usePosts } from "../hooks/usePosts";
import { getTags } from "../services/api";
import { useState, useEffect } from "react";
import type { TagCount } from "../types/post";
import { site, nowItems } from "../config";

export function Home() {
  const { posts, loading, error } = usePosts({ per_page: 4 });
  const [tags, setTags] = useState<TagCount[]>([]);

  useEffect(() => {
    getTags().then(setTags).catch(() => setTags([]));
  }, []);

  return (
    <Layout>
      <SEO path="/" />

      {/* Hero */}
      <div className="pt-10 pb-0">
        <h1 className="text-[34px] font-bold text-text-primary leading-tight tracking-tight mb-1">
          {site.fullName}
        </h1>
        <div className="text-[16px] text-accent font-medium mb-3.5">{site.handle}</div>
        <p className="text-[19px] text-text-secondary max-w-[520px] leading-relaxed">
          {site.tagline}
        </p>
      </div>

      <hr className="border-t border-border my-9" />

      {/* Recent Posts */}
      <div>
        <h2 className="font-semibold text-[13px] uppercase tracking-widest text-text-secondary mb-4">
          Recent posts
        </h2>
        {loading && <p className="text-text-secondary">Loading…</p>}
        {error && <p className="text-accent">Failed to load posts.</p>}
        {!loading && !error && posts.map((post) => (
          <PostCard key={post.id} post={post} />
        ))}
        <div className="mt-5">
          <Link to="/blog" className="text-accent hover:text-accent-alt no-underline text-[15px]">
            all posts →
          </Link>
        </div>
      </div>

      <hr className="border-t border-border my-9" />

      {/* Now */}
      <div>
        <h2 className="font-semibold text-[13px] uppercase tracking-widest text-text-secondary mb-4">
          Now
        </h2>
        <div className="flex flex-col gap-2.5">
          {nowItems.map((item) => (
            <div key={item} className="flex gap-2.5 text-text-secondary leading-relaxed">
              <span className="text-accent font-semibold shrink-0">→</span>
              <span>{item}</span>
            </div>
          ))}
        </div>
      </div>

      <hr className="border-t border-border my-9" />

      {/* Tags */}
      <div>
        <h2 className="font-semibold text-[13px] uppercase tracking-widest text-text-secondary mb-4">
          Tags
        </h2>
        <div className="flex flex-wrap gap-1.5">
          {tags.map((t) => (
            <TagPill key={t.name} tag={t.name} />
          ))}
        </div>
      </div>
    </Layout>
  );
}
