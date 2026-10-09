package build

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"swiss/utils"

	"github.com/Tyy47/clibox/argbin"
)

var (

	// generic errors

	ErrUnknownLanguage = errors.New("unknown language")

	// build object errors

	ErrNilLang      = errors.New("lang argument cannot be nil")
	ErrNilBuildTool = errors.New("build tool argument cannot be nil")
	ErrNilBuildArgs = errors.New("build args argument cannot be nil")
	ErrNilBuildFile = errors.New("build file argument cannot be nil")

	// program type errors

	ErrNilProgram = errors.New("program return cannot be nil")
	ErrNilBuild   = errors.New("build in program type cannot be nil")
	ErrNilRun     = errors.New("run in program type cannot be nil")
)

// Language safely adds language guard rails
type Language string

// Available languages for swiss to build & run
const (
	Go   Language = "go"
	Rust Language = "rust"
)

// build stores all the needed information to build a program
type build struct {
	Tool           string   // Tool stores the build tool of a language (e.g cargo, bun, etc)
	BuildArguments []string // BuildArguments stores optional flags for the build Tool
	BuildFile      string   // File thats associated with a certain language I.E main.go, Cargo.toml, etc.
}

// run stores all the needed information to run a program
type run struct {
	Tool           string   // Tool stores the run tool of a language (e.g cargo, bun, etc)
	BuildArguments []string // BuildArguments stores optional flags for the run Tool
	BuildFile      string   // File thats associated with a certain language I.E main.go, Cargo.toml, etc.
}

// program stores context between functions to get build and run objects for later function executions
type program struct {
	Build *build
	Run   *run
}

// buildMap contains all of the available languages that can be built via swiss.
var buildMap = map[Language]build{
	Go: {
		Tool:           "go",
		BuildArguments: []string{"build"},
		BuildFile:      "go.mod",
	},

	Rust: {
		Tool:           "cargo",
		BuildArguments: []string{"build"},
		BuildFile:      "Cargo.toml",
	},
}

// runMap contains all of the available languages that can be ran via swiss
var runMap = map[Language]run{
	Go: {
		Tool:           "go",
		BuildArguments: []string{"run", "."},
		BuildFile:      "go.mod",
	},

	Rust: {
		Tool:           "cargo",
		BuildArguments: []string{"run"},
		BuildFile:      "Cargo.toml",
	},
}

// manualFindLanguage grabs all files in the current directory
// and looks for a BuildFile to determine what language is being used and what needs to be ran or built.
func manualFindLanguage(ctx *argbin.Context) (*program, error) {
	// Grab all files in the current directory
	files, err := os.ReadDir("./")
	if err != nil {
		return nil, err
	}

	langs := &program{}

	// Loop over each file
	for _, file := range files {
		// If the file is a directory, then skip
		if file.Type().IsDir() {
			continue
		}

		// Loop over the build map to get the build object
		for k, v := range buildMap {
			if file.Name() != v.BuildFile {
				continue
			}

			ctx.Values["language"] = k
			langs.Build = &v
		}

		// Loop over the run map to get the run object
		for k, v := range runMap {
			if file.Name() != v.BuildFile {
				continue
			}

			ctx.Values["language"] = k
			langs.Run = &v
		}
	}
	return langs, nil
}

// getLanguage takes a language input and returns a build object thats related to that language.
// If no languages are found, ErrUnknownLanguage is returned.
func getLanguage(lang string) (*program, error) {
	// Nil check for lang
	if lang == "" {
		return nil, ErrNilLang
	}

	// Lowercase the input
	fixed := strings.ToLower(lang)

	// Gather the build object
	buildObject, ok := buildMap[Language(fixed)]
	if !ok {
		return nil, ErrUnknownLanguage
	}

	// Gather the run object
	runObject, ok := runMap[Language(fixed)]
	if !ok {
		return nil, ErrUnknownLanguage
	}

	// Create a program and assign the build and run objects to store for context
	langs := program{
		Build: &buildObject,
		Run:   &runObject,
	}

	// Return the build object
	return &langs, nil
}

// buildLanguage handles the execution of the build of a program.
func buildLanguage(p *program) error {
	// Nil check program
	if p == nil {
		return ErrNilProgram
	}

	// Nil check for build
	if p.Build == nil {
		return ErrNilBuild
	}

	// Assigning build to b
	b := p.Build

	// Takes in the tool and arguments to create a command
	cmd := exec.Command(b.Tool, b.BuildArguments...)

	// Assigns cmd's output to stdout and stderr
	utils.ToggleOutputForCMD(cmd)

	// Runs the command and returns the error if it fails
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// runLanguage handles the "running" of a program.
// It will take a run object from program and gather all the information required
// to run the program.
func runLanguage(p *program) error {
	// Nil check for program
	if p == nil {
		return ErrNilProgram
	}

	// Nil check for stored build
	if p.Build == nil {
		return ErrNilRun
	}

	// Assigning run to r
	r := p.Run

	// Takes in the tool and arguments to create a command
	cmd := exec.Command(r.Tool, r.BuildArguments...)

	// Assigns cmd's output to stdout and stderr
	utils.ToggleOutputForCMD(cmd)

	// Runs the command and returns the error if it fails
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// listFlag returns a flag for build and run to print the list of valid languages
// that swiss can build and run.
func listFlag() *argbin.Flag {
	return &argbin.Flag{
		Execute: func(ctx *argbin.Context) error {
			fmt.Println("Buildable and runable languages:")

			// Loop over build map and print entries
			for k, v := range buildMap {
				fmt.Printf("%s: required tool = %s.\n", k, v.Tool)
			}

			return nil
		},
		Terminal: true,
	}
}

func showArgsFlag() *argbin.Flag {
	return &argbin.Flag{
		Execute: func(ctx *argbin.Context) error {
			ctx.Values["show_args"] = true
			return nil
		},
	}
}

// SwissInstall installs swiss for Linux users when ran inside of the cloned swiss repo.
func SwissInstall() *argbin.Command {
	return &argbin.Command{
		Name: "install",
		Execute: func(ctx *argbin.Context) error {
			// Check if the user is running linux
			if utils.GetOperatingSystem() != "linux" {
				utils.Output.Warning("unable to install swiss. unsupported operating system")
				return nil
			}

			// Install message
			utils.Output.Info("installing swiss.")

			// Get language
			lang, err := getLanguage("go")
			if err != nil {
				return err
			}

			// Build swiss
			if err := buildLanguage(lang); err != nil {
				return err
			}

			// Create destination path
			desPath := fmt.Sprintf("/home/%s/.local/bin/swiss", utils.GetUsersName())

			// Build the move command
			moveCmd := exec.Command("mv", "swiss", desPath)

			// Execute the move command
			if err := moveCmd.Run(); err != nil {
				return err
			}

			// Success message after install
			utils.Output.Success("swiss has been installed! make sure .local/bin is added to your path")

			return nil
		},
	}
}

// getLanguageName gathers the language name to build or run. 
// *NOTE* Only used in BuildCommand & RunCommand
func getProjectAndName(ctx argbin.Context) (*program, string, error) {
	var proj *program
	var title string
	var err error

	if ctx.ParsedValue == "" {
		proj, err = manualFindLanguage(&ctx)
		if err != nil {
			return nil, "", err
		}

		temp, ok := ctx.Values["language"].(Language)
		if !ok {
			return nil, "", fmt.Errorf("swiss doesn't support this language.")
		}
		title = string(temp)

	} else {
		proj, err = getLanguage(ctx.ParsedValue)
		if err != nil {
			return nil, "", err
		}
		title = ctx.ParsedValue
	}

	return proj, title, nil
}

// showProjectArguments displays the arguments being ran for build and run depending on the command ran.
func showProjectArguments(ctx argbin.Context, title string, args ...string) error {
	if show, exists := ctx.Values["show_args"].(bool); exists && show {
		utils.Output.Infof("%s command args: %s.", title, strings.Join(args, ", "))
		if cont := utils.ContinuePrompt(); !cont {
			return utils.ErrExitingSwissGracefully
		}
	}

	return nil
}



// BuildCommand creates the "build" command for swiss.
// Command is added in main.go in main function.
func BuildCommand() *argbin.Command {
	return &argbin.Command{
		Name: "build",
		Execute: func(ctx *argbin.Context) error {
			utils.Output.Info("searching for language.")

			b, title, err := getProjectAndName(*ctx)
			if err != nil {
				return err
			}

			if err := showProjectArguments(*ctx, title, b.Build.BuildArguments...); err != nil {
				return err
			}

			utils.Output.Info("building program.")
			if err := buildLanguage(b); err != nil {
				return err
			}

			utils.Output.Successf("%s program has been built", title)
			return nil
		},
		Flags: argbin.Flags{
			"-a":     showArgsFlag(),
			"--args": showArgsFlag(),
			"-h":     utils.HelpFlag(),
			"--help": utils.HelpFlag(),
			"-l":     listFlag(),
			"--list": listFlag(),
		},
		HelpMenu: `
╭───────────────────  Swiss  ────────────────────╮
│                                                │
│       The army knife of CLI applications       │
│                                                │
╰────────────────────────────────────────────────╯
Build & Run module - Builds or Runs a program based on the language provided.

Commands:
	build <string>: Builds a program based on the language you input.
	run <string>: Runs a program based on the language you input.
Flags:
	-h --help: Opens the help menu.
	-l --list: Prints a list of available languages to build and run with their respective build tools available in swiss.
	-a --args: Shows all arguments used when building or running a program via swiss.`,
	}
}

// RunCommand create the "run" command for swiss.
// Command is added in main.go in main function.
func RunCommand() *argbin.Command {
	return &argbin.Command{
		Name: "run",
		Execute: func(ctx *argbin.Context) error {
			utils.Output.Info("searching for language.")

			r, title, err := getProjectAndName(*ctx)
			if err != nil {
				return err
			}

			if err := showProjectArguments(*ctx, title, r.Run.BuildArguments...); err != nil {
				return err
			}


			utils.Output.Info("running program.")
			if err := runLanguage(r); err != nil {
				return err
			}

			utils.Output.Successf("%s program has been ran.", title)
			return nil
		},
		Flags: argbin.Flags{
			"-a":     showArgsFlag(),
			"--args": showArgsFlag(),
			"-h":     utils.HelpFlag(),
			"--help": utils.HelpFlag(),
			"-l":     listFlag(),
			"--list": listFlag(),
		},
		HelpMenu: `
╭───────────────────  Swiss  ────────────────────╮
│                                                │
│       The army knife of CLI applications       │
│                                                │
╰────────────────────────────────────────────────╯
Build & Run module - Builds or Runs a program based on the language inputted.

Commands:
	build <string>: Builds a program based on the language you input.
	run <string>: Runs a program based on the language you input.
Flags:
	-h --help: Opens the help menu.
	-l --list: Prints a list of available languages to build and run with their respective build tools available in Swiss.`,
	}
}
