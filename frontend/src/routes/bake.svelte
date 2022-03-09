<script lang="ts">
	import { goto } from '$app/navigation';
	import prettyPrintInterval from '$lib/utility/prettyPrint';

	import Button from '$lib/Button.svelte';
	import InputItem from '$lib/InputItem.svelte';

	let interval = 3600;
	let name = '';
	let times = 21;
	let amount = 21000;

	let error = '';

	let errors = {};

	// function isNumber(str) {
	//     var pattern = /^\d+$/;
	//     return pattern.test(str);  // returns a boolean
	// }

	async function handleSubmit(e) {
		e.preventDefault();

		error = '';
		errors = {};

		const data = {
			// 1000 ms in a sec
			interval: interval * 1000,
			name,
			times,
			amount
		};

		if (!data.name) {
			errors['name'] = "The name can't be blank";
		}
		if (data.interval < 0) {
			errors['interval'] = "The interval can't be negative nice try";
		}
		if (data.interval === 0 && times !== 1) {
			errors['interval'] = "You can't have an interval of zero if you're paying more than once";
		}

		if (data.times < 0) {
			errors['times'] = "Times can't be negative nice try";
		}
		if (data.amount < 1) {
			errors['amount'] = 'The minimum is one sat';
		}

		// Make sure the numbers are numbers
		if (typeof data.interval !== 'number') {
			errors['interval'] = 'Interval should be a number I bet';
		}
		if (typeof data.times !== 'number') {
			errors['times'] = 'Times should be a number I bet';
		}
		if (typeof data.amount !== 'number') {
			errors['amount'] = 'Interval should be a number I bet';
		}

		if (Object.keys(errors).length !== 0) {
			return;
		}

		try {
			const res = await fetch('http://localhost:8080/bake', {
				method: 'POST',
				body: JSON.stringify(data),
				headers: {
					'Content-Type': 'application/json'
				}
			});

			const json = await res.json();
			if (json.id) {
				goto(`/mac/${json.id}`);
				// console.debug(json);
			} else {
				error = 'Something went wrong';
			}
		} catch (e) {
			error = 'Something went wrong';
			console.error(e);
		}
	}

	// $: prettyInterval = prettyPrintInterval(interval);
</script>

<h2 class="text-2xl font-bold mb-4">Step One: What would you like to bake?</h2>
<form
	class="flex flex-col space-y-2 text-gray-300"
	on:submit|preventDefault={handleSubmit}
	action="/endpoints/bake"
	method="POST"
>
	<InputItem
		error={errors['name']}
		label="Name"
		type="text"
		bind:value={name}
		id="name"
		placeholder="name-of-macaroon"
	/>
	<p>
		What do you want to call your Macaroon? This is so you can keep track of this subscription
		later.
	</p>

	<InputItem
		error={errors['interval']}
		label="Interval (seconds)"
		type="number"
		bind:value={interval}
		id="interval"
		placeholder={3600}
	/>
	<p>How often do you want to pay?</p>
	<p>
		Right now it's set to recur every {prettyPrintInterval(interval)}.
	</p>

	<InputItem
		error={errors['times']}
		label="Times"
		type="number"
		bind:value={times}
		id="times"
		placeholder={100}
	/>
	<p>How many times would you like to pay?</p>
	<p>
		With an interval of {interval ? interval.toLocaleString() : '0'} seconds you'll be subscribed for
		{times ? prettyPrintInterval(interval * times) : 'forever'}.
	</p>
	<p>Set this to zero to pay indefinitely.</p>

	<InputItem
		error={errors['amount']}
		label="Amount (sats)"
		type="number"
		bind:value={amount}
		id="amount"
		placeholder={21000}
	/>
	<p>How many sats do you want to pay each time?</p>

	<div class="flex items-center w-full space-x-4 py-4">
		<Button onClick={handleSubmit}>Let's Bake!</Button>
		<Button onClick={() => goto('/')} secondary>Nevermind</Button>
		{#if error}
			<div class="text-red-400">
				{error}
			</div>
		{/if}
	</div>
</form>
