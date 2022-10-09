<script lang="ts">
	import SearchBar from './components/SearchBar.svelte';
	import Tailwindcss from './Tailwind.svelte';
	import { Router, Route } from 'svelte-navigator';
	import { GetJapaneseWords } from './api/generated';
	import WordBrowser from './components/WordBrowser.svelte';

	let query: string = '';
	let wordsQuery: ReturnType<typeof GetJapaneseWords>;
	$: if (query) {
		wordsQuery = GetJapaneseWords({
			variables: {
				query: query
			}
		});
	} else {
		wordsQuery = null;
	}
</script>

<Tailwindcss />

<main class="bg-white space-y-4 p-4 md:max-w-4xl mx-auto">
	<Router>
		<Route path="/*">
			<SearchBar bind:query />

			{#if wordsQuery}
				{#if $wordsQuery.loading}
					<h1>Loading</h1>
				{:else if $wordsQuery.error}
					<h1>Error</h1>
				{:else}
					<WordBrowser words={$wordsQuery.data.japaneseWords?.words} />
				{/if}
			{:else}
				<h1>type kanji</h1>
			{/if}
		</Route>
	</Router>
</main>
