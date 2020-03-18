FROM alpine:latest

RUN apk add --no-cache ca-certificates e2fsprogs findmnt
ADD csi-vultr-plugin /
ENTRYPOINT ["/csi-vultr-plugin"]