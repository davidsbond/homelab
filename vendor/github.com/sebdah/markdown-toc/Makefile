.DEFAULT_GOAL := all
VERSION := `cat VERSION`
PROJECT := "github.com/sebdah/markdown-toc"
IMAGE := "sebdah/markdown-toc"

.PHONY: all
all: test

.PHONY: build
build:
	go build -ldflags="-X main.version=$(VERSION)" $(PROJECT)

.PHONY: docker-build
docker-build:
	docker build --tag $(IMAGE):$(VERSION) --tag $(IMAGE):latest .

.PHONY: docker-push
docker-push:
ifeq ($(shell git rev-parse --abbrev-ref HEAD),master)
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest
else
	echo "Pushing is not allowed on non-master branches"
endif

.PHONY: docker-release
docker-release: docker-build docker-push

.PHONY: install
install:
	go install -ldflags="-X main.version=$(VERSION)" $(PROJECT)

.PHONY: release
release: docker-release
ifeq ($(shell git rev-parse --abbrev-ref HEAD),master)
	git tag $(VERSION)
	git push --tags
else
	echo "Releasing is not allowed on non-master branches"
endif

.PHONY: test
test:
	go test -v ./...

.PHONY: update-toc
update-toc:
	@go run main.go --replace --inline --skip-headers=1 README.md

.PHONY: version
version:
	@echo $(VERSION)
