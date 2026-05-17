package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/dinarasaurae/path-of-exile-cli/internal/dungeon"
)

func main() {
	os.Exit(execute(os.Args[1:], os.Stdout, os.Stderr))
}

func execute(args []string, stdout io.Writer, stderr io.Writer) int {
	if err := run(args, stdout, stderr); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}

		_, _ = fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func run(args []string, stdout io.Writer, stderr io.Writer) (err error) {
	flags := flag.NewFlagSet("path-of-exile-cli", flag.ContinueOnError)
	flags.SetOutput(stderr)

	configPath := flags.String("config", "config.json", "path to config file")
	eventsPath := flags.String("events", "events", "path to events file")

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

	eventsFile, err := os.Open(*eventsPath)
	if err != nil {
		return fmt.Errorf("open events: %w", err)
	}
	defer func() {
		if closeErr := eventsFile.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close events: %w", closeErr)
		}
	}()

	events, err := dungeon.ParseEvents(eventsFile)
	if err != nil {
		return fmt.Errorf("parse events: %w", err)
	}

	result, err := dungeon.Process(settings, events)
	if err != nil {
		return fmt.Errorf("process events: %w", err)
	}

	_, err = io.WriteString(stdout, dungeon.FormatResult(result))
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
