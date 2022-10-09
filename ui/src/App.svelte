<script lang="ts">
	import SearchBar from './components/SearchBar.svelte';
	import Tailwindcss from './Tailwind.svelte';
	import { Router, Route } from 'svelte-navigator';
	import { GetJapaneseWords } from './api/generated';

	let query: string = '';
	let wordsQuery: ReturnType<typeof GetJapaneseWords>;
	$: if (query) {
		wordsQuery = GetJapaneseWords({
			variables: {
				query: query
			}
		});
	}
</script>

<Tailwindcss />

<main class="bg-white space-y-4 p-4 md:max-w-4xl mx-auto">
	<Router>
		<Route path="/*">
			<SearchBar bind:query />

			{#if $wordsQuery}
				{JSON.stringify($wordsQuery.data?.japaneseWords)}
			{/if}
		</Route>
	</Router>
</main>
