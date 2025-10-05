import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  const jwt = cookies.get('jwt');

  // If there is no JWT, the user is not logged in.
  // The hook handles the redirect, but we can double-check here.
  if (!jwt) {
    // No need to fetch a feed
    return { posts: [] };
  }

  try {
    const res = await fetch('http://localhost:8080/api/feed?limit=10', {
      headers: { 'Cookie': `jwt=${jwt}` }
    });

    if (res.ok) {
      const posts = await res.json();
      return { posts };
    }
    return { posts: [] };
  } catch {
    return { posts: [] };
  }
};
