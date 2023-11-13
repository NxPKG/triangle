// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

package observe

import (
	"testing"
	"time"

	observerpb "github.com/khulnasoft/shipyard/api/v1/observer"
	"github.com/khulnasoft/triangle/pkg/defaults"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_getDebugEventsRequest(t *testing.T) {
	selectorOpts.since = ""
	selectorOpts.until = ""
	req, err := getDebugEventsRequest()
	assert.NoError(t, err)
	assert.Equal(t, &observerpb.GetDebugEventsRequest{Number: defaults.EventsPrintCount}, req)
	selectorOpts.since = "2021-04-26T01:00:00Z"
	selectorOpts.until = "2021-04-26T01:01:00Z"
	req, err = getDebugEventsRequest()
	assert.NoError(t, err)
	since, err := time.Parse(time.RFC3339, selectorOpts.since)
	assert.NoError(t, err)
	until, err := time.Parse(time.RFC3339, selectorOpts.until)
	assert.NoError(t, err)
	assert.Equal(t, &observerpb.GetDebugEventsRequest{
		Number: defaults.EventsPrintCount,
		Since:  timestamppb.New(since),
		Until:  timestamppb.New(until),
	}, req)
}

func Test_getDebugEventsRequestWithoutSince(t *testing.T) {
	selectorOpts.since = ""
	selectorOpts.until = ""
	req, err := getDebugEventsRequest()
	assert.NoError(t, err)
	assert.Equal(t, &observerpb.GetDebugEventsRequest{Number: defaults.EventsPrintCount}, req)
	selectorOpts.until = "2021-04-26T01:01:00Z"
	req, err = getDebugEventsRequest()
	assert.NoError(t, err)
	until, err := time.Parse(time.RFC3339, selectorOpts.until)
	assert.NoError(t, err)
	assert.Equal(t, &observerpb.GetDebugEventsRequest{
		Number: defaults.EventsPrintCount,
		Until:  timestamppb.New(until),
	}, req)
}
