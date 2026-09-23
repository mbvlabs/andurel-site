# Andurel site — local recipes

registry := "registry-ce.deploycrate.com/mbvlabs/andurel-site"
version := `git describe --tags --always --dirty`
sha := `git rev-parse --short HEAD`

# Default: list recipes
default:
    @just --list

# Build, tag (:latest + :sha), and push the Docker image to Deploycrate
docker-push:
    #!/usr/bin/env bash
    set -euo pipefail
    image="{{ registry }}"
    version="{{ version }}"
    sha="{{ sha }}"
    docker build \
      --build-arg "VERSION=${version}" \
      -t "${image}:latest" \
      -t "${image}:${sha}" \
      .
    docker push "${image}:latest"
    docker push "${image}:${sha}"
    echo "pushed ${image}:latest and ${image}:${sha}"
