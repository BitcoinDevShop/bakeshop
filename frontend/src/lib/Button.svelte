<script lang="ts">
	export let onClick;
	export let secondary = false;
	export let error = false;
	export let small = false;
</script>

<button class="bg-[#9DD1F1]" class:secondary class:error class:small on:click={onClick} {...$$props}
	><slot /></button
>

<style lang="postcss">
	.small {
		@apply text-base py-2 rounded-sm;
	}
	.secondary {
		@apply bg-red-400;
	}

	.error {
		@apply bg-[#0C1713] text-white;
	}
	button {
		@apply text-black rounded shadow transition;
		@apply text-xl font-bold p-4;
	}
	button:hover {
		@apply drop-shadow-[0px_0px_5px_rgba(192,87,70,0.75)];
	}
	button:active {
		opacity: 0.75;
	}
</style>
