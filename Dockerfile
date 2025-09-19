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
FROM --platform=$BUILDPLATFORM mcr.microsoft.com/dotnet/sdk:10.0-preview-alpine AS server-build
ARG TARGETARCH
WORKDIR /src/MySvelteApp.Server

# restore
COPY MySvelteApp.Server/*.csproj ./
RUN dotnet restore -a $TARGETARCH

# copy code + static assets
COPY MySvelteApp.Server/. .
COPY --from=client-build /src/MySvelteApp.Client/.svelte-kit/output/client ./wwwroot

# publish (no --no-restore)
RUN dotnet publish -c Release -a $TARGETARCH -o /app/publish

# 3) Runtime
FROM mcr.microsoft.com/dotnet/aspnet:10.0-preview-alpine
WORKDIR /app
COPY --from=server-build /app/publish .

EXPOSE 8080
ENTRYPOINT ["dotnet", "MySvelteApp.Server.dll"]
