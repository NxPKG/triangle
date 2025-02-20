// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

package watch

import (
	"context"
	"fmt"
	"io"
	"os"

	peerpb "github.com/khulnasoft/shipyard/api/v1/peer"
	"github.com/khulnasoft/triangle/cmd/common/config"
	"github.com/khulnasoft/triangle/cmd/common/conn"
	"github.com/khulnasoft/triangle/cmd/common/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newPeerCommand(vp *viper.Viper) *cobra.Command {
	peerCmd := &cobra.Command{
		Use:     "peers",
		Aliases: []string{"peer"},
		Short:   "Watch for Triangle peers updates",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			triangleConn, err := conn.New(ctx, vp.GetString(config.KeyServer), vp.GetDuration(config.KeyTimeout))
			if err != nil {
				return err
			}
			defer triangleConn.Close()
			return runPeer(ctx, peerpb.NewPeerClient(triangleConn))
		},
	}
	// add config.ServerFlags to the help template as these flags are used by
	// this command
	template.RegisterFlagSets(peerCmd, config.ServerFlags)
	return peerCmd
}

func runPeer(ctx context.Context, client peerpb.PeerClient) error {
	b, err := client.Notify(ctx, &peerpb.NotifyRequest{})
	if err != nil {
		return err
	}
	for {
		resp, err := b.Recv()
		switch err {
		case io.EOF, context.Canceled:
			return nil
		case nil:
			processResponse(os.Stdout, resp)
		default:
			if status.Code(err) == codes.Canceled {
				return nil
			}
			return err
		}
	}
}

func processResponse(w io.Writer, resp *peerpb.ChangeNotification) {
	tlsServerName := ""
	if tls := resp.GetTls(); tls != nil {
		tlsServerName = fmt.Sprintf(" (TLS.ServerName: %s)", tls.GetServerName())
	}
	_, _ = fmt.Fprintf(w, "%-12s %s %s%s\n", resp.GetType(), resp.GetAddress(), resp.GetName(), tlsServerName)
}
