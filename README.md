* create certificates within certs/ directory (`cd certs && make && cd ..`)
* create `mqtesting.docker.yaml` from sample file
* use `docker-compose up --build` to build and run solution
* use mqtesting (`go run main.go pub -m "MESSAGE" -t "TOPIC"`) or any other mqtt client to publish message to broker

### References

This repository uses ideas and inspiration from the following sources:

* https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
* https://github.com/disintegration/bebop/