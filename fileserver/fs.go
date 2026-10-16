package fileserver

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"swiss/utils"

	"github.com/Tyy47/clibox/argbin"
	"github.com/Tyy47/clibox/colorbin"
	"github.com/Tyy47/clibox/inputbin"
)

func setPortFlag() *argbin.Flag {
	return &argbin.Flag{
		TakesValue: true,
		Execute: func(cmdCtx *argbin.CommandContext) error {
			// Add port to context
			cmdCtx.Values["port"] = cmdCtx.ParsedFlagValue

			return nil
		},
	}
}

func FileServerCommand() *argbin.Command {
	return &argbin.Command{
		Name:            "fileserver",
		AdditionalNames: []string{"fs"},
		Execute: func(cmdCtx *argbin.CommandContext) error {
			// Host port
			var port string = "3030"

			// Assign port to gathered input
			if gatheredPort, exists := cmdCtx.Values["port"].(string); exists {
				port = gatheredPort
			}

			// Get current working directory
			dir, err := os.Getwd()
			if err != nil {
				return err
			}

			// Create context for http server
			serverCtx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Create server object
			serve := &http.Server{
				Addr:    ":" + port,
				Handler: http.FileServer(http.Dir(dir)),
			}

			// Create an error channel, host the server
			// and receive any potential errors.
			errorChan := make(chan error, 1)
			go func() {
				utils.Output.Infof("serving %s on port %s", dir, port)
				errorChan <- serve.ListenAndServe()
			}()

			// Goroutine to prompt the user if
			// they want to close the file server.
			go func() {
				// Input options
				input := inputbin.InputOptions{
					Prefix: "Close fileserver? [y/n]",
					PrefixColor: &colorbin.ColorOptions{
						ChosenColor: colorbin.ColorYellow,
					},
				}

				// Sleep for 1 seconds to have confirmation prompt
				// be the last print
				time.Sleep(1 * time.Second)
				for {
					// Present prompt
					if inputbin.Confirm(&input, true) {
						cancel()
						utils.Output.Infof("closing fileserver on port %s", port)
						break
					}
				}
			}()

			// Server shutdown process
			select {
			// Check if the context has finished
			case <-serverCtx.Done():
				// Grab shutdown context
				shutdownCtx, stop := context.WithTimeout(context.Background(), 5*time.Second)
				defer stop()

				// Attempt to shutdown server with given context
				if err := serve.Shutdown(shutdownCtx); err != nil {
					return err
				}

				// Returns a received error if it isn't the
				// http server being closed
				if err := <-errorChan; err != nil && !errors.Is(err, http.ErrServerClosed) {
					return err
				}

			// Returns a received error if it isn't the
			// http server being closed
			case err := <-errorChan:
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					return err
				}
			}

			return nil
		},
		Flags: argbin.Flags{
			"-p":     setPortFlag(),
			"--port": setPortFlag(),
		},
	}
}
