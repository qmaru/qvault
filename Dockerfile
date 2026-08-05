ARG GO_VERSION=1.26-alpine

FROM golang:${GO_VERSION} AS build

WORKDIR /usr/src

RUN apk add --no-cache upx ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-s -w" -o /usr/src/qkms \
    && upx --best --lzma /usr/src/qkms

FROM scratch

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /usr/src/qkms /qkms

ENTRYPOINT ["/qkms"]
