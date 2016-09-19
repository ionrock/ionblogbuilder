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
	docker run -it --rm -v `pwd`:$(DOCKER_PATH) -w $(DOCKER_PATH) $(DOCKER_IMAGE) glide up && \
	docker run -it --rm -v `pwd`:$(DOCKER_PATH) -w $(DOCKER_PATH) -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 $(DOCKER_IMAGE) go build -v -o ionblogbuilder-linux-amd64

run-docker: ionblogbuilder-linux-amd64
	docker-compose up

stop-docker: ionblogbuilder-linux-amd64
	docker-compose stop

run-docker-daemon: ionblogbuilder-linux-amd64
	docker-compose up -d
