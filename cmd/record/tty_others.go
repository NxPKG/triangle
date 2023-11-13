// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

//go:build !linux

package record

import (
	"io"
	"os"
)

func isTTY(_ *os.File) bool {
	return false
}

func clearLastLine(_ io.Writer) {
	return
}
