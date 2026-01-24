import { getApiBase } from '$lib/apiBase';

async function fetchJson(fetch, url) {
	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}
	return response.json();
}

export const load = async ({ fetch }) => {
	let topAlliances = [];
	let topAlliancesError = null;

	try {
		const apiUrl = `${getApiBase()}/alliances/top/americas`;
		topAlliances = await fetchJson(fetch, apiUrl);
	} catch (err) {
		topAlliancesError = err instanceof Error ? err.message : 'Failed to load top alliances';
	}

	return {
		topAlliances,
		topAlliancesError
	};
};
