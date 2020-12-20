package mature

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/tenntenn/natureremo"
)

type mockCloudWatchAPI struct {
	cloudwatchiface.CloudWatchAPI
	putMetricDataWithContextFunc func(
		ctx aws.Context,
		input *cloudwatch.PutMetricDataInput,
		opts ...request.Option,
	) (*cloudwatch.PutMetricDataOutput, error)
}

func (m *mockCloudWatchAPI) PutMetricDataWithContext(
	ctx aws.Context,
	input *cloudwatch.PutMetricDataInput,
	opts ...request.Option,
) (*cloudwatch.PutMetricDataOutput, error) {
	return m.putMetricDataWithContextFunc(ctx, input, opts...)
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

func TestPutSensorValues(t *testing.T) {
	testcases := []struct {
		deviceID  string
		timestamp time.Time
		values    map[natureremo.SensorType]natureremo.SensorValue
	}{
		{
			deviceID:  "016e00cc-a76a-43be-bdfc-b305722fe4fd",
			timestamp: time.Date(2020, 12, 11, 10, 9, 8, 7, time.UTC),
			values:    map[natureremo.SensorType]natureremo.SensorValue{},
		},
	}

	for _, tc := range testcases {
		ctx := context.Background()
		cw := &cloudWatchImpl{
			client: &mockCloudWatchAPI{
				putMetricDataWithContextFunc: func(
					ctx aws.Context,
					input *cloudwatch.PutMetricDataInput,
					opts ...request.Option,
				) (*cloudwatch.PutMetricDataOutput, error) {
					return &cloudwatch.PutMetricDataOutput{}, nil
				},
			},
		}

		if err := cw.PutSensorValues(ctx, tc.timestamp, tc.deviceID, tc.values); err != nil {
			t.Errorf("want no error. got: %s", err)
		}
	}
}

func TestPutSensorValues_error(t *testing.T) {
	testcases := []struct {
		subtitle  string
		deviceID  string
		timestamp time.Time
		values    map[natureremo.SensorType]natureremo.SensorValue
		err       string
		want      string
	}{
		{
			subtitle:  "api error",
			deviceID:  "9347d9c2-cde0-43be-b931-40b229a6645d",
			timestamp: time.Date(2020, 12, 11, 10, 9, 8, 7, time.UTC),
			values:    map[natureremo.SensorType]natureremo.SensorValue{},
			err:       "unexpected error",
			want:      "cannot put metric data: unexpected error",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.subtitle, func(t *testing.T) {
			ctx := context.Background()
			cw := &cloudWatchImpl{
				client: &mockCloudWatchAPI{
					putMetricDataWithContextFunc: func(
						ctx aws.Context,
						input *cloudwatch.PutMetricDataInput,
						opts ...request.Option,
					) (*cloudwatch.PutMetricDataOutput, error) {
						return nil, fmt.Errorf(tc.err)
					},
				},
			}

			err := cw.PutSensorValues(ctx, tc.timestamp, tc.deviceID, tc.values)
			if err == nil {
				t.Fatal("want error. got nil")
			}
			if err.Error() != tc.want {
				t.Errorf("want: %q. got: %q", tc.want, err.Error())
			}
		})
	}
}
