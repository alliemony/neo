import { tagPalette } from "../config";

export function getTagColor(tag: string): { bg: string; text: string } {
  let h = 0;
  for (let i = 0; i < tag.length; i++) h = tag.charCodeAt(i) + ((h << 5) - h);
  const idx = Math.abs(h) % tagPalette.length;
  return tagPalette[idx] as { bg: string; text: string };
}
