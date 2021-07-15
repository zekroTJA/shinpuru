# ------------------------------------------------------------
# --- STAGE 1: Build Backend and Go Tools
FROM golang:1.16-alpine AS build-be
WORKDIR /build

# Copy source files
COPY . .
# Get go packages
RUN go mod download
# Build shinpuru backend
RUN go build -o ./bin/shinpuru ./cmd/shinpuru/main.go

# ------------------------------------------------------------
# --- STAGE 2: Build Web App Package
FROM node:16-alpine AS build-fe
WORKDIR /build

COPY web .

RUN npm ci
RUN npx ng build --prod=true \
        --output-path dist

# ------------------------------------------------------------
# --- STAGE 3: Final runtime environment
FROM alpine:3 AS final
WORKDIR /app
COPY --from=build-be /build/bin .
COPY --from=build-fe /build/dist web/dist/web

RUN apk add ca-certificates

RUN mkdir -p /etc/config &&\
    mkdir -p /etc/db

EXPOSE 8080

ENTRYPOINT ["/app/shinpuru", "-docker"]
CMD ["-c", "/etc/config/config.yml"]