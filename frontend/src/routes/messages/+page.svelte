<script lang="ts">
  import type { PageData } from './$types';
  import { page } from '$app/stores';
  import { debounce } from '$lib/utils/debounce';
  import { goto } from '$app/navigation';
  import type { User } from '$lib/types';

  export let data: PageData;

  let isSearching = false;
  let searchTerm = '';
  let searchResults: User[] = [];
  let isSearchLoading = false;

  const search = debounce(async () => {
    if (searchTerm.trim().length < 2) {
      searchResults = [];
      return;
    }
    isSearchLoading = true;
    const res = await fetch(`/api/users/search?q=${searchTerm}`);
    if (res.ok) {
      searchResults = await res.json();
    }
    isSearchLoading = false;
  }, 300);

  function startConversationWith(user: User) {
    isSearching = false;
    searchTerm = '';
    searchResults = [];
    goto(`/messages/${user.id}`);
  }

  $: if (isSearching) search();
</script>

<div class="w-1/3 border-r flex flex-col bg-white">
  <div class="p-4 border-b flex justify-between items-center">
    <h2 class="font-semibold">Conversations</h2>
    <button
      on:click={() => isSearching = !isSearching}
      class="p-1 rounded-full hover:bg-gray-200"
      title="New message"
    >
      <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
    </button>
  </div>
  
  {#if isSearching}
    <div class="p-2 border-b relative">
      <input
        type="text"
        bind:value={searchTerm}
        placeholder="Search for a user..."
        class="w-full p-2 border rounded-md"
      />
      {#if searchTerm.trim().length > 1}
        <div class="absolute top-full left-0 w-full bg-white border shadow-lg rounded-b-md z-10">
          {#if isSearchLoading}
            <div class="p-2 text-gray-500">Searching...</div>
          {:else if searchResults.length > 0}
            {#each searchResults as user (user.id)}
              <button on:click={() => startConversationWith(user)} class="w-full text-left p-2 hover:bg-gray-100">
                {user.username}
              </button>
            {/each}
          {:else}
            <div class="p-2 text-gray-500">No users found.</div>
          {/if}
        </div>
      {/if}
    </div>
  {/if}

  <div class="flex-1 overflow-y-auto">
    {#if data.conversations && data.conversations.length > 0}
      {#each data.conversations as convo (convo.id)}
        <a
          href="/messages/{convo.id}"
          class="block p-4 border-b hover:bg-gray-50"
          class:bg-indigo-50={$page.params.userId == convo.id.toString()}
        >
          <div class="flex items-center">
            <div class="w-10 h-10 rounded-full bg-gray-300 mr-3"></div>
            <span class="font-semibold">{convo.username}</span>
          </div>
        </a>
      {/each}
    {:else}
        <p class="p-4 text-center text-gray-500">No conversations yet.</p>
    {/if}
  </div>
</div>

<div class="w-2/3 flex items-center justify-center text-gray-500 bg-gray-50">
  Select a conversation to start chatting.
</div>
