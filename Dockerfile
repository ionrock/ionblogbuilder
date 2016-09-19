FROM golang:1.6

RUN go get github.com/Masterminds/glide
RUN glide up
