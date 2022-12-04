<script>
	import { fade } from 'svelte/transition';

	export let notification = {};
	// `onRemove` function will be passed into your component.
	export let onRemove = null;

	const handleButtonClick = () => {
		onRemove();
	};
</script>

<div class="p-4 m-4 rounded-lg bg-green text-white text-lg" in:fade out:fade>
	<button on:click={handleButtonClick}>
		<h4>{notification.text}</h4>
	</button>
</div>
