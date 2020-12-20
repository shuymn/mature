package notify

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

type mockCloudWatchAPI struct {
	cloudwatchiface.CloudWatchAPI
	getMetricWidgetImageWithContextFunc func(
		ctx aws.Context,
		input *cloudwatch.GetMetricWidgetImageInput,
		opts ...request.Option,
	) (*cloudwatch.GetMetricWidgetImageOutput, error)
}

func (m *mockCloudWatchAPI) GetMetricWidgetImageWithContext(
	ctx aws.Context,
	input *cloudwatch.GetMetricWidgetImageInput,
	opts ...request.Option,
) (*cloudwatch.GetMetricWidgetImageOutput, error) {
	return m.getMetricWidgetImageWithContextFunc(ctx, input, opts...)
}

func TestNewCloudWatch(t *testing.T) {
	client := &mockCloudWatchAPI{}
	got := NewCloudWatch(client)

	if got == nil {
		t.Error("got nil")
	} else if got.client != client {
		t.Error("cloudwatch api client does not match")
	}
}

func TestGetMetricWidgetImage(t *testing.T) {
	testcases := []struct {
		widget string
	}{
		{widget: ""},
		{widget: "dummy-widget"},
	}

	for _, tc := range testcases {
		ctx := context.Background()
		cw := &cloudWatchImpl{
			client: &mockCloudWatchAPI{
				getMetricWidgetImageWithContextFunc: func(
					ctx aws.Context,
					input *cloudwatch.GetMetricWidgetImageInput,
					opts ...request.Option,
				) (*cloudwatch.GetMetricWidgetImageOutput, error) {
					return &cloudwatch.GetMetricWidgetImageOutput{}, nil
				},
			},
		}

		if _, err := cw.GetMetricWidgetImage(ctx, tc.widget); err != nil {
			t.Errorf("want no error. got: %s", err)
		}
	}
}

func TestGetMetricWidgetImage_error(t *testing.T) {
	testcases := []struct {
		subtitle string
		widget   string
		err      string
		want     string
	}{
		{
			subtitle: "api error",
			widget:   "",
			err:      "unexpected error",
			want:     "failed to get metric widget image: unexpected error",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.subtitle, func(t *testing.T) {
			ctx := context.Background()
			cw := &cloudWatchImpl{
				client: &mockCloudWatchAPI{
					getMetricWidgetImageWithContextFunc: func(
						ctx aws.Context,
						input *cloudwatch.GetMetricWidgetImageInput,
						opts ...request.Option,
					) (*cloudwatch.GetMetricWidgetImageOutput, error) {
						return nil, fmt.Errorf(tc.err)
					},
				},
			}

			_, err := cw.GetMetricWidgetImage(ctx, tc.widget)
			if err == nil {
				t.Fatal("want error. got nil")
			}
			if err.Error() != tc.want {
				t.Errorf("want: %q. got: %q", tc.want, err.Error())
			}
		})
	}
}
