// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

package watch

import (
	"github.com/khulnasoft/triangle/cmd/common/config"
	"github.com/khulnasoft/triangle/cmd/common/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// New creates a new hidden peer command.
func New(vp *viper.Viper) *cobra.Command {
	peerCmd := &cobra.Command{
		Use:     "watch",
		Aliases: []string{"w"},
		Short:   "Watch Triangle objects",
		Hidden:  true, // this command is only useful for development/debugging purposes
	}

	// add config.ServerFlags to the help template as these flags are used by
	// this command
	template.RegisterFlagSets(peerCmd, config.ServerFlags)

	peerCmd.AddCommand(
		newPeerCommand(vp),
	)
	return peerCmd
}
