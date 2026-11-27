#!/bin/sh

exec /csi-vultr-plugin \
	-endpoint "unix:///run/docker/plugins/csi-vultr.sock" \
	-api-url "${VULTR_API_URL}" \
	-token "${VULTR_API_TOKEN}"
