import { error } from '@sveltejs/kit';
import { getApiBase } from '$lib/apiBase';
import { validRegions } from '$lib/utils';

async function fetchJson(fetch, url) {
	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}
	return response.json();
}

export const load = async ({ params, fetch }) => {
	const region = params.region;
	const guildName = params.name;
	let decodedName = '';
	if (guildName) {
		try {
			decodedName = decodeURIComponent(guildName);
		} catch {
			decodedName = '';
		}
	}

	const validRegion = validRegions.has(region);
	let guildData = null;
	let guildError = null;

	if (!validRegion) {
		throw error(404, 'Invalid region');
	} else if (!decodedName) {
		throw error(404, 'Guild not found');
	} else {
		try {
			const apiUrl = `${getApiBase()}/guilds/${region}/${encodeURIComponent(decodedName)}`;
			guildData = await fetchJson(fetch, apiUrl);
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to load guild data';
			if (message.includes('status: 404')) {
				throw error(404, 'Guild not found');
			}
			guildError = message;
		}
	}

	return {
		region,
		validRegion,
		guildData,
		guildError
	};
};
