import { useState, useEffect, useRef } from 'react';
import { Layout } from '../components/layout/Layout';
import { SEO } from '../components/SEO';
import { getMusic } from '../services/api';
import type { MusicRec } from '../types/post';

function LyricsSection({ title, artist }: { title: string; artist: string }) {
  const [open, setOpen] = useState(false);
  const [lyrics, setLyrics] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);
  const fetched = useRef(false);

  async function handleToggle() {
    if (!open && !fetched.current) {
      fetched.current = true;
      setLoading(true);
      try {
        const res = await fetch(
          `https://api.lyrics.ovh/v1/${encodeURIComponent(artist)}/${encodeURIComponent(title)}`
        );
        if (!res.ok) throw new Error('not found');
        const data = await res.json();
        setLyrics(data.lyrics || null);
        if (!data.lyrics) setError(true);
      } catch {
        setError(true);
      } finally {
        setLoading(false);
      }
    }
    setOpen((o) => !o);
  }

  return (
    <div className="mt-2">
      <button
        onClick={handleToggle}
        className="text-xs text-text-secondary hover:text-accent transition-colors flex items-center gap-1"
      >
        <span>{open ? '▾' : '▸'}</span>
        <span>Lyrics</span>
      </button>
      {open && (
        <div className="mt-2 text-sm text-text-secondary">
          {loading && <p>Loading lyrics…</p>}
          {error && !loading && <p className="italic">Lyrics unavailable.</p>}
          {lyrics && !loading && (
            <pre className="whitespace-pre-wrap font-body leading-relaxed">{lyrics}</pre>
          )}
        </div>
      )}
    </div>
  );
}

function StreamingPill({ href, label }: { href: string; label: string }) {
  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className="inline-flex items-center text-xs px-2.5 py-0.5 bg-surface border border-border text-text-secondary no-underline hover:text-accent hover:border-accent transition-colors"
    >
      {label}
    </a>
  );
}

function MusicCard({ rec }: { rec: MusicRec }) {
  return (
    <div className="py-4 border-t border-border first:border-t-0">
      <div className="flex gap-4">
        {rec.cover_url && (
          <img
            src={rec.cover_url}
            alt={`${rec.title} cover`}
            className="w-16 h-16 object-cover shrink-0"
          />
        )}
        <div className="min-w-0 flex-1">
          <p className="font-bold text-text-primary text-[17px] leading-tight">{rec.title}</p>
          <p className="text-text-secondary text-sm mt-0.5">
            {rec.artist}
            {rec.album && <span className="ml-1">· {rec.album}</span>}
          </p>
          {rec.note && (
            <p className="text-text-secondary text-sm mt-1 italic">{rec.note}</p>
          )}
          {(rec.spotify_url || rec.apple_url) && (
            <div className="flex gap-2 mt-2 flex-wrap">
              {rec.spotify_url && <StreamingPill href={rec.spotify_url} label="Spotify" />}
              {rec.apple_url && <StreamingPill href={rec.apple_url} label="Apple Music" />}
            </div>
          )}
          <LyricsSection title={rec.title} artist={rec.artist} />
        </div>
      </div>
    </div>
  );
}

export function Music() {
  const [recs, setRecs] = useState<MusicRec[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    getMusic()
      .then((data) => setRecs(data.recs ?? []))
      .catch(() => setError(true))
      .finally(() => setLoading(false));
  }, []);

  return (
    <Layout>
      <SEO title="Music" path="/music" />

      <div className="mb-9">
        <h1 className="text-[30px] font-bold text-text-primary tracking-tight mb-1.5">
          Music
        </h1>
        <p className="text-text-secondary text-[18px]">
          What's been on repeat.
        </p>
      </div>

      {loading && <p className="text-text-secondary">Loading…</p>}
      {error && <p className="text-text-secondary">Could not load music. Try again later.</p>}
      {!loading && !error && recs.length === 0 && (
        <p className="text-text-secondary">No music recs yet.</p>
      )}

      <div>
        {recs.map((rec) => (
          <MusicCard key={rec.id} rec={rec} />
        ))}
      </div>
    </Layout>
  );
}
