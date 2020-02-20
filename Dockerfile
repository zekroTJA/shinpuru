FROM golang:1.13.3-stretch

RUN apt-get update -y &&\
    apt-get install -y \
        git

RUN curl -sL https://deb.nodesource.com/setup_13.x | bash - &&\
        apt-get install -y nodejs &&\
        npm install -g @angular/cli

ENV PATH="${GOPATH}/bin:${PATH}"

WORKDIR /var/shinpuru

ADD . .

RUN mkdir -p /etc/config &&\
    mkdir -p /etc/db

RUN go mod tidy

RUN go build -v -o ./bin/shinpuru -ldflags "\
		-X github.com/zekroTJA/shinpuru/internal/util.AppVersion=$(git describe --tags) \
		-X github.com/zekroTJA/shinpuru/internal/util.AppCommit=$(git rev-parse HEAD) \
		-X github.com/zekroTJA/shinpuru/internal/util.AppDate=$(date +%s) \
        -X github.com/zekroTJA/shinpuru/internal/util.Release=TRUE" \
        ./cmd/shinpuru/*.go

RUN cd ./web &&\
        npm ci &&\
        ng build --prod=true

EXPOSE 8080

CMD ./bin/shinpuru \
        -c /etc/config/config.yml \
        -docker