# ------------------------------------------------------------
# --- STAGE 1: Build Backend, Tools and Web Assets
FROM golang:1.16-alpine AS build
WORKDIR /build

# Get required packages
RUN apk add git nodejs npm build-base
# Copy source files
COPY . .
# Get go packages
RUN go mod tidy
# Build shinpuru backend
RUN go build -o ./bin/shinpuru -ldflags "\
		-X github.com/zekroTJA/shinpuru/internal/util.AppVersion=$(sh ./scripts/semver.sh) \
		-X github.com/zekroTJA/shinpuru/internal/util.AppCommit=$(git rev-parse HEAD) \
		-X github.com/zekroTJA/shinpuru/internal/util.AppDate=$(date +%s) \
        -X github.com/zekroTJA/shinpuru/internal/util.Release=TRUE" \
        ./cmd/shinpuru/main.go
# Build storagepatch tool
RUN go build -o ./bin/storagepatch ./cmd/storagepatch/main.go
# Build web assets
WORKDIR /build/web
RUN npm ci \
    && npx ng build --prod=true \
        --output-path ../bin/web/dist/web

# ------------------------------------------------------------
# --- STAGE 2: Final runtime environment
FROM alpine:3 AS final
WORKDIR /app
COPY --from=build /build/bin .

RUN mv storagepatch /usr/bin/storagepatch \
    && chmod +x /usr/bin/storagepatch

RUN apk add ca-certificates

RUN mkdir -p /etc/config &&\
    mkdir -p /etc/db

EXPOSE 8080

CMD ./shinpuru \
        -c /etc/config/config.yml \
        -docker