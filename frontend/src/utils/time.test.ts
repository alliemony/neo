import { describe, it, expect, vi, afterEach } from 'vitest';
import { formatRelativeTime } from './time';

describe('formatRelativeTime', () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  it('shows "just now" for very recent times', () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-04-10T12:00:00Z'));

    const result = formatRelativeTime('2026-04-10T11:59:30Z');
    expect(result).toBe('just now');
  });

  it('shows minutes ago', () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-04-10T12:00:00Z'));

    const result = formatRelativeTime('2026-04-10T11:55:00Z');
    expect(result).toBe('5 minutes ago');
  });

  it('shows 1 minute ago (singular)', () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-04-10T12:00:00Z'));

    const result = formatRelativeTime('2026-04-10T11:59:00Z');
    expect(result).toBe('1 minute ago');
  });

  it('shows hours ago', () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-04-10T12:00:00Z'));

    const result = formatRelativeTime('2026-04-10T10:00:00Z');
    expect(result).toBe('2 hours ago');
  });

  it('shows days ago', () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-04-10T12:00:00Z'));

    const result = formatRelativeTime('2026-04-07T12:00:00Z');
    expect(result).toBe('3 days ago');
  });

  it('shows absolute date for posts older than 7 days', () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-04-10T12:00:00Z'));

    const result = formatRelativeTime('2026-03-09T12:00:00Z');
    expect(result).toBe('Mar 9, 2026');
  });
});
