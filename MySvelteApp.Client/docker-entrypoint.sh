# MySvelteApp.Client/docker-entrypoint.sh
#!/bin/sh
set -e
cd /app

# If node_modules isn’t present (volume), install once
if [ ! -d node_modules ] || [ ! -f node_modules/.modules-stamp ]; then
  echo "[entrypoint] Installing dependencies..."
  npm ci
  mkdir -p node_modules && touch node_modules/.modules-stamp
fi

exec npm run dev -- --host 0.0.0.0 --port 5173
