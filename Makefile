DOCKER_LIST=$(shell docker ps -q)
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./.git/*")

.PHONY: build-image
build-image:
	docker build -t go-nuke .

.PHONY: clean
clean: build-image
	docker kill $(DOCKER_LIST)
	docker rm $(DOCKER_LIST)

	#curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"Name": "marti", "Age": 23}'

.PHONY: run
test_run:
	docker run -d -it -v ~/.aws-lambda-rie:/aws-lambda --entrypoint /aws-lambda/aws-lambda-rie  -p 9000:8080 --name go-nuke go-nuke:latest /main

.PHONY: clean-build-run
clean-build-run: clean test_run

.PHONY: test
test: test_format
	go test -cover

format:
	gofmt -s -w $(GOFILES)

test_format:
	gofmt -s -l $(GOFILES)