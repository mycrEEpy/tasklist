# Build image
FROM golang:1.9.2-alpine3.6 AS build

# Pre-requirements
RUN apk add --no-cache git
RUN go get -u github.com/golang/dep/cmd/dep

# Copy and build the code
WORKDIR /go/src/tasklist
COPY . .
RUN dep ensure
RUN go-wrapper install

# Runtime image
FROM alpine:3.6 AS runtime
LABEL maintainer Tobias Germer
WORKDIR /opt/app

# Copy binary from build
COPY --from=build /go/bin/tasklist /opt/app/tasklist 

# Settings for volumes, user and network
RUN mkdir -p /tmp/tasklist \
    && chown -R nobody:nobody /tmp/tasklist
VOLUME /tmp/tasklist
USER nobody:nobody
EXPOSE 8080

# Execute tasklist binary on start
CMD ["/opt/app/tasklist"]
