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
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS server-build
WORKDIR /src/MySvelteApp.Server

# copy go mod files
COPY MySvelteApp.Server/go.* ./

# download dependencies
RUN go mod download

# copy code + static assets
COPY MySvelteApp.Server/. .
COPY --from=client-build /src/MySvelteApp.Client/.svelte-kit/output/client ./static

# build the binary
RUN go build -o /app/server ./cmd/server

# 3) Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=server-build /app/server .

EXPOSE 8080
ENTRYPOINT ["./server"]
