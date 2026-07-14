import { describe, expect, it } from 'vitest';
import { buildQuery, isTaskTerminal, taskStreamId } from '../src/lib/utils.js';

describe('buildQuery', () => {
  it('appends pagination params', () => {
    expect(buildQuery('/tasks', { page: 2, per_page: 20 })).toBe('/tasks?page=2&per_page=20');
  });

  it('returns path when no params', () => {
    expect(buildQuery('/agents', {})).toBe('/agents');
  });
});

describe('taskStreamId', () => {
  it('prefers correlation_id', () => {
    expect(taskStreamId({ id: 'rec1', correlation_id: 'task-123' })).toBe('task-123');
  });

  it('falls back to id', () => {
    expect(taskStreamId({ id: 'rec1' })).toBe('rec1');
  });
});

describe('isTaskTerminal', () => {
  it('detects terminal statuses', () => {
    expect(isTaskTerminal('completed')).toBe(true);
    expect(isTaskTerminal('running')).toBe(false);
  });
});