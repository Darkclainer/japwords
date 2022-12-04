<script lang="ts">
	import SearchBar from './components/SearchBar.svelte';
	import Tailwindcss from './Tailwind.svelte';
	import { Router, Route } from 'svelte-navigator';
	import { GetLemmas } from './api/generated';
	import LemmaBrowser from './components/LemmaBrowser.svelte';

	let query: string = '';
	let lemmasQuery: ReturnType<typeof GetLemmas>;
	$: if (query) {
		lemmasQuery = GetLemmas({
			variables: {
				query: query
			}
		});
	} else {
		lemmasQuery = null;
	}
</script>

<Tailwindcss />

<main class="bg-white space-y-4 p-4 md:max-w-4xl mx-auto">
	<Router>
		<Route path="/*">
			<SearchBar bind:query />

			{#if lemmasQuery}
				{#if $lemmasQuery.loading}
					<h1>Loading</h1>
				{:else if $lemmasQuery.error}
					<h1>Error</h1>
				{:else}
					<LemmaBrowser lemmas={$lemmasQuery.data.Lemmas.lemmas} />
				{/if}
			{:else}
				<h1>type kanji</h1>
			{/if}
		</Route>
	</Router>
</main>
