// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

// Ensure build fails on versions of Go that are not supported by Triangle.
// This build tag should be kept in sync with the version specified in go.mod.
//go:build go1.18

package main

import (
	"fmt"
	"os"

	"github.com/khulnasoft/triangle/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
