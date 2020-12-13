package mature

import "testing"

func TestNewConfig(t *testing.T) {
	testcases := []struct {
		region, accessTokenKey, deviceIDKey string
	}{
		{
			region:         "dummy-region",
			accessTokenKey: "dummy-access-token-key",
			deviceIDKey:    "dummy-device-id-key",
		},
		{
			region:         "",
			accessTokenKey: "dummy-access-token-key",
			deviceIDKey:    "dummy-device-id-key",
		},
		{
			region:         "",
			accessTokenKey: "",
			deviceIDKey:    "dummy-device-id-key",
		},
		{
			region:         "",
			accessTokenKey: "",
			deviceIDKey:    "",
		},
	}

	for _, tc := range testcases {
		got := NewConfig(tc.region, tc.accessTokenKey, tc.deviceIDKey)
		if got == nil {
			t.Error("got nil")
		}
	}
}
