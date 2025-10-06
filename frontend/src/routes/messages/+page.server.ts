import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  const jwt = cookies.get('jwt');
  if (!jwt) throw redirect(307, '/login');

  const res = await fetch('http://localhost:8080/api/chat/conversations', {
    headers: { 'Cookie': `jwt=${jwt}` }
  });
  
  if (res.ok) {
    const conversations = await res.json();
    return { conversations };
  }
  return { conversations: [] };
};
