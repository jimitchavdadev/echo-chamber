<script lang="ts">
  import { goto } from '$app/navigation';

  let username = '';
  let email = '';
  let password = '';
  let error = '';
  let message = '';

  async function handleSubmit() {
    error = '';
    message = '';
    try {
      const res = await fetch('http://localhost:8080/api/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password }),
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.error || 'Registration failed');
      }
      
      message = data.message;
      setTimeout(() => goto('/login'), 2000); // Redirect to login after 2 seconds

    } catch (err: any) {
      error = err.message;
    }
  }
</script>

<div class="flex items-center justify-center min-h-screen bg-gray-100">
  <div class="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">
    <h1 class="text-2xl font-bold text-center">Create an Account</h1>
    <form on:submit|preventDefault={handleSubmit} class="space-y-6">
      <div>
        <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
        <input id="username" type="text" bind:value={username} required class="w-full px-3 py-2 mt-1 border rounded-md"/>
      </div>
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
        <input id="email" type="email" bind:value={email} required class="w-full px-3 py-2 mt-1 border rounded-md"/>
      </div>
      <div>
        <label for="password" class="block text-sm font-medium text-gray-700">Password (min 6 chars)</label>
        <input id="password" type="password" bind:value={password} required class="w-full px-3 py-2 mt-1 border rounded-md"/>
      </div>
      {#if error}
        <p class="text-sm text-red-600">{error}</p>
      {/if}
       {#if message}
        <p class="text-sm text-green-600">{message}</p>
      {/if}
      <div>
        <button type="submit" class="w-full px-4 py-2 font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700">
          Register
        </button>
      </div>
    </form>
		<p class="text-sm text-center">
			Already have an account? <a href="/login" class="text-indigo-600 hover:underline">Login here</a>
		</p>
  </div>
</div>
