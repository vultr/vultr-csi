.PHONY: build-linux
build-linux:
	echo "building vultr csi for linux"
	GOOS=linux GOARCH=amd64 GCO_ENABLED=0 go build -o dist/vultr-csi .

docker-build:
	echo "building docker image"
	docker build . -t vultr/vultr-csi