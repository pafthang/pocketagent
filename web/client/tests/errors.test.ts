import { describe, expect, it } from 'vitest';
import { PocketAgentError, isForbidden, isUnauthorized } from '../src/lib/errors.js';

describe('PocketAgentError', () => {
  it('exposes status helpers', () => {
    const err = new PocketAgentError('nope', 401);
    expect(err.isUnauthorized).toBe(true);
    expect(isUnauthorized(err)).toBe(true);
    expect(isForbidden(err)).toBe(false);
  });
});