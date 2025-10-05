<script lang="ts">
  import type { PageData } from './$types';
  import { page, navigating } from '$app/stores';
  import { chatMessages, sendChatMessage } from '$lib/stores/websocket';
  import { authStore } from '$lib/stores/auth';
  import { onMount, afterUpdate } from 'svelte';

  export let data: PageData;

  let currentChat: any[] = [];
  let newMessage = '';
  let chatContainer: HTMLElement;
  let partnerId: number;
  let partner: any;

  // This reactive block merges server-loaded history with live messages from the store
  $: {
    partnerId = parseInt($page.params.userId, 10);
    partner = data.conversations?.find(c => c.id === partnerId);
    
    const history = data.messages || [];
    const live = $chatMessages[partnerId] || [];
    
    // Combine and deduplicate
    const all = [...history, ...live];
    const uniqueMessages = Array.from(new Map(all.map(item => [item.ID, item])).values());
    currentChat = uniqueMessages.sort((a, b) => new Date(a.CreatedAt).getTime() - new Date(b.CreatedAt).getTime());
  }

  function handleSend() {
    if (!newMessage.trim()) return;
    sendChatMessage(partnerId, newMessage);
    // Optimistic update
    currentChat = [...currentChat, {
      ID: Math.random(),
      SenderID: $authStore.user?.id,
      Content: newMessage,
      CreatedAt: new Date().toISOString()
    }];
    newMessage = '';
  }

  // Auto-scroll logic
  let shouldScroll = true;
  afterUpdate(() => {
    if (shouldScroll && chatContainer) {
      chatContainer.scrollTo(0, chatContainer.scrollHeight);
    }
  });

  // Keep auto-scroll active when new data comes in, but not from user scrolling up
  $: if (currentChat, () => shouldScroll = true);
</script>

<div class="w-1/3 border-r flex flex-col bg-white">
  <div class="p-4 border-b font-semibold">Conversations</div>
  <div class="flex-1 overflow-y-auto">
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
  </div>
</div>

<div class="w-2/3 flex flex-col h-full">
  {#if navigating}
    <div class="flex-1 flex items-center justify-center"><p>Loading...</p></div>
  {:else if partner}
    <div class="p-4 border-b font-semibold bg-white">{partner.username}</div>
    <div 
      bind:this={chatContainer} 
      on:scroll={() => {
        // Disable auto-scroll if user scrolls up
        if (chatContainer.scrollHeight - chatContainer.scrollTop > chatContainer.clientHeight + 100) {
          shouldScroll = false;
        }
      }}
      class="flex-1 overflow-y-auto p-4 bg-gray-100"
    >
      {#each currentChat as msg (msg.ID)}
        <div class="flex mb-3" class:justify-end={msg.SenderID === $authStore.user?.id}>
          <div 
            class="rounded-lg py-2 px-3 max-w-sm"
            class:bg-indigo-500={msg.SenderID === $authStore.user?.id}
            class:text-white={msg.SenderID === $authStore.user?.id}
            class:bg-gray-200={msg.SenderID !== $authStore.user?.id}
          >
            {msg.Content}
          </div>
        </div>
      {/each}
    </div>
    <div class="p-4 bg-white border-t">
      <form on:submit|preventDefault={handleSend} class="flex space-x-3">
        <input
          bind:value={newMessage}
          type="text"
          class="flex-1 p-2 border rounded-md focus:ring-indigo-500 focus:border-indigo-500"
          placeholder="Type a message..."
        />
        <button type="submit" class="px-4 py-2 bg-indigo-600 text-white font-semibold rounded-md">Send</button>
      </form>
    </div>
  {:else}
     <div class="flex-1 flex items-center justify-center text-gray-500">
        Select a conversation to start chatting.
     </div>
  {/if}
</div>
