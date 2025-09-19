import type { CreateClientConfig } from '$api/schema/client/client.gen';
import { config } from '$api/config';

export const createClientConfig: CreateClientConfig = (override) => ({
	...override,
	baseUrl: config.apiEndpoint
});
