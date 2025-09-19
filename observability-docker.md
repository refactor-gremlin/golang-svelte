### Observability: Docker stack + local apps

The SvelteKit client and the ASP.NET Core API run locally while Grafana, Loki,
Tempo, and Promtail run inside Docker. Both apps push logs to Promtail (which
forwards them to Loki) and send traces directly to Tempo.

All commands below run from the repository root. Docker can be rootful or
rootless.

---

### 1. Docker Compose definition
Create `observability.compose.yml` alongside this document.

```yaml
services:
  loki:
    image: grafana/loki:3.4.1
    command: -config.file=/etc/loki/local-config.yaml
    volumes:
      - ./observability/loki/loki-config.yaml:/etc/loki/local-config.yaml:ro
      - loki-data:/loki
    ports:
      - "3100:3100"

  tempo:
    image: grafana/tempo:2.5.0
    command: ["-config.file=/etc/tempo.yaml"]
    ports:
      - "4317:4317"  # OTLP gRPC
      - "4318:4318"  # OTLP HTTP
    volumes:
      - ./observability/tempo/tempo.yaml:/etc/tempo.yaml:ro
      - tempo-data:/var/tempo

  promtail:
    image: grafana/promtail:3.1.1
    command: ["-config.file=/etc/promtail/config.yml"]
    ports:
      - "3101:3101"  # Promtail push API
    volumes:
      - ./observability/promtail/config.yml:/etc/promtail/config.yml:ro
      - /var/log:/var/log:ro

  grafana:
    image: grafana/grafana-oss:10.4.3
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_AUTH_ANONYMOUS_ENABLED=true
      - GF_AUTH_ANONYMOUS_ORG_ROLE=Admin
      - GF_AUTH_DISABLE_LOGIN_FORM=true
    volumes:
      - grafana-data:/var/lib/grafana
      - ./observability/grafana/provisioning/datasources:/etc/grafana/provisioning/datasources:ro

volumes:
  loki-data:
  tempo-data:
  grafana-data:
```

---

### 2. Service configuration files
Create the following files under `observability/` (already committed).

`observability/loki/loki-config.yaml`
```yaml
auth_enabled: false

server:
  http_listen_port: 3100
  grpc_listen_port: 9095

common:
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
  replication_factor: 1
  ring:
    kvstore:
      store: inmemory

ingester:
  lifecycler:
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
  chunk_idle_period: 5m
  chunk_retain_period: 30s
  wal:
    enabled: false

schema_config:
  configs:
    - from: 2024-01-01
      store: boltdb-shipper
      object_store: filesystem
      schema: v13
      index:
        prefix: index_
        period: 24h

storage_config:
  filesystem:
    directory: /loki/chunks
  boltdb_shipper:
    active_index_directory: /loki/index
    cache_location: /loki/boltdb-cache

ruler:
  storage:
    type: local
    local:
      directory: /loki/rules

limits_config:
  ingestion_rate_mb: 16
  ingestion_burst_size_mb: 32
  allow_structured_metadata: false
  max_global_streams_per_user: 0
```

`observability/tempo/tempo.yaml`
```yaml
server:
  http_listen_port: 3200
  http_listen_address: 0.0.0.0
distributor:
  receivers:
    otlp:
      protocols:
        http:
          endpoint: 0.0.0.0:4318
        grpc:
          endpoint: 0.0.0.0:4317
storage:
  trace:
    backend: local
    wal:
      path: /var/tempo/wal
    local:
      path: /var/tempo/traces
compactor:
  compaction:
    block_retention: 48h
metrics_generator:
  registry:
    external_labels:
      source: tempo
```

`observability/promtail/config.yml`
```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

clients:
  - url: http://loki:3100/loki/api/v1/push

positions:
  filename: /tmp/positions.yaml

scrape_configs:
  - job_name: app-push
    loki_push_api:
      server:
        http_listen_port: 3101
      labels:
        source: app
        env: dev
```

`observability/grafana/provisioning/datasources/datasources.yml`
```yaml
apiVersion: 1
datasources:
  - name: Loki
    type: loki
    access: proxy
    url: http://loki:3100
    isDefault: true
  - name: Tempo
    type: tempo
    access: proxy
    url: http://tempo:3200
    jsonData:
      httpMethod: GET
      tracesToLogs:
        datasourceUid: Loki
        tags: ["job", "service"]
        mappedTags:
          - key: service.name
            value: service
```

---

### 3. Start the stack
```bash
docker compose -f observability.compose.yml up -d --remove-orphans
```

Services:
- Grafana → <http://localhost:3000> (admin/admin, anonymous access enabled)
- Loki HTTP API → <http://localhost:3100>
- Tempo OTLP → `http://localhost:4318` (HTTP) or `localhost:4317` (gRPC)
- Promtail push API → `http://localhost:3101/loki/api/v1/push`

---

### 4. SvelteKit client instrumentation

`MySvelteApp.Client/src/instrumentation.server.js` already boots the
OpenTelemetry Node SDK and exports traces via OTLP/HTTP to Tempo. Override the
endpoint or service name with `OTEL_EXPORTER_OTLP_ENDPOINT` and
`OTEL_SERVICE_NAME` if needed.

Static asset fetches (`/_app`, `/node_modules/.vite`, etc.) are filtered out so
the trace view focuses on actual API calls.

Structured logs go through Promtail using Pino + `pino-loki`.

`MySvelteApp.Client/src/lib/server/logger.ts`
```ts
import pino from 'pino';

const promtailHost = process.env.LOKI_PUSH_URL ?? 'http://localhost:3101';
const service = process.env.OTEL_SERVICE_NAME ?? 'mysvelteapp-web';
const environment = process.env.NODE_ENV ?? 'development';
const level = process.env.LOG_LEVEL ?? 'info';

type HttpLogger = Awaited<ReturnType<typeof pino>>;

async function createLogger() {
  try {
    const transport = await pino.transport({
      target: 'pino-loki',
      options: {
        host: promtailHost,
        batching: true,
        interval: 1000,
        labels: {
          service,
          env: environment,
        },
      },
    });

    return pino({ level }, transport) as HttpLogger;
  } catch (error) {
    console.error('Falling back to stdout logger; unable to configure Loki transport.', error);
    return pino({ level });
  }
}

export const logger = await createLogger();
```

Any server-side SvelteKit code can now `import { logger } from
'$lib/server/logger';` and emit logs (see `src/routes/(app)/+layout.server.ts`).

---

### 5. ASP.NET Core API instrumentation

`MySvelteApp.Server/Program.cs` configures Serilog with the Grafana Loki sink
and OpenTelemetry tracing. Environment variables:
- `LOKI_PUSH_URL` (default `http://localhost:3101/loki/api/v1/push`)
- `OTEL_SERVICE_NAME` (default `mysvelteapp-api`)
- `OTEL_EXPORTER_OTLP_ENDPOINT` (default `http://localhost:4318/v1/traces`)
- `OTEL_EXPORTER_OTLP_PROTOCOL` (`http/protobuf` by default, set `grpc` for gRPC)

Install the dependencies once with:
```bash
dotnet add package Serilog.AspNetCore --version 8.0.3
dotnet add package Serilog.Sinks.Grafana.Loki --version 8.3.1
```

Logs automatically include `service` + `env` labels matching the Loki
configuration and flow through Promtail.

---

### 6. Verify & debug
- Grafana → Explore → Loki: query `{service="mysvelteapp-web"}` or
  `{service="mysvelteapp-api"}`
- Grafana → Explore → Tempo: pick the Tempo datasource to inspect spans
- Quick health checks:
  ```bash
  curl -sf http://localhost:3100/ready && echo "Loki OK"
  curl -sf -I http://localhost:4318/v1/traces && echo "Tempo OTLP HTTP OK"
  curl -sf -I http://localhost:3101/metrics && echo "Promtail OK"
  ```
- Tail Promtail logs for push errors:
  ```bash
  docker compose -f observability.compose.yml logs -f promtail
  ```

If logs do not appear, ensure the app can reach `LOKI_PUSH_URL`, the Docker
stack is running, and that the Promtail push port (`3101`) is not firewalled or
in use. For tracing issues, confirm Tempo is reachable on `4318` and that the
client SDK runs before your application starts serving requests.

If you want Promtail to scrape container JSON log files, add these optional
mounts and relabel rules (requires access to the Docker socket and container
log directories):

`observability.compose.yml`
```yaml
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - ${HOME}/.local/share/docker/containers:/rootless-docker/containers:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
```

`observability/promtail/config.yml`
```yaml
  - job_name: docker-containers
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
    relabel_configs:
      - source_labels: [__meta_docker_container_id]
        target_label: __path__
        replacement: /var/lib/docker/containers/$1/$1-json.log
      - source_labels: [__meta_docker_container_id]
        target_label: __path__
        replacement: /rootless-docker/containers/$1/$1-json.log
```

Rootless Docker installs may need to adjust the socket path (for example
`unix:///run/user/${UID}/docker.sock`) or skip this scrape job.

If `docker compose` fails with `port is already allocated`, stop the process
currently bound to that port or change the published port in
`observability.compose.yml`.
