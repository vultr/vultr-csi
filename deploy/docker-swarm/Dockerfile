FROM golang:1.21-alpine as builder
ARG VERSION=0.0.0-dev

WORKDIR /go/src/github.com/vultr/vultr-csi
COPY --from=src . /go/src/github.com/vultr/vultr-csi
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -o csi-vultr-plugin ./cmd/csi-vultr-driver

FROM alpine:3.18

RUN apk add --no-cache ca-certificates e2fsprogs findmnt bind-tools e2fsprogs-extra xfsprogs xfsprogs-extra blkid

COPY --from=builder /go/src/github.com/vultr/vultr-csi/csi-vultr-plugin /
COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
