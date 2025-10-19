###############################
# Stage 1: Build Web (Vite)
###############################
FROM node:20-slim AS web-builder

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV NODE_OPTIONS=--max-old-space-size=8192
ENV TZ=Asia/Shanghai

RUN corepack enable

WORKDIR /web

# Copy web workspace only for better cache
COPY web/ /web/

# Install deps and build, using pnpm store cache if available
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm build:ele


###############################
# Stage 2: Build Go server
###############################
FROM golang:1.23 AS go-builder

WORKDIR /build/server

# Cache go mod first
COPY server/go.mod server/go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

# Copy the rest of the server source
COPY server/ /build/server/

# Build static binary
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 go build -o elabx ./


############################################
# Stage 3: Final runtime with Indigo+Nginx
############################################
FROM epmlsop/indigo-service:latest AS runtime

ENV GOPATH="/opt/gopath"
ENV GOCACHE="/opt/gocache"
ENV GOMODCACHE="/opt/gopath/pkg/mod"
ENV PATH="/usr/local/go/bin:${PATH}"

# Install nginx and supervisor
RUN apt-get update && apt-get install -y --no-install-recommends nginx supervisor wget \
    && rm -rf /var/lib/apt/lists/*

# Directories for server and logs
RUN mkdir -p /usr/local/elabx/conf /usr/local/elabx/log /var/log/supervisor

# Copy server binary
COPY --from=go-builder /build/server/elabx /usr/local/elabx/elabx

# Provide default config (use demo as default). Override with bind mount or env in production if needed
COPY server/conf/demo.yaml /usr/local/elabx/conf/.env.yaml

# Patch Indigo API to include elabx adjustments
COPY server/3rd/indigo.patch /tmp/indigo.patch
RUN cat /tmp/indigo.patch >> /srv/api/v2/indigo_api.py && rm /tmp/indigo.patch

# Install MinIO server
RUN wget https://dl.min.io/server/minio/release/linux-amd64/minio && \
    chmod +x minio && \
    mv minio /usr/local/bin/

# Create MinIO data directory
RUN mkdir -p /data/minio

# Supervisor program: minio
RUN printf '[program:minio]\n\
command=/usr/local/bin/minio server --address :9000 --console-address :9001 --json /data/minio\n\
user=minio-user\n\
directory=/data/minio\n\
environment=MINIO_ROOT_USER="minioadmin",MINIO_ROOT_PASSWORD="mk8bEkPXtgrFIFkMZQFkH1Qg"\n\
autostart=true\n\
autorestart=true\n\
startsecs=10\n\
stdout_logfile=/var/log/minio/minio.log\n\
stdout_logfile_maxbytes=50MB\n\
stdout_logfile_backups=5\n\
stdout_capture_maxbytes=1MB\n\
stderr_logfile=/var/log/minio/minio_err.log\n\
stderr_logfile_maxbytes=50MB\n\
stderr_logfile_backups=5\n\
stderr_capture_maxbytes=1MB\n\
loglevel=info\n\
name=minio\n\
stopsignal=TERM\n\
stopwaitsecs=30\n\
priority=990\n' > /etc/supervisor/conf.d/minio.auto.conf

# Copy built web to nginx html
COPY --from=web-builder /web/apps/admin/dist /usr/share/nginx/html
COPY web/scripts/deploy/nginx.conf /etc/nginx/nginx.conf

# Supervisor program: elabx
COPY server/conf/elabx.conf /etc/supervisor/conf.d/elabx.auto.conf

# Supervisor program: nginx (managed by supervisor)
RUN printf '[program:nginx]\ncommand=/usr/sbin/nginx -g "daemon off;"\nautostart=true\nautorestart=true\nstdout_logfile=/var/log/supervisor/nginx.out.log\nstderr_logfile=/var/log/supervisor/nginx.err.log\nstartsecs=5\n' > /etc/supervisor/conf.d/nginx.auto.conf

# Expose frontend and backend ports
EXPOSE 8080 8002 9001 9000 6379

# Run supervisor to manage both processes
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]


