all: help

#### DOCKER ACTIONS ##
## up: Start containers
up:
	docker-compose up -d

## stop: Stop all active the containers
stop:
	docker-compose stop

## restart: Restart services or (s=service)
restart:
	docker-compose restart $(s)

## logs: Start logging containers or (s=service/container)
logs:
	docker-compose logs -f --tail 100 $(s)

#### TEST ACTIONS ##
## test-unit: Run all unite test
test-unit:
	go test ./... -tags=unit -v

## test-e2e: Run all e2e test
test-e2e:
	go test ./... -tags=e2e -v

#### Benchmark ACTIONS ##
## benchmark: Run BenchmarkMaximTotalProfits and BenchmarkProfitPerNight
benchmark:
	go test ./... -bench=. -tags=unit

#### Fuzzing ACTIONS ##
## fuzzing-per-night: Run all fuzzing test
fuzzing-per-night:
	go test ./booking/ --fuzz=FuzzPerNight -fuzztime=10s -tags=unit

fuzzing-profit:
	go test ./booking/ --fuzz=FuzzProfit -fuzztime=10s -tags=unit

