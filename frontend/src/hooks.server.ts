import { redirect, type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
  // --- This part fetches the user and is working correctly ---
  const jwt = event.cookies.get('jwt');
  if (jwt) {
    try {
      const response = await fetch('http://localhost:8080/api/me', {
        headers: { 'Cookie': `jwt=${jwt}` }
      });
      if (response.ok) {
        const { user } = await response.json();
        event.locals.user = user;
      } else {
        event.locals.user = null;
      }
    } catch {
      event.locals.user = null;
    }
  } else {
    event.locals.user = null;
  }

  // --- NEW: Global Route Protection ---
  const protectedRoutes = ['/profile']; // Add any other path prefixes that need protecting
  const publicRoutes = ['/', '/login', '/register'];

  const isProtectedRoute = protectedRoutes.some(path => event.url.pathname.startsWith(path));

  if (!event.locals.user && !publicRoutes.includes(event.url.pathname) && isProtectedRoute) {
    // If user is not logged in AND is trying to access a protected page, redirect to homepage.
    throw redirect(307, '/'); 
  }

  return resolve(event);
};
