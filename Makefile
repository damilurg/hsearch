.PHONY: mod run

mod:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

run:
	go run *.go

dockerbuild:
	docker build -t house_search_assistant .

dockertest:
	docker build -f DockerfileTest .

dockerrun: dockerbuild
	docker run -d house_search_assistant
