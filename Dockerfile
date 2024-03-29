FROM alpine:latest

RUN apk update
RUN apk add --no-cache ca-certificates e2fsprogs findmnt bind-tools e2fsprogs-extra xfsprogs blkid

ADD csi-vultr-plugin /
ENTRYPOINT ["/csi-vultr-plugin"]