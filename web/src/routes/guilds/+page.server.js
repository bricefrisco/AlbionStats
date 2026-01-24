import { getApiBase } from '$lib/apiBase';

async function fetchJson(fetch, url) {
	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}
	return response.json();
}

export const load = async ({ fetch }) => {
	let topGuilds = [];
	let topGuildsError = null;

	try {
		const apiUrl = `${getApiBase()}/guilds/top/americas`;
		topGuilds = await fetchJson(fetch, apiUrl);
	} catch (err) {
		topGuildsError = err instanceof Error ? err.message : 'Failed to load top guilds';
	}

	return {
		topGuilds,
		topGuildsError
	};
};
