import { describe, expect, it, vi } from 'vitest';
import { parseSSEFrame } from '../src/lib/streams.js';
import type { TaskStreamHandlers } from '../src/lib/streams.js';

// Re-import dispatch via stream test - test parseSSEFrame directly
describe('parseSSEFrame', () => {
  it('parses token event', () => {
    const frame = parseSSEFrame('event: token\ndata: {"task_id":"t1","delta":"hi"}\n');
    expect(frame).toEqual({
      event: 'token',
      data: '{"task_id":"t1","delta":"hi"}',
    });
  });

  it('ignores heartbeat comments', () => {
    expect(parseSSEFrame(': ping')).toBeNull();
  });

  it('ignores empty chunks', () => {
    expect(parseSSEFrame('')).toBeNull();
  });
});

describe('SSE dispatch integration', () => {
  it('invokes onToken handler', async () => {
    const { streamTaskSSE } = await import('../src/lib/streams.js');

    const onToken = vi.fn();
    const handlers: TaskStreamHandlers = { onToken };

    const body = new ReadableStream({
      start(controller) {
        controller.enqueue(
          new TextEncoder().encode(
            'event: token\ndata: {"task_id":"t1","delta":"Hello"}\n\n',
          ),
        );
        controller.close();
      },
    });

    await streamTaskSSE(
      {
        baseUrl: 'http://localhost:8080',
        fetch: async () => new Response(body, { status: 200 }),
      },
      't1',
      handlers,
    );

    expect(onToken).toHaveBeenCalledWith({ task_id: 't1', delta: 'Hello' });
  });
});