<script context="module">
	/** @type {import('@sveltejs/kit').ErrorLoad} */
	export function load({ error, status }) {
		console.debug(error);
		return {
			props: {
				status,
				error
			}
		};
	}
</script>

<script>
	export let status;
	export let error;

	console.debug(error);
</script>

<h2>{status} error</h2>
<p class="mb-4">Did you remember to fire up the backend?</p>
<code class="mt-4">{error.message}</code>
