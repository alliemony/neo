const MINUTE = 60;
const HOUR = 3600;
const DAY = 86400;
const WEEK = 7 * DAY;

export function formatRelativeTime(dateStr: string): string {
  const date = new Date(dateStr);
  const now = new Date();
  const diffSec = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (diffSec < MINUTE) return 'just now';
  if (diffSec < HOUR) {
    const mins = Math.floor(diffSec / MINUTE);
    return `${mins} ${mins === 1 ? 'minute' : 'minutes'} ago`;
  }
  if (diffSec < DAY) {
    const hrs = Math.floor(diffSec / HOUR);
    return `${hrs} ${hrs === 1 ? 'hour' : 'hours'} ago`;
  }
  if (diffSec < WEEK) {
    const days = Math.floor(diffSec / DAY);
    return `${days} ${days === 1 ? 'day' : 'days'} ago`;
  }

  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}
