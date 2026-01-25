function buildUrlset(origin) {
	const lastmod = new Date().toISOString();
	const urls = [
		`${origin}/`
	];

	const entries = urls
		.map(
			(loc) =>
				`  <url>\n    <loc>${loc}</loc>\n    <lastmod>${lastmod}</lastmod>\n  </url>`
		)
		.join('\n');

	return `<?xml version="1.0" encoding="UTF-8"?>\n` +
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n` +
		`${entries}\n` +
		`</urlset>\n`;
}

export const GET = ({ url }) => {
	const xml = buildUrlset(url.origin);
	return new Response(xml, {
		headers: {
			'Content-Type': 'application/xml; charset=utf-8',
			'Cache-Control': 'public, max-age=3600'
		}
	});
};
