export const routes = [
	{ label: 'Battle Boards', base: 'battle-boards' },
	{ label: 'Alliances', base: 'alliances' },
	{ label: 'Guilds', base: 'guilds' },
	{ label: 'Players', base: 'players' }
];

export const regionRoutes = routes.map((route) => route.base);
