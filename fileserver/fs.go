package fileserver

import (
	"net/http"
	"os"

	"swiss/utils"

	"github.com/Tyy47/clibox/argbin"
)

func setPortFlag() *argbin.Flag {
	return &argbin.Flag{
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
			// Add port to context
			ctx.Values["port"] = ctx.ParsedFlagValue

			return nil
		},
	}
}

func FileServerCommand() *argbin.Command {
	return &argbin.Command{
		Name:            "fileserver",
		AdditionalNames: []string{"fs"},
		Execute: func(ctx *argbin.Context) error {
			// Default port if one isn't provided
			const defaultPort string = "3030"

			// Host port
			var port string

			// Assign port to gathered input
			if gatheredPort, exists := ctx.Values["port"].(string); exists {
				port = gatheredPort
			} else {
				port = defaultPort
			}

			// Get current working directory
			dir, err := os.Getwd()
			if err != nil {
				return err
			}

			server := http.FileServer(http.Dir(dir))

			http.Handle(dir, server)

			utils.Output.Infof("starting fileserver on port %s", port)

			if err := http.ListenAndServe(":" + port, nil); err != nil {
				return err
			}

			return nil
		},
		Flags: argbin.Flags{
			"-p":     setPortFlag(),
			"--port": setPortFlag(),
		},
	}
}
