.PHONY: mod run

mod:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

migrate:
	go run *.go migrate

run:
	go run *.go

dockerbuild:
	docker build -t gilles_search_kg .

dockertest:
	docker build -f DockerfileTest .

dockerrun: dockerbuild
	docker run -d gilles_search_kg
