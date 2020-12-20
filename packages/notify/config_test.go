package notify

import "testing"

func TestNewConfig(t *testing.T) {
	testcases := []struct {
		region, slackTokenKey, channelIDKey string
	}{
		{
			region:        "dummy-region",
			slackTokenKey: "dummy-slack-token-key",
			channelIDKey:  "dummy-channel-id-key",
		},
		{
			region:        "",
			slackTokenKey: "dummy-slack-token-key",
			channelIDKey:  "dummy-channel-id-key",
		},
		{
			region:        "",
			slackTokenKey: "",
			channelIDKey:  "dummy-channel-id-key",
		},
		{
			region:        "",
			slackTokenKey: "",
			channelIDKey:  "",
		},
	}

	for _, tc := range testcases {
		got := NewConfig(tc.region, tc.slackTokenKey, tc.channelIDKey)
		if got == nil {
			t.Error("got nil")
		}
	}
}
