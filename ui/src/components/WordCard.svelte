<script lang="ts">
	import type { JapaneseWord } from '../api/generated';
	export let word: JapaneseWord;
</script>

<div class="flex shadow-md rounded-md bg-sky-50 mb-10 p-4">
	<div class="divide-y-2 flex-1">
		<div>
			<h1 class="text-4xl">
				{word.kanji}
				<ruby>
					明<rp>[</rp><rt>あし</rt><rp>]</rp>日<rp>[</rp><rt>たこ</rt><rp>[</rp>
				</ruby>
			</h1>
		</div>
		<div>
			<p>{word.meaning}</p>
		</div>
		<div />
	</div>
	<div class="flex-initial">
		<button class="p-2 bg-blue-100"> Add to Anki </button>
	</div>
</div>

<style>
</style>
