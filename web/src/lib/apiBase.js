import { dev } from '$app/environment';

export function getApiBase() {
	return dev ? 'https://albionstats.bricefrisco.com/api' : 'http://localhost:8080/api';
}
