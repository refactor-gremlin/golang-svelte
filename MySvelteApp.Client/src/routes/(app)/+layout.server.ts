// if the user is not authenticated, redirect to the login page
import { redirect } from '@sveltejs/kit';
import { logger } from '$lib/server/logger';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies }) => {
	const token = cookies.get('auth_token');
	const log = logger;

	if (!token) {
		log.warn('User is not authenticated, redirecting to login page');
		return redirect(302, '/login');
	}

    // No validation endpoint yet; treat presence of token as authenticated
    return {
        user: {
            id: 'authenticated',
            name: 'User'
        }
    };
};
