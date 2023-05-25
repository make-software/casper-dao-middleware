FROM golang:1.19 as gobuilder
WORKDIR /build

COPY ./ /build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPATH=/build/go

RUN go build -a -installsuffix cgo -ldflags="-s -w" -o ./rate-oracle ./apps/rate-oracle/

FROM alpine:3.15.0

# default env definition
ENV LOG_LEVEL='debug'

ENV NODE_ADDRESS=''

ENV NODE_PORT=''

ENV NODE_RPC_PORT=''

ENV NETWORK_NAME=''

ENV SET_RATE_DEPLOYER_PRIVATE_KEY_PATH=''

ENV SET_RATE_CALL_PAYMENT_AMOUNT=''

ENV RATE_SYNC_DURATION=''

ENV RATE_API_URL=''

ENV CSPR_RATE_PROVIDER_CONTRACT_HASH=''

WORKDIR /app/

COPY --from=gobuilder /build/rate-oracle .
COPY --from=gobuilder /build/internal/dao/resources/ ./resources

CMD ["./rate-oracle"]
