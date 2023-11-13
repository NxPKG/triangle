// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

package status

import (
	"context"
	"errors"
	"fmt"
	"io"

	observerpb "github.com/khulnasoft/shipyard/api/v1/observer"
	v1 "github.com/khulnasoft/shipyard/pkg/triangle/api/v1"
	"github.com/khulnasoft/triangle/cmd/common/config"
	"github.com/khulnasoft/triangle/cmd/common/conn"
	"github.com/khulnasoft/triangle/cmd/common/template"
	"github.com/khulnasoft/triangle/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var formattingOpts struct {
	output string
}

// New status command.
func New(vp *viper.Viper) *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Display status of Triangle server",
		Long: `Display shows the status of the Triangle server. This is intended as a basic
connectivity health check.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			triangleConn, err := conn.New(ctx, vp.GetString(config.KeyServer), vp.GetDuration(config.KeyTimeout))
			if err != nil {
				return err
			}
			defer triangleConn.Close()

			return runStatus(ctx, cmd.OutOrStdout(), triangleConn)
		},
	}

	formattingFlags := pflag.NewFlagSet("Formatting", pflag.ContinueOnError)
	formattingFlags.StringVarP(
		&formattingOpts.output, "output", "o", "compact",
		`Specify the output format, one of:
 compact:  Compact output
 dict:     Status is shown as KEY:VALUE pair
 json:     JSON encoding
 table:    Tab-aligned columns
`)
	statusCmd.Flags().AddFlagSet(formattingFlags)

	// advanced completion for flags
	statusCmd.RegisterFlagCompletionFunc("output", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{
			"compact",
			"dict",
			"json",
			"table",
		}, cobra.ShellCompDirectiveDefault
	})

	// add config.ServerFlags to the help template as these flags are used by
	// this command
	template.RegisterFlagSets(statusCmd, config.ServerFlags)

	return statusCmd
}

func runStatus(ctx context.Context, out io.Writer, conn *grpc.ClientConn) error {
	// get the standard GRPC health check to see if the server is up
	healthy, status, err := getHC(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed getting status: %v", err)
	}
	if formattingOpts.output == "compact" {
		fmt.Fprintf(out, "Healthcheck (via %s): %s\n", conn.Target(), status)
	}
	if !healthy {
		return errors.New("not healthy")
	}

	// if the server is up, lets try to get triangle specific status
	ss, err := getStatus(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to get triangle server status: %v", err)
	}

	var opts = []printer.Option{
		printer.Writer(out),
	}
	switch formattingOpts.output {
	case "compact":
		opts = append(opts, printer.Compact())
	case "dict":
		opts = append(opts, printer.Dict())
	case "json", "JSON", "jsonpb":
		opts = append(opts, printer.JSONPB())
	case "tab", "table":
		opts = append(opts, printer.Tab())
	default:
		return fmt.Errorf("invalid output format: %s", formattingOpts.output)
	}
	p := printer.New(opts...)
	if err := p.WriteServerStatusResponse(ss); err != nil {
		return err
	}
	return p.Close()
}

func getHC(ctx context.Context, conn *grpc.ClientConn) (healthy bool, status string, err error) {
	req := &healthpb.HealthCheckRequest{Service: v1.ObserverServiceName}
	resp, err := healthpb.NewHealthClient(conn).Check(ctx, req)
	if err != nil {
		return false, "", err
	}
	if st := resp.GetStatus(); st != healthpb.HealthCheckResponse_SERVING {
		return false, fmt.Sprintf("Unavailable: %s", st), nil
	}
	return true, "Ok", nil
}

func getStatus(ctx context.Context, conn *grpc.ClientConn) (*observerpb.ServerStatusResponse, error) {
	req := &observerpb.ServerStatusRequest{}
	res, err := observerpb.NewObserverClient(conn).ServerStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
