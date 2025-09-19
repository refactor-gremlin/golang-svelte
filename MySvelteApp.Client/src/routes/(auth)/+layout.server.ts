import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async () => {
	// Auth routes are publicly accessible - no authentication required
	return {};
};
