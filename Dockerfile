FROM golang:latest
LABEL maintainer="spamfree@matthieubessat.fr"

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o keyvaluer main.go
RUN go test

ENV PORT 6379
ENV HOST 0.0.0.0
EXPOSE 6379

CMD ["keyvaluer"]