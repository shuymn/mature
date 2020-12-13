package mature

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/tenntenn/natureremo"
	"golang.org/x/xerrors"
)

const (
	metricsNamespace  = "NatureRemo/RoomMetrics"
	dimensionDeviceID = "DeviceID"
)

type CloudWatch interface {
	PutSensorValues(
		ctx context.Context,
		timestamp time.Time,
		deviceID string,
		vs map[natureremo.SensorType]natureremo.SensorValue,
	) error
}

type cloudWatchImpl struct {
	client cloudwatchiface.CloudWatchAPI
}

func NewCloudWatch(client cloudwatchiface.CloudWatchAPI) *cloudWatchImpl {
	return &cloudWatchImpl{client: client}
}

func (cw *cloudWatchImpl) PutSensorValues(
	ctx context.Context,
	timestamp time.Time,
	deviceID string,
	vs map[natureremo.SensorType]natureremo.SensorValue,
) error {
	dims := []*cloudwatch.Dimension{
		new(cloudwatch.Dimension).
			SetName(dimensionDeviceID).
			SetValue(deviceID),
	}
	data := make([]*cloudwatch.MetricDatum, 0, len(vs))
	for t, v := range vs {
		name := strings.Title(SensorTypeString[t])
		datum := new(cloudwatch.MetricDatum).
			SetMetricName(name).
			SetTimestamp(timestamp).
			SetValue(v.Value).
			SetDimensions(dims)
		data = append(data, datum)
	}
	input := new(cloudwatch.PutMetricDataInput).
		SetNamespace(metricsNamespace).
		SetMetricData(data)
	_, err := cw.client.PutMetricDataWithContext(ctx, input)
	if err != nil {
		return xerrors.Errorf("cannot put metric data: %w", err)
	}
	return nil
}
