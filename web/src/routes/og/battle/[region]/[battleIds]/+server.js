import { readFile } from 'node:fs/promises';
import path from 'node:path';
import satori from 'satori';
import { Resvg } from '@resvg/resvg-js';
import { getApiBase } from '$lib/apiBase.js';
import { formatFame } from '$lib/utils.js';

const OG_WIDTH = 1200;
const OG_HEIGHT = 630;

function truncateName(value, maxLength = 22) {
	if (!value) return '-';
	if (value.length <= maxLength) return value;
	return `${value.slice(0, Math.max(0, maxLength - 3))}...`;
}

function formatValue(value) {
	if (typeof value === 'number') return value.toLocaleString();
	return value ?? '-';
}

function buildMarkup({ alliances, guilds }) {
	const hasAlliances = alliances.length > 0;
	const hasGuilds = guilds.length > 0;
	const showFallback = !hasAlliances && !hasGuilds;
	return {
		type: 'div',
		props: {
			style: {
				width: '100%',
				height: '100%',
				display: 'flex',
				flexDirection: 'column',
				justifyContent: 'flex-start',
				padding: '0px',
				backgroundColor: '#0f172a',
				color: '#e2e8f0',
				fontFamily: 'Inter'
			},
			children: [
				{
					type: 'div',
					props: {
						style: {
							display: 'flex',
							flexDirection: 'column',
							gap: '40px',
							alignItems: 'stretch',
							width: '100%',
							height: '100%',
							padding: '20px'
						},
						children: [
							hasAlliances
								? {
										type: 'div',
										props: {
											style: {
												display: 'flex',
												flexDirection: 'column',
												gap: '12px',
												backgroundColor: '#0f172a',
												width: '100%'
											},
											children: [
												{
													type: 'div',
													props: {
														style: {
															display: 'flex',
															justifyContent: 'space-between',
															fontSize: '28px',
															textTransform: 'uppercase',
															letterSpacing: '0.08em',
															color: '#94a3b8'
														},
														children: [
															{
																type: 'div',
																props: { style: { width: '36%' }, children: 'Alliance' }
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#94a3b8' },
																	children: 'Players'
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#94a3b8' },
																	children: 'Kills'
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#94a3b8' },
																	children: 'Deaths'
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#94a3b8' },
																	children: 'Kill Fame'
																}
															}
														]
													}
												},
												...alliances.map((alliance) => ({
													type: 'div',
													props: {
														style: {
															display: 'flex',
															justifyContent: 'space-between',
															alignItems: 'center',
															fontSize: '36px',
															color: '#e2e8f0'
														},
														children: [
															{
																type: 'div',
																props: {
																	style: { width: '36%', fontWeight: 600 },
																	children: truncateName(alliance.AllianceName)
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#93c5fd' },
																	children: formatValue(alliance.PlayerCount)
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#fca5a5' },
																	children: formatValue(alliance.Kills)
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#f0abfc' },
																	children: formatValue(alliance.Deaths)
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#facc15' },
																	children: formatFame(alliance.KillFame)
																}
															}
														]
													}
												}))
											]
										}
									}
								: null,
							hasGuilds
								? {
										type: 'div',
										props: {
											style: {
												display: 'flex',
												flexDirection: 'column',
												gap: '12px',
												backgroundColor: '#0f172a',
												width: '100%'
											},
											children: [
												{
													type: 'div',
													props: {
														style: {
															display: 'flex',
															justifyContent: 'space-between',
															fontSize: '28px',
															textTransform: 'uppercase',
															letterSpacing: '0.08em',
															color: '#94a3b8'
														},
														children: [
															{
																type: 'div',
																props: { style: { width: '36%' }, children: 'Guild' }
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#94a3b8' },
																	children: 'Players'
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#94a3b8' },
																	children: 'Kills'
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#94a3b8' },
																	children: 'Deaths'
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#94a3b8' },
																	children: 'Kill Fame'
																}
															}
														]
													}
												},
												...guilds.map((guild) => ({
													type: 'div',
													props: {
														style: {
															display: 'flex',
															justifyContent: 'space-between',
															alignItems: 'center',
															fontSize: '36px',
															color: '#e2e8f0'
														},
														children: [
															{
																type: 'div',
																props: {
																	style: { width: '36%', fontWeight: 600 },
																	children: truncateName(guild.GuildName)
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#93c5fd' },
																	children: formatValue(guild.PlayerCount)
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#fca5a5' },
																	children: formatValue(guild.Kills)
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#f0abfc' },
																	children: formatValue(guild.Deaths)
																}
															},
															{
																type: 'div',
																props: {
																	style: { width: '16%', textAlign: 'right', color: '#facc15' },
																	children: formatFame(guild.KillFame)
																}
															}
														]
													}
												}))
											]
										}
									}
								: null,
							showFallback
								? {
										type: 'div',
										props: {
											style: {
												padding: '32px',
												backgroundColor: '#0f172a',
												textAlign: 'center',
												fontSize: '28px',
												color: '#e2e8f0',
												width: '100%'
											},
											children: 'No alliance or guild data found.'
										}
									}
								: null
						].filter(Boolean)
					}
				}
			]
		}
	};
}

async function loadFonts() {
	const regularPath = path.join(
		process.cwd(),
		'node_modules',
		'@fontsource',
		'inter',
		'files',
		'inter-latin-400-normal.woff'
	);
	const semiboldPath = path.join(
		process.cwd(),
		'node_modules',
		'@fontsource',
		'inter',
		'files',
		'inter-latin-600-normal.woff'
	);

	const [regular, semibold] = await Promise.all([
		readFile(regularPath),
		readFile(semiboldPath)
	]);

	return [
		{
			name: 'Inter',
			data: regular,
			weight: 400,
			style: 'normal'
		},
		{
			name: 'Inter',
			data: semibold,
			weight: 600,
			style: 'normal'
		}
	];
}

export const GET = async ({ params, fetch }) => {
	let alliances = [];
	let guilds = [];
	try {
		const response = await fetch(
			`${getApiBase()}/battles/${params.region}/${params.battleIds}`
		);
		if (response.ok) {
			const payload = await response.json();
			const rawAlliances = payload?.Alliances || [];
			const rawGuilds = payload?.Guilds || [];
			alliances = rawAlliances
				.filter((alliance) => alliance?.AllianceName)
				.sort((a, b) => (b.Kills || 0) - (a.Kills || 0))
				.slice(0, 4);
			guilds = rawGuilds
				.filter((guild) => guild?.GuildName)
				.sort((a, b) => (b.Kills || 0) - (a.Kills || 0))
				.slice(0, 4);
		}
	} catch {
		alliances = [];
		guilds = [];
	}

	const fonts = await loadFonts();
	const svg = await satori(buildMarkup({ alliances, guilds }), {
		width: OG_WIDTH,
		height: OG_HEIGHT,
		fonts
	});

	const resvg = new Resvg(svg, {
		fitTo: { mode: 'width', value: OG_WIDTH }
	});
	const pngData = resvg.render();
	const pngBuffer = pngData.asPng();

	return new Response(pngBuffer, {
		headers: {
			'Content-Type': 'image/png',
			'Cache-Control': 'public, max-age=3600'
		}
	});
};
