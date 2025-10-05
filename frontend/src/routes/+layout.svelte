<script lang="ts">
  import { page } from '$app/stores';
  import { authStore } from '$lib/stores/auth';
  import Header from '$lib/components/layout/Header.svelte';
  import '../app.css';

  // This reactive statement is the key to persistence.
  // It runs whenever the page data from the server changes.
  // It ensures our client-side store is always in sync with the server's session state.
  $: {
    if ($page.data.user) {
      authStore.set({ user: $page.data.user });
    } else {
      authStore.set({ user: null });
    }
  }
</script>

<div class="min-h-screen bg-gray-50">
  <Header />
  <main>
    <slot />
  </main>
</div>
