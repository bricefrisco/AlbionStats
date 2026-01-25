import { redirect } from '@sveltejs/kit';

export const load = async ({ url }) => {
	throw redirect(308, `/players/americas${url.search}`);
};
