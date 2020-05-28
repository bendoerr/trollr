package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bendoerr/trollr"
	"github.com/bendoerr/trollr/util"
	"github.com/heetch/confita"
	confitaenv "github.com/heetch/confita/backend/env"
	confitaflags "github.com/heetch/confita/backend/flags"
)

// Version is a constant that stores the version information.
// nolint:deadcode // These values are populated by 'govvv' see https://github.com/troian/govvv
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
	for {
		cfg := trollr.AppConfig{
			Listen: ":7891",
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

		// Create the parts
		px := util.NewPoolExecutor(util.Run)
		troll := trollr.NewTroll(cfg.TrollBin, px.Run)
		http := trollr.NewAPI(troll)

		// Start the HTTP Server
		http.Start()

		// Confirm configuration values
		fmt.Printf("  Running with Configuration:\n")
		fmt.Printf("    Listen:    '%s'\n", cfg.Listen)
		fmt.Printf("    Mosmllib:  '%s'\n", cfg.Mosmllib)
		fmt.Printf("    Troll Bin: '%s'\n", cfg.TrollBin)
		fmt.Printf("\n\n")

		// Signal and exit handling from here down
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

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
