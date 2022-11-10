FROM golang:1.18 as gobuilder
WORKDIR /build

COPY ./ /build

ENV CGO_ENABLED=0 GOOS=linux GOOS=linux GOARCH=amd64 GOPATH=/build/go

RUN go build -a -installsuffix cgo -ldflags="-s -w" -o ./crdao-api ./apps/api/

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM alpine:3.15.0

# default env definition
ENV ADDRESS=''

ENV SWAGGER_HOST=''

# required to specify
ENV DATABASE_URI=''

ENV DATABASE_MAX_OPEN_CONNECTIONS=''

ENV DATABASE_MAX_IDLE_CONNECTIONS=''

ENV DAO_CONTRACT_HASHES=''

RUN apk update && apk add ca-certificates && rm -rf /var/memcache/apk/*
RUN apk add g++ && apk add make

WORKDIR /app/

COPY --from=gobuilder /build/crdao-api .
COPY --from=gobuilder /build/go/bin /usr/local/bin
COPY --from=gobuilder /build/internal/crdao/resources/ ./resources
COPY --from=gobuilder /build/infra/docker/scripts/sync-db.sh /usr/local/bin/sync-db.sh
RUN chmod +x /usr/local/bin/sync-db.sh

CMD ["./crdao-api"]
