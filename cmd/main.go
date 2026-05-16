package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/dinarasaurae/path-of-exile-cli/internal/dungeon"
)

func run(args []string, stdout io.Writer, stderr io.Writer) (err error) {
	flags := flag.NewFlagSet("path-of-exile-cli", flag.ContinueOnError)
	flags.SetOutput(stderr)

	configPath := flags.String("config", "config.json", "path to config file")

	if err := flags.Parse(args); err != nil {
		return err
	}
	configFile, err := os.Open(*configPath)
	if err != nil {
		return fmt.Errorf("open config: %w", err)
	}
	defer func() {
		if closeErr := configFile.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close config: %w", closeErr)
		}
	}()

	settings, err := dungeon.LoadSettings(configFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	return nil
}
