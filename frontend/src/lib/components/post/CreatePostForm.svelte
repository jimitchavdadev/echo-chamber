<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Post } from '$lib/types';

  const dispatch = createEventDispatcher();

  let content = '';
  let isLoading = false;
  let error = '';

  async function handleSubmit() {
    if (!content.trim() || isLoading) return;
    isLoading = true;
    error = '';

    try {
      const res = await fetch('http://localhost:8080/api/posts', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content }),
        credentials: 'include'
      });

      const newPost: Post = await res.json();
      if (!res.ok) {
        throw new Error('Failed to create post');
      }

      dispatch('postCreated', newPost);
      content = ''; // Clear the form
    } catch (err: any) {
      error = err.message;
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="bg-white p-4 rounded-lg shadow-sm border mb-8">
  <form on:submit|preventDefault={handleSubmit}>
    <textarea
      bind:value={content}
      class="w-full p-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
      rows="3"
      placeholder="What's on your mind?"
    ></textarea>
    <div class="flex justify-end items-center mt-2">
       {#if error}<p class="text-sm text-red-500 mr-4">{error}</p>{/if}
      <button
        type="submit"
        disabled={isLoading || !content.trim()}
        class="px-4 py-2 bg-indigo-600 text-white font-semibold rounded-md disabled:bg-indigo-300 hover:bg-indigo-700 transition"
      >
        {isLoading ? 'Posting...' : 'Post'}
      </button>
    </div>
  </form>
</div>
