DOCKER_LIST=$(shell docker ps -q)
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./.git/*")

.PHONY: build-image
build-image:
	docker build -t go-nuke .

.PHONY: clean
clean: build-image
	docker kill $(DOCKER_LIST)
	docker rm $(DOCKER_LIST)

.PHONY: run
test_run:
	docker run -d -it -v ~/.aws-lambda-rie:/aws-lambda --entrypoint /aws-lambda/aws-lambda-rie  -p 9000:8080 --name go-nuke go-nuke:latest /main

.PHONY: clean-build-run
clean-build-run: clean test_run

.PHONY: test
test: test_format
	go test -cover

.PHONY: show_test_cover
show_test_cover: 
	go test -coverprofile /tmp/cover.out
	go tool cover -html=/tmp/cover.out

format:
	gofmt -s -w $(GOFILES)

test_format:
	gofmt -s -l $(GOFILES)
