<script context="module">
	/** @type {import('@sveltejs/kit').Load} */
	export async function load({ params, fetch, session, stuff }) {
		const url = `http://localhost:8080/list`;
		const res = await fetch(url);

		const list = await res.json();

		if (res.ok) {
			return {
				props: {
					list
				}
			};
		}

		return {
			status: res.status,
			error: new Error(`Could not load ${url}`)
		};
	}
</script>

<script lang="ts">
	import { goto } from '$app/navigation';

	export let list;

	import Button from '$lib/Button.svelte';
</script>

<p class="mb-4">
	This is an easy way to <em class="text-red-400">dangerously </em>and
	<em class="text-red-400">recklessly</em> send subscription payments over Lightning.
</p>
<p class="mb-4">
	Bake Shop prepares a special Macaroon which you can give to the person you'd like to pay on a
	recurring basis. The Macaroon embeds the time interval, the amount, and the number of times this
	subscription is valid for.
</p>
<p class="mb-4">
	Equipped with this Macaroon, and the knowledge of how to connect to your Lightning node, your
	counterparty can dial into your node and command it to pay invoices it supplies.
</p>
<p class="mb-4">
	If everything goes accoring to plan, Bake Shop's Macaroon interceptor will validate every command
	from your counterparty and verify that it's within the scope of contract you set up
	<em> (e.g. only once a month, only for 2,000 sats, etc.). </em>
</p>
<p class="mb-4">Does that sound like a good idea?</p>
<Button
	onClick={() => {
		console.log('heyo');
		goto('/bake');
	}}>Begin!</Button
>

{#if list?.items?.length > 0}
	<h1 class="mt-4">Previous Bakes</h1>

	{#each list.items as mac}
		<div class="rounded-xl border p-4 space-y-4 flex flex-col mb-4">
			<h2>{mac.name}</h2>
			<ul>
				<li>
					<h3>Created</h3>
					<h4>{new Date(mac.created_at).toLocaleString()}</h4>
				</li>
				<li>
					<h3>Status</h3>
					<h4>{mac.status}</h4>
				</li>
			</ul>
			<Button primary onClick={() => goto(`/mac/${mac.id}`)}>View</Button>
		</div>
		<!-- 
		<code class="break-all">
			{JSON.stringify(mac, null, 2)}
		</code> -->
	{/each}
{/if}

<style lang="postcss">
	h3 {
		@apply mr-4 w-1/5;
	}

	h4 {
		@apply w-4/5;
	}

	li {
		@apply flex space-x-2 items-start p-1 py-2;
	}

	li:nth-child(odd) {
		@apply rounded;
		background-color: rgba(255, 255, 255, 0.1);
	}
</style>
