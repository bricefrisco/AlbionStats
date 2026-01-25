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
	const battleIds = params.battleIds;
	const validRegion = validRegions.has(region);

	let battleData = null;
	let pageError = null;
	let loading = false;

	if (!validRegion || !battleIds) {
		throw error(404, 'Invalid region or battle IDs');
	} else {
		try {
			battleData = await fetchJson(
				fetch,
				`${getApiBase()}/battles/${region}/${battleIds}`
			);
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to load battle data';
			throw error(404, message);
		}
	}

	return {
		region,
		battleIds,
		battleData,
		error: pageError,
		loading
	};
};
