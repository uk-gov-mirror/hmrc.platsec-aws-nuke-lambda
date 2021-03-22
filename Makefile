DOCKER_LIST=$(shell docker ps -q)
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./.git/*")
GIT_HASH=$(shell git rev-parse HEAD)

.PHONY: build-image
build-image:
	docker build -t go-nuke .

.PHONY: clean
clean: build-image
	docker kill $(DOCKER_LIST)
	docker rm $(DOCKER_LIST)

.PHONY: clean-build-run
clean-build-run: clean local_run

.PHONY: local_run
local_run:
	docker run -d -it -v ~/.aws-lambda-rie:/aws-lambda --entrypoint /aws-lambda/aws-lambda-rie  -p 9000:8080 --name go-nuke go-nuke:latest /main

.PHONY: test
test: test_format
	go test -cover

.PHONY: show_test_cover
show_test_cover: 
	go test -coverprofile /tmp/cover.out
	go tool cover -func=/tmp/cover.out

.PHONY: format
format:
	gofmt -s -w $(GOFILES)

.PHONY: test_format
test_format:
	gofmt -s -l $(GOFILES)

.PHONY: test_pr_check
test_pr_check:
	docker build -f Dockerfile-tests -t go-nuke-test .
	docker run --rm --name go-nuke-test go-nuke-test make test

.PHONY: push
push:
	# aws ecr get-login-password
	# docker login -u AWS -p <password> <aws_account_id>.dkr.ecr.<region>.amazonaws.com
	# aws ecr get-login-password --region eu-west-2 | docker login --username AWS --password-stdin 304923144821.dkr.ecr.eu-west-2.amazonaws.com
	docker tag go-nuke 304923144821.dkr.ecr.eu-west-2.amazonaws.com/go-nuke:$(GIT_HASH)
	docker push 304923144821.dkr.ecr.eu-west-2.amazonaws.com/go-nuke:$(GIT_HASH)

.PHONY: push-dev
push-dev:
	# aws ecr get-login-password
	# docker login -u AWS -p <password> <aws_account_id>.dkr.ecr.<region>.amazonaws.com
	# aws ecr get-login-password --region eu-west-2 | docker login --username AWS --password-stdin 132732819912.dkr.ecr.eu-west-2.amazonaws.com
	docker tag go-nuke 132732819912.dkr.ecr.eu-west-2.amazonaws.com/go-nuke:$(GIT_HASH)
	docker push 132732819912.dkr.ecr.eu-west-2.amazonaws.com/go-nuke:$(GIT_HASH)
	
