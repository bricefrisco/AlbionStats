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
	let dailyActiveUsers = { timestamps: [], americas: [], europe: [], asia: [], error: null };

	const [playersResult, snapshotsResult, dauResult] = await Promise.allSettled([
		fetchJson(fetch, `${getApiBase()}/metrics/players_total?granularity=1w`),
		fetchJson(fetch, `${getApiBase()}/metrics/snapshots?granularity=1w`),
		fetchJson(fetch, `${getApiBase()}/metrics/dau`)
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

	if (dauResult.status === 'fulfilled') {
		dailyActiveUsers = {
			timestamps: dauResult.value.timestamps || [],
			americas: dauResult.value.americas || [],
			europe: dauResult.value.europe || [],
			asia: dauResult.value.asia || [],
			error: null
		};
	} else {
		dailyActiveUsers.error = dauResult.reason?.message || 'Failed to load DAU metrics';
	}

	return { playersTracked, totalDataPoints, dailyActiveUsers };
};
