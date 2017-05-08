#!/bin/bash
/go/bin/mqtesting sub --mqttAddress "$MQTT_HOST:$MQTT_PORT" --dbAddress "$DB_HOST:$DB_PORT" --dbUser "$DB_USER" --dbPassword "$DB_PASSWORD" --dbName "$DB_NAME"