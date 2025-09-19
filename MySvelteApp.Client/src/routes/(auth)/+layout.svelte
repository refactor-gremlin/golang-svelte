<script lang="ts">
	import '../../app.css';
	import { Toaster } from 'svelte-sonner';

	let { children } = $props();
</script>

<svelte:head>
	<title>Authentication - MySvelteApp</title>
</svelte:head>

{@render children?.()}
<Toaster />
