// src/api/authService.ts
class AuthService {
	private token: string | null = null;

	constructor() {
		// In browser environment, token will be managed via cookies
		// Server-side rendering will handle token via load functions
	}

	// Get the current auth token (this would typically come from a load function)
	getToken(): string | null {
		return this.token;
	}

	// Check if user is authenticated (this would typically come from a load function)
	isAuthenticated(): boolean {
		return !!this.token;
	}

	// Set the auth token (primarily for server-side use)
	setToken(token: string): void {
		this.token = token;
	}

	// Clear the auth token (primarily for server-side use)
	clearToken(): void {
		this.token = null;
	}

	// Get auth headers for API requests
	getAuthHeaders(): HeadersInit {
		if (this.token) {
			return {
				Authorization: `Bearer ${this.token}`,
				'Content-Type': 'application/json'
			};
		}
		return {
			'Content-Type': 'application/json'
		};
	}

	// Intercept API calls to add auth headers
	async fetchWithAuth(url: string, options: RequestInit = {}): Promise<Response> {
		const headers = {
			...this.getAuthHeaders(),
			...options.headers
		};

		return fetch(url, {
			...options,
			headers
		});
	}
}

export const authService = new AuthService();
