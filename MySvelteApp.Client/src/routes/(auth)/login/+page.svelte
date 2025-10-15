<script lang="ts">
	import { goto } from '$app/navigation';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { login } from '$src/routes/(auth)/auth.remote';
	import { resolveAuthErrorMessage } from '$lib/auth/error-messages';
	import { toast } from 'svelte-sonner';
	import { z } from 'zod';
	import { createFormValidation } from '$lib/composables/client-form-validation';

	const DEFAULT_ERROR_MESSAGE = 'Login failed. Please check your credentials.';

	const loginSchema = z.object({
		username: z.string().trim().min(1, 'Username is required'),
		password: z.string().min(1, 'Password is required')
	});

	const form = createFormValidation(loginSchema);
</script>

<div
	class="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8 dark:bg-gray-900"
>
	<div class="w-full max-w-md space-y-8">
		<div class="text-center">
			<h2 class="mt-6 text-3xl font-bold text-gray-900 dark:text-white">Sign in to your account</h2>
			<p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
				Don't have an account?
				<a
					href="/register"
					class="font-medium text-blue-600 hover:text-blue-500 dark:text-blue-400"
				>
					Sign up
				</a>
			</p>
		</div>

		<Card>
			<CardHeader class="space-y-1">
				<CardTitle class="text-center text-2xl font-bold">Welcome back</CardTitle>
				<CardDescription class="text-center">
					Enter your credentials to access your account
				</CardDescription>
			</CardHeader>
			<CardContent>
				<form
					{...login.enhance(async ({ submit }) => {
						if (!form.validateForm()) return;
						
						try {
							await submit();
							toast.success('Login successful!');
							goto('/');
						} catch (error: unknown) {
							toast.error(resolveAuthErrorMessage(error, DEFAULT_ERROR_MESSAGE));
						}
					})}
					class="space-y-4"
				>
					<div class="space-y-2">
						<Label for="username">Username</Label>
						<Input
							id="username"
							name="username"
							type="text"
							placeholder="Enter your username"
							required
							bind:value={form.formData.username}
							oninput={() => form.validateField('username', form.formData.username)}
							onblur={() => form.validateField('username', form.formData.username)}
							class={form.errors.username && form.touched.username ? 'border-red-500 focus:ring-red-500' : ''}
						/>
						{#if form.errors.username && form.touched.username}
							<p class="text-sm text-red-600 mt-1">{form.errors.username}</p>
						{/if}
					</div>

					<div class="space-y-2">
						<Label for="password">Password</Label>
						<Input
							id="password"
							name="password"
							type="password"
							placeholder="Enter your password"
							required
							bind:value={form.formData.password}
							oninput={() => form.validateField('password', form.formData.password)}
							onblur={() => form.validateField('password', form.formData.password)}
							class={form.errors.password && form.touched.password ? 'border-red-500 focus:ring-red-500' : ''}
						/>
						{#if form.errors.password && form.touched.password}
							<p class="text-sm text-red-600 mt-1">{form.errors.password}</p>
						{/if}
					</div>

					<Button type="submit" class="w-full" disabled={login.pending > 0 || !form.isValid}>
						{#if login.pending > 0}
							Signing in...
						{:else}
							Sign in
						{/if}
					</Button>

					<div class="text-center">
						<a
							href="/forgot-password"
							class="text-sm text-blue-600 hover:text-blue-500 dark:text-blue-400"
						>
							Forgot your password?
						</a>
					</div>
				</form>
			</CardContent>
		</Card>
	</div>
</div>
