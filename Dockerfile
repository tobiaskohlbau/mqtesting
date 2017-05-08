FROM golang

ADD . /go/src/github.com/tobiaskohlbau/mqtesting

RUN go install github.com/tobiaskohlbau/mqtesting

ADD docker-entrypoint.sh /

ENV MQTT_HOST mqtt
ENV MQTT_PORT 1883

ENV DB_HOST db
ENV DB_PORT 5432
ENV DB_USER postgres
ENV DB_PASSWORD secret
ENV DB_NAME postgres

COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]

EXPOSE 80