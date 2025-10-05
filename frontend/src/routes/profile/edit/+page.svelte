<script lang="ts">
  import { authStore } from '$lib/stores/auth';

  let bio = $authStore.user?.bio || '';
  let error = '';
  let successMessage = '';

  async function handleSubmit() {
    error = '';
    successMessage = '';

    try {
      const res = await fetch('http://localhost:8080/api/profile', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ bio }),
        credentials: 'include', // Important for sending the auth cookie
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.error || 'Failed to update profile');
      }

      // Update the local store with the new user data from the response
      authStore.set({ user: data.user });
      successMessage = 'Profile updated successfully!';
      setTimeout(() => (successMessage = ''), 3000); // Clear message after 3 seconds
    } catch (err: any) {
      error = err.message;
    }
  }
</script>

<div class="container mx-auto max-w-2xl p-8">
  <h1 class="text-3xl font-bold mb-6">Edit Your Profile</h1>

  {#if $authStore.user}
    <form on:submit|preventDefault={handleSubmit} class="space-y-6 bg-white p-8 rounded-lg shadow-md">
      <div>
        <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
        <input
          type="text"
          id="username"
          value={$authStore.user.username}
          disabled
          class="w-full px-3 py-2 mt-1 border bg-gray-100 border-gray-300 rounded-md"
        />
      </div>
       <div>
        <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
        <input
          type="email"
          id="email"
          value={$authStore.user.email}
          disabled
          class="w-full px-3 py-2 mt-1 border bg-gray-100 border-gray-300 rounded-md"
        />
      </div>
      <div>
        <label for="bio" class="block text-sm font-medium text-gray-700">Bio</label>
        <textarea
          id="bio"
          bind:value={bio}
          rows="4"
          class="w-full px-3 py-2 mt-1 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          placeholder="Tell us a little about yourself..."
        ></textarea>
      </div>

      {#if error}
        <p class="text-sm text-red-600">{error}</p>
      {/if}

      {#if successMessage}
        <p class="text-sm text-green-600">{successMessage}</p>
      {/if}

      <div>
        <button
          type="submit"
          class="w-full px-4 py-2 font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          Save Changes
        </button>
      </div>
    </form>
  {/if}
</div>
