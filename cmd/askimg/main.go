package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/igolaizola/askimg"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/peterbourgon/ff/v3/ffyaml"
)

func main() {
	// Create signal based context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Launch command
	cmd := newCommand()
	if err := cmd.ParseAndRun(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func newCommand() *ffcli.Command {
	fs := flag.NewFlagSet("askimg", flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	var cfg askimg.Config
	fs.StringVar(&cfg.Token, "token", "", "authentication token")
	fs.StringVar(&cfg.Image, "image", "", "url of the image")
	fs.StringVar(&cfg.Question, "question", "", "question to ask, if empty it will be captioned")
	fs.IntVar(&cfg.Temperature, "temperature", 1, "temperature to use with nucleus sampling")
	fs.BoolVar(&cfg.UseNucleusSampling, "nucleus", false, "use nucleus sampling")
	fs.DurationVar(&cfg.Timeout, "timeout", 30*time.Second, "timeout of the request")

	return &ffcli.Command{
		ShortUsage: fmt.Sprintf("askimg [flags] <key> <value data...>"),
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ffyaml.Parser),
			ff.WithEnvVarPrefix("ASKIMG"),
		},
		ShortHelp: fmt.Sprintf("askimg asks a question about an image."),
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			output, err := askimg.Ask(ctx, &cfg)
			if err != nil {
				return err
			}
			fmt.Println(output)
			return nil
		},
	}
}
