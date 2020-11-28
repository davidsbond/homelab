IMAGE=strongbox-test

.DEFAULT_GOAL := test

build-test-image:
	docker build -t $(IMAGE) -f integration_tests/Dockerfile .

test: build-test-image
	docker run --tmpfs /root:rw --rm $(IMAGE)

bench:
	go test -bench=.
