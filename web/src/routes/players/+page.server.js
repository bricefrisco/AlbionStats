import { getApiBase } from '$lib/apiBase';

async function fetchJson(fetch, url) {
	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}
	return response.json();
}

export const load = async ({ fetch }) => {
	let topPlayers = [];
	let topPlayersError = null;

	try {
		const apiUrl = `${getApiBase()}/players/top/americas`;
		topPlayers = await fetchJson(fetch, apiUrl);
	} catch (err) {
		topPlayersError = err instanceof Error ? err.message : 'Failed to load top players';
	}

	return {
		topPlayers,
		topPlayersError
	};
};
