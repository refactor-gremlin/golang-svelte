import { execSync } from 'child_process';

// Use a hardcoded development URL since this script runs outside SvelteKit context
const apiEndpoint = 'http://localhost:7216';
const openApiPath = `${apiEndpoint}/swagger/v1/swagger.json`;
const zodOut = './src/api/api-zod-client.ts';

console.log('Generating Zod client at', openApiPath);
execSync(`npx openapi-zod-client "${openApiPath}" -o "${zodOut}"`, { stdio: 'inherit' });
console.log('Generated Zod client at', zodOut);
