<script context="module">
	/** @type {import('@sveltejs/kit').Load} */
	export async function load({ params, fetch, session, stuff }) {
		// const url = `/endpoints/details/${params.uuid}.json`;
		const url = `http://localhost:8080/details/${params.uuid}`;
		const res = await fetch(url);

		const mac = await res.json();

		if (res.ok) {
			return {
				props: {
					mac
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
	import Button from '$lib/Button.svelte';
	import * as macaroon from 'macaroon';
	import { Buffer } from 'buffer';

	export let mac;

	// TODO display errors
	let error = '';

	const encodedMacaroon = mac.macaroon;
	let decodedMac;
	try {
		// const Buffer = bitcoin.Buffer;
		const buffer = Buffer.from(encodedMacaroon.replace(/\s*/gi, ''), 'hex');
		decodedMac = macaroon.importMacaroon(buffer).exportJSON();
		console.debug(decodedMac);
	} catch (e) {
		console.error(e);
	}

	// 'lnd-custom subscribe ms:3600 amount:21000 times:21'
	function parseCaveats(caveats: string) {
		const arr = caveats.split(' ');
		return {
			interval: arr[2].split(':')[1],
			amount: arr[3].split(':')[1],
			times: arr[4].split(':')[1]
		};
	}

	$: caveats = decodedMac?.c && decodedMac.c.length === 1 ? parseCaveats(decodedMac?.c[0].i) : null;

	async function copyToClipboard(text: string) {
		try {
			await navigator.clipboard.writeText(text);
		} catch (e) {
			console.error(e);
		}
	}

	function downloadTlsCert(cert: string) {
		const link = document.createElement('a');
		link.href = `data:application/octet-stream;base64,${cert}`;
		link.download = `tls.cert`;
		link.click();
	}

	let show = false;

	async function cancel() {
		try {
			const res = await fetch('http://localhost:8080/cancel', {
				method: 'POST',
				body: JSON.stringify({ id: mac.id }),
				headers: {
					'Content-Type': 'application/json'
				}
			});

			// If we cancelled correctly we should reload the macaroon info
			if (res.ok) {
				const url = `http://localhost:8080/details/${mac.id}`;
				const res = await fetch(url);
				mac = await res.json();
			}
		} catch (e) {
			error = 'Something went wrong';
			console.error(e);
		}
	}
</script>

<div class="rounded-xl border p-4 space-y-4 flex flex-col">
	<h2>{mac.name}</h2>
	<ul>
		<li>
			<h3>Created</h3>
			<h4>{new Date(mac.created_at).toLocaleString()}</h4>
		</li>
		<li>
			<h3>Uuid</h3>
			<div class="w-4/5 space-y-2">
				<h4 class="break-all text-sm">{mac.id}</h4>
				<Button small onClick={() => copyToClipboard(mac.id)}>Copy</Button>
			</div>
		</li>
		<li>
			<h3>Macaroon</h3>
			<div class="w-4/5 space-y-2">
				<h4 class="break-all text-sm">{mac.macaroon}</h4>
				<Button small onClick={() => copyToClipboard(mac.macaroon)}>Copy</Button>
			</div>
		</li>
		<li>
			<h3>Status</h3>
			<h4>{mac.status}</h4>
		</li>
		{#if caveats}
			<li>
				<h3>Amount</h3>
				<h4>{caveats.amount} sats</h4>
			</li>
			<li>
				<h3>Times</h3>
				<h4>{caveats.times}x</h4>
			</li>
			<li>
				<h3>Interval</h3>
				<h4>Every {caveats.interval} milliseconds</h4>
			</li>
		{/if}
	</ul>
	{#if mac.status === 'active'}
		<Button primary onClick={cancel}>Cancel</Button>
	{/if}
	<Button secondary onClick={() => (show = !show)}>{show ? 'Hide' : 'Show'} the scary stuff</Button>
	{#if show}
		<figure class="flex flex-col space-y-2">
			<h3 class="w-full whitespace-nowrap">Connect String</h3>
			<code class="text-xs break-all">
				{mac.lndConnect}
			</code>
			<div class="flex items-center space-x-4">
				<Button small onClick={() => copyToClipboard(mac.lndConnect)}>Copy Lndconnect</Button>
				<!-- TODO do I even need "download tls" ? -->
				<!-- <Button small onClick={() => downloadTlsCert('abcdefg')}>Download TLS</Button> -->
			</div>
			<figcaption class="text-xs">
				^ The lndconnect string and TLS cert are what you give to the payee. This gives them access
				to your node under the conditions of this Macaroon. Copy and paste responsibly!
			</figcaption>
		</figure>
	{/if}
</div>

<code class="break-all">
	{JSON.stringify(mac, null, 2)}
	{JSON.stringify(decodedMac.c[0], null, 2)}
</code>

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
