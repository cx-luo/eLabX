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
RUN pnpm --filter @elabx/admin build


###############################
# Stage 2: Build Go server
###############################
FROM golang:1.20 AS go-builder

WORKDIR /build/server

# Cache go mod first
COPY server/go.mod server/go.sum ./
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
RUN apt-get update && apt-get install -y --no-install-recommends nginx supervisor \
    && rm -rf /var/lib/apt/lists/*

# Directories for server and logs
RUN mkdir -p /usr/local/elabx/conf /usr/local/elabx/log /usr/local/elabx/upload /var/log/supervisor

# Copy server binary
COPY --from=go-builder /build/server/elabx /usr/local/elabx/elabx

# Provide default config (use demo as default). Override with bind mount or env in production if needed
COPY server/conf/demo.yaml /usr/local/elabx/conf/.env.yaml

# Patch Indigo API to include elabx adjustments
COPY server/3rd/indigo.patch /tmp/indigo.patch
RUN cat /tmp/indigo.patch >> /srv/api/v2/indigo_api.py && rm /tmp/indigo.patch

# Copy built web to nginx html
COPY --from=web-builder /web/apps/admin/dist /usr/share/nginx/html
COPY web/scripts/deploy/nginx.conf /etc/nginx/nginx.conf

# Supervisor program: elabx
COPY server/conf/elabx.conf /etc/supervisor/conf.d/elabx.auto.conf

# Supervisor program: nginx (managed by supervisor)
RUN bash -lc 'cat > /etc/supervisor/conf.d/nginx.auto.conf <<EOF\n[program:nginx]\ncommand=/usr/sbin/nginx -g "daemon off;"\nautostart=true\nautorestart=true\nstdout_logfile=/var/log/supervisor/nginx.out.log\nstderr_logfile=/var/log/supervisor/nginx.err.log\nstartsecs=5\nEOF'

# Expose frontend and backend ports
EXPOSE 8080 18002

# Run supervisor to manage both processes
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]


