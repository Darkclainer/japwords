<script lang="ts">
	import type { Lemma } from '../api/generated';
	import LemmaCard from './LemmaCard.svelte';
	export let lemmas: Array<Lemma> = [];
</script>

<div>
	{#each lemmas as lemma}
		<LemmaCard {lemma} />
	{/each}
</div>
