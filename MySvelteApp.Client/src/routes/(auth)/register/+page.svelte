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
	import { register } from '$src/routes/(auth)/auth.remote';
	import { resolveAuthErrorMessage } from '$lib/auth/error-messages';

	let error = $state<string | null>(null);
	let success = $state<string | null>(null);

	const DEFAULT_ERROR_MESSAGE = 'Registration failed. Please try again.';
</script>

<div
	class="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8 dark:bg-gray-900"
>
	<div class="w-full max-w-md space-y-8">
		<div class="text-center">
			<h2 class="mt-6 text-3xl font-bold text-gray-900 dark:text-white">Create your account</h2>
			<p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
				Already have an account?
				<a href="/login" class="font-medium text-blue-600 hover:text-blue-500 dark:text-blue-400">
					Sign in
				</a>
			</p>
		</div>

		<Card>
			<CardHeader class="space-y-1">
				<CardTitle class="text-center text-2xl font-bold">Get started</CardTitle>
				<CardDescription class="text-center">
					Create a new account to access all features
				</CardDescription>
			</CardHeader>
			<CardContent>
				<!-- Error and Success Messages -->
				{#if error}
					<div class="mb-4 rounded border border-red-400 bg-red-100 p-3 text-red-700">
						{error}
					</div>
				{/if}

				{#if success}
					<div class="mb-4 rounded border border-green-400 bg-green-100 p-3 text-green-700">
						{success}
					</div>
				{/if}

				<form
					{...register.enhance(async ({ submit }) => {
						error = null;
						success = null;

						try {
							await submit();
							success = 'Registration successful! Please log in.';
							setTimeout(() => goto('/login'), 2000);
						} catch (err: unknown) {
							error = resolveAuthErrorMessage(err, DEFAULT_ERROR_MESSAGE);
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
						/>
					</div>

					<div class="space-y-2">
						<Label for="email">Email address</Label>
						<Input id="email" name="email" type="email" placeholder="Enter your email" required />
					</div>

					<div class="space-y-2">
						<Label for="password">Password</Label>
						<Input
							id="password"
							name="password"
							type="password"
							placeholder="Create a password"
							required
						/>
					</div>

					<div class="space-y-2">
						<Label for="confirmPassword">Confirm password</Label>
						<Input
							id="confirmPassword"
							name="confirmPassword"
							type="password"
							placeholder="Confirm your password"
							required
						/>
					</div>

					<Button type="submit" class="w-full" disabled={register.pending > 0}>
						{register.pending > 0 ? 'Creating account...' : 'Create account'}
					</Button>

					<div class="text-center text-xs text-gray-500 dark:text-gray-400">
						By creating an account, you agree to our
						<a href="/terms" class="text-blue-600 hover:text-blue-500 dark:text-blue-400">
							Terms of Service
						</a>
						and
						<a href="/privacy" class="text-blue-600 hover:text-blue-500 dark:text-blue-400">
							Privacy Policy
						</a>
					</div>
				</form>
			</CardContent>
		</Card>
	</div>
</div>
