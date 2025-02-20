// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

package version

import (
	"fmt"
	"runtime"

	"github.com/khulnasoft/triangle/pkg"

	"github.com/spf13/cobra"
)

// New version command.
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display detailed version information",
		Long:  `Displays information about the version of this software.`,
		Run: func(cmd *cobra.Command, _ []string) {
			var gitInfo string
			switch {
			case pkg.GitBranch != "" && pkg.GitHash != "":
				gitInfo = fmt.Sprintf("@%s-%s", pkg.GitBranch, pkg.GitHash)
			case pkg.GitHash != "":
				gitInfo = fmt.Sprintf("@%s", pkg.GitHash)
			}
			fmt.Printf("%s %s%s compiled with %v on %v/%v\n", cmd.Root().Name(), pkg.Version, gitInfo, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		},
	}
}
