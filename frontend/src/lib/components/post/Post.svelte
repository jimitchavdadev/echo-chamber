<script lang="ts">
  import type { Post } from '$lib/types';
  import { createEventDispatcher } from 'svelte';
  import { authStore } from '$lib/stores/auth';

  export let post: Post;
  
  const dispatch = createEventDispatcher();
  let isLikeLoading = false;

  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit'
    });
  }

  async function handleLikeToggle() {
    if (isLikeLoading) return;
    isLikeLoading = true;

    // Optimistic update
    const originalIsLiked = post.isLiked;
    const originalLikeCount = post.likeCount;
    post.isLiked = !post.isLiked;
    post.likeCount += post.isLiked ? 1 : -1;
    
    const method = originalIsLiked ? 'DELETE' : 'POST';
    
    try {
      const res = await fetch(`http://localhost:8080/api/posts/${post.ID}/like`, {
        method,
        credentials: 'include'
      });
      if (!res.ok) throw new Error('Like failed');
    } catch {
      // Revert on failure
      post.isLiked = originalIsLiked;
      post.likeCount = originalLikeCount;
    } finally {
      isLikeLoading = false;
    }
  }

  async function handleDelete() {
    if (confirm('Are you sure you want to delete this post?')) {
      try {
        const res = await fetch(`http://localhost:8080/api/posts/${post.ID}`, {
          method: 'DELETE',
          credentials: 'include'
        });
        if (!res.ok) throw new Error('Delete failed');
        dispatch('postDeleted', post.ID);
      } catch (err) {
        alert('Failed to delete post.');
      }
    }
  }
</script>

<div class="bg-white border border-gray-200 rounded-lg p-4 my-4 shadow-sm">
  <div class="flex justify-between items-start">
    <div class="flex items-center mb-3">
      <div class="w-10 h-10 rounded-full bg-gray-300 mr-3"></div>
      <div>
        <a href="/users/{post.author.username}" class="font-bold text-gray-800 hover:underline">
          {post.author.username}
        </a>
        <p class="text-sm text-gray-500">{formatDate(post.CreatedAt)}</p>
      </div>
    </div>
    {#if $authStore.user && $authStore.user.id === post.author.id}
       <button on:click={handleDelete} class="text-gray-400 hover:text-red-500 text-lg">&times;</button>
    {/if}
  </div>

  <a href="/posts/{post.ID}" class="block">
    <p class="text-gray-700 whitespace-pre-wrap">{post.content}</p>
  </a>

  <div class="flex items-center space-x-4 mt-4 pt-2 border-t">
    <button
      on:click={handleLikeToggle}
      disabled={!$authStore.user || isLikeLoading}
      class="flex items-center space-x-1 text-gray-500 hover:text-red-500 disabled:text-gray-300 disabled:cursor-not-allowed"
    >
      <svg class:text-red-500={post.isLiked} class="w-5 h-5" fill={post.isLiked ? 'currentColor' : 'none'} stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 016.364 0L12 7.636l1.318-1.318a4.5 4.5 0 016.364 6.364L12 20.364l-7.682-7.682a4.5 4.5 0 010-6.364z"></path></svg>
      <span>{post.likeCount}</span>
    </button>
    <a href="/posts/{post.ID}" class="flex items-center space-x-1 text-gray-500 hover:text-indigo-500">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path></svg>
      </a>
  </div>
</div>
