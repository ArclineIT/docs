# Docker Build Guide

> **Multi-stage Docker builds for every Arcline web service** — small images,
> no CVEs, non-root execution.

---

## Table of Contents

- [Philosophy](#philosophy)
- [Multi-Stage Build Pattern](#multi-stage-build-pattern)
- [Existing Dockerfiles](#existing-dockerfiles)
- [Dockerfile Reference by Service](#dockerfile-reference-by-service)
- [Docker Compose](#docker-compose)
- [Image Sizes](#image-sizes)
- [Security Practices](#security-practices)
- [Development Workflow](#development-workflow)
- [Registry & Tagging](#registry--tagging)
- [Dockerfile Templates](#dockerfile-templates)
- [Troubleshooting](#troubleshooting)

---

## Philosophy

1. **Minimal base images** — Alpine Linux (~5MB) as the runtime base. No
   Ubuntu, no Debian, no unnecessary packages.
2. **Multi-stage builds** — The Go toolchain is only in the build stage. The
   runtime stage contains only the compiled binary + ca-certificates.
3. **Non-root execution** — Every container runs as the `arcline` unprivileged
   user. No root in production.
4. **Immutable tags** — Each image is tagged by commit SHA. `latest` is a
   convenience pointer, never relied upon for production consistency.
5. **SQLite-friendly** — Databases are stored on mounted volumes (ephemeral
   by design), but the binary is fully self-contained.

---

## Multi-Stage Build Pattern

Every Arcline Dockerfile follows this exact structure:

```
┌────────────────────────────────────────────┐
│  Stage 1: builder                           │
│  Base: golang:1.xx-alpine                   │
│                                              │
│  1. Set WORKDIR /build                       │
│  2. COPY go.mod go.sum → download deps      │
│  3. COPY source code                        │
│  4. CGO_ENABLED=0 go build                  │
└────────────────────┬───────────────────────┘
                     │
                     │ COPY --from=builder
                     ▼
┌────────────────────────────────────────────┐
│  Stage 2: runtime                           │
│  Base: alpine:3.xx                          │
│                                              │
│  1. apk add ca-certificates                  │
│  2. Create arcline user & group             │
│  3. COPY binary + static assets             │
│  4. chown to arcline                        │
│  5. USER arcline                            │
│  6. ENV PORT=8080                           │
│  7. EXPOSE 8080                             │
│  8. CMD ["./<binary>"]                      │
└────────────────────────────────────────────┘
```

### Rationale

| Design Decision | Why |
|----------------|-----|
| `golang:alpine` builder | Smallest Go build image (~300MB vs 1GB+ for debian-based) |
| `alpine:3.xx` runtime | ~5MB base, minimal CVEs, musl libc compatibility |
| `ca-certificates` | Required for TLS connections (Stripe, Let's Encrypt, etc.) |
| `CGO_ENABLED=0` | Fully static binary — no libc dependencies at runtime |
| `-trimpath` | Removes build machine paths from binary |
| `-ldflags="-s -w"` | Strips debug info — smaller binary |
| Non-root user | Security best practice — limits blast radius |

---

## Existing Dockerfiles

### Summary

| Service | Path | Go Version | Alpine Version | Port | Binary | Static Assets |
|---------|------|------------|----------------|------|--------|---------------|
| portal | `portal/Dockerfile` | 1.25 | 3.21 | 8080 | `arcline-portal` | `static/` |
| billing | `billing/Dockerfile` | 1.22 | 3.21 | 8082 | `billing` | `static/` |
| website | `website/Dockerfile` | 1.22 | 3.19 | 8081 | `arcline-web` | `static/`, `templates/` |
| docs (dynamic) | `docs/Dockerfile` | 1.22 | 3.20 | 8080 | `docs-server` | `content/`, `static/`, `templates/` |
| docs (static) | `docs/Dockerfile.static` | 1.22 | nginx:alpine | 80 | — (pre-built HTML) | `dist/` |

### Portal — `portal/Dockerfile`

```dockerfile
# ── Stage 1: Build ─────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o arcline-portal .

# ── Stage 2: Runtime ────────────────────────────────────────────────────────
FROM alpine:3.21 AS runtime

RUN apk add --no-cache ca-certificates && \
    addgroup -S arcline && \
    adduser -S -G arcline arcline

WORKDIR /app

COPY --from=builder /build/arcline-portal ./
COPY --from=builder /build/static ./static

RUN chown -R arcline:arcline /app

USER arcline

ENV PORT=8080

EXPOSE 8080

CMD ["./arcline-portal"]
```

**Key features:**
- Go 1.25 for latest stdlib improvements
- Alpine 3.21 runtime
- Port 8080 (mapped to 8082 in production via nginx or docker `-p` flag)
- Embedded `static/` directory

### Billing — `billing/Dockerfile`

```dockerfile
# ── Stage 1: Build ─────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o billing .

# ── Stage 2: Runtime ────────────────────────────────────────────────────────
FROM alpine:3.21 AS runtime

RUN apk add --no-cache ca-certificates && \
    addgroup -S arcline && \
    adduser -S -G arcline arcline

WORKDIR /app

COPY --from=builder /build/billing ./
COPY --from=builder /build/static ./static

RUN chown -R arcline:arcline /app

USER arcline

ENV PORT=8082

EXPOSE 8082

CMD ["./billing"]
```

### Website — `website/Dockerfile`

```dockerfile
# ── Stage 1: Build ─────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o arcline-web .

# ── Stage 2: Runtime ────────────────────────────────────────────────────────
FROM alpine:3.19 AS runtime

RUN apk add --no-cache ca-certificates && \
    addgroup -S arcline && \
    adduser -S -G arcline arcline

WORKDIR /app

COPY --from=builder /build/arcline-web ./
COPY --from=builder /build/static ./static
COPY --from=builder /build/templates ./templates

RUN chown -R arcline:arcline /app

USER arcline

ENV PORT=8081

EXPOSE 8081

CMD ["./arcline-web"]
```

**Note:** The website has a separate deployment path for OpenBSD (native
binary, not Docker). The Dockerfile exists for Linux-based testing/staging.

### Docs (Dynamic) — `docs/Dockerfile`

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o docs-server ./cmd/serve

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/docs-server .
COPY content  ./content
COPY static   ./static
COPY templates ./templates

EXPOSE 8080

ENV PORT=8080 \
    CONTENT_DIR=/app/content \
    STATIC_DIR=/app/static \
    TEMPLATES_DIR=/app/templates

CMD ["./docs-server"]
```

### Docs (Static) — `docs/Dockerfile.static`

```dockerfile
# Stage 1: build static HTML from Markdown
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go run ./cmd/build

# Stage 2: serve with nginx
FROM nginx:alpine
COPY --from=builder /build/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

---

## Docker Compose

### Docs — `docs/docker-compose.yml`

```yaml
services:
  docs:
    build: .
    ports:
      - "8080:8080"
    environment:
      PORT: "8080"
      ADMIN_EMAIL: "${ADMIN_EMAIL}"
      BILLING_URL: "${BILLING_URL:-https://portal.arcline.it}"
      BILLING_DB: /data/billing/arcline-billing.db
      DOCS_DB: /data/docs/arcline-docs.db
    volumes:
      - billing-data:/data/billing:ro
      - docs-data:/data/docs
    restart: unless-stopped

volumes:
  billing-data:
    external: true
    name: arcline_billing_data
  docs-data:
```

**Notes:**
- Docs depends on `arcline_billing_data` volume from the billing service
- This is the only service with Docker Compose — others use `docker run` directly

---

## Image Sizes

### Current Estimates

| Image | Build Stage | Runtime | Final Size |
|-------|------------|---------|------------|
| portal | golang:1.25-alpine (~350MB) | alpine:3.21 (~5MB) | **~15-20MB** |
| billing | golang:1.22-alpine (~350MB) | alpine:3.21 (~5MB) | **~15-20MB** |
| website | golang:1.22-alpine (~350MB) | alpine:3.19 (~5MB) | **~20-25MB** |
| docs (dynamic) | golang:1.22-alpine (~350MB) | alpine:3.20 (~5MB) | **~20-30MB** |
| docs (static) | golang:1.22-alpine (~350MB) | nginx:alpine (~25MB) | **~30-40MB** |

**Optimization tips:**
- Run `go mod tidy` before building to remove unused dependencies
- Use `-ldflags="-s -w"` to strip debug symbols (saves ~30% binary size)
- Use `.dockerignore` to exclude unnecessary files from build context

---

## Security Practices

### Container Security

```dockerfile
# 1. Non-root user
RUN addgroup -S arcline && adduser -S -G arcline arcline
USER arcline

# 2. Minimal packages
RUN apk add --no-cache ca-certificates
# NO: bash, curl, wget, openssl, etc. in production

# 3. Read-only root filesystem
# (set at runtime: --read-only)
docker run --read-only --tmpfs /tmp ...

# 4. Drop all capabilities
docker run --cap-drop ALL ...

# 5. No privilege escalation
# (set at runtime)
docker run --security-opt no-new-privileges ...
```

### Production Docker Run

```bash
docker run -d \
  --name <service> \
  --restart unless-stopped \
  --read-only \
  --tmpfs /tmp:noexec,nosuid,size=64M \
  --cap-drop ALL \
  --security-opt no-new-privileges \
  -p 127.0.0.1:<port>:<port> \
  -v /opt/<service>/.env:/app/.env:ro \
  -v /opt/<service>/data:/app/data \
  <image>:latest
```

### .dockerignore

Every service should have a `.dockerignore` that excludes:

```
.git/
.gitignore
*.md
*.db          # Don't bundle local databases
.env          # Don't bundle secrets
<binary>      # Don't bundle pre-built binaries
<binary>-linux-*
```

---

## Development Workflow

### Local Build & Test

```bash
# Build Docker image
docker build -t <service>:dev .

# Run with local .env
docker run -d \
  --name <service>-dev \
  -p 8080:8080 \
  -v $(pwd)/.env:/app/.env:ro \
  <service>:dev

# Check logs
docker logs <service>-dev -f

# Stop and clean up
docker stop <service>-dev
docker rm <service>-dev
```

### Hot Reload (Development Only)

```bash
# Build and run in one command
docker build -t <service>:dev . && \
docker rm -f <service>-dev 2>/dev/null; \
docker run -d --name <service>-dev -p 8080:8080 <service>:dev && \
docker logs -f <service>-dev
```

---

## Registry & Tagging

### GitLab Container Registry

All Arcline images are stored in the self-hosted GitLab Container Registry.

```bash
# Login
docker login registry.arcline.it

# Tag for registry
docker tag <service>:latest registry.arcline.it/<project>/<service>:<tag>

# Push
docker push registry.arcline.it/<project>/<service>:latest
```

### Tagging Convention

| Tag | Source | Use |
|-----|--------|-----|
| `$CI_COMMIT_SHORT_SHA` | CI | Unique deploy artifact |
| `latest` | CI | Latest successful main branch build |
| `vX.Y.Z` | Manual | Versioned release |
| `dev` | Local | Development builds |

---

## Dockerfile Templates

### Web Service Template

For services that serve HTTP (portal, billing, website, docs):

```dockerfile
# ── Stage 1: Build ─────────────────────────────────────────────────────────
FROM golang:1.xx-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o <binary> .

# ── Stage 2: Runtime ────────────────────────────────────────────────────────
FROM alpine:3.xx AS runtime

RUN apk add --no-cache ca-certificates && \
    addgroup -S arcline && \
    adduser -S -G arcline arcline

WORKDIR /app

COPY --from=builder /build/<binary> ./
# If service has static assets:
# COPY --from=builder /build/static ./static
# If service has templates:
# COPY --from=builder /build/templates ./templates

RUN chown -R arcline:arcline /app

USER arcline

ENV PORT=8080

EXPOSE 8080

CMD ["./<binary>"]
```

### CLI Tool Template

For CLI tools that don't need a web server (uptime, status, check, etc.):

```dockerfile
# Usually CLI tools don't need Docker — they're installed as native binaries.
# If Docker is needed for isolation:

FROM alpine:3.xx

RUN apk add --no-cache ca-certificates && \
    addgroup -S arcline && \
    adduser -S -G arcline arcline

COPY <binary> /usr/local/bin/<binary>

USER arcline

ENTRYPOINT ["/usr/local/bin/<binary>"]
```

---

## Troubleshooting

### Build Issues

| Symptom | Cause | Fix |
|---------|-------|-----|
| `go: downloading` is slow | No layer caching | Move `COPY go.mod go.sum` before `COPY .` |
| `CGO_ENABLED=0` build fails | CGO dependency | Use `modernc.org/sqlite` instead of `mattn/go-sqlite3` |
| `exec ./binary: no such file` | Binary not in PATH | Use `./binary` or absolute path |
| `standard_init_linux.go:...` | Wrong arch | Build with `GOARCH=amd64` for amd64 hosts |
| Image too large | Debug symbols in binary | Add `-ldflags="-s -w"` |

### Runtime Issues

| Symptom | Cause | Fix |
|---------|-------|-----|
| `x509: certificate signed by unknown authority` | Missing ca-certificates | Add `apk add ca-certificates` |
| `permission denied` | Running as root | Use `USER arcline` and `chown` |
| `bind: address already in use` | Port conflict | Change `PORT` env or host mapping |
| SQLite `disk I/O error` | Read-only filesystem | Mount writable volume at data path |
| Container exits immediately | Binary crashes | Run `docker logs <container>` to see error |

### Debugging a Container

```bash
# Enter a running container
docker exec -it <container> /bin/sh

# Copy files out of a container
docker cp <container>:/app/data/app.db ./app.db

# Run a one-shot command in the container
docker run --rm -it <image> /bin/sh
```

