import { getApiBase } from '$lib/apiBase';
import { buildBattleBoardsUrl, mapBattleBoardsData } from '$lib/battleBoards';
const validTypes = new Set(['alliance', 'guild', 'player']);

async function fetchJson(fetch, url) {
	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}
	return response.json();
}

export const load = async ({ url, fetch }) => {
	const q = url.searchParams.get('q') || '';
	const type = validTypes.has(url.searchParams.get('type')) ? url.searchParams.get('type') : 'alliance';
	const p = url.searchParams.get('p') || '10';
	const region = 'americas';

	let initialBattles = [];
	let initialHasMore = false;
	let initialError = null;

	try {
		const apiUrl = buildBattleBoardsUrl({
			base: getApiBase(),
			region,
			type,
			q,
			p,
			offset: 0
		});

		const data = await fetchJson(fetch, apiUrl.toString());
		initialBattles = mapBattleBoardsData(data);
		initialHasMore = initialBattles.length >= 20;
	} catch (err) {
		initialError = err instanceof Error ? err.message : 'Failed to load battles';
	}

	return {
		q,
		type,
		p,
		initialBattles,
		initialHasMore,
		initialError
	};
};
