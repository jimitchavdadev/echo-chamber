<script lang="ts">
	import { authStore } from "$lib/stores/auth";
</script>

<div class="container mx-auto p-8">
	<h1 class="text-3xl font-bold">Welcome to Echo-Chamber</h1>
	{#if $authStore.user}
		<p class="mt-4">You are logged in as {$authStore.user.username}.</p>
	{:else}
		<p class="mt-4">Please log in or register to continue.</p>
	{/if}
</div>
