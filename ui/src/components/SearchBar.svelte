<script lang="ts">
	import { useLocation, useNavigate } from 'svelte-navigator';

	const location = useLocation();
	const navigate = useNavigate();

	let currentQuery = decodeURI($location.pathname.slice(1));

	export let query: string;
	query = currentQuery;
</script>

<div class="flex items-center justify-between">
	<h2 class="font-semibold text-xl text-slate-900 mb-2 focus:outline-none">Japanese</h2>
</div>
<form
	class="group relative"
	on:submit|preventDefault={(e) => {
		if (!(e.target instanceof HTMLFormElement)) {
			return;
		}
		const formData = new FormData(e.target);
		query = formData.get('query').toString();
		navigate(encodeURI(query));
	}}
>
	<svg
		width="20"
		height="20"
		fill="currentColor"
		class="absolute left-3 top-1/2 -mt-2.5 text-slate-400 pointer-events-none group-focus-within:text-blue-500"
		aria-hidden="true"
	>
		<path
			fill-rule="evenodd"
			clip-rule="evenodd"
			d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
		/>
	</svg>
	<input
		name="query"
		class="focus:ring-2 focus:ring-green focus:outline-none appearance-none w-full text-sm leading-6 text-slate-900 placeholder-slate-400 rounded-md py-2 pl-10 ring-1 ring-blue shadow-sm"
		type="text"
		aria-label="Enter japanese word"
		placeholder="Enter japanese word..."
		bind:value={currentQuery}
	/>
</form>
