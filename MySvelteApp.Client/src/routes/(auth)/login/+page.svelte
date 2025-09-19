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

	const DEFAULT_ERROR_MESSAGE = 'Login failed. Please check your credentials.';
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
						/>
					</div>

					<div class="space-y-2">
						<Label for="password">Password</Label>
						<Input
							id="password"
							name="password"
							type="password"
							placeholder="Enter your password"
							required
						/>
					</div>

					<Button type="submit" class="w-full" disabled={login.pending > 0}>
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
