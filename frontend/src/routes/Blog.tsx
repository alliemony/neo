import { useState, useEffect, useMemo } from "react";
import { useSearchParams } from "react-router-dom";
import { Layout } from "../components/layout/Layout";
import { PostCard } from "../components/blog/PostCard";
import { TagPill } from "../components/blog/TagPill";
import { SEO } from "../components/SEO";
import { getPosts } from "../services/api";
import type { Post } from "../types/post";

export function Blog() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [allPosts, setAllPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [search, setSearch] = useState("");
  const [sort, setSort] = useState<"newest" | "oldest">("newest");

  const activeTag = searchParams.get("tag") || "all";

  function setActiveTag(tag: string) {
    if (tag === "all") {
      setSearchParams({});
    } else {
      setSearchParams({ tag });
    }
  }

  useEffect(() => {
    setLoading(true);
    getPosts({ per_page: 100 })
      .then((data) => { setAllPosts(data.posts); setLoading(false); })
      .catch(() => { setError(true); setLoading(false); });
  }, []);

  const tagCounts = useMemo(() => {
    const counts: Record<string, number> = {};
    allPosts.forEach((p) => p.tags.forEach((t) => { counts[t] = (counts[t] || 0) + 1; }));
    return Object.entries(counts).sort((a, b) => b[1] - a[1]);
  }, [allPosts]);

  const filtered = useMemo(() => {
    let posts = allPosts;
    if (activeTag !== "all") posts = posts.filter((p) => p.tags.includes(activeTag));
    if (search.trim()) {
      const q = search.toLowerCase();
      posts = posts.filter(
        (p) =>
          p.title.toLowerCase().includes(q) ||
          p.content.toLowerCase().includes(q) ||
          p.tags.some((t) => t.includes(q))
      );
    }
    return [...posts].sort((a, b) => {
      const da = new Date(a.created_at).getTime();
      const db = new Date(b.created_at).getTime();
      return sort === "newest" ? db - da : da - db;
    });
  }, [allPosts, activeTag, search, sort]);

  return (
    <Layout>
      <SEO title="Blog" path="/blog" />

      <div className="mb-8">
        <h1 className="text-[30px] font-bold text-text-primary tracking-tight mb-1.5">Blog</h1>
        <p className="text-text-secondary text-[18px]">notes, experiments, and intermittent thoughts</p>
      </div>

      {/* Filters */}
      <div className="flex flex-col gap-3.5 mb-7">
        <div className="relative">
          <span className="absolute left-3 top-1/2 -translate-y-1/2 text-text-secondary text-base pointer-events-none">⌕</span>
          <input
            type="text"
            placeholder="search posts…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full pl-9 pr-3.5 py-2.5 bg-bg border border-border rounded-md text-[16px] text-text-primary placeholder:text-text-secondary outline-none focus:border-accent transition-colors"
          />
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          <div className="flex gap-1.5 flex-wrap flex-1">
            <button
              onClick={() => setActiveTag("all")}
              className={`inline-flex items-center text-[13px] font-medium px-2.5 py-0.5 rounded-full transition-all cursor-pointer border-0 ${
                activeTag === "all"
                  ? "bg-accent text-white"
                  : "bg-surface text-text-secondary hover:brightness-95"
              }`}
            >
              all ({allPosts.length})
            </button>
            {tagCounts.map(([tag, count]) => (
              <TagPill
                key={tag}
                tag={`${tag} (${count})`}
                active={activeTag === tag}
                onClick={() => setActiveTag(activeTag === tag ? "all" : tag)}
                asSpan
              />
            ))}
          </div>
          <div className="ml-auto">
            <select
              value={sort}
              onChange={(e) => setSort(e.target.value as "newest" | "oldest")}
              className="text-[14px] py-1 px-3 border border-border rounded-md bg-bg text-text-secondary outline-none focus:border-accent appearance-none cursor-pointer"
            >
              <option value="newest">newest first</option>
              <option value="oldest">oldest first</option>
            </select>
          </div>
        </div>
      </div>

      {loading && <p className="text-text-secondary">Loading…</p>}
      {error && <p className="text-accent">Failed to load posts.</p>}
      {!loading && !error && (
        <>
          <p className="text-[14px] text-text-secondary font-medium mb-2">
            {filtered.length} post{filtered.length !== 1 ? "s" : ""}
          </p>
          {filtered.length === 0 ? (
            <div className="text-center py-12 text-text-secondary border border-dashed border-border rounded-lg">
              <p className="font-medium">No posts found</p>
              <p className="text-sm mt-1">try a different tag or search term</p>
            </div>
          ) : (
            <div>
              {filtered.map((post) => (
                <PostCard key={post.id} post={post} />
              ))}
            </div>
          )}
        </>
      )}
    </Layout>
  );
}
