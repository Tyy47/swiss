package initialize

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"swiss/utils"

	"github.com/Tyy47/clibox/argbin"
	"github.com/Tyy47/clibox/inputbin"
	"github.com/Tyy47/clibox/outbin"
)

var (
	ErrUnknownProject = errors.New("unknown project name")
)

const initProjectList = `Languages:
Rust - Cargo: swiss init rust
Go - Go: swiss init go
C - Swiss: swiss init c
Zig - Zig: swiss init zig
Python - uv: swiss init python
Typescript - bun: swiss init ts`

var output = outbin.NewOutput(os.Stdout, os.Stderr)

// findProject reads through a map of projects. If a project is found,
// It'll return. If not, it will return an error.
func findProject(s string) (*project, error) {
	lower := strings.ToLower(s)

	proj, ok := initMap[Language(lower)]
	if !ok {
		return nil, fmt.Errorf("%w %s", ErrUnknownProject, lower)
	}

	return &proj, nil

}

// createProjectDirectory creates the files and folders required for the project init
func createProjectDirectory(p *project) error {

	// Create folders based on the project request
	if p.Folders != nil {
		if err := utils.MakeFolder(p.Folders...); err != nil {
			return err
		}
	}

	// Create files based on the project request
	if p.Files != nil {
		if err := utils.MakeFile(p.Files...); err != nil {
			return err
		}
	}

	return nil
}

// createExecuteCommand builds the command to create the project and returns it as a
// pointer to an exec Cmd.
func createExecuteCommand(proj *project) *exec.Cmd {

	if proj.NeedsProjectName {

		ops := inputbin.InputOptions{
			Question: "enter project name: ",
		}

		// Gather user input
		userInput, err := inputbin.Text(&ops)
		if err != nil {
			// If input is blank it will get current working directory name
			userInput, err = func() (string, error) {
				// Grab working directory
				wd, _ := os.Getwd()

				// Split the directory string
				split := strings.Split(wd, "/")

				// Grab last entry
				if len(split) == 1 {
					return split[0], nil
				}

				return split[len(split)-1], nil
			}()

			// Defaults the project name to "project" if name cannot be gotten elsewhere
			if err != nil {
				userInput = "project"
			}
		}

		// Adds project name to arguments list
		proj.Arguments = append(proj.Arguments, userInput)
	}

	// Create the command
	cmd := exec.Command(proj.Tool, proj.Arguments...)

	return cmd
}

// initialize is the main function on the init process for swiss.
// It grabs the project from findProject and feeds it into createProjectDirectory.
// It returns the project and err if the two functions fail.
func initialize(projectName string) (*project, error) {

	// Finds the project based on the users input
	proj, err := findProject(projectName)
	if err != nil {
		return nil, err
	}

	// Creates the project folders & files.
	if err := createProjectDirectory(proj); err != nil {
		return nil, err
	}

	return proj, nil
}

// gitInit runs "git init" in the current directory
func gitInit() error {
	// Initing git command
	init := exec.Command("git", "init")

	// Runs git init and returns an error if unable to
	if err := init.Run(); err != nil {
		return err
	}

	return nil
}

// jjInit runs "jj git init" in the current directory
func jjInit() error {
	// Init jujutsu command
	init := exec.Command("jj", "git", "init")

	// Runs the init command and returns the error if unsuccessful
	if err := init.Run(); err != nil {
		return err
	}

	return nil
}

func InitCommand() *argbin.Command {
	return &argbin.Command{
		Name:       "init",
		TakesValue: true,
		Execute: func(cmdCtx *argbin.CommandContext) error {
			// Toggle for silent output
			silent, _ := cmdCtx.Values["silent"].(bool)

			utils.Output.Infof("initializing %s project", cmdCtx.ParsedValue)

			// Finds the project
			proj, err := initialize(cmdCtx.ParsedValue)
			if err != nil {
				return err
			}

			// Create the command to initalize a project.
			initCmd := createExecuteCommand(proj)

			// Toggles stdout & stderr based on silent flag
			if !silent {
				// Put cmds output to stdout and stderr
				utils.ToggleOutputForCMD(initCmd)
			}

			// Execute init command
			if !proj.ManualInit {
				if err := initCmd.Run(); err != nil {
					return err
				}
			}
			
			// Displays success message if silent isnt toggled
			if !silent { output.Successf("%s project has been created.", cmdCtx.ParsedValue) }
			return nil
		},
		Flags: argbin.Flags{
			"-h":     utils.HelpFlag(),
			"--help": utils.HelpFlag(),
			"--silent": {
				Execute: func(cmdCtx *argbin.CommandContext) error {
					cmdCtx.Values["silent"] = true
					return nil
				},
			},
			"-s": {
				Execute: func(cmdCtx *argbin.CommandContext) error {
					cmdCtx.Values["silent"] = true
					return nil
				},
			},
			"-g": {
				Execute: func(cmdCtx *argbin.CommandContext) error {
					silent, _ := cmdCtx.Values["silent"].(bool)

					if err := gitInit(); err != nil {
						output.Error("unable to init git")
						return nil
					}

					if !silent {
						output.Success("git initialized successfully")
					}

					return nil
				},
			},
			"--git": {
				Execute: func(cmdCtx *argbin.CommandContext) error {
					silent, _ := cmdCtx.Values["silent"].(bool)

					if err := gitInit(); err != nil {
						output.Error("unable to init git")
						return nil
					}

					if !silent {
						output.Success("git initialized successfully")
					}

					return nil
				},
			},
			"-j": {
				Execute: func(cmdCtx *argbin.CommandContext) error {
					silent, _ := cmdCtx.Values["silent"].(bool)

					if err := jjInit(); err != nil {
						output.Error("unable to init jujutsu")
						return nil
					}

					if !silent {
						output.Success("jujutsu initialized successfully")
					}

					return nil
				},
			},
			"--jujutsu": {
				Execute: func(cmdCtx *argbin.CommandContext) error {
					silent, _ := cmdCtx.Values["silent"].(bool)

					if err := jjInit(); err != nil {
						output.Error("unable to init jujutsu")
						return nil
					}

					if !silent {
						output.Success("jujutsu initialized successfully")
					}

					return nil
				},
			},
			"-e": {
				Execute: func(cmdCtx *argbin.CommandContext) error {
					files := []string{"TODO.md", "README.md"}

					if err := utils.MakeFile(files...); err != nil {
						return err
					}

					return nil
				},
			},
			"--extras": {
				Execute: func(cmdCtx *argbin.CommandContext) error {
					files := []string{"TODO.md", "README.md"}

					if err := utils.MakeFile(files...); err != nil {
						return err
					}

					return nil
				},
			},
		},
		HelpMenu: `
╭───────────────────  Swiss  ────────────────────╮
│                                                │
│       The army knife of CLI applications       │
│                                                │
╰────────────────────────────────────────────────╯
Init module - Initialize a project using Swiss.

Commands:
	init <string> [name: string]: Inits a project based on the given input.

Flags:
	-h --help: Opens the help menu.
	-l --list: Prints a list of projects that can be initialized and if they are supported with additional arguments for names.
	-g --git: Inits git alongside your project.
	-j --jujutsu: Inits jj alongside your project.
	-e, --extras: Creates extra files for your project. (TODO.md & README.md)
`,
	}
}
