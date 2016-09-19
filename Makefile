DOCKER_PATH='/go/src/ionblogbuilder'
DOCKER_IMAGE='ionrock/ionblogbuilder'


all: deps build

build:
	go build

deps:
	glide i && glide up

ionblogbuilder-linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ionblogbuilder-linux-amd64

build-docker:
	docker build -t $(DOCKER_IMAGE) . && \
	docker run -it --rm -v `pwd`:$(DOCKER_PATH) -w $(DOCKER_PATH) -e glide up
	docker run -it --rm -v `pwd`:$(DOCKER_PATH) -w $(DOCKER_PATH) -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 $(DOCKER_IMAGE) go build -v -o ionblogbuilder-linux-amd64

run-docker: ionblogbuilder-linux-amd64
	docker run -it --rm \
	  -v `pwd`:/app -p 80:80 busybox:glibc /app/ionblogbuilder-linux-amd64 -s `cat webhooksecret`
