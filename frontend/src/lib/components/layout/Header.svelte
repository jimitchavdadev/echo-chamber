<script lang="ts">
  import { authStore } from '$lib/stores/auth';

  async function handleLogout() {
    // 1. Tell the backend to clear the HttpOnly cookie
    await fetch('http://localhost:8080/api/logout', { method: 'POST' });

    // 2. Clear the user from our client-side store
    authStore.set({ user: null });

    // 3. Force a full page navigation to the login page to ensure a clean state
    window.location.href = '/login';
  }
</script>

<header class="bg-white shadow-md">
  <nav class="container mx-auto px-6 py-3 flex justify-between items-center">
    <a href="/" class="text-xl font-bold text-gray-800">Echo-Chamber</a>
    <div class="flex items-center space-x-4">
      {#if $authStore.user}
        <span>Welcome, {$authStore.user.username}</span>
        <button on:click={handleLogout} class="px-4 py-2 text-sm text-white bg-red-500 rounded hover:bg-red-600">
          Logout
        </button>
      {:else}
        <a href="/login" class="text-gray-600 hover:text-gray-800">Login</a>
        <a href="/register" class="px-4 py-2 text-sm text-white bg-indigo-500 rounded hover:bg-indigo-600">Register</a>
      {/if}
    </div>
  </nav>
</header>
