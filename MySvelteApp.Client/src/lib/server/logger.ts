import pino from 'pino';
import type { Logger } from 'pino';

const promtailHost = process.env.LOKI_PUSH_URL ?? 'http://localhost:3101';
const service = process.env.OTEL_SERVICE_NAME ?? 'mysvelteapp-web';
const environment = process.env.NODE_ENV ?? 'development';
const level = process.env.LOG_LEVEL ?? 'info';

async function createLogger(): Promise<Logger> {
	try {
		const transport = await pino.transport({
			target: 'pino-loki',
			options: {
				host: promtailHost,
				batching: true,
				interval: 1000,
				labels: {
					service,
					env: environment
				}
			}
		});

		return pino({ level }, transport);
	} catch (error) {
		console.error('Falling back to stdout logger; unable to configure Loki transport.', error);
		return pino({ level });
	}
}

export const logger: Logger = await createLogger();
