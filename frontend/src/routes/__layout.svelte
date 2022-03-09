<script>
	import '../app.css';
</script>

<svelte:head>
	<title>Bake Shop</title>
</svelte:head>

<main class="mx-auto p-8 max-w-prose">
	<a href="/" class="hover:underline">
		<h1>Bake Shop</h1>
	</a>
	<slot />
</main>
