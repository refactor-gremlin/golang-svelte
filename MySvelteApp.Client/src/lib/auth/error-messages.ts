import type { ActionFailure } from '@sveltejs/kit';
import { zAuthErrorResponse } from '$api/schema/zod.gen';
import type { AuthErrorResponse } from '$api/schema/types.gen';

const extractAuthMessage = (payload: unknown): string | undefined => {
	const parsed = zAuthErrorResponse.safeParse(payload);
	if (!parsed.success) return undefined;
	const message = parsed.data.message;
	return typeof message === 'string' && message.trim().length > 0 ? message : undefined;
};

export const resolveAuthErrorMessage = (err: unknown, fallback: string): string => {
	if (typeof err === 'object' && err !== null) {
		if ('body' in err) {
			const message = extractAuthMessage((err as { body?: unknown }).body);
			if (message) return message;
		}

		if ('data' in err) {
			const failure = err as ActionFailure<AuthErrorResponse>;
			const message = extractAuthMessage(failure.data);
			if (message) return message;
		}
	}

	if (err instanceof Error) {
		const message = err.message?.trim();
		if (message) return message;
	}

	return fallback;
};
