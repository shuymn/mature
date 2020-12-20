package mature

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"golang.org/x/xerrors"
)

type SSM interface {
	GetParameters(ctx context.Context, keys ...string) (map[string]string, error)
}

type ssmImpl struct {
	client ssmiface.SSMAPI
}

func NewSSM(client ssmiface.SSMAPI) *ssmImpl {
	return &ssmImpl{client: client}
}

func (s *ssmImpl) GetParameters(ctx context.Context, keys ...string) (map[string]string, error) {
	names := make([]*string, 0, len(keys))
	for _, k := range keys {
		if k == "" {
			continue
		}
		names = append(names, aws.String(k))
	}

	if len(names) == 0 {
		return nil, xerrors.New("invalid keys given")
	}

	input := new(ssm.GetParametersInput).SetNames(names).SetWithDecryption(true)
	output, err := s.client.GetParametersWithContext(ctx, input)
	if err != nil {
		return nil, xerrors.Errorf("cannot get parameters: %w", err)
	}

	values := make(map[string]string, len(keys))
	for _, p := range output.Parameters {
		values[aws.StringValue(p.Name)] = aws.StringValue(p.Value)
	}

	return values, nil
}
