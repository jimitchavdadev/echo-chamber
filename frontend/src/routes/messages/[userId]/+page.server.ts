import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies, fetch, params }) => {
  const jwt = cookies.get('jwt');
  if (!jwt) throw redirect(307, '/login');

  // Fetch both conversations list and history for the selected user
  const convosRes = fetch('http://localhost:8080/api/chat/conversations', {
    headers: { 'Cookie': `jwt=${jwt}` }
  });
  const historyRes = fetch(`http://localhost:8080/api/chat/history/${params.userId}`, {
    headers: { 'Cookie': `jwt=${jwt}` }
  });

  const [convosResult, historyResult] = await Promise.all([convosRes, historyRes]);

  return {
    conversations: convosResult.ok ? await convosResult.json() : [],
    messages: historyResult.ok ? await historyResult.json() : [],
  };
};
