import { describe, it, expect, vi, afterEach } from 'vitest';
import { useFetch } from '../FetchData.js';

const jsonResponse = (body, status) => new Response(JSON.stringify(body), {
  status,
  headers: { 'Content-Type': 'application/json' },
});

describe('useFetch', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('returns the data and the status of a successful response', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => jsonResponse({ id: 42 }, 200)));

    const { data, hasError, status } = await useFetch('/trees/42', {});

    expect(hasError.value).toBe(false);
    expect(status.value).toBe(200);
    expect(data.value).toEqual({ id: 42 });
  });

  it('keeps the message of a JSON string error body', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => jsonResponse('current user has no editor privilege on this object', 403)));

    const { hasError, status, errorMessage } = await useFetch('/trees/42', {});

    expect(hasError.value).toBe(true);
    expect(status.value).toBe(403);
    expect(errorMessage.value).toBe('current user has no editor privilege on this object');
  });

  it('keeps the message of a {"message"} error body', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => jsonResponse({ message: 'invalid or expired jwt' }, 401)));

    const { status, errorMessage } = await useFetch('/trees/42', {});

    expect(status.value).toBe(401);
    expect(errorMessage.value).toBe('invalid or expired jwt');
  });

  it('falls back to the status when the error body is not JSON', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response('oops', { status: 502 })));

    const { hasError, errorMessage } = await useFetch('/trees/42', {});

    expect(hasError.value).toBe(true);
    expect(errorMessage.value).toBe('HTTP 502');
  });

  it('reports status 0 when the server cannot be reached', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => { throw new TypeError('Failed to fetch'); }));

    const { hasError, status, errorMessage } = await useFetch('/trees/42', {});

    expect(hasError.value).toBe(true);
    expect(status.value).toBe(0);
    expect(errorMessage.value).toBe('Failed to fetch');
  });
});
