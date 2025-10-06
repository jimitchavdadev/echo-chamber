<script lang="ts">
  import { notifications } from '$lib/stores/websocket';
  import { slide } from 'svelte/transition';

  export let show = false;
</script>

{#if show}
  <div
    transition:slide={{ duration: 200 }}
    class="absolute right-0 mt-2 w-80 bg-white rounded-md shadow-lg border z-20"
  >
    <div class="p-3 border-b font-semibold">Notifications</div>
    <div class="max-h-96 overflow-y-auto">
      {#if $notifications.length > 0}
        {#each $notifications as notification (notification.ID)}
          <a
            href="/posts/{notification.EntityID}"
            class="block p-3 hover:bg-gray-50 border-b text-sm"
          >
            <span class="font-bold">{notification.Actor.username}</span>
            {notification.Type === 'like' ? 'liked your post.' : 'commented on your post.'}
          </a>
        {/each}
      {:else}
        <p class="text-center text-gray-500 p-4">No new notifications.</p>
      {/if}
    </div>
  </div>
{/if}
