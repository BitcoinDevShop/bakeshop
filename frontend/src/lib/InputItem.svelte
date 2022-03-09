<script lang="ts">
	export let error;
	export let value;
	export let label;
	export let id;
	export let placeholder = null;
	export let type;
</script>

<div class="flex flex-col space-y-2 items-start">
	<div class="pt-4 flex space-x-2">
		<label class="font-bold" for={id}>{label}</label>
		{#if error}
			<div class="text-red-400 transition">
				{error}
			</div>
		{/if}
	</div>

	{#if type === 'text'}
		<input
			{id}
			bind:value
			type="text"
			{placeholder}
			class="text-xl text-black rounded transition hover:drop-shadow-[0px_0px_5px_rgba(192,87,70,0.75)]"
		/>
	{:else}
		<input
			{id}
			bind:value
			type="number"
			{placeholder}
			class="text-xl text-black rounded transition hover:drop-shadow-[0px_0px_5px_rgba(192,87,70,0.75)]"
		/>
	{/if}
</div>
