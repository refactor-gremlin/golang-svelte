<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardHeader, CardTitle, CardFooter } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { Avatar, AvatarImage, AvatarFallback } from '$lib/components/ui/avatar';
	import { getRandomPokemonData } from './data.remote';
	import SparklesIcon from '@lucide/svelte/icons/sparkles';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import ZapIcon from '@lucide/svelte/icons/zap';
	import SkullIcon from '@lucide/svelte/icons/skull';
	import { cn } from '$lib/utils';
	import { fade, slide, fly, scale } from 'svelte/transition';

	// Button fade state
	let isButtonFading = $state(false);
	let currentPromise = $state<Promise<unknown> | null>(null);

	// Reset button fade state when the current promise settles
	$effect(() => {
		if (currentPromise) {
			currentPromise
				.finally(() => {
					// Reset after a short delay to allow for smooth transition
					setTimeout(() => {
						isButtonFading = false;
						currentPromise = null;
					}, 300);
				})
				.catch(() => {
					// Handle error case - still reset the button
					setTimeout(() => {
						isButtonFading = false;
						currentPromise = null;
					}, 300);
				});
		}
	});
</script>

<div class="min-h-screen">
	<div class="container mx-auto max-w-4xl px-4 py-12">
		<!-- Header Section -->
		<div class="mb-8 text-center" in:fade>
			<div class="mb-4 inline-flex items-center gap-2" in:fly={{ y: 20 }}>
				<SparklesIcon class="h-8 w-8 text-yellow-500" />
				<h1 class="text-4xl font-bold">Pokemon Explorer</h1>
				<SparklesIcon class="h-8 w-8 text-yellow-500" />
			</div>
			<p class="mb-6 text-lg text-muted-foreground" in:slide>
				Discover amazing Pokemon with a single click! ⚡
			</p>
		</div>

		<!-- Action Button -->
		<div class="mb-8 flex justify-center" in:fade>
			<div transition:fade={{ duration: 300 }}>
				<Button
					onclick={() => {
						isButtonFading = true;
						currentPromise = getRandomPokemonData().refresh();
					}}
					size="lg"
					class={cn(
						'px-8 py-4 text-lg font-semibold',
						'transform shadow-lg hover:scale-105 hover:shadow-xl',
						isButtonFading && 'cursor-not-allowed opacity-50'
					)}
					disabled={isButtonFading}
				>
					<SparklesIcon class="mr-2 h-5 w-5" />
					{isButtonFading ? 'Discovering...' : 'Discover Pokemon'}
				</Button>
			</div>
		</div>

		<!-- Pokemon Display -->
		<div class="mb-12 flex justify-center" in:fade>
			<svelte:boundary>
				{#await getRandomPokemonData() then pokemon}
					<div class="flex flex-col items-center">
						<Card
							class="w-full max-w-md border-0 bg-white/80 shadow-2xl backdrop-blur-sm dark:bg-gray-800/80"
						>
							<CardHeader class="pb-4 text-center">
								<div class="mb-4 flex justify-center">
									<div in:scale={{ delay: 200 }}>
										<Avatar class="h-32 w-32 border-4 border-white shadow-lg">
											<AvatarImage
												src={pokemon?.image || ''}
												alt={pokemon?.name || 'Pokemon'}
												class="object-cover"
											/>
											<AvatarFallback
												class="bg-gradient-to-br from-blue-400 to-purple-500 text-2xl font-bold text-white"
											>
												{pokemon?.name?.charAt(0)?.toUpperCase() || '?'}
											</AvatarFallback>
										</Avatar>
									</div>
								</div>
								<div in:slide={{ delay: 400 }}>
									<CardTitle class="text-center text-3xl font-bold">
										{pokemon?.name || 'Unknown Pokemon'}
									</CardTitle>
								</div>
							</CardHeader>

							<CardContent class="pb-4 text-center">
								{#if pokemon?.type}
									<div class="flex justify-center">
										<div in:fly={{ y: 10, delay: 600 }}>
											<Badge
												variant="default"
												class="flex items-center gap-2 px-4 py-2 text-sm font-semibold"
											>
												<ZapIcon class="h-4 w-4" />
												{pokemon.type}
											</Badge>
										</div>
									</div>
								{:else}
									<div class="flex justify-center">
										<div in:fade={{ delay: 600 }}>
											<Badge variant="outline" class="px-4 py-2 text-sm font-semibold">
												<SparklesIcon class="mr-1 h-4 w-4" />
												Mystery Type
											</Badge>
										</div>
									</div>
								{/if}
							</CardContent>

							<CardFooter class="border-t border-gray-100 pt-4 text-center dark:border-gray-700">
								<div in:fade={{ delay: 800 }}>
									<div class="text-sm text-muted-foreground">
										✨ Caught a wild {pokemon?.name || 'Pokemon'}! ✨
									</div>
								</div>
							</CardFooter>
						</Card>
					</div>
				{/await}

				{#snippet pending()}
					<div class="flex flex-col items-center">
						<Card
							class="mx-auto w-full max-w-md border-0 bg-white/80 shadow-2xl backdrop-blur-sm dark:bg-gray-800/80"
						>
							<CardHeader class="pb-4 text-center">
								<div class="mb-4 flex justify-center">
									<div in:scale>
										<Skeleton class="h-32 w-32 rounded-full" />
									</div>
								</div>
								<div in:slide>
									<Skeleton class="mx-auto h-8 w-48" />
								</div>
							</CardHeader>

							<CardContent class="pb-4 text-center">
								<div in:fade>
									<Skeleton class="mx-auto h-8 w-24" />
								</div>
							</CardContent>

							<CardFooter class="border-t border-gray-100 pt-4 text-center dark:border-gray-700">
								<div in:fade>
									<Skeleton class="mx-auto h-4 w-40" />
								</div>
							</CardFooter>
						</Card>

						<!-- Fun Stats Section - Loading State -->
						<div class="mt-8" in:fade={{ delay: 1000 }}>
							<!-- Add a subtle loading hint -->
							<div
								class="mt-4 text-center text-sm text-muted-foreground opacity-70"
								in:fade={{ delay: 1600 }}
							>
								🔍 Searching for Pokemon...
							</div>
						</div>
					</div>
				{/snippet}

				{#snippet onerror(err: unknown)}
					<div class="flex flex-col items-center">
						<Card class="mx-auto w-full max-w-md border-0 bg-red-50 shadow-2xl dark:bg-red-900/20">
							<CardContent class="py-8 text-center">
								<div class="mb-4 text-red-500">
									<div in:scale>
										<SkullIcon class="mx-auto h-12 w-12" />
									</div>
								</div>
								<div in:fade>
									<h3 class="mb-2 text-lg font-semibold text-red-800 dark:text-red-200">
										Oops! Pokemon escaped!
									</h3>
								</div>
								<div in:slide>
									<p class="mb-4 text-sm text-red-600 dark:text-red-300">
										{err instanceof Error
											? err.message
											: 'Something went wrong while catching Pokemon'}
									</p>
								</div>
								<Button
									onclick={() => getRandomPokemonData().refresh()}
									variant="outline"
									class="border-red-300 text-red-700 hover:bg-red-50"
								>
									<RefreshCwIcon class="mr-2 h-4 w-4" />
									Try Again
								</Button>
							</CardContent>
						</Card>
					</div>
				{/snippet}
			</svelte:boundary>
		</div>
	</div>
</div>
