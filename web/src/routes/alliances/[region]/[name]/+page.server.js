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
		allianceError = 'Invalid region';
	} else if (!decodedName) {
		allianceError = 'Alliance not found';
	} else {
		try {
			const apiUrl = `${getApiBase()}/alliances/${region}/${encodeURIComponent(decodedName)}`;
			allianceData = await fetchJson(fetch, apiUrl);
		} catch (err) {
			allianceError = err instanceof Error ? err.message : 'Failed to load alliance data';
		}
	}

	return {
		region,
		decodedName,
		validRegion,
		allianceData,
		allianceError
	};
};
