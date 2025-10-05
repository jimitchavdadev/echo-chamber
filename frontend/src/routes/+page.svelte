<script lang="ts">
  import { authStore } from '$lib/stores/auth';
  import Post from '$lib/components/post/Post.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<div class="container mx-auto max-w-2xl p-4 md:p-8">
  {#if $authStore.user}
    <div class="bg-white p-4 rounded-lg shadow-sm border mb-8">
      <h2 class="text-lg font-semibold">Home</h2>
    </div>

    {#if data.posts && data.posts.length > 0}
      {#each data.posts as post (post.ID)}
        <Post {post} />
      {/each}
    {:else}
      <div class="text-center text-gray-500 p-8">
        <h3 class="text-xl font-semibold">Your feed is empty</h3>
        <p>Follow some users to see their posts here.</p>
      </div>
    {/if}

  {:else}
    <div class="text-center p-8">
      <h1 class="text-4xl font-bold">Welcome to Echo-Chamber</h1>
      <p class="text-gray-600 mt-4">Log in or register to see the feed and connect with others.</p>
    </div>
  {/if}
</div>
