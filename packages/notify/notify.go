package notify

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/rakyll/statik/fs"
	_ "github.com/shuymn/mature/packages/notify/statik"
	"github.com/slack-go/slack"
	"golang.org/x/xerrors"
)

const (
	ExitCodeOK = iota
	ExitCodeErr
)

type App struct {
	CloudWatch CloudWatch
	ChannelID  string
	Config     *Config
	Slack      Slack
}

//go:generate statik -src=widget -include=*.json

func New(conf *Config) *App {
	return &App{Config: conf}
}

func (app *App) initLazy(ctx context.Context) error {
	if app.CloudWatch != nil && app.Slack != nil && app.ChannelID != "" {
		return nil
	}

	appConf := app.Config
	awsConf := new(aws.Config)
	if appConf.Region != "" {
		awsConf.Region = aws.String(appConf.Region)
	}
	sessOpt := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *awsConf,
	}
	sess := session.Must(session.NewSessionWithOptions(sessOpt))

	if app.CloudWatch == nil {
		app.CloudWatch = NewCloudWatch(cloudwatch.New(sess))
	}

	if app.Slack != nil && app.ChannelID != "" {
		return nil
	}

	vs, err := NewSSM(ssm.New(sess)).GetParameters(ctx, appConf.SlackTokenKey, appConf.ChannelIDKey)
	if err != nil {
		return xerrors.Errorf("failed to get parameters: %w", err)
	}

	if app.Slack == nil {
		app.Slack = NewSlack(slack.New(vs[appConf.SlackTokenKey]))
	}

	if app.ChannelID == "" {
		app.ChannelID = vs[appConf.ChannelIDKey]
	}

	return nil
}

func (app *App) Handle(ctx context.Context) error {
	return app.RunWithContext(ctx)
}

func (app *App) Run() error {
	return app.RunWithContext(context.Background())
}

func (app *App) RunWithContext(ctx context.Context) error {
	if err := app.initLazy(ctx); err != nil {
		return xerrors.Errorf("failed to init lazy: %w", err)
	}

	statikFS, err := fs.New()
	if err != nil {
		return xerrors.Errorf("failed to get statik fs: %w", err)
	}

	r, err := statikFS.Open("/natureremo.json")
	if err != nil {
		return xerrors.Errorf("failed to open natureremo.json: %w", err)
	}
	defer r.Close()

	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return xerrors.Errorf("failed to read natureremo.json: %w", err)
	}

	b, err := app.CloudWatch.GetMetricWidgetImage(ctx, string(contents))
	if err != nil {
		return xerrors.Errorf("failed to get metric widget image: %w", err)
	}

	if err = app.Slack.UploadImage(ctx, app.ChannelID, bytes.NewReader(b)); err != nil {
		return xerrors.Errorf("failed to upload image: %w", err)
	}

	return nil
}
