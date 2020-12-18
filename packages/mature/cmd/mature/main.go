package main

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hashicorp/logutils"
	"github.com/mature/packages/mature"
	"gopkg.in/alecthomas/kingpin.v2"
)

var app *mature.App

func main() {
	debug := kingpin.Flag("debug", "enable debug log").Bool()

	opt := &mature.Option{
		Region: kingpin.Flag("region", "AWS region").
			OverrideDefaultFromEnvar("AWS_REGION").String(),
		AccessTokenKey: kingpin.Flag("access-token-key", "AWS SSM key of Nature Remo access token").
			OverrideDefaultFromEnvar("MATURE_ACCESS_TOKEN_KEY").Required().String(),
		DeviceIDKey: kingpin.Flag("device-id-key", "AWS SSM key of Nature Remo device id").
			OverrideDefaultFromEnvar("MATURE_DEVICE_ID_KEY").Required().String(),
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

	conf := mature.NewConfig(*opt.Region, *opt.AccessTokenKey, *opt.DeviceIDKey)
	app = mature.New(conf)

	if strings.HasPrefix(os.Getenv("AWS_EXECUTION_ENV"), "AWS_Lambda") || os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		lambda.Start(app.Handle)
	} else {
		if err := app.Run(); err != nil {
			log.Printf("[ERROR] %+v\n", err)
			os.Exit(mature.ExitCodeErr)
		}
		os.Exit(mature.ExitCodeOK)
	}
}
