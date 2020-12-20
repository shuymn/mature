package notify

import (
	"context"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"golang.org/x/xerrors"
)

type CloudWatch interface {
	GetMetricWidgetImage(ctx context.Context, widget string) ([]byte, error)
}

type cloudWatchImpl struct {
	client cloudwatchiface.CloudWatchAPI
}

func NewCloudWatch(client cloudwatchiface.CloudWatchAPI) *cloudWatchImpl {
	return &cloudWatchImpl{client: client}
}

func (c *cloudWatchImpl) GetMetricWidgetImage(ctx context.Context, widget string) ([]byte, error) {
	input := new(cloudwatch.GetMetricWidgetImageInput).SetMetricWidget(widget)
	output, err := c.client.GetMetricWidgetImageWithContext(ctx, input)
	if err != nil {
		return nil, xerrors.Errorf("failed to get metric widget image: %w", err)
	}
	return output.MetricWidgetImage, nil
}
