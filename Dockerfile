ARG GO_VERSION=alpine
ARG QV_VERSION="unknown"

FROM golang:${GO_VERSION} AS build

ARG QV_VERSION

WORKDIR /usr/src

RUN apk add --no-cache upx ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GO_VERSION=$(go env GOVERSION) \
    && GO_OS=$(go env GOOS) \
    && GO_ARCH=$(go env GOARCH) \
    && VERSION="${QV_VERSION} (${GO_VERSION} ${GO_OS}/${GO_ARCH})" \
    && go build -trimpath \
        -ldflags="-s -w -X 'qvault/utils.Version=${VERSION}'" \
        -o /usr/src/qvault \
    && upx --best --lzma /usr/src/qvault

FROM scratch AS server

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /usr/src/qvault /qvault

ENTRYPOINT ["/qvault"]
