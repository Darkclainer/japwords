<script lang="ts">
	import { getNotificationsContext } from 'svelte-notifications';
	import { Lemma, PitchType } from '../api/generated';
	import PlayIcon from './PlayIcon.svelte';

	export let lemma: Lemma;

	const { addNotification } = getNotificationsContext();
	const copy_kanji = () => {
		navigator.clipboard.writeText(lemma.slug.word);
		addNotification({
			text: 'Copied ' + lemma.slug.word,
			position: 'bottom-right',
			removeAfter: 2000
		});
	};

	let audio: HTMLAudioElement;
	const play_mouse_up = () => {
		if (!audio) {
			return;
		}
		audio.play();
	};
</script>

<div class="flex justify-between shadow-md rounded-md bg-gray my-4 px-8 pt-10 pb-7">
	<div class="grow divide-y divide-blue flex-1">
		<div>
			<button
				on:mouseup={copy_kanji}
				class="hover:text-blue active:text-dark-blue transition-colors duration-300"
			>
				<h1 class="text-5xl pb-4">
					{#if lemma.slug.furigana.length == 0}
						{lemma.slug.word}
					{:else}
						<ruby>
							{#each lemma.slug.furigana as furigana}
								{#if furigana.kanji}
									{furigana.kanji}<rp>[</rp><rt class="text-dark-gray text-xl"
										>{furigana.hiragana}</rt
									><rp>]</rp>
								{:else}
									{furigana.hiragana}<rp>[</rp><rt /><rp>]</rp>
								{/if}
							{/each}
						</ruby>
					{/if}
				</h1>
			</button>
		</div>
		<div class="text-xl py-4">
			{#if lemma.slug.pitch.length == 0}
				{#each lemma.slug.furigana as furigana}
					{furigana.hiragana}
				{/each}
			{:else}
				{#each lemma.slug.pitch as pitch}
					<span
						class:border-t={pitch.pitch.includes(PitchType.Up)}
						class:border-b={pitch.pitch.includes(PitchType.Down)}
						class:border-r={pitch.pitch.includes(PitchType.Right)}
						class="border-blue">{pitch.hiragana}</span
					>
				{/each}
			{/if}
		</div>
		{#if lemma.senses.length != 0}
			<div class="text-xl py-4">
				{#each lemma.senses as sense, i}
					<div class="mb-4 last:mb-0">
						<p class="text-base text-blue">{sense.partOfSpeech.join(', ')}</p>
						<p>{i + 1}. {sense.definition.join('; ')}</p>
						<p class="text-base text-dark-gray">{sense.tags.join(', ')}</p>
					</div>
				{/each}
			</div>
		{/if}
		{#if lemma.forms.length != 0}
			<div class="text-xl pt-4">
				<p class="text-base text-blue">Other Forms</p>
				<div class="flex flex-col">
					{#each lemma.forms as form}
						<p>
							{form.word}
							{#if form.hiragana}「{form.hiragana}」{/if}
						</p>
					{/each}
				</div>
			</div>
		{/if}
	</div>
	<div class="shrink basis-14" />
	<div class="flex-none flex flex-col justify-items-stretch basis-40">
		<button
			class="px-2 py-4 text-xl rounded-md bg-blue hover:bg-green active:bg-dark-green transition-colors duration-300 text-white"
			>Add to Anki</button
		>
		<div class="divide-y divide-blue">
			{#if lemma.tags.length != 0}
				<div class="flex flex-col justify-items-center gap-4 py-4 text-lg text-blue">
					{#each lemma.tags as tag}
						<p class="text-center">{tag}</p>
					{/each}
				</div>
			{/if}
			{#if lemma.audio.length != 0}
				<div class="flex flex-row pt-4">
					<button class="m-2" on:mouseup={play_mouse_up}>
						<PlayIcon />
					</button>
					<p class="m-1 text-dark-gray text-sm leading-5">Listen to pronunciation</p>
					<audio bind:this={audio}>
						{#each lemma.audio as audio}
							<source src={audio.source} type={audio.type} />
						{/each}
					</audio>
				</div>
			{/if}
		</div>
	</div>
</div>
