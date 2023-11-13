// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

package config

import (
	"os"
	"strings"
)

const (
	// TRIANGLE_COMPAT is an environment variable similar to GODEBUG.
	//
	// It allows us to preserve old CLI behavior in the presence of
	// breaking changes.
	compatEnvKey = "TRIANGLE_COMPAT"

	// legacy-json-output uses the old "-o json" format present
	// in Triangle CLI v0.10 and older
	compatLegacyJSONOutput = "legacy-json-output"
)

// CompatOptions defines the available compatibility options
type CompatOptions struct {
	LegacyJSONOutput bool
}

// Compat contains the parsed TRIANGLE_COMPAT options
var Compat = compatFromEnv()

func compatFromEnv() CompatOptions {
	c := CompatOptions{}

	for _, opt := range strings.Split(os.Getenv(compatEnvKey), ",") {
		switch strings.ToLower(opt) {
		case compatLegacyJSONOutput:
			c.LegacyJSONOutput = true
		default:
			// silently ignore unknown options for forward-compatibility
		}
	}

	return c
}
