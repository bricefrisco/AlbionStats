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

	let topAlliances = [];
	let topAlliancesError = null;

	if (!validRegion) {
		throw error(404, 'Invalid region');
	} else {
		try {
			const apiUrl = `${getApiBase()}/alliances/top/${region}`;
			topAlliances = await fetchJson(fetch, apiUrl);
		} catch (err) {
			topAlliancesError = err instanceof Error ? err.message : 'Failed to load top alliances';
		}
	}

	return {
		region,
		topAlliances,
		topAlliancesError
	};
};
