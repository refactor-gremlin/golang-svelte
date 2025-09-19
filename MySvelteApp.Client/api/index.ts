// src/api/index.ts

import { client } from './schema/client.gen'; // Generated client instance
import * as sdk from './schema/sdk.gen'; // All your SDK methods
import * as schemas from './schema/schemas.gen'; // JSON Schema objects
import * as types from './schema/types.gen'; // TypeScript types

// Export for use in your app
export const api = sdk; // Call sdk.getWeatherForecast(), etc.
export { client, schemas, types };
