PROG=ionblogbuilder


build: $(PROG)
	go build -o $(PROG) main.go

deps:
	go get
