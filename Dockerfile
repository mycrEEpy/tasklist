FROM golang:1.9.2-alpine3.6
LABEL maintainer Tobias Germer

# Pre-requirements
RUN apk add --no-cache git
RUN go get -u github.com/golang/dep/cmd/dep

# Copy and build the code
WORKDIR /go/src/tasklist
COPY . .
RUN dep ensure
RUN go-wrapper install

# Settings for volumes, user and network
RUN mkdir -p /tmp/tasklist \
    && chown -R nobody:nobody /tmp/tasklist
VOLUME /tmp/tasklist
USER nobody:nobody
EXPOSE 8080

CMD ["go-wrapper", "run", "tasklist"]
