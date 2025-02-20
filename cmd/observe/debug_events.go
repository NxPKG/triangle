// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

package observe

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	observerpb "github.com/khulnasoft/shipyard/api/v1/observer"
	"github.com/khulnasoft/triangle/cmd/common/config"
	"github.com/khulnasoft/triangle/cmd/common/conn"
	"github.com/khulnasoft/triangle/cmd/common/template"
	"github.com/khulnasoft/triangle/pkg/defaults"
	"github.com/khulnasoft/triangle/pkg/logger"
	hubtime "github.com/khulnasoft/triangle/pkg/time"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newDebugEventsCommand(vp *viper.Viper) *cobra.Command {
	debugEventsCmd := &cobra.Command{
		Use:   "debug-events",
		Short: "Observe Shipyard debug events",
		RunE: func(cmd *cobra.Command, _ []string) error {
			debug := vp.GetBool(config.KeyDebug)
			if err := handleEventsArgs(cmd.OutOrStdout(), debug); err != nil {
				return err
			}
			req, err := getDebugEventsRequest()
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
			defer cancel()

			triangleConn, err := conn.New(ctx, vp.GetString(config.KeyServer), vp.GetDuration(config.KeyTimeout))
			if err != nil {
				return err
			}
			defer triangleConn.Close()
			client := observerpb.NewObserverClient(triangleConn)
			logger.Logger.Debug("Sending GetDebugEvents request", "request", req)
			if err := getDebugEvents(ctx, client, req); err != nil {
				msg := err.Error()
				// extract custom error message from failed grpc call
				if s, ok := status.FromError(err); ok && s.Code() == codes.Unknown {
					msg = s.Message()
				}
				return errors.New(msg)
			}
			return nil
		},
	}

	flagSets := []*pflag.FlagSet{selectorFlags, formattingFlags, config.ServerFlags, otherFlags}
	for _, fs := range flagSets {
		debugEventsCmd.Flags().AddFlagSet(fs)
	}
	template.RegisterFlagSets(debugEventsCmd, flagSets...)
	return debugEventsCmd
}

func getDebugEventsRequest() (*observerpb.GetDebugEventsRequest, error) {
	// convert selectorOpts.since into a param for GetDebugEvents
	var since, until *timestamppb.Timestamp
	if selectorOpts.since != "" {
		st, err := hubtime.FromString(selectorOpts.since)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the since time: %v", err)
		}

		since = timestamppb.New(st)
		if err := since.CheckValid(); err != nil {
			return nil, fmt.Errorf("failed to convert `since` timestamp to proto: %v", err)
		}
	}
	// Set the until field if --until option is specified and --follow
	// is not specified. If --since is specified but --until is not, the server sets the
	// --until option to the current timestamp.
	if selectorOpts.until != "" && !selectorOpts.follow {
		ut, err := hubtime.FromString(selectorOpts.until)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the until time: %v", err)
		}
		until = timestamppb.New(ut)
		if err := until.CheckValid(); err != nil {
			return nil, fmt.Errorf("failed to convert `until` timestamp to proto: %v", err)
		}
	}

	if since == nil && until == nil {
		switch {
		case selectorOpts.all:
			// all is an alias for last=uint64_max
			selectorOpts.last = ^uint64(0)
		case selectorOpts.last == 0:
			// no specific parameters were provided, just a vanilla `triangle events debug`
			selectorOpts.last = defaults.EventsPrintCount
		}
	}

	return &observerpb.GetDebugEventsRequest{
		Number: selectorOpts.last,
		Follow: selectorOpts.follow,
		Since:  since,
		Until:  until,
	}, nil
}

func getDebugEvents(ctx context.Context, client observerpb.ObserverClient, req *observerpb.GetDebugEventsRequest) error {
	b, err := client.GetDebugEvents(ctx, req)
	if err != nil {
		return err
	}

	defer printer.Close()

	for {
		resp, err := b.Recv()
		switch err {
		case io.EOF, context.Canceled:
			return nil
		case nil:
		default:
			if status.Code(err) == codes.Canceled {
				return nil
			}
			return err
		}

		if err = printer.WriteProtoDebugEvent(resp); err != nil {
			return err
		}
	}
}
