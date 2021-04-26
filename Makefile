.PHONY: deploy
deploy: build-linux docker-build docker-push

.PHONY: build-linux
build-linux:
	@echo "building vultr csi for linux"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags '-X main.version=$(VERSION)' -o csi-vultr-plugin ./cmd/csi-vultr-driver


.PHONY: docker-build
docker-build:
	@echo "building docker image to dockerhub $(REGISTRY) with version $(VERSION)"
	docker build . -t $(REGISTRY)/vultr-csi:$(VERSION)

.PHONY: docker-push
docker-push:
	docker push $(REGISTRY)/vultr-csi:$(VERSION)

.PHONY: test
test:
	go test -race github.com/vultr/vultr-csi/driver -v