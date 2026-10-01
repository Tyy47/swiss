package replace

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Tyy47/clibox/argbin"
)

type FileInfo struct {
	Name string
	Replacement string
	Seperator string
	Count int
}

func (fi *FileInfo) AssignValues(ctx *argbin.Context) error {
	// Grab file name from command context
	name, ok := ctx.Values["file_name"].(string)
	if !ok {
		return fmt.Errorf("need file name from -f or --file")
	}

	// Grab replacement string from -w or --with flag
	replace, ok := ctx.Values["with"].(string)
	if !ok {
		return fmt.Errorf("need replacement string from -w or --with")
	}

	// Get count from command context if supplied
	count, ok := ctx.Values["amount"].(int)
	if ok {
		fi.Count = count
	}
	
	// Assigning a name to file fi
	fi.Name = name
	fi.Replacement = replace
	fi.Seperator = " "

	return nil
}

func getFileContents(info *FileInfo) ([]byte, error) {
	file, err := os.ReadFile(info.Name)
	return file, err
}

func searchFileContents(content []byte, info *FileInfo) ([]int, error) {
}


func findStringFlag() *argbin.Flag {
	return &argbin.Flag{
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
			ctx.Values["with"] = ctx.ParsedFlagValue
			return nil
		},
	}
}

func replaceAmountFlag() *argbin.Flag {
	return &argbin.Flag{
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
	
			// Convert user input into a int
			convert, err := strconv.Atoi(ctx.ParsedValue)
			if err != nil {
				return fmt.Errorf("unable to convert %s to an int", ctx.ParsedValue)
			}

			ctx.Values["amount"] = convert
			return nil
		},
	}
}

func getFileNameFlag() *argbin.Flag {
	return &argbin.Flag{
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
			// Grab file name value
			ctx.Values["file_name"] = ctx.ParsedFlagValue
			return nil
		},
	}
}

func ReplaceCommand() *argbin.Command {
	return &argbin.Command{
		Name: "replace",
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
			
			// Init file info
			info := &FileInfo{}

			// Assign values to info
			if err := info.AssignValues(ctx); err != nil {
				return err
			}

			content, err := getFileContents(info)
			if err != nil { return err }

			return nil
		},
		Flags: argbin.Flags{
			"-a": replaceAmountFlag(),
			"--amount": replaceAmountFlag(),
			"-w": findStringFlag(),
			"--with": findStringFlag(),
			"-f": getFileNameFlag(),
			"--file": getFileNameFlag(),
		},
	}
}
