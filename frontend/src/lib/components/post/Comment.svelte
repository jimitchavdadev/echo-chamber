<script lang="ts">
  import type { Comment } from '$lib/types';

  export let comment: Comment;
  
  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleString('en-US', { month: 'short', day: 'numeric' });
  }
</script>

<div class="flex space-x-3 py-3 border-t">
  <div class="w-8 h-8 rounded-full bg-gray-300 flex-shrink-0"></div>
  <div class="flex-1">
    <div class="flex items-baseline space-x-2">
      <a href="/users/{comment.author.username}" class="font-bold text-sm hover:underline">{comment.author.username}</a>
      <span class="text-xs text-gray-500">{formatDate(comment.CreatedAt)}</span>
    </div>
    <p class="text-gray-700 text-sm mt-1">{comment.content}</p>
  </div>
</div>
