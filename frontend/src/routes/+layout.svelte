<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { authStore } from '$lib/stores/auth';
  import { connectWebSocket } from '$lib/stores/websocket';
  import Header from '$lib/components/layout/Header.svelte';
  import '../app.css';

  $: {
    if ($page.data.user) {
      authStore.set({ user: $page.data.user });
    } else {
      authStore.set({ user: null });
    }
  }

  // Connect to WebSocket on client-side when user is logged in
  onMount(() => {
    const unsubscribe = authStore.subscribe(state => {
      if (state.user) {
        connectWebSocket();
      }
    });
    return unsubscribe;
  });
</script>

<div class="min-h-screen bg-gray-50">
  <Header />
  <main>
    <slot />
  </main>
</div>
