// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

package observe

import (
	"testing"
	"time"

	observerpb "github.com/khulnasoft/shipyard/api/v1/observer"
	monitorAPI "github.com/khulnasoft/shipyard/pkg/monitor/api"
	"github.com/khulnasoft/triangle/pkg/defaults"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAgentEventSubTypeMap(t *testing.T) {
	// Make sure to keep agent event sub-types maps in sync. See
	// agentEventSubtypes godoc for details.
	require.Len(t, agentEventSubtypes, len(monitorAPI.AgentNotifications))
	for _, v := range agentEventSubtypes {
		require.Contains(t, monitorAPI.AgentNotifications, v)
	}
	agentEventSubtypesContainsValue := func(an monitorAPI.AgentNotification) bool {
		for _, v := range agentEventSubtypes {
			if v == an {
				return true
			}
		}
		return false
	}
	for k := range monitorAPI.AgentNotifications {
		require.True(t, agentEventSubtypesContainsValue(k))
	}
}

func Test_getAgentEventsRequest(t *testing.T) {
	selectorOpts.since = ""
	selectorOpts.until = ""
	req, err := getAgentEventsRequest()
	assert.NoError(t, err)
	assert.Equal(t, &observerpb.GetAgentEventsRequest{Number: defaults.EventsPrintCount}, req)
	selectorOpts.since = "2021-04-26T00:00:00Z"
	selectorOpts.until = "2021-04-26T00:01:00Z"
	req, err = getAgentEventsRequest()
	assert.NoError(t, err)
	since, err := time.Parse(time.RFC3339, selectorOpts.since)
	assert.NoError(t, err)
	until, err := time.Parse(time.RFC3339, selectorOpts.until)
	assert.NoError(t, err)
	assert.Equal(t, &observerpb.GetAgentEventsRequest{
		Number: defaults.EventsPrintCount,
		Since:  timestamppb.New(since),
		Until:  timestamppb.New(until),
	}, req)
}

func Test_getAgentEventsRequestWithoutSince(t *testing.T) {
	selectorOpts.since = ""
	selectorOpts.until = ""
	req, err := getAgentEventsRequest()
	assert.NoError(t, err)
	assert.Equal(t, &observerpb.GetAgentEventsRequest{Number: defaults.EventsPrintCount}, req)
	selectorOpts.until = "2021-04-26T00:01:00Z"
	req, err = getAgentEventsRequest()
	assert.NoError(t, err)
	assert.NoError(t, err)
	until, err := time.Parse(time.RFC3339, selectorOpts.until)
	assert.NoError(t, err)
	assert.Equal(t, &observerpb.GetAgentEventsRequest{
		Number: defaults.EventsPrintCount,
		Until:  timestamppb.New(until),
	}, req)
}
