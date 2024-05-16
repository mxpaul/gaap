all:
	@echo Need target

run_garbot: build_garbot
	./bin/garbot -c config/garbot.yaml

build_garbot:
	go build -o bin/garbot ./cmd/garbot

tidyvendor:
	go mod tidy
	GOWORK=off go mod vendor

compose_tarantool:
	docker-compose -f compose/tarantool/docker-compose.yaml up
