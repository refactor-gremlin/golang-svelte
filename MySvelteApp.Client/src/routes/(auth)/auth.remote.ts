import { form, command, query } from '$app/server';
import { getRequestEvent } from '$app/server';
import { error } from '@sveltejs/kit';
import { postAuthLogin, postAuthRegister, getTestAuth } from '$api/schema/sdk.gen';
import { zAuthErrorResponse } from '$api/schema/zod.gen';
import { z } from 'zod';

// Stricter UI-side validation schemas for immediate feedback
const zLoginForm = z.object({
	username: z.string().trim().min(1, 'Username is required'),
	password: z.string().min(1, 'Password is required')
});

const zRegisterForm = z.object({
	username: z.string().trim().min(1, 'Username is required'),
	email: z.email('Valid email required'),
	password: z.string().min(8, 'Password must be at least 8 characters')
});

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

const getString = (value: FormDataEntryValue | null) => (typeof value === 'string' ? value : '');

// Login form handler with automatic validation
export const login = form(async (formData) => {
	const parsed = zLoginForm.safeParse({
		username: getString(formData.get('username')),
		password: getString(formData.get('password'))
	});

	if (!parsed.success) {
		throw error(400, { message: 'Invalid login data' });
	}

	const { username, password } = parsed.data;
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
	const parsed = zRegisterFormWithConfirm.safeParse({
		username: getString(formData.get('username')),
		email: getString(formData.get('email')),
		password: getString(formData.get('password')),
		confirmPassword: getString(formData.get('confirmPassword'))
	});

	if (!parsed.success) {
		const message = parsed.error.issues[0]?.message ?? 'Invalid registration data';
		throw error(400, { message });
	}

	const { username, email, password } = parsed.data;
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

// Get current user query
export const getCurrentUser = query(async () => {
	const { cookies } = getRequestEvent();

	const token = cookies.get('auth_token');
	if (!token) {
		error(401, 'Not authenticated');
	}

	// Validate token with backend using generated TestAuth client
	try {
		await getTestAuth({
			headers: {
				Authorization: `Bearer ${token}`
			},
			throwOnError: true as const
		});
		return {
			id: 'user123',
			email: 'user@example.com',
			name: 'Test User',
			token
		};
	} catch (err) {
		console.error('Token validation error:', err);
		error(401, 'Authentication failed');
	}
});

// Check if user is authenticated
export const isAuthenticated = query(async () => {
	const { cookies } = getRequestEvent();
	const token = cookies.get('auth_token');
	return !!token;
});
