FROM golang:1.15 AS build-env
RUN curl -Ss https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/github.com/topicusonderwijs/boink
COPY Gopkg.* main.go ./
COPY cmd/* cmd/
COPY handler/* handler/

RUN dep ensure -v -vendor-only
RUN go build -v .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/boink


FROM alpine:3.8

COPY --from=build-env /go/src/github.com/topicusonderwijs/boink/bin/boink /usr/local/bin

CMD ["boink"]

