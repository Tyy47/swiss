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
	TargetWord string
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
	fi.TargetWord = ctx.ParsedValue
	fi.Replacement = replace
	fi.Seperator = " "

	return nil
}

func getFileContents(info *FileInfo) ([]byte, error) {
	file, err := os.ReadFile(info.Name)
	return file, err
}

func replaceFileContents(content []byte, info *FileInfo) error {
	
	search := func() []string {

		x := strings.Split(string(content), info.Seperator)

		strings.Split(string(content), info.Seperator)

		for i, single := range x {
			x[i] = strings.TrimSpace(single)
		}

		return x
	}()

	if len(search) == 0 {
		return fmt.Errorf("%s is an empty file.", info.Name)
	}

	replaceList := search

	for i, single := range search {
		if single == info.TargetWord {
			replaceList[i] = info.Replacement
		}
	}

	byteArray := func(stringArray []string) []byte {
		newString := strings.Join(stringArray, " ") + "\n"

		bytes := []byte(newString)

		return bytes
	}(replaceList)


	if err := os.WriteFile(info.Name, byteArray, 0755); err != nil {
		return fmt.Errorf("Unable to write to file %s due to error %w", info.Name, err)
	}

	return nil
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

			replaceFileContents(content, info)

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
