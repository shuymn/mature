package mature

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ssm"
	"golang.org/x/xerrors"
)

const (
	ExitCodeOK = iota
	ExitCodeErr
)

type App struct {
	CloudWatch CloudWatch
	Config     *Config
	DeviceID   string
	NatureRemo NatureRemo
}

func New(conf *Config) *App {
	return &App{Config: conf}
}

func (app *App) initLazy(ctx context.Context) error {
	if app.CloudWatch != nil && app.NatureRemo != nil && app.DeviceID != "" {
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

	if app.NatureRemo != nil && app.DeviceID != "" {
		return nil
	}

	vs, err := NewSSM(ssm.New(sess)).GetParameters(ctx, appConf.AccessTokenKey, appConf.DeviceIDKey)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	if app.NatureRemo == nil {
		app.NatureRemo = NewNatureRemo(vs[appConf.AccessTokenKey])
	}

	if app.DeviceID == "" {
		app.DeviceID = vs[appConf.DeviceIDKey]
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
		return xerrors.Errorf("%w", err)
	}

	vs, err := app.NatureRemo.FetchAllNewestSensorValue(ctx, app.DeviceID)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	if err = app.CloudWatch.PutSensorValues(ctx, time.Now(), app.DeviceID, vs); err != nil {
		return xerrors.Errorf("%w", err)
	}

	return nil
}
