PROG=ionblogbuilder


build:
	go build -o $(PROG) main.go

deps:
	go get
