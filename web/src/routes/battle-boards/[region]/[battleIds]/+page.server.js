import { getApiBase } from '$lib/apiBase';
const validRegions = new Set(['americas', 'europe', 'asia']);

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
	let error = null;
	let loading = false;

	if (!validRegion || !battleIds) {
		error = 'Invalid region or battle IDs';
	} else {
		try {
			battleData = await fetchJson(
				fetch,
				`${getApiBase()}/battles/${region}/${battleIds}`
			);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load battle data';
		}
	}

	return {
		region,
		battleIds,
		battleData,
		error,
		loading
	};
};
