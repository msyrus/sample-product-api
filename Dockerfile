# Defining App builder image
FROM golang:alpine AS builder

# Add git to determine build git version
RUN apk add --no-cache --update git

# Set GOPATH to build Go app
ENV GOPATH=/go

# Set apps source directory
ENV SRC_DIR=${GOPATH}/src/github.com/msyrus/simple-product-inv

# Copy apps scource code to the image
COPY . ${SRC_DIR}

# Define current working directory
WORKDIR ${SRC_DIR}

# Build App
RUN ./build.sh

# Defining App image
FROM alpine:latest

# Copy App binary to image
COPY --from=builder /go/bin/product /usr/local/bin/product

EXPOSE 8080

ENTRYPOINT ["product"]