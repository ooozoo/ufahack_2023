package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"ufahack_2023/internal/config"
)

type Args struct {
	ConfigPath string
}

func processArgs() Args {
	var a Args

	f := flag.NewFlagSet("Server", flag.ExitOnError)
	f.StringVar(&a.ConfigPath, "c", "config.yml", "Path to configuration file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := config.Description()
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a
}

func run() error {
	args := processArgs()

	cfg, err := config.Load(args.ConfigPath)
	if err != nil {
		return err
	}

	_ = cfg

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
