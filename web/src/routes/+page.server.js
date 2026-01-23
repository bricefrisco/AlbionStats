import { getApiBase } from '$lib/apiBase';

async function fetchJson(fetch, url) {
	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}
	return response.json();
}

export const load = async ({ fetch }) => {
	let playersTracked = { timestamps: [], values: [], error: null };
	let totalDataPoints = { timestamps: [], values: [], error: null };

	const [playersResult, snapshotsResult] = await Promise.allSettled([
		fetchJson(fetch, `${getApiBase()}/metrics/players_total?granularity=1w`),
		fetchJson(fetch, `${getApiBase()}/metrics/snapshots?granularity=1w`)
	]);

	if (playersResult.status === 'fulfilled') {
		playersTracked = {
			timestamps: playersResult.value.timestamps || [],
			values: playersResult.value.values || [],
			error: null
		};
	} else {
		playersTracked.error = playersResult.reason?.message || 'Failed to load player metrics';
	}

	if (snapshotsResult.status === 'fulfilled') {
		totalDataPoints = {
			timestamps: snapshotsResult.value.timestamps || [],
			values: (snapshotsResult.value.values || []).map((value) => value * 40),
			error: null
		};
	} else {
		totalDataPoints.error =
			snapshotsResult.reason?.message || 'Failed to load snapshots metrics';
	}

	return { playersTracked, totalDataPoints };
};
