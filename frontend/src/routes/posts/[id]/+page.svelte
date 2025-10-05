<script lang="ts">
  import type { PageData } from './$types';
  import type { Comment as CommentType } from '$lib/types';
  import Post from '$lib/components/post/Post.svelte';
  import Comment from '$lib/components/post/Comment.svelte';
  import { authStore } from '$lib/stores/auth';

  export let data: PageData;

  let post = data.post;
  let comments: CommentType[] = data.comments || [];
  let newCommentContent = '';
  let isSubmitting = false;

  async function handleCommentSubmit() {
    if (!newCommentContent.trim() || isSubmitting) return;
    isSubmitting = true;

    try {
      const res = await fetch(`http://localhost:8080/api/posts/${post.ID}/comments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content: newCommentContent }),
        credentials: 'include'
      });
      const newComment = await res.json();
      if (!res.ok) throw new Error('Failed to post comment');

      comments = [...comments, newComment]; // Add new comment to the list
      newCommentContent = ''; // Clear form
    } catch (err) {
      alert('Could not post comment.');
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="container mx-auto max-w-2xl p-4 md:p-8">
  {#if post}
    <Post bind:post />

    <div class="bg-white border border-gray-200 rounded-lg p-4 my-4 shadow-sm">
      <h3 class="text-lg font-semibold mb-2">Comments</h3>
      
      {#if $authStore.user}
        <form on:submit|preventDefault={handleCommentSubmit} class="flex items-start space-x-3 mb-4">
          <div class="w-8 h-8 rounded-full bg-gray-300 flex-shrink-0"></div>
          <div class="flex-1">
            <textarea
              bind:value={newCommentContent}
              class="w-full p-2 border border-gray-300 rounded-md text-sm"
              rows="2"
              placeholder="Add a comment..."
            ></textarea>
            <div class="flex justify-end mt-2">
              <button 
                type="submit" 
                disabled={isSubmitting || !newCommentContent.trim()}
                class="px-3 py-1 bg-indigo-600 text-white text-sm font-semibold rounded-md disabled:bg-indigo-300"
              >
                {isSubmitting ? '...' : 'Comment'}
              </button>
            </div>
          </div>
        </form>
      {/if}

      <div class="space-y-2">
        {#if comments.length > 0}
          {#each comments as comment (comment.ID)}
            <Comment {comment} />
          {/each}
        {:else}
          <p class="text-sm text-gray-500 text-center py-4">No comments yet.</p>
        {/if}
      </div>
    </div>
  {:else}
    <p class="text-red-500 text-center">Post not found.</p>
  {/if}
</div>
