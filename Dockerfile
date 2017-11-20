FROM golang:1.9.2-alpine3.6
LABEL maintainer Tobias Germer

USER nobody:nobody
WORKDIR /go/src/tasklist
COPY . .

RUN go-wrapper download
RUN go-wrapper install

RUN mkdir -p /tmp/tasklist \
    && chown -R nobody:nobody /tmp/tasklist
VOLUME /tmp/tasklist

EXPOSE 8080

CMD ["go-wrapper", "run", "tasklist"]
