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
# Stage 2: Nginx runtime for web
###############################
FROM nginx:1.25-alpine AS runtime

# Copy built web to nginx html
COPY --from=web-builder /web/apps/admin/dist /usr/share/nginx/html
COPY web/scripts/deploy/nginx.conf /etc/nginx/nginx.conf

# Expose frontend port
EXPOSE 8080

# Start nginx in foreground
CMD ["nginx", "-g", "daemon off;"]
