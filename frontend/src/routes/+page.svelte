<script lang="ts">
  import { authStore } from '$lib/stores/auth';
  import Post from '$lib/components/post/Post.svelte';
  import type { PageData } from './$types';
  import { intersectionObserver } from '$lib/actions/intersectionObserver';
  import Loader from '$lib/components/ui/Loader.svelte';

  export let data: PageData;
  
  let posts = data.posts || [];
  let isLoading = false;
  let allPostsLoaded = posts.length < 10; // If initial fetch is less than a full page, we're done

  async function loadMorePosts() {
    if (isLoading || allPostsLoaded) return;
    isLoading = true;

    const offset = posts.length;
    const res = await fetch(`/api/feed?limit=10&offset=${offset}`);

    if (res.ok) {
      const newPosts = await res.json();
      if (newPosts.length > 0) {
        posts = [...posts, ...newPosts];
      }
      if (newPosts.length < 10) {
        allPostsLoaded = true;
      }
    } else {
      // Handle error or stop trying
      allPostsLoaded = true;
    }
    isLoading = false;
  }
</script>

<div class="container mx-auto max-w-2xl p-4 md:p-8">
  {#if $authStore.user}
    <div class="bg-white p-4 rounded-lg shadow-sm border mb-8">
      <h2 class="text-lg font-semibold">Home</h2>
    </div>

    {#if posts.length > 0}
      {#each posts as post (post.ID)}
        <Post {post} />
      {/each}

      {#if !allPostsLoaded}
        <div use:intersectionObserver on:intersect={loadMorePosts}>
          {#if isLoading}
            <Loader />
          {/if}
        </div>
      {:else}
        <p class="text-center text-gray-500 py-4">You've reached the end!</p>
      {/if}
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
