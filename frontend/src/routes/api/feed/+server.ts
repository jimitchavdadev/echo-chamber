import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ fetch, cookies, url }) => {
  const limit = url.searchParams.get('limit') || '10';
  const offset = url.searchParams.get('offset') || '0';
  const jwt = cookies.get('jwt');

  if (!jwt) {
    return json({ error: 'Not authenticated' }, { status: 401 });
  }

  try {
    const res = await fetch(`http://localhost:8080/api/feed?limit=${limit}&offset=${offset}`, {
      headers: { 'Cookie': `jwt=${jwt}` }
    });

    if (res.ok) {
      const posts = await res.json();
      return json(posts);
    }
    return json({ error: 'Failed to fetch feed' }, { status: res.status });
  } catch (err) {
    return json({ error: 'Server error' }, { status: 500 });
  }
};
