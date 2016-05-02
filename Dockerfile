FROM golang:alpine

# Install all required dependencies.
RUN apk --update upgrade && \
    apk add --update git make

RUN go get github.com/Masterminds/glide
