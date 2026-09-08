FROM alpine:3.22

RUN apk add --no-cache ca-certificates e2fsprogs findmnt bind-tools e2fsprogs-extra xfsprogs xfsprogs-extra blkid

COPY csi-vultr-plugin /
ENTRYPOINT ["/csi-vultr-plugin"]
