const SITEMAPS = ['sitemap-static.xml', 'sitemap-sections.xml'];

function buildSitemapIndex(origin) {
	const lastmod = new Date().toISOString();
	const entries = SITEMAPS.map(
		(path) =>
			`  <sitemap>\n    <loc>${origin}/${path}</loc>\n    <lastmod>${lastmod}</lastmod>\n  </sitemap>`
	).join('\n');

	return `<?xml version="1.0" encoding="UTF-8"?>\n` +
		`<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n` +
		`${entries}\n` +
		`</sitemapindex>\n`;
}

export const GET = ({ url }) => {
	const xml = buildSitemapIndex(url.origin);
	return new Response(xml, {
		headers: {
			'Content-Type': 'application/xml; charset=utf-8',
			'Cache-Control': 'public, max-age=3600'
		}
	});
};
