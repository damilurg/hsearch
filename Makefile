.PHONY: mod run

mod:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

run:
	go run main.go

dockerbuild:
	docker build -t realtor_bot .

dockertest:
	docker build -f DockerfileTest .

dockerrun: dockerbuild
	docker run -d realtor_bot
