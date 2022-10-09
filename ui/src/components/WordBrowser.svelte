<script lang="ts">
	import type { JapaneseWord } from '../api/generated';
	import WordCard from './WordCard.svelte';
	export let words: Array<JapaneseWord> = [];
</script>

<div>
	{#each words as word}
		<WordCard {word} />
	{/each}
</div>
