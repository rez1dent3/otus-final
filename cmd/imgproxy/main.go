package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rez1dent3/otus-final/internal/imgprev"
	"github.com/rez1dent3/otus-final/internal/server"
	"gopkg.in/yaml.v2"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/imgproxy/config.yaml", "Path to configuration file")
}

func newConfig(reader io.Reader) (*imgprev.Config, error) {
	config := imgprev.Config{}
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	flag.Parse()

	file, err := os.Open(configFile)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	config, err := newConfig(file)
	if err != nil {
		log.Println(err)
		return
	}

	app := imgprev.New(config)
	httpServ := server.New(app)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	defer func() {
		app.Logger().Info("http server stopping")
		if err := httpServ.Stop(ctx); err != nil {
			app.Logger().Error(err.Error())
		}
	}()

	app.Logger().Info("http server starting")
	if err := httpServ.Start(ctx); err != nil {
		app.Logger().Error(err.Error())
	}
}
