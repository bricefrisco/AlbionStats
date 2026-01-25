import { error } from '@sveltejs/kit';
import { getApiBase } from '$lib/apiBase';
import { validRegions } from '$lib/utils';

export const load = async ({ params, fetch }) => {
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
		throw error(404, 'Invalid region');
	} else if (!decodedName) {
		throw error(404, 'Player not found');
	} else {
		try {
			const response = await fetch(
				`${getApiBase()}/players/${region}/${encodeURIComponent(decodedName)}`
			);

			if (response.status === 404) {
				throw error(404, 'Player not found');
			} else if (!response.ok) {
				throw error(500, `HTTP error! status: ${response.status}`);
			} else {
				const payload = await response.json();
				playerData = payload?.Player || null;

				if (payload) {
					metrics = {
						pvp: {
							data: {
								timestamps: payload.Timestamps || [],
								kill_fame: payload.Pvp?.KillFame || [],
								death_fame: payload.Pvp?.DeathFame || [],
								fame_ratio: payload.Pvp?.FameRatio || []
							},
							error: null
						},
						pve: {
							data: {
								timestamps: payload.Timestamps || [],
								total: payload.Pve?.Total || [],
								royal: payload.Pve?.Royal || [],
								outlands: payload.Pve?.Outlands || [],
								avalon: payload.Pve?.Avalon || [],
								hellgate: payload.Pve?.Hellgate || [],
								corrupted: payload.Pve?.Corrupted || [],
								mists: payload.Pve?.Mists || []
							},
							error: null
						},
						gathering: {
							data: {
								timestamps: payload.Timestamps || [],
								total: payload.Gathering?.Total || [],
								royal: payload.Gathering?.Royal || [],
								outlands: payload.Gathering?.Outlands || [],
								avalon: payload.Gathering?.Avalon || []
							},
							error: null
						},
						crafting: {
							data: {
								timestamps: payload.Timestamps || [],
								total: payload.Crafting?.Total || []
							},
							error: null
						}
					};
				}
			}
		} catch (err) {
			if (err?.status === 404) {
				throw err;
			}
			playerError = err instanceof Error ? err.message : 'Failed to load player data';
		}
	}

	return {
		region,
		decodedName,
		validRegion,
		playerData,
		playerError,
		loading,
		metrics
	};
};
