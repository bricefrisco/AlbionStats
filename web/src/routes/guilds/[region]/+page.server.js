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

	let topGuilds = [];
	let topGuildsError = null;

	if (!validRegion) {
		throw error(404, 'Invalid region');
	} else {
		try {
			const apiUrl = `${getApiBase()}/guilds/top/${region}`;
			topGuilds = await fetchJson(fetch, apiUrl);
		} catch (err) {
			topGuildsError = err instanceof Error ? err.message : 'Failed to load top guilds';
		}
	}

	return {
		region,
		topGuilds,
		topGuildsError
	};
};
