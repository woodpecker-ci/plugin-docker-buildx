package main

import (
	"os"

	"codeberg.org/woodpecker-plugins/plugin-docker-buildx/plugin"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"github.com/drone-plugins/drone-plugin-lib/errors"
	"github.com/drone-plugins/drone-plugin-lib/urfave"
)

var version = "unknown"

func main() {
	settings := &plugin.Settings{}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		godotenv.Overload("/run/drone/env")
	}
	
	if envFile, set := os.LookupEnv("PLUGIN_ENV_FILE"); set {
		godotenv.Overload(envFile)
	}

	app := &cli.App{
		Name:    "docker-buildx",
		Usage:   "build docker container with DinD and buildx",
		Version: version,
		Flags:   append(settingsFlags(settings), urfave.Flags()...),
		Action:  run(settings),
	}

	if err := app.Run(os.Args); err != nil {
		errors.HandleExit(err)
	}
}

func run(settings *plugin.Settings) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		urfave.LoggingFromContext(ctx)

		plugin := plugin.New(
			*settings,
			urfave.PipelineFromContext(ctx),
			urfave.NetworkFromContext(ctx),
		)

		if err := plugin.Validate(); err != nil {
			if e, ok := err.(errors.ExitCoder); ok {
				return e
			}

			return errors.ExitMessagef("validation failed: %w", err)
		}

		if err := plugin.Execute(); err != nil {
			if e, ok := err.(errors.ExitCoder); ok {
				return e
			}

			return errors.ExitMessagef("execution failed: %w", err)
		}

		return nil
	}
}
