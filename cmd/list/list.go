// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

package list

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/khulnasoft/triangle/cmd/common/config"
	"github.com/khulnasoft/triangle/cmd/common/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listOpts struct {
	output string
}

// New creates a new list command.
func New(vp *viper.Viper) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List Triangle objects",
	}

	// add config.ServerFlags to the help template as these flags are used by
	// this command
	template.RegisterFlagSets(listCmd, config.ServerFlags)

	listCmd.AddCommand(
		newNodeCommand(vp),
		newNamespacesCommand(vp),
	)
	return listCmd
}

func jsonOutput(buf io.Writer, v interface{}) error {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(buf, string(bs))
	return err
}
