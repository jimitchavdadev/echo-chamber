import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, params, cookies }) => {
  const { id } = params;
  const jwt = cookies.get('jwt');
  const headers: HeadersInit = jwt ? { 'Cookie': `jwt=${jwt}` } : {};

  // Fetch post and comments in parallel
  const postRes = fetch(`http://localhost:8080/api/posts/${id}`, { headers });
  const commentsRes = fetch(`http://localhost:8080/api/posts/${id}/comments`, { headers });

  const [postResult, commentsResult] = await Promise.all([postRes, commentsRes]);

  if (!postResult.ok) {
    return { status: postResult.status, error: new Error('Post not found') };
  }

  const post = await postResult.json();
  const comments = commentsResult.ok ? await commentsResult.json() : [];
  
  return { post, comments };
};
