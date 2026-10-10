package fileserver

import (
	"net/http"
	"strconv"

	"github.com/Tyy47/clibox/argbin"
)



func setPortFlag() *argbin.Flag {
	return &argbin.Flag{
		Execute: func(ctx *argbin.Context) error {
			// Gather the port and convert to int
			port, err := strconv.Atoi(ctx.ParsedFlagValue)
			if err != nil {
				return err
			}

			// Add port to context
			ctx.Values["port"] = port

			return nil
		},

	}
}

func FileServerCommand() *argbin.Command {
	return &argbin.Command{
		Name: "fileserver",
		AdditionalNames: []string{"fs"},
		Execute: func(ctx *argbin.Context) error {
			// Default port if one isn't provided
			const defaultPort int = 3030

			// Host port
			var port int
			
			// Assign port to gathered input
			if gatheredPort, exists := ctx.Values["port"].(int); exists {
				port = gatheredPort
			} else {
				port = defaultPort
			}



			return nil
		},
		Flags: argbin.Flags{
			"-p": setPortFlag(),
			"--port": setPortFlag(),
		},
	}
}
