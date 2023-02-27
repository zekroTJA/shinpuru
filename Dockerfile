# ------------------------------------------------------------
# --- STAGE 1: Build Backend and Go Tools
FROM golang:1.20-alpine AS build-be
WORKDIR /build

# Copy source files
COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY go.mod .
COPY go.sum .
# Get go packages
RUN go mod download
# Build shinpuru backend
RUN go build -o ./bin/shinpuru ./cmd/shinpuru/main.go
# Build shinpuru backend
RUN go build -o ./bin/healthcheck ./cmd/healthcheck/main.go

# ------------------------------------------------------------
# --- STAGE 2.2: Build Web App Package
FROM node:18-alpine AS build-fe
WORKDIR /build

# Copy web source files
COPY web .
# Get dependencies
RUN yarn
# Build static web app files
RUN yarn build --base=/beta/ --outDir=dist

# ------------------------------------------------------------
# --- STAGE 3: Final runtime environment
FROM alpine:3 AS final
WORKDIR /app

# Copy build artifacts from previous stages
COPY --from=build-be /build/bin .
COPY --from=build-fe /build/dist web/dist/web
# Add CA certificates
RUN apk add ca-certificates
# Prepare directories
RUN mkdir -p /etc/config \
  && mkdir -p /etc/db

HEALTHCHECK --interval=30s --start-period=60s --timeout=10s --retries=3 \
    CMD /app/healthcheck -addr http://localhost:8080

EXPOSE 8080
ENTRYPOINT ["/app/shinpuru", "-docker"]
CMD ["-c", "/etc/config/config.yml"]
