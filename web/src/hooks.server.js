const RATE_LIMIT = 4;
const WINDOW_MS = 1000;

// ip -> { count, windowStart }
const ipBuckets = new Map();

export const handle = async ({ event, resolve }) => {
	const now = Date.now();

	const clientIp =
		event.request.headers.get('cf-connecting-ip') || event.getClientAddress?.() || 'unknown';

	const bucket = ipBuckets.get(clientIp);

	if (!bucket || now - bucket.windowStart >= WINDOW_MS) {
		ipBuckets.set(clientIp, { count: 1, windowStart: now });
		return await resolve(event);
	}

	if (bucket.count >= RATE_LIMIT) {
		return new Response('Too Many Requests', {
			status: 429,
			headers: {
				'Content-Type': 'text/plain',
				'Retry-After': Math.ceil(WINDOW_MS / 1000).toString()
			}
		});
	}

	bucket.count++;
	return await resolve(event);
};
