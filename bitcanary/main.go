package bitcanary

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/n25a/BitCanary/internal/log"
	"go.uber.org/zap"

	"github.com/n25a/BitCanary/internal/config"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "The path of config file")

	// Parse the terminal flags
	flag.Parse()
}

func main() {
	// load the config
	config.LoadConfig(configPath)

	// validate the config
	if config.C.PrimaryAddress == "" {
		log.Logger.Fatal("primary address can not be empty")
	}

	if config.C.CanaryAddress == "" {
		log.Logger.Fatal("canary address can not be empty")
	}

	// start the server
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)

	// TODO: add metrics handler

	server := &http.Server{
		Addr:         config.C.HTTP.Bind,
		Handler:      mux,
		ReadTimeout:  config.C.HTTP.ReadTimeout,
		WriteTimeout: config.C.HTTP.WriteTimeout,
	}

	go func() {
		log.Logger.Info("Starting HTTP server",
			zap.String("address", config.C.HTTP.Bind),
			zap.String("primary address", config.C.PrimaryAddress),
			zap.String("canary address", config.C.CanaryAddress),
		)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Logger.Fatal("HTTP server ListenAndServe Error", zap.Error(err))
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	log.Logger.Debug("Closing HTTP connections")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Logger.Error("Error in shutting down the HTTP server", zap.Error(err))
	}

	log.Logger.Info("HTTP server is shut down")
}
