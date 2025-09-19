# MySvelteApp.Client/Dockerfile.dev
FROM node:22-alpine

WORKDIR /app

# Install deps up front for faster rebuilds
COPY package*.json ./
RUN npm ci

ENV HOST=0.0.0.0
EXPOSE 5173

# Ensure deps present when node_modules is a named volume
COPY docker-entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

CMD ["/usr/local/bin/entrypoint.sh"]
