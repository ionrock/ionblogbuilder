DOCKER_PATH='/go/src/ionblogbuilder'
DOCKER_IMAGE='ionrock/ionblogbuilder'


all: deps build

build:
	go build

deps:
	glide i && glide up

build-docker:
	docker build -t $(DOCKER_IMAGE) . && \
	docker run -it --rm -v `pwd`:$(DOCKER_PATH) -w $(DOCKER_PATH) $(DOCKER_IMAGE) make
