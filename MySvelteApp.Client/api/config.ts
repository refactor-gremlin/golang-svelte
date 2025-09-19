// src/api/config.ts
import { dev } from '$app/environment';
import { PUBLIC_API_ENDPOINT } from '$env/static/public';

const defaultUrl = 'http://localhost:8080';
export const config = {
	apiEndpoint:
	PUBLIC_API_ENDPOINT && PUBLIC_API_ENDPOINT !== ''
		? PUBLIC_API_ENDPOINT
		: defaultUrl
};
