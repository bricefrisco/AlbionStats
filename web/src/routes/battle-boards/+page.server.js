const apiBase = 'https://albionstats.bricefrisco.com/api';
const validTypes = new Set(['alliance', 'guild', 'player']);

function mapEntries(list = []) {
	return (list || []).map((entry) => {
		const match = entry?.match(/^(.*?)\s*\((\d+)\)$/);
		return {
			label: match ? match[1].trim() : entry,
			count: match ? match[2] : null
		};
	});
}

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
		let apiUrl;
		if (type === 'alliance' && q) {
			apiUrl = new URL(`${apiBase}/boards/alliance/${region}/${encodeURIComponent(q)}`);
			apiUrl.searchParams.set('playerCount', p || '10');
		} else if (type === 'guild' && q) {
			apiUrl = new URL(`${apiBase}/boards/guild/${region}/${encodeURIComponent(q)}`);
			apiUrl.searchParams.set('playerCount', p || '10');
		} else if (type === 'player' && q) {
			apiUrl = new URL(`${apiBase}/boards/player/${region}/${encodeURIComponent(q)}`);
			apiUrl.searchParams.set('playerCount', p || '10');
		} else {
			apiUrl = new URL(`${apiBase}/boards/${region}`);
			apiUrl.searchParams.set('totalPlayers', p || '10');
		}

		const data = await fetchJson(fetch, apiUrl.toString());
		if (Array.isArray(data)) {
			initialBattles = data.map((battle) => ({
				...battle,
				AllianceEntries: mapEntries(battle.AllianceNames),
				GuildEntries: mapEntries(battle.GuildNames)
			}));
			initialHasMore = initialBattles.length >= 20;
		}
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
