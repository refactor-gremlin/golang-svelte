# 1) Client build
FROM node:24-alpine AS client-build
WORKDIR /src/MySvelteApp.Client

# install deps
COPY MySvelteApp.Client/package*.json ./
RUN npm install

# build client
COPY MySvelteApp.Client/. .
RUN npm run build

# 2) Server build
FROM python:3.11-slim AS server-build
WORKDIR /src/MySvelteApp.Server

# Install uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /usr/local/bin/uv

# copy requirements
COPY MySvelteApp.Server/requirements.txt ./
RUN uv pip install --system --no-cache -r requirements.txt

# copy code + static assets
COPY MySvelteApp.Server/. .
COPY --from=client-build /src/MySvelteApp.Client/.svelte-kit/output/client ./static

# 3) Runtime
FROM python:3.11-slim
WORKDIR /app

# Copy Python dependencies
COPY --from=server-build /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=server-build /usr/local/bin /usr/local/bin

# Copy application code
COPY --from=server-build /src/MySvelteApp.Server /app

EXPOSE 8080
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8080"]
