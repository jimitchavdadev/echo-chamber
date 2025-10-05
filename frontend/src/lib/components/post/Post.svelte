<script lang="ts">
  import type { Post } from '$lib/types'; // We will create this type definition next

  export let post: Post;

  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit'
    });
  }
</script>

<div class="bg-white border border-gray-200 rounded-lg p-4 my-4 shadow-sm">
  <div class="flex items-center mb-3">
    <div class="w-10 h-10 rounded-full bg-gray-300 mr-3"></div>
    <div>
      <a href="/users/{post.author.username}" class="font-bold text-gray-800 hover:underline">
        {post.author.username}
      </a>
      <p class="text-sm text-gray-500">{formatDate(post.CreatedAt)}</p>
    </div>
  </div>
  <p class="text-gray-700 whitespace-pre-wrap">{post.content}</p>
</div>
