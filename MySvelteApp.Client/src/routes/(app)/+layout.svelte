<script lang="ts">
	import '../../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Home, Settings, User, Plus, ChevronUp } from '@lucide/svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { logout } from '$src/routes/(auth)/auth.remote';
	import { goto } from '$app/navigation';
	import { Toaster } from 'svelte-sonner';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Sidebar.Provider>
	<div class="flex h-screen w-full">
		<Sidebar.Root>
			<Sidebar.Header>
				<h2 class="text-lg font-semibold">My App</h2>
			</Sidebar.Header>

			<Sidebar.Content>
				<Sidebar.Group>
					<Sidebar.GroupLabel>Application</Sidebar.GroupLabel>
					<Sidebar.GroupAction>
						<Plus size={16} />
						<span class="sr-only">Add Project</span>
					</Sidebar.GroupAction>
					<Sidebar.GroupContent>
						<Sidebar.Menu>
							<Sidebar.MenuItem>
								<Sidebar.MenuButton>
									<Home size={16} />
									<span>Dashboard</span>
								</Sidebar.MenuButton>
							</Sidebar.MenuItem>
							<Sidebar.MenuItem>
								<Sidebar.MenuButton>
									<User size={16} />
									<span>Profile</span>
								</Sidebar.MenuButton>
							</Sidebar.MenuItem>
						</Sidebar.Menu>
					</Sidebar.GroupContent>
				</Sidebar.Group>

				<Sidebar.Separator />

				<Sidebar.Group>
					<Sidebar.GroupLabel>Settings</Sidebar.GroupLabel>
					<Sidebar.GroupContent>
						<Sidebar.Menu>
							<Sidebar.MenuItem>
								<Sidebar.MenuButton>
									<Settings size={16} />
									<span>Settings</span>
								</Sidebar.MenuButton>
							</Sidebar.MenuItem>
						</Sidebar.Menu>
					</Sidebar.GroupContent>
				</Sidebar.Group>
			</Sidebar.Content>

			<Sidebar.Footer>
				<Sidebar.Menu>
					<Sidebar.MenuItem>
						<DropdownMenu.Root>
							<DropdownMenu.Trigger>
								<Sidebar.MenuButton
									class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
								>
									<User size={16} />
									<span class="ml-2">Username</span>
									<ChevronUp class="ml-auto" size={16} />
								</Sidebar.MenuButton>
							</DropdownMenu.Trigger>
							<DropdownMenu.Content side="top" class="w-(--bits-dropdown-menu-anchor-width)">
								<DropdownMenu.Item>
									<span>Account</span>
								</DropdownMenu.Item>
								<DropdownMenu.Item>
									<span>Billing</span>
								</DropdownMenu.Item>
								<DropdownMenu.Separator />
								<DropdownMenu.Item>
									<button
										type="button"
										class="w-full text-left hover:bg-transparent"
										onclick={async () => {
											try {
												await logout();
												goto('/login');
											} catch (error) {
												console.error('Logout failed:', error);
											}
										}}
									>
										Sign out
									</button>
								</DropdownMenu.Item>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</Sidebar.MenuItem>
				</Sidebar.Menu>
			</Sidebar.Footer>

			<Sidebar.Rail />
		</Sidebar.Root>

		<Sidebar.Inset>
			<header class="flex h-12 items-center justify-between border-b px-4">
				<Sidebar.Trigger />
				<div class="flex items-center space-x-2">
					<h1 class="text-xl font-semibold">Welcome</h1>
				</div>
			</header>
			<main class="flex-1 p-4">
				{@render children?.()}
			</main>
		</Sidebar.Inset>
	</div>
	<Toaster />
</Sidebar.Provider>
