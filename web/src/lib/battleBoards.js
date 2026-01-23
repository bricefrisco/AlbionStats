export function mapEntries(list = []) {
	return (list || []).map((entry) => {
		const match = entry?.match(/^(.*?)\s*\((\d+)\)$/);
		return {
			label: match ? match[1].trim() : entry,
			count: match ? match[2] : null
		};
	});
}

export function mapBattleBoardsData(data) {
	if (!Array.isArray(data)) return [];
	return data.map((battle) => ({
		...battle,
		AllianceEntries: mapEntries(battle.AllianceNames),
		GuildEntries: mapEntries(battle.GuildNames)
	}));
}

export function buildBattleBoardsUrl({ base, region, type, q, p, offset }) {
	let url;
	const query = q || '';

	if (type === 'alliance' && query) {
		url = new URL(`${base}/boards/alliance/${region}/${encodeURIComponent(query)}`);
		url.searchParams.set('playerCount', p || '10');
	} else if (type === 'guild' && query) {
		url = new URL(`${base}/boards/guild/${region}/${encodeURIComponent(query)}`);
		url.searchParams.set('playerCount', p || '10');
	} else if (type === 'player' && query) {
		url = new URL(`${base}/boards/player/${region}/${encodeURIComponent(query)}`);
		url.searchParams.set('playerCount', p || '10');
	} else {
		url = new URL(`${base}/boards/${region}`);
		url.searchParams.set('totalPlayers', p || '10');
	}

	if (offset > 0) {
		url.searchParams.set('offset', offset.toString());
	}

	return url;
}
