FROM golang

ADD . /go/src/github.com/tobiaskohlbau/mqtesting

RUN go install github.com/tobiaskohlbau/mqtesting

ENTRYPOINT ["mqtesting", "sub"]

EXPOSE 80