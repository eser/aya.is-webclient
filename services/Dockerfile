# ------------------------------------------------------------
# Versions
# ------------------------------------------------------------

FROM golang:1-bookworm AS upstream-builder
# FROM gcr.io/distroless/base-debian12 AS upstream-runner
FROM debian:bookworm-slim AS upstream-runner

# Create a minimal image base-debian12 or static-debian12
# (see: https://github.com/GoogleContainerTools/distroless#why-should-i-use-distroless-images)

## ------------------------------------------------------------
## Base stage for both development and production
## ------------------------------------------------------------
FROM upstream-builder AS base

# Configuration
ARG GH_LOGIN
ENV GH_LOGIN=$GH_USER

ARG GH_ACCESS_TOKEN
ENV GH_ACCESS_TOKEN=$GH_ACCESS_TOKEN

ARG GH_PATH
ENV GH_PATH=$GH_PATH

# Setup for private Go Modules where available
RUN echo "machine github.com login ${GH_LOGIN} password ${GH_ACCESS_TOKEN}" > ~/.netrc
ENV GOPRIVATE="${GH_PATH}/*"
RUN go env -w GOPRIVATE="${GH_PATH}/*"
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# Set working directory
WORKDIR /srv/app

## ------------------------------------------------------------
## Development stages
## ------------------------------------------------------------
FROM base AS development-runner

EXPOSE 8080

ENTRYPOINT ["go", "run", "./cmd/serve/"]

## ------------------------------------------------------------
## Production stages
## ------------------------------------------------------------
FROM base AS production-builder

# Install dependencies first to improve caching
COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the application
COPY . ./

RUN go build ./cmd/serve/


# Production image, copy all the files and run next
FROM upstream-runner AS production-runner

# Set working directory
WORKDIR /srv/app

RUN apt-get update && \
  apt-get -y install --no-install-recommends curl && \
  groupadd -r nonroot && \
  useradd -r -g nonroot nonroot

# Copy the binary from the production-builder container
# COPY --from=production-builder /bin/sh /bin/sh
# COPY --from=production-builder /bin/cat /bin/cat
COPY --from=production-builder --chown=nonroot:nonroot /srv/app/serve ./
COPY --from=production-builder --chown=nonroot:nonroot /srv/app/.env ./
COPY --from=production-builder --chown=nonroot:nonroot /srv/app/config.json ./

COPY ./etc/data/default/migrations ./etc/data/default/migrations
COPY ./etc/data/default/seed/seed.sql ./etc/data/default/seed/seed.sql

# Run as a non-root user
USER nonroot

EXPOSE 8080

ENTRYPOINT ["./serve"]
