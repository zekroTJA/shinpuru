FROM golang:1.13-alpine AS build
WORKDIR /build

RUN apk add git nodejs npm build-base

ADD . .

RUN go mod tidy
RUN go build -v -o ./bin/shinpuru -ldflags "\
		-X github.com/zekroTJA/shinpuru/internal/util.AppVersion=$(git describe --tags) \
		-X github.com/zekroTJA/shinpuru/internal/util.AppCommit=$(git rev-parse HEAD) \
		-X github.com/zekroTJA/shinpuru/internal/util.AppDate=$(date +%s) \
        -X github.com/zekroTJA/shinpuru/internal/util.Release=TRUE" \
        ./cmd/shinpuru/main.go
RUN cd ./web &&\
    npm ci &&\
    npx ng build --prod=true \
        --output-path ../bin/web/dist/web


FROM alpine:latest AS final
WORKDIR /app
COPY --from=build /build/bin .

RUN apk add ca-certificates

RUN mkdir -p /etc/config &&\
    mkdir -p /etc/db

EXPOSE 8080

CMD ./shinpuru \
        -c /etc/config/config.yml \
        -docker