import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, params, cookies }) => {
  try {
    const jwt = cookies.get('jwt');
    const headers: HeadersInit = jwt ? { 'Cookie': `jwt=${jwt}` } : {};

    const res = await fetch(`http://localhost:8080/api/users/${params.username}`, { headers });
    
    if (!res.ok) {
      return { status: res.status, error: new Error('Could not fetch user profile') };
    }

    const profileData = await res.json();
    return { profile: profileData };
  } catch (error) {
    return { status: 500, error: new Error('Server error') };
  }
};
