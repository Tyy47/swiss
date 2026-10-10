package main

import (
	"errors"
	"fmt"

	"swiss/build"
	"swiss/fileserver"
	"swiss/gen"
	"swiss/initialize"
	"swiss/replace"
	"swiss/shortcuts"
	"swiss/utils"

	"github.com/Tyy47/clibox/argbin"
	"github.com/Tyy47/clibox/colorbin"
)

// Creating the root object of the application
var root = argbin.Root{
	AppName:     "swiss",
	AppVersion:  "1.2",
	CommandList: make([]*argbin.Command, 0),
	HelpMenu: `
╭───────────────────  Swiss  ────────────────────╮
│                                                │
│       The army knife of CLI applications       │
│                                                │
╰────────────────────────────────────────────────╯

usage: swiss [command] [additional_arguments] <flags>

Commands:
	build: Builds a program with a native build tool and swiss as a shorthand wrapper.
	run: Runs a program with a native build tool and swiss as a shorthand wrapper.
	init: Initializes a programming based project in current folder.
	gen: Generates different codes that are most commonly used in development
	sc: Command shortcuts for various CLI utilities to make development faster
	replace: Replaces words in a file.

Flags:
	-h, --help: Displays the swiss help menu
	-v, --version: Displays the current swiss version number`,
}

// helpCommand creates the "help" command for the root.
func helpCommand() *argbin.Command {
	return &argbin.Command{
		Name:            "help",
		AdditionalNames: []string{"--help", "-h"},
		Execute: func(ctx *argbin.Context) error {
			fmt.Println(root.HelpMenu)
			return nil
		},
	}
}

// versionCommand creates the "version" command for the root.
func versionCommand() *argbin.Command {
	return &argbin.Command{
		Name:            "version",
		AdditionalNames: []string{"--version", "-v"},
		Execute: func(ctx *argbin.Context) error {
			// Color app version to green
			version := colorbin.Green(root.AppVersion).ToHighIntensityBold().String()

			// Profit
			fmt.Printf("swiss: version %s\n", version)
			return nil
		},
	}
}

func main() {
	// Command storage to add to root.AddCommand
	commands := []*argbin.Command{
		helpCommand(),
		versionCommand(),
		build.BuildCommand(),
		build.RunCommand(),
		build.SwissInstall(),
		gen.GenerateCommand(),
		initialize.InitCommand(),
		shortcuts.ShortcutCommand(),
		replace.ReplaceCommand(),
		fileserver.FileServerCommand(),
	}

	// Adds all commands to app
	if err := root.AddCommand(commands...); err != nil {
		// Panic for potential programmer errors
		panic(err)
	}

	// Starts the project using argbin
	if err := root.Run(); err != nil {
		// Shows help menu if arguments are missing from swiss
		if errors.Is(err, argbin.ErrMissingArguments) {
			fmt.Println(root.HelpMenu)
			return
		}

		// Checks if swiss was exited gracefully, unwraps error to string to print it as a note 
		// instead of an error.
		if errors.Is(err, utils.ErrExitingSwissGracefully) {
			utils.Output.Info(err.Error())
			return
		}
		
		// Prints if no specific errors are specified
		utils.Output.Error(err)
	}
}
