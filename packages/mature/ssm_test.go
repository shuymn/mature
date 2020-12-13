package mature

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type mockSSMAPI struct {
	ssmiface.SSMAPI
	getParametersWithContextFunc func(
		ctx aws.Context,
		input *ssm.GetParametersInput,
		opts ...request.Option,
	) (*ssm.GetParametersOutput, error)
}

func (m *mockSSMAPI) GetParametersWithContext(
	ctx aws.Context,
	input *ssm.GetParametersInput,
	opts ...request.Option,
) (*ssm.GetParametersOutput, error) {
	return m.getParametersWithContextFunc(ctx, input, opts...)
}

func TestNewSSM(t *testing.T) {
	client := &mockSSMAPI{}
	got := NewSSM(client)

	if got == nil {
		t.Error("got nil")
	} else if got.client != client {
		t.Error("ssm api client does not match")
	}
}

func TestGetParameters(t *testing.T) {
	testcases := []struct {
		accessTokenKey, deviceIDKey, unknownKey string
	}{
		{
			accessTokenKey: "/mature/test/access-token-key",
			deviceIDKey:    "/mature/test/device-id-key",
			unknownKey:     "/mature/test/unknown-key",
		},
		{
			deviceIDKey: "/mature/test/device-id-key",
			unknownKey:  "/mature/test/unknown-key",
		},
		{
			unknownKey: "/mature/test/unknown-key",
		},
	}

	for _, tc := range testcases {
		ctx := context.Background()
		ssm := &ssmImpl{
			client: &mockSSMAPI{
				getParametersWithContextFunc: func(
					ctx aws.Context,
					input *ssm.GetParametersInput,
					opts ...request.Option,
				) (*ssm.GetParametersOutput, error) {
					return &ssm.GetParametersOutput{}, nil
				},
			},
		}

		_, err := ssm.GetParameters(ctx, tc.accessTokenKey, tc.deviceIDKey, tc.unknownKey)
		if err != nil {
			t.Errorf("want no error. got: %s", err)
		}
	}
}

func TestGetParameters_error(t *testing.T) {
	testcases := []struct {
		subtitle       string
		err            string
		accessTokenKey string
		deviceIDKey    string
		unknownKey     string
		want           string
	}{
		{
			subtitle:       "api error",
			accessTokenKey: "/mature/test/access-token-key",
			deviceIDKey:    "/mature/test/device-id-key",
			unknownKey:     "/mature/test/unknown-key",
			err:            "unexpected error",
			want:           "cannot get parameters: unexpected error",
		},
		{
			subtitle: "invalid keys error",
			want:     "invalid keys given",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.subtitle, func(t *testing.T) {
			ctx := context.Background()
			ssm := &ssmImpl{
				client: &mockSSMAPI{
					getParametersWithContextFunc: func(
						ctx aws.Context,
						input *ssm.GetParametersInput,
						opts ...request.Option,
					) (*ssm.GetParametersOutput, error) {
						return nil, fmt.Errorf(tc.err)
					},
				},
			}

			_, err := ssm.GetParameters(ctx, tc.accessTokenKey, tc.deviceIDKey, tc.unknownKey)
			if err == nil {
				t.Fatal("want error. got nil")
			}
			if err.Error() != tc.want {
				t.Errorf("want: %q. got: %q", tc.want, err.Error())
			}
		})
	}
}
