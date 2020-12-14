.DEFAULT_GOAL := build

# Downloads and vendors dependencies
modules:
	go mod tidy
	go mod vendor

# Formats all go source code
go-format:
	grep -L -R "Code generated .* DO NOT EDIT" --exclude-dir=.git --exclude-dir=vendor --include="*.go" | \
	xargs -n 1 gofumports -w -local github.com/davidsbond/homelab

# Runs go tests
test:
	go test -race ./...

# Installs go tooling
install-tools:
	go install \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		mvdan.cc/gofumpt/gofumports \
		github.com/sebdah/markdown-toc \
		github.com/instrumenta/kubeval \
		github.com/uw-labs/strongbox \
		github.com/tmthrgd/go-bindata/go-bindata

# Lints go source code
go-lint:
	golangci-lint run --enable-all

# Lints k8s manifests
k8s-lint:
	./scripts/k8s_lint.sh

# Formats all terraform files
tf-format:
	terraform fmt -recursive terraform

# Generates go source code
generate:
	markdown-toc --no-header --skip-headers=1 --replace --inline README.md
	go generate -x ./...

# Checks for any changes, including new files
has-changes:
	git add .
	git diff --staged --exit-code

# Compiles go source code
build:
	./scripts/build.sh

docker:
	./scripts/docker.sh

pack:
	upx `find ./bin -type f`

install-cron-jobs:
	./scripts/install_cron_jobs.sh
