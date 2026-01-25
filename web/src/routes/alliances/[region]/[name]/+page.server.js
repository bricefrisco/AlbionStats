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
	const allianceName = params.name;
	let decodedName = '';
	if (allianceName) {
		try {
			decodedName = decodeURIComponent(allianceName);
		} catch {
			decodedName = '';
		}
	}

	const validRegion = validRegions.has(region);
	let allianceData = null;
	let allianceError = null;

	if (!validRegion) {
		throw error(404, 'Invalid region');
	} else if (!decodedName) {
		throw error(404, 'Alliance not found');
	} else {
		try {
			const apiUrl = `${getApiBase()}/alliances/${region}/${encodeURIComponent(decodedName)}`;
			allianceData = await fetchJson(fetch, apiUrl);
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to load alliance data';
			if (message.includes('status: 404')) {
				throw error(404, 'Alliance not found');
			}
			allianceError = message;
		}
	}

	return {
		region,
		validRegion,
		allianceData,
		allianceError
	};
};
