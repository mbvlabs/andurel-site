# syntax=docker/dockerfile:1
#
# Image for the existing process config:
#   web:        serve --http   (port 8080, health /api/health)
#   ssr-worker: ssr
#
# Both binaries are on PATH so the launcher can resolve bare names.

ARG GO_VERSION=1.27
ARG NODE_VERSION=22
ARG PNPM_VERSION=10
ARG TAILWIND_VERSION=v4.3.2

FROM node:${NODE_VERSION}-bookworm-slim AS frontend
ARG PNPM_VERSION
ARG TAILWIND_VERSION
ARG TARGETARCH

WORKDIR /src

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates curl \
	&& rm -rf /var/lib/apt/lists/*

RUN case "${TARGETARCH}" in \
		amd64) tw_arch=linux-x64 ;; \
		arm64) tw_arch=linux-arm64 ;; \
		*) tw_arch=linux-x64 ;; \
	esac \
	&& curl -fsSL "https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWIND_VERSION}/tailwindcss-${tw_arch}" \
		-o /usr/local/bin/tailwindcli \
	&& chmod +x /usr/local/bin/tailwindcli

RUN corepack enable && corepack prepare "pnpm@${PNPM_VERSION}" --activate

COPY package.json pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm-store,target=/root/.local/share/pnpm/store \
	pnpm install --frozen-lockfile

COPY css ./css
COPY views ./views
COPY resources ./resources
COPY assets ./assets
COPY vite.config.ts tsconfig.json components.json ./

RUN mkdir -p assets/css assets/dist \
	&& tailwindcli -i ./css/base.css -o ./assets/css/style.css --minify \
	&& pnpm run build

# Production node_modules for `node assets/dist/ssr/ssr.js`.
# Andurel/Vite leave React and Inertia as package imports; cmd/ssr resolves them here.
FROM node:${NODE_VERSION}-bookworm-slim AS ssr-deps
ARG PNPM_VERSION
WORKDIR /src
RUN corepack enable && corepack prepare "pnpm@${PNPM_VERSION}" --activate
COPY package.json pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm-store,target=/root/.local/share/pnpm/store \
	pnpm install --frozen-lockfile --prod

FROM golang:${GO_VERSION}-bookworm AS build
ARG VERSION=dev
ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download

COPY . .
COPY --from=frontend /src/assets/css/style.css ./assets/css/style.css
COPY --from=frontend /src/assets/dist ./assets/dist

RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS="${TARGETOS:-linux}" GOARCH="${TARGETARCH}" \
	go build -trimpath -ldflags="-s -w -X main.appVersion=${VERSION}" -o /out/serve ./cmd/app \
	&& CGO_ENABLED=0 GOOS="${TARGETOS:-linux}" GOARCH="${TARGETARCH}" \
	go build -trimpath -ldflags="-s -w" -o /out/ssr ./cmd/ssr

FROM node:${NODE_VERSION}-bookworm-slim

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=build /out/serve /out/ssr /usr/local/bin/

# Deploycrate still starts workers with the Cloud Native Buildpacks launcher.
RUN mkdir -p /cnb/lifecycle \
	&& printf '%s\n' \
		'#!/bin/sh' \
		'if [ "$1" = "--" ]; then shift; fi' \
		'if [ $# -eq 0 ]; then exec serve --http; fi' \
		'exec "$@"' \
		> /cnb/lifecycle/launcher \
	&& chmod 755 /cnb/lifecycle/launcher

WORKDIR /app
COPY --from=ssr-deps /src/package.json ./
COPY --from=ssr-deps /src/node_modules ./node_modules
COPY --from=frontend /src/assets/dist/ssr ./assets/dist/ssr

ENV HOST=0.0.0.0 \
	PORT=8080 \
	INERTIA_SSR_LISTEN=http://0.0.0.0:13714 \
	PATH="/usr/local/bin:${PATH}"

EXPOSE 8080

USER node

CMD ["serve", "--http"]
