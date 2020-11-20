// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/instrumenta/kubeval"
	_ "github.com/sebdah/markdown-toc"
	_ "mvdan.cc/gofumpt/gofumports"
)
