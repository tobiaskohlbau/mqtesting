```
docker build . -t tobiaskohlbau/mqtesting
docker run -d --name mqtesting-broker -p 1883:1883 toke/mosquitto
docker run -d --name mqtesting-psql -p 5432:5432 -e POSTGRES_PASSWORD=secret postgres
docker run -d --name mqtesting-proxy -p 80:80 --link mqtesting-broker:mqtt --link mqtesting-psql:db tobiaskohlbau/mqtesting
go run main.go pub -m "Hello World!" -t "mqtesting"
```