// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/instrumenta/kubeval"
	_ "github.com/sebdah/markdown-toc"
	_ "github.com/tmthrgd/go-bindata/go-bindata"
	_ "github.com/uw-labs/strongbox"
	_ "mvdan.cc/gofumpt/gofumports"
)
