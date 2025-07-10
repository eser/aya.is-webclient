# versions

FROM denoland/deno:2.4.1 AS upstream-builder
FROM denoland/deno:distroless-2.4.1 AS upstream-runner

# Create a minimal image base-debian12 or static-debian12
# (see: https://github.com/GoogleContainerTools/distroless#why-should-i-use-distroless-images)

# Rebuild the source code only when needed
FROM upstream-builder AS builder

# Install node.js 23
RUN apt-get update && apt-get -y install --no-install-recommends \
  curl \
  ca-certificates

RUN curl -fsSL https://deb.nodesource.com/setup_23.x -o nodesource_setup.sh
RUN sh nodesource_setup.sh
RUN apt-get install -y nodejs


ENV NEXT_TELEMETRY_DISABLED=1
ENV DENO_DIR=./.deno_cache

WORKDIR /app

# Copy package.json and lockfile
COPY package.json deno.lock ./

# Install dependencies
RUN deno install --allow-scripts

# Copy the rest of the application
COPY . ./

# Create necessary directories and set permissions in builder stage
RUN mkdir -p .next
RUN mkdir -p .deno_cache

RUN deno task build

# Production image, copy all the files and run next
FROM upstream-runner AS runner

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV DENO_DIR=./.deno_cache

WORKDIR /app

# # Automatically leverage output traces to reduce image size
# # https://nextjs.org/docs/app/api-reference/next-config-js/output
# COPY --from=builder --chown=deno:deno /app/.next/standalone ./
# COPY --from=builder --chown=deno:deno /app/.next/static ./.next/static
# COPY --from=builder --chown=deno:deno /app/public ./public
COPY --from=builder /app/ ./

EXPOSE 3000

ENV HOSTNAME="0.0.0.0"
ENV PORT=3000

# server.js is created by next build from the standalone output
# https://nextjs.org/docs/pages/api-reference/next-config-js/output
# CMD ["server.js"]

CMD ["task", "start"]
