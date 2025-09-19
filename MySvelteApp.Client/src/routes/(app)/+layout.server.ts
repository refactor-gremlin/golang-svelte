// if the user is not authenticated, redirect to the login page
import { redirect } from '@sveltejs/kit';
import { logger } from '$lib/server/logger';
import { getTestAuth } from '$api/schema/sdk.gen';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies }) => {
	const token = cookies.get('auth_token');
	const log = logger;

	if (!token) {
		log.warn('User is not authenticated, redirecting to login page');
		return redirect(302, '/login');
	}

	// Verify token with backend using generated TestAuth client
	try {
		await getTestAuth({
			headers: {
				Authorization: `Bearer ${token}`
			},
			throwOnError: true as const
		});

		// Since the TestAuth endpoint doesn't return user data, we'll return a basic user object
		return {
			user: {
				id: 'authenticated',
				name: 'User'
			}
		};
	} catch (error) {
		log.error({ err: error }, 'Token validation failed, redirecting to login page');
		cookies.delete('auth_token', { path: '/' });
		return redirect(302, '/login');
	}
};
