<script lang="ts">
  import type { PageData } from './$types';
  import { authStore } from '$lib/stores/auth';

  export let data: PageData;
  let profile = data.profile;
  let isLoading = false;

  async function handleFollowToggle() {
    if (!$authStore.user || isLoading) return;
    isLoading = true;

    const method = profile.isFollowing ? 'DELETE' : 'POST';
    
    // Optimistic UI update
    const originalFollowState = profile.isFollowing;
    profile.isFollowing = !profile.isFollowing;

    try {
      const res = await fetch(`http://localhost:8080/api/users/${profile.id}/follow`, {
        method,
        credentials: 'include',
      });

      if (!res.ok) {
        // Revert on failure
        profile.isFollowing = originalFollowState;
        alert('Action failed. Please try again.');
      }
    } catch {
      // Revert on failure
      profile.isFollowing = originalFollowState;
      alert('Action failed. Please try again.');
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="container mx-auto max-w-2xl p-8">
  {#if profile}
    <div class="flex justify-between items-start">
      <div>
        <h1 class="text-3xl font-bold">{profile.username}</h1>
        <p class="text-gray-600 mt-2">{profile.bio || 'No bio yet.'}</p>
      </div>
      {#if $authStore.user && $authStore.user.id !== profile.id}
        <button
          on:click={handleFollowToggle}
          disabled={isLoading}
          class="px-4 py-2 rounded-md font-semibold text-sm transition
            {profile.isFollowing 
              ? 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50' 
              : 'bg-black text-white hover:bg-gray-800'}"
        >
          {isLoading ? '...' : (profile.isFollowing ? 'Following' : 'Follow')}
        </button>
      {/if}
    </div>
    <div class="mt-8 border-t pt-4">
      <h2 class="text-xl font-semibold">Posts</h2>
      <p class="text-gray-500 mt-4">User posts will appear here.</p>
    </div>
  {:else}
    <p class="text-red-500 text-center">Could not load profile.</p>
  {/if}
</div>
