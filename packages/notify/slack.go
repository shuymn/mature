package notify

import (
	"context"
	"io"

	"github.com/slack-go/slack"
	"golang.org/x/xerrors"
)

type Message struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type string `json:"type"`
	Text Text   `json:"text"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Slack interface {
	UploadImage(ctx context.Context, channelID string, image io.Reader) error
}

type slackImpl struct {
	client *slack.Client
}

func NewSlack(client *slack.Client) *slackImpl {
	return &slackImpl{client: client}
}

func (s *slackImpl) UploadImage(ctx context.Context, channelID string, image io.Reader) error {
	params := slack.FileUploadParameters{
		Reader:   image,
		Filename: "report",
		Channels: []string{channelID},
	}
	_, err := s.client.UploadFileContext(ctx, params)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	return nil
}
