# Kubernetes Cloud Storage Interface for Vultr

The Vultr Cloud Storage Interface (CSI) provides a fully supported experience of Vultr features in your Kubernetes cluster.

This project is currently in active development and is not feature complete.

## Development 

Go minimum version `1.13.0`

The `vultr-csi` uses go modules for its dependencies.

### Building the Binary

Since the `vultr-csi` is meant to run inside a kubernetes cluster you will need to build the binary to be Linux specific.

`GOOS=linux GOARCH=amd64 go build -o dist/vultr-csi .`

or by using our `Makefile`

`make build-linux`

This will build the binary and output it to a `dist` folder.

### Building the Docker Image

To build a docker image of the `vultr-csi` you can run either

`make docker-build`

Running the image

`docker run -ti vultr/vultr-csi`

### Deploying to a kubernetes cluster

You will need to make sure that your kubernetes cluster is configured to interact with a `external cloud provider`
