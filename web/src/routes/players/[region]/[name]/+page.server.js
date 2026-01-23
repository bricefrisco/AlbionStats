import { getApiBase } from '$lib/apiBase';
const validRegions = new Set(['americas', 'europe', 'asia']);

async function fetchJson(fetch, url) {
	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}
	return response.json();
}

async function fetchMetric(fetch, path) {
	try {
		const data = await fetchJson(fetch, `${getApiBase()}${path}`);
		return { data, error: null };
	} catch (err) {
		return { data: null, error: err instanceof Error ? err.message : 'Failed to load metric' };
	}
}

export const load = async ({ params, url, fetch }) => {
	const region = params.region;
	const playerName = params.name;
	let decodedName = '';
	if (playerName) {
		try {
			decodedName = decodeURIComponent(playerName);
		} catch {
			decodedName = '';
		}
	}
	const validRegion = validRegions.has(region);
	const activeTab = url.searchParams.get('tab') || 'pvp';

	let playerData = null;
	let playerError = null;
	let loading = false;
	let metrics = {
		pvp: { data: null, error: null },
		pve: { data: null, error: null },
		gathering: { data: null, error: null },
		crafting: { data: null, error: null }
	};

	if (!validRegion) {
		playerError = 'Invalid region';
	} else if (!decodedName) {
		playerError = 'Player not found';
	} else {
		try {
	const response = await fetch(
		`${getApiBase()}/players/${region}/${encodeURIComponent(decodedName)}`
	);

			if (response.status === 404) {
				playerError = 'Player not found';
			} else if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			} else {
				playerData = await response.json();
			}
		} catch (err) {
			playerError = err instanceof Error ? err.message : 'Failed to load player data';
		}

		if (playerData?.PlayerID) {
			const [pvp, pve, gathering, crafting] = await Promise.all([
				fetchMetric(fetch, `/metrics/pvp/${region}/${playerData.PlayerID}`),
				fetchMetric(fetch, `/metrics/pve/${region}/${playerData.PlayerID}`),
				fetchMetric(fetch, `/metrics/gathering/${region}/${playerData.PlayerID}`),
				fetchMetric(fetch, `/metrics/crafting/${region}/${playerData.PlayerID}`)
			]);

			metrics = { pvp, pve, gathering, crafting };
		}
	}

	return {
		region,
		decodedName,
		validRegion,
		playerData,
		playerError,
		loading,
		activeTab,
		metrics
	};
};
