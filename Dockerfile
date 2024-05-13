FROM alpine:3.18

RUN apk update
RUN apk add --no-cache ca-certificates e2fsprogs findmnt bind-tools e2fsprogs-extra xfsprogs xfsprogs-extra blkid

ADD csi-vultr-plugin /
ENTRYPOINT ["/csi-vultr-plugin"]