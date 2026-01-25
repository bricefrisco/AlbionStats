import { error } from '@sveltejs/kit';
import { getApiBase } from '$lib/apiBase.js';
import { validRegions } from '$lib/utils';

async function fetchJson(fetch, url) {
	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}
	return response.json();
}

export const load = async ({ fetch, params }) => {
	const region = params.region;
	const validRegion = validRegions.has(region);

	let topPlayers = [];
	let topPlayersError = null;

	if (!validRegion) {
		throw error(404, 'Invalid region');
	} else {
		try {
			const apiUrl = `${getApiBase()}/players/top/${region}`;
			topPlayers = await fetchJson(fetch, apiUrl);
		} catch (err) {
			topPlayersError = err instanceof Error ? err.message : 'Failed to load top players';
		}
	}

	return {
		region,
		topPlayers,
		topPlayersError
	};
};
