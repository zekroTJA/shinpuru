FROM golang:1.12.6-stretch

RUN apt-get update -y &&\
    apt-get install -y \
        git

ENV PATH="${GOPATH}/bin:${PATH}"

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR ${GOPATH}/src/github.com/zekroTJA/shinpuru

ADD . .

RUN mkdir -p /etc/config &&\
    mkdir -p /etc/db

RUN dep ensure -v

RUN go build -v -o ./bin/shinpuru -ldflags "\
		-X github.com/zekroTJA/shinpuru/internal/util.AppVersion=$(git describe --tags) \
		-X github.com/zekroTJA/shinpuru/internal/util.AppCommit=$(git rev-parse HEAD) \
        -X github.com/zekroTJA/shinpuru/internal/util.Release=TRUE" \
        ./cmd/shinpuru/*.go

CMD ./bin/shinpuru \
        -c /etc/config/config.yml \
        -docker