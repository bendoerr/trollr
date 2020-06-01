package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bendoerr/trollr/exec"
	"github.com/heetch/confita"
	confitaenv "github.com/heetch/confita/backend/env"
	confitaflags "github.com/heetch/confita/backend/flags"
	"go.uber.org/zap"
)

// These values are populated by 'govvv' see https://github.com/troian/govvv
// nolint:deadcode
var Version, GitSummary, BuildDate string
var BuildInfo string

func init() {
	if len(Version) < 1 {
		Version = "dev"
	}

	if len(GitSummary) < 1 {
		GitSummary = "snapshot"
	}

	BuildInfo = fmt.Sprintf("Src Version: %s\n   Build Date: %s\n", GitSummary, BuildDate)
}

func main() {
	// Setup the logger
	var zapLogger zap.Logger
	defer func() {
		_ = zapLogger.Sync()
	}()

	for {
		cfg := AppConfig{
			Listen: ":6789",
		}

		// Load the configuration
		loader := confita.NewLoader(
			confitaenv.NewBackend(),
			confitaflags.NewBackend())

		err := loader.Load(context.Background(), &cfg)
		if err != nil {
			fmt.Printf("configruation error: %s\n", err)
			os.Exit(1)
		}

		// Print out a fancy logo!
		fmt.Printf(` 
     _____         _ _          .-------.    ______
    |_   _|       | | |        /   o   /|   /\     \
      | |_ __ ___ | | |_ __   /_______/o|  /o \  o  \
      | | '__/ _ \| | | '__|  | o     | | /   o\_____\
      | | | | (_) | | | |     |   o   |o/ \o   /o    /
      | |_|  \___/|_|_|_|     |     o |/   \ o/  o  /
      \_/ %-18s  '-------'     \/____o/
`+"\n", Version)

		fmt.Printf("  %s\n", BuildInfo)

		// Setup logging
		zapConfig := zap.NewProductionConfig()
		if len(cfg.LogFile) > 0 {
			zapConfig.OutputPaths = []string{cfg.LogFile}
		}
		zapLogger, _ := zapConfig.Build()
		logger := zapLogger.Named("main")

		// Setup Services
		tx := exec.NewTimingExecutor(exec.Run)
		lx := exec.NewLoggingExecutor(tx.Run, logger)
		px := exec.NewPoolExecutor(lx.Run)
		troll := NewTroll(cfg.TrollBin, px.Run)
		http := NewAPI(cfg.Listen, troll, logger.Named("http"))

		// Start the HTTP Server
		http.Start()

		// Confirm configuration values
		fmt.Printf("  Running with Configuration:\n")
		fmt.Printf("    Listen:    '%s'\n", cfg.Listen)
		fmt.Printf("    Log File:  '%s'\n", cfg.LogFile)
		fmt.Printf("    Mosmllib:  '%s'\n", cfg.Mosmllib)
		fmt.Printf("    Troll Bin: '%s'\n", cfg.TrollBin)
		fmt.Printf("\n\n")

		// Signal and exit handling from here down
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		time.Sleep(time.Second)
		logger.Info("started")

		sig := <-sigchan
		num := int(sig.(syscall.Signal))

		fmt.Printf("signal: %s\n", sig)

		if sig == os.Kill {
			os.Exit(128 + num)
		}

		// Clean up
		if sig != os.Kill {
			err = http.Stop()
			if err != nil {
				fmt.Printf("error shutting down serevr, %s\n", err)
			}
		}

		switch sig {
		case syscall.SIGTERM:
			// Exit Normally.
			os.Exit(0)
		case syscall.SIGHUP:
			// Reload
			continue
		default:
			// Exit with error
			os.Exit(128 + num)
		}
	}
}
