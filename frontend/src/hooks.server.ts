// frontend/src/hooks.server.ts
import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
  console.log('\n--- SVELTEKIT HOOK RUNNING ---');

  const jwt = event.cookies.get('jwt');
  console.log('1. Reading cookie from browser:', jwt ? `Found a cookie` : 'No cookie found');

  if (jwt) {
    try {
      // Forward the cookie to the backend's /api/me endpoint
      const response = await fetch('http://localhost:8080/api/me', {
        headers: {
          // Manually construct the Cookie header for the server-side fetch
          'Cookie': `jwt=${jwt}`
        }
      });

      console.log('2. Backend /api/me response status:', response.status);

      if (response.ok) {
        const data = await response.json();
        event.locals.user = data.user;
        console.log('3. Successfully fetched user:', event.locals.user.username);
      } else {
        // This will happen if the token is invalid or expired
        event.locals.user = null;
        console.log('3. Backend returned an error, clearing user.');
      }
    } catch (error) {
      console.error('4. CRITICAL: Failed to fetch user from backend.', error);
      event.locals.user = null;
    }
  } else {
    event.locals.user = null;
  }

  console.log('5. Final user object for this request:', event.locals.user);
  console.log('----------------------------\n');

  return resolve(event);
};