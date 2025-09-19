// src/api/index.ts

import { config } from './config';
import { client } from './schema/client.gen'; // Generated client instance :contentReference[oaicite:9]{index=9}
import * as sdk from './schema/sdk.gen'; // All your SDK methods
import * as schemas from './schema/schemas.gen'; // JSON Schema objects
import * as types from './schema/types.gen'; // TypeScript types

// 1) Apply your runtime base URL from SvelteKit config/env
client.setConfig({ baseUrl: config.apiEndpoint }); // Applies to Fetch client :contentReference[oaicite:10]{index=10}

// 2) Export for use in your app
export const api = sdk; // Call sdk.getWeatherForecast(), etc.
export { client, schemas, types };
