<script lang="ts">
  import { authStore } from '$lib/stores/auth';
  import { unreadNotificationCount } from '$lib/stores/websocket';
  import NotificationDropdown from '$lib/components/notifications/NotificationDropdown.svelte';

  let showNotifications = false;

  async function handleLogout() {
    await fetch('http://localhost:8080/api/logout', { 
      method: 'POST',
      credentials: 'include',
    });
    authStore.set({ user: null });
    window.location.href = '/';
  }
</script>

<header class="bg-white shadow-md sticky top-0 z-10">
  <nav class="container mx-auto px-6 py-3 flex justify-between items-center">
    <a href="/" class="text-xl font-bold text-gray-800">Echo-Chamber</a>
    <div class="flex items-center space-x-4">
      {#if $authStore.user}
        <a href="/messages" class="text-gray-600 hover:text-gray-800">Messages</a>

        <div class="relative">
          <button on:click={() => {
            showNotifications = !showNotifications;
            if (showNotifications) unreadNotificationCount.set(0);
          }} class="relative text-gray-600 hover:text-gray-800">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"></path></svg>
            {#if $unreadNotificationCount > 0}
              <span class="absolute top-0 right-0 block h-2 w-2 rounded-full bg-red-500 ring-2 ring-white"></span>
            {/if}
          </button>
          <NotificationDropdown bind:show={showNotifications} />
        </div>
        
        <a href="/profile/edit" class="text-gray-600 hover:text-gray-800 transition">Edit Profile</a>
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
