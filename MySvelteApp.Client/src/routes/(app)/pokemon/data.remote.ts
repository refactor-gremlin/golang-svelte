import { query } from '$app/server';
// sdk.gen automatically generates the client for the API + automatically parses the response with zod
import { getRandomPokemon } from '$api/schema/sdk.gen';

export const getRandomPokemonData = query(async () => {
	const response = await getRandomPokemon();
	if (!response.data) {
		throw new Error('Invalid response from getRandomPokemon');
	}
	return response.data;
});
