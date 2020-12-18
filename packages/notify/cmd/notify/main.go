package main

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hashicorp/logutils"
	"github.com/shuymn/mature/packages/notify"
	"gopkg.in/alecthomas/kingpin.v2"
)

var app *notify.App

func main() {
	debug := kingpin.Flag("debug", "enable debug log").Bool()

	opt := &notify.Option{
		Region: kingpin.Flag("region", "AWS region").
			OverrideDefaultFromEnvar("AWS_REGION").String(),
		SlackTokenKey: kingpin.Flag("slack-token-key", "AWS SSM key of Slack token").
			OverrideDefaultFromEnvar("MATURE_SLACK_TOKEN_KEY").Required().String(),
		ChannelIDKey: kingpin.Flag("channel-id-key", "AWS SSM key of Slack channel ID").
			OverrideDefaultFromEnvar("MATURE_SLACK_CHANNEL_ID_KEY").Required().String(),
	}

	kingpin.Parse()

	logLevel := "INFO"
	if *debug {
		logLevel = "DEBUG"
	}
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(logLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	conf := notify.NewConfig(*opt.Region, *opt.SlackTokenKey, *opt.ChannelIDKey)
	app = notify.New(conf)

	if strings.HasPrefix(os.Getenv("AWS_EXECUTION_ENV"), "AWS_Lambda") || os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		lambda.Start(app.Handle)
	} else {
		if err := app.Run(); err != nil {
			log.Printf("[ERROR] %+v\n", err)
			os.Exit(notify.ExitCodeErr)
		}
		os.Exit(notify.ExitCodeOK)
	}
}
