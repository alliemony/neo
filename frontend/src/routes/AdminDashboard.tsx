import { useState, useEffect, useCallback } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Layout } from "../components/layout/Layout";
import { useAuth } from "../hooks/useAuth";
import {
  adminGetPosts,
  adminGetPages,
  adminDeletePost,
  adminDeletePage,
  adminGetRecs,
  adminCreateRec,
  adminUpdateRec,
  adminDeleteRec,
  adminGetMusic,
  adminCreateMusic,
  adminUpdateMusic,
  adminDeleteMusic,
} from "../services/api";
import type {
  Post,
  Page,
  Rec,
  CreateRecInput,
  UpdateRecInput,
  MusicRec,
  CreateMusicRecInput,
  UpdateMusicRecInput,
} from "../types/post";

const EMPTY_REC_FORM: CreateRecInput = {
  name: "", section: "", tag: "", tag_bg: "#eaeff8", tag_color: "#2e58a0",
  description: "", href: "", published: true, sort_order: 0,
};

const EMPTY_MUSIC_FORM: CreateMusicRecInput = {
  title: "", artist: "", album: "", cover_url: "", spotify_url: "", apple_url: "",
  note: "", published: true, sort_order: 0,
};

export function AdminDashboard() {
  const { credentials, user, authMode, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  const [posts, setPosts] = useState<Post[]>([]);
  const [pages, setPages] = useState<Page[]>([]);
  const [recs, setRecs] = useState<Rec[]>([]);
  const [musicRecs, setMusicRecs] = useState<MusicRec[]>([]);
  const [loading, setLoading] = useState(true);

  // Recs form state
  const [recForm, setRecForm] = useState<CreateRecInput>(EMPTY_REC_FORM);
  const [editingRecId, setEditingRecId] = useState<number | null>(null);
  const [recFormOpen, setRecFormOpen] = useState(false);
  const [recFormError, setRecFormError] = useState("");

  // Music form state
  const [musicForm, setMusicForm] = useState<CreateMusicRecInput>(EMPTY_MUSIC_FORM);
  const [editingMusicId, setEditingMusicId] = useState<number | null>(null);
  const [musicFormOpen, setMusicFormOpen] = useState(false);
  const [musicFormError, setMusicFormError] = useState("");

  // For OAuth mode, use null credentials (cookies handle it); for basic, use credentials.
  const authCreds = authMode === "oauth" ? null : credentials;

  const loadData = useCallback(async () => {
    try {
      const [postData, pageData, recsData, musicData] = await Promise.all([
        adminGetPosts(authCreds),
        adminGetPages(authCreds),
        adminGetRecs(authCreds),
        adminGetMusic(authCreds),
      ]);
      setPosts(postData.posts || []);
      setPages(pageData || []);
      setRecs(recsData.recs || []);
      setMusicRecs(musicData.recs || []);
    } catch {
      // unauthorized or network error — list stays empty
    }
    setLoading(false);
  }, [authCreds]);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/admin/login");
      return;
    }
    loadData();
  }, [isAuthenticated, loadData, navigate]);

  const handleDeletePost = async (slug: string) => {
    if (!confirm(`Delete post "${slug}"?`)) return;
    try {
      await adminDeletePost(authCreds, slug);
      setPosts((prev) => prev.filter((p) => p.slug !== slug));
    } catch {
      alert("Failed to delete post.");
    }
  };

  const handleDeletePage = async (slug: string) => {
    if (!confirm(`Delete page "${slug}"?`)) return;
    try {
      await adminDeletePage(authCreds, slug);
      setPages((prev) => prev.filter((p) => p.slug !== slug));
    } catch {
      alert("Failed to delete page.");
    }
  };

  // --- Recs ---
  const existingSections = [...new Set(recs.map((r) => r.section))];

  function openNewRecForm() {
    setEditingRecId(null);
    setRecForm(EMPTY_REC_FORM);
    setRecFormError("");
    setRecFormOpen(true);
  }

  function openEditRecForm(rec: Rec) {
    setEditingRecId(rec.id);
    setRecForm({
      name: rec.name,
      section: rec.section,
      tag: rec.tag,
      tag_bg: rec.tag_bg,
      tag_color: rec.tag_color,
      description: rec.description,
      href: rec.href ?? "",
      published: rec.published,
      sort_order: rec.sort_order,
    });
    setRecFormError("");
    setRecFormOpen(true);
  }

  async function handleRecSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!recForm.name || !recForm.section || !recForm.tag || !recForm.tag_bg || !recForm.tag_color || !recForm.description) {
      setRecFormError("Name, section, tag, tag colours, and description are required.");
      return;
    }
    const payload = { ...recForm, href: recForm.href || null };
    try {
      if (editingRecId !== null) {
        await adminUpdateRec(authCreds, editingRecId, payload as UpdateRecInput);
      } else {
        await adminCreateRec(authCreds, payload);
      }
      setRecFormOpen(false);
      loadData();
    } catch {
      setRecFormError("Failed to save. Please try again.");
    }
  }

  async function handleRecPublishToggle(rec: Rec) {
    try {
      await adminUpdateRec(authCreds, rec.id, { published: !rec.published });
      loadData();
    } catch {
      // ignore
    }
  }

  async function handleDeleteRec(rec: Rec) {
    if (!confirm(`Delete rec "${rec.name}"?`)) return;
    try {
      await adminDeleteRec(authCreds, rec.id);
      setRecs((prev) => prev.filter((r) => r.id !== rec.id));
    } catch {
      alert("Failed to delete rec.");
    }
  }

  // --- Music ---
  function openNewMusicForm() {
    setEditingMusicId(null);
    setMusicForm(EMPTY_MUSIC_FORM);
    setMusicFormError("");
    setMusicFormOpen(true);
  }

  function openEditMusicForm(rec: MusicRec) {
    setEditingMusicId(rec.id);
    setMusicForm({
      title: rec.title,
      artist: rec.artist,
      album: rec.album ?? "",
      cover_url: rec.cover_url ?? "",
      spotify_url: rec.spotify_url ?? "",
      apple_url: rec.apple_url ?? "",
      note: rec.note ?? "",
      published: rec.published,
      sort_order: rec.sort_order,
    });
    setMusicFormError("");
    setMusicFormOpen(true);
  }

  async function handleMusicSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!musicForm.title || !musicForm.artist) {
      setMusicFormError("Title and artist are required.");
      return;
    }
    const nullify = (v: string | null | undefined) => (v === "" ? null : v);
    const payload = {
      ...musicForm,
      album: nullify(musicForm.album as string),
      cover_url: nullify(musicForm.cover_url as string),
      spotify_url: nullify(musicForm.spotify_url as string),
      apple_url: nullify(musicForm.apple_url as string),
      note: nullify(musicForm.note as string),
    };
    try {
      if (editingMusicId !== null) {
        await adminUpdateMusic(authCreds, editingMusicId, payload as UpdateMusicRecInput);
      } else {
        await adminCreateMusic(authCreds, payload);
      }
      setMusicFormOpen(false);
      loadData();
    } catch {
      setMusicFormError("Failed to save. Please try again.");
    }
  }

  async function handleMusicPublishToggle(rec: MusicRec) {
    try {
      await adminUpdateMusic(authCreds, rec.id, { published: !rec.published });
      loadData();
    } catch {
      // ignore
    }
  }

  async function handleDeleteMusic(rec: MusicRec) {
    if (!confirm(`Delete music rec "${rec.title}"?`)) return;
    try {
      await adminDeleteMusic(authCreds, rec.id);
      setMusicRecs((prev) => prev.filter((m) => m.id !== rec.id));
    } catch {
      alert("Failed to delete music rec.");
    }
  }

  if (!isAuthenticated) return null;

  const displayName = authMode === "oauth" && user ? user.username : "Admin";
  const avatarUrl = authMode === "oauth" && user ? user.avatar_url : null;

  const inputClass = "w-full border-2 border-border bg-bg px-2 py-1 text-sm";
  const createBtnClass =
    "border-2 border-border bg-accent text-white! px-3 py-1 text-sm font-heading hover:opacity-90 no-underline";

  return (
    <Layout>
      <div className="flex justify-between items-center mb-8">
        <h1 className="font-heading text-2xl font-bold">Admin Dashboard</h1>
        <div className="flex items-center gap-3">
          {avatarUrl && (
            <img
              src={avatarUrl}
              alt={displayName}
              className="w-6 h-6 border border-border"
            />
          )}
          <span className="text-sm text-text-secondary font-heading">
            {displayName}
          </span>
          <button
            onClick={() => {
              logout();
              navigate("/admin/login");
            }}
            className="border-2 border-border px-3 py-1 text-sm font-heading hover:border-accent"
          >
            Logout
          </button>
        </div>
      </div>

      {loading ? (
        <p className="text-text-secondary">Loading…</p>
      ) : (
        <>
          {/* Posts Section */}
          <section className="mb-10">
            <div className="flex justify-between items-center mb-4">
              <h2 className="font-heading text-xl font-bold">
                Posts ({posts.length})
              </h2>
              <Link to="/admin/posts/new" className={createBtnClass}>
                + New Post
              </Link>
            </div>

            <div className="space-y-2">
              {posts.map((post) => (
                <div
                  key={post.slug}
                  className="flex items-center justify-between border-2 border-border bg-surface p-3"
                >
                  <div className="flex items-center gap-3">
                    <span
                      className={`inline-block w-2 h-2 rounded-full ${
                        post.published ? "bg-success" : "bg-text-secondary"
                      }`}
                      title={post.published ? "Published" : "Draft"}
                    />
                    <span className="font-heading text-sm">{post.title}</span>
                    {!post.published && (
                      <span className="text-xs text-text-secondary">
                        (draft)
                      </span>
                    )}
                  </div>
                  <div className="flex gap-2">
                    <Link
                      to={`/admin/posts/${post.slug}/edit`}
                      className="border border-border px-2 py-1 text-xs font-heading hover:border-accent no-underline"
                    >
                      Edit
                    </Link>
                    <button
                      onClick={() => handleDeletePost(post.slug)}
                      className="border border-border px-2 py-1 text-xs font-heading hover:border-accent text-accent"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))}
              {posts.length === 0 && (
                <p className="text-text-secondary text-sm">No posts yet.</p>
              )}
            </div>
          </section>

          {/* Pages Section */}
          <section className="mb-10">
            <div className="flex justify-between items-center mb-4">
              <h2 className="font-heading text-xl font-bold">
                Pages ({pages.length})
              </h2>
              <Link to="/admin/pages/new" className={createBtnClass}>
                + New Page
              </Link>
            </div>

            <div className="space-y-2">
              {pages.map((page) => (
                <div
                  key={page.slug}
                  className="flex items-center justify-between border-2 border-border bg-surface p-3"
                >
                  <div className="flex items-center gap-3">
                    <span
                      className={`inline-block w-2 h-2 rounded-full ${
                        page.published ? "bg-success" : "bg-text-secondary"
                      }`}
                      title={page.published ? "Published" : "Draft"}
                    />
                    <span className="font-heading text-sm">{page.title}</span>
                    <span className="text-xs text-text-secondary">
                      #{page.sort_order}
                    </span>
                    {!page.published && (
                      <span className="text-xs text-text-secondary">
                        (draft)
                      </span>
                    )}
                  </div>
                  <div className="flex gap-2">
                    <Link
                      to={`/admin/pages/${page.slug}/edit`}
                      className="border border-border px-2 py-1 text-xs font-heading hover:border-accent no-underline"
                    >
                      Edit
                    </Link>
                    <button
                      onClick={() => handleDeletePage(page.slug)}
                      className="border border-border px-2 py-1 text-xs font-heading hover:border-accent text-accent"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))}
              {pages.length === 0 && (
                <p className="text-text-secondary text-sm">No pages yet.</p>
              )}
            </div>
          </section>

          {/* Recs Section */}
          <section className="mb-10">
            <div className="flex justify-between items-center mb-4">
              <h2 className="font-heading text-xl font-bold">
                Recs ({recs.length})
              </h2>
              <button onClick={openNewRecForm} className={createBtnClass}>
                + New Rec
              </button>
            </div>

            {recFormOpen && (
              <form onSubmit={handleRecSubmit} className="mb-4 border-2 border-border bg-surface p-4 space-y-3">
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-xs font-bold mb-1">Name *</label>
                    <input className={inputClass} value={recForm.name} onChange={(e) => setRecForm({ ...recForm, name: e.target.value })} />
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Section *</label>
                    <input list="sections-list" className={inputClass} value={recForm.section} onChange={(e) => setRecForm({ ...recForm, section: e.target.value })} />
                    <datalist id="sections-list">
                      {existingSections.map((s) => <option key={s} value={s} />)}
                    </datalist>
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Tag *</label>
                    <input className={inputClass} value={recForm.tag} onChange={(e) => setRecForm({ ...recForm, tag: e.target.value })} />
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Link (href)</label>
                    <input type="url" className={inputClass} value={recForm.href ?? ""} onChange={(e) => setRecForm({ ...recForm, href: e.target.value })} />
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Tag BG colour *</label>
                    <div className="flex gap-2 items-center">
                      <input type="color" value={recForm.tag_bg} onChange={(e) => setRecForm({ ...recForm, tag_bg: e.target.value })} className="w-8 h-8 border-0 p-0 cursor-pointer" />
                      <input className={`${inputClass} font-mono`} value={recForm.tag_bg} onChange={(e) => setRecForm({ ...recForm, tag_bg: e.target.value })} />
                    </div>
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Tag text colour *</label>
                    <div className="flex gap-2 items-center">
                      <input type="color" value={recForm.tag_color} onChange={(e) => setRecForm({ ...recForm, tag_color: e.target.value })} className="w-8 h-8 border-0 p-0 cursor-pointer" />
                      <input className={`${inputClass} font-mono`} value={recForm.tag_color} onChange={(e) => setRecForm({ ...recForm, tag_color: e.target.value })} />
                    </div>
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Sort order</label>
                    <input type="number" className={inputClass} value={recForm.sort_order} onChange={(e) => setRecForm({ ...recForm, sort_order: Number(e.target.value) })} />
                  </div>
                  <div className="flex items-center gap-2 pt-5">
                    <input type="checkbox" id="rec-published" checked={recForm.published} onChange={(e) => setRecForm({ ...recForm, published: e.target.checked })} />
                    <label htmlFor="rec-published" className="text-sm font-bold">Published</label>
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-bold mb-1">Description *</label>
                  <textarea className={inputClass} rows={3} value={recForm.description} onChange={(e) => setRecForm({ ...recForm, description: e.target.value })} />
                </div>
                {recFormError && <p className="text-accent text-sm">{recFormError}</p>}
                <div className="flex gap-2 justify-end">
                  <button type="button" onClick={() => setRecFormOpen(false)} className="px-3 py-1 text-sm border-2 border-border">Cancel</button>
                  <button type="submit" className="px-3 py-1 text-sm bg-accent text-white! font-bold border-2 border-border">{editingRecId !== null ? "Save" : "Create"}</button>
                </div>
              </form>
            )}

            <div className="space-y-2">
              {recs.map((rec) => (
                <div key={rec.id} className="flex items-center justify-between border-2 border-border bg-surface p-3 gap-2">
                  <div className="flex items-center gap-3 min-w-0">
                    <span className="font-heading text-sm truncate">{rec.name}</span>
                    <span className="text-xs text-text-secondary shrink-0">{rec.section}</span>
                    <button
                      onClick={() => handleRecPublishToggle(rec)}
                      className={`shrink-0 text-xs px-2 py-0.5 border border-border ${rec.published ? "text-success" : "text-text-secondary"}`}
                      title="Toggle published"
                    >
                      {rec.published ? "Published" : "Draft"}
                    </button>
                  </div>
                  <div className="flex gap-2 shrink-0">
                    <button onClick={() => openEditRecForm(rec)} className="border border-border px-2 py-1 text-xs font-heading hover:border-accent">Edit</button>
                    <button onClick={() => handleDeleteRec(rec)} className="border border-border px-2 py-1 text-xs font-heading hover:border-accent text-accent">Delete</button>
                  </div>
                </div>
              ))}
              {recs.length === 0 && (
                <p className="text-text-secondary text-sm">No recs yet.</p>
              )}
            </div>
          </section>

          {/* Music Section */}
          <section>
            <div className="flex justify-between items-center mb-4">
              <h2 className="font-heading text-xl font-bold">
                Music ({musicRecs.length})
              </h2>
              <button onClick={openNewMusicForm} className={createBtnClass}>
                + New Music Rec
              </button>
            </div>

            {musicFormOpen && (
              <form onSubmit={handleMusicSubmit} className="mb-4 border-2 border-border bg-surface p-4 space-y-3">
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-xs font-bold mb-1">Title *</label>
                    <input className={inputClass} value={musicForm.title} onChange={(e) => setMusicForm({ ...musicForm, title: e.target.value })} />
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Artist *</label>
                    <input className={inputClass} value={musicForm.artist} onChange={(e) => setMusicForm({ ...musicForm, artist: e.target.value })} />
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Album</label>
                    <input className={inputClass} value={musicForm.album ?? ""} onChange={(e) => setMusicForm({ ...musicForm, album: e.target.value })} />
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Cover URL</label>
                    <input type="url" className={inputClass} value={musicForm.cover_url ?? ""} onChange={(e) => setMusicForm({ ...musicForm, cover_url: e.target.value })} />
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Spotify URL</label>
                    <input type="url" className={inputClass} value={musicForm.spotify_url ?? ""} onChange={(e) => setMusicForm({ ...musicForm, spotify_url: e.target.value })} />
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Apple Music URL</label>
                    <input type="url" className={inputClass} value={musicForm.apple_url ?? ""} onChange={(e) => setMusicForm({ ...musicForm, apple_url: e.target.value })} />
                  </div>
                  <div>
                    <label className="block text-xs font-bold mb-1">Sort order</label>
                    <input type="number" className={inputClass} value={musicForm.sort_order} onChange={(e) => setMusicForm({ ...musicForm, sort_order: Number(e.target.value) })} />
                  </div>
                  <div className="flex items-center gap-2 pt-5">
                    <input type="checkbox" id="music-published" checked={musicForm.published} onChange={(e) => setMusicForm({ ...musicForm, published: e.target.checked })} />
                    <label htmlFor="music-published" className="text-sm font-bold">Published</label>
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-bold mb-1">Personal note</label>
                  <textarea className={inputClass} rows={2} value={musicForm.note ?? ""} onChange={(e) => setMusicForm({ ...musicForm, note: e.target.value })} />
                </div>
                {musicFormError && <p className="text-accent text-sm">{musicFormError}</p>}
                <div className="flex gap-2 justify-end">
                  <button type="button" onClick={() => setMusicFormOpen(false)} className="px-3 py-1 text-sm border-2 border-border">Cancel</button>
                  <button type="submit" className="px-3 py-1 text-sm bg-accent text-white! font-bold border-2 border-border">{editingMusicId !== null ? "Save" : "Create"}</button>
                </div>
              </form>
            )}

            <div className="space-y-2">
              {musicRecs.map((rec) => (
                <div key={rec.id} className="flex items-center justify-between border-2 border-border bg-surface p-3 gap-2">
                  <div className="flex items-center gap-3 min-w-0">
                    <span className="font-heading text-sm truncate">{rec.title}</span>
                    <span className="text-xs text-text-secondary shrink-0">{rec.artist}</span>
                    <button
                      onClick={() => handleMusicPublishToggle(rec)}
                      className={`shrink-0 text-xs px-2 py-0.5 border border-border ${rec.published ? "text-success" : "text-text-secondary"}`}
                      title="Toggle published"
                    >
                      {rec.published ? "Published" : "Draft"}
                    </button>
                  </div>
                  <div className="flex gap-2 shrink-0">
                    <button onClick={() => openEditMusicForm(rec)} className="border border-border px-2 py-1 text-xs font-heading hover:border-accent">Edit</button>
                    <button onClick={() => handleDeleteMusic(rec)} className="border border-border px-2 py-1 text-xs font-heading hover:border-accent text-accent">Delete</button>
                  </div>
                </div>
              ))}
              {musicRecs.length === 0 && (
                <p className="text-text-secondary text-sm">No music recs yet.</p>
              )}
            </div>
          </section>
        </>
      )}
    </Layout>
  );
}
