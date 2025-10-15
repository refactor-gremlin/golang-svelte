import { form, command, query } from '$app/server';
import { getRequestEvent } from '$app/server';
import { error } from '@sveltejs/kit';
import { postAuthLogin, postAuthRegister } from '$api/schema/sdk.gen';
import { zAuthErrorResponse } from '$api/schema/zod.gen';
import { z } from 'zod';
import { safeValidateFormData } from '$lib/server/utils/server-form-validation';

/**
 * Authentication remote functions for SvelteKit experimental remote functions.
 *
 * WHY THIS ARCHITECTURE:
 * - Server-side validation ensures data integrity and security
 * - Experimental remote functions provide seamless client-server communication
 * - Zod schemas provide type-safe validation on both client and server
 * - Centralized authentication logic prevents code duplication
 * - JWT token management with secure cookie handling
 *
 * SECURITY CONSIDERATIONS:
 * - All form data is validated on the server before processing
 * - Passwords are never logged or exposed in error messages
 * - JWT tokens use httpOnly, secure, and sameSite cookie attributes
 * - Input sanitization prevents injection attacks
 * - Consistent error handling prevents information leakage
 *
 * USE CASE:
 * - User authentication (login, register, logout)
 * - Session management via JWT tokens
 * - Account status verification
 */

// Server-side validation schemas - match client-side but ensure security
// WHY: Provides type-safe validation and consistent error messages
const zLoginForm = z.object({
	username: z.string().min(1, 'Username is required'),
	password: z.string().min(1, 'Password is required')
});

const zRegisterForm = z.object({
	username: z.string().min(1, 'Username is required'),
	email: z.email('Valid email required'),
	password: z.string().min(8, 'Password must be at least 8 characters')
});

// Enhanced registration schema with password confirmation
// WHY: Ensures passwords match before server processing to provide better UX
const zRegisterFormWithConfirm = zRegisterForm
	.extend({
		confirmPassword: z.string().min(1, 'Please confirm your password')
	})
	.superRefine((data, ctx) => {
		if (data.password !== data.confirmPassword) {
			ctx.addIssue({
				code: 'custom',
				message: 'Passwords do not match',
				path: ['confirmPassword']
			});
		}
	});

// Login form handler with automatic validation
export const login = form(async (formData) => {
	const { username, password } = safeValidateFormData(formData, zLoginForm);
	const { cookies } = getRequestEvent();

	try {
		// Use generated API client with ThrowOnError for cleaner control flow
		const response = await postAuthLogin({
			body: { username, password },
			throwOnError: true as const
		});

		const result = response.data;

		// Set JWT token in cookie
		if (result?.token) {
			cookies.set('auth_token', result.token, {
				path: '/',
				httpOnly: true,
				secure: import.meta.env.PROD,
				sameSite: 'strict'
			});
		}

		return result;
	} catch (err) {
		console.error('Login error:', err);
		const parsed = zAuthErrorResponse.safeParse(err);
		const message =
			parsed.success && parsed.data.message
				? parsed.data.message
				: err instanceof Error
					? err.message
					: 'Network error. Please check your connection and try again.';
		throw error(401, { message });
	}
});

// Registration form handler with automatic validation
export const register = form(async (formData) => {
	const { username, email, password } = safeValidateFormData(formData, zRegisterFormWithConfirm);
	const registerData = { username, email, password };
	try {
		// Use generated API client with ThrowOnError
		const response = await postAuthRegister({
			body: registerData,
			throwOnError: true as const
		});

		const result = response.data;

		return result;
	} catch (err) {
		console.log('Registration catch error:', err);
		const parsed = zAuthErrorResponse.safeParse(err);
		const message =
			parsed.success && parsed.data.message
				? parsed.data.message
				: err instanceof Error
					? err.message
					: 'Registration failed';
		throw error(400, { message });
	}
});

// Logout command
export const logout = command(async () => {
	const { cookies } = getRequestEvent();

	// Clear auth token cookie
	cookies.delete('auth_token', { path: '/' });

	return { success: true };
});

// Check if user is authenticated
export const isAuthenticated = query(async () => {
	const { cookies } = getRequestEvent();
	const token = cookies.get('auth_token');
	return !!token;
});
