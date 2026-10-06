package replace

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Tyy47/clibox/argbin"
)

type FileInfo struct {
	Name        string
	TargetWord  string
	Replacement string
	Seperator   string
	Count       int
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

	// Search closure to grab a stringed, trimmed content array
	search := func() []string {

		x := strings.Split(string(content), info.Seperator)

		strings.Split(string(content), info.Seperator)

		for i, single := range x {
			x[i] = strings.TrimSpace(single)
		}

		return x
	}()

	// Error check to see if it's an empty file
	if len(search) == 0 {
		return fmt.Errorf("%s is an empty file.", info.Name)
	}

	// Replacement list to switch changed words
	replaceList := search

	// Counter to count string changes in a file
	var changeCounter int

	// Loop over the stringed search array
	for i, single := range search {

		// If the change counter meets the Count in fileinfo, then it breaks out of the loop
		if info.Count != 0 {
			if changeCounter == info.Count {
				break
			}
		}

		// If the target word is found, it is changed and the change counter is incremented
		if info.TargetWord == single {
			trimmedString := strings.Trim(info.Replacement, "")

			punc := func(s string) string {
				// All available punctuation
				puncList := []string{
					".", ",", ";", ":", "{", "[", "]", `}`,
					"!", "@", "#", "$", "%", "^", "&", "*",
					"-", "--", "_", "__", "+", "=", "`", "~",
				}

				// Search the string to see if it has punctuation
				for _, suffix := range puncList {
					if strings.HasSuffix(s, suffix) {
						return suffix
					}
				}

				return ""
			}(trimmedString)

			// Punctuation nil check
			if punc != "" {
				replaceList[i] = trimmedString + punc
				changeCounter += 1
				continue
			}

			// Replace word in array
			replaceList[i] = trimmedString
			changeCounter += 1
		}
	}

	// byteArray closure to convert the string array back into bytes
	byteArray := func(stringArray []string) []byte {
		newString := strings.Join(stringArray, " ") + "\n"

		bytes := []byte(newString)

		return bytes
	}(replaceList)

	// Write the byteArray to the file to replace
	if err := os.WriteFile(info.Name, byteArray, 0755); err != nil {
		return fmt.Errorf("Unable to write to file %s due to error %w", info.Name, err)
	}

	return nil
}

func findStringFlag() *argbin.Flag {
	return &argbin.Flag{
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
			// Assign parsed flag value to command context
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
			convert, err := strconv.Atoi(strings.TrimSpace(ctx.ParsedFlagValue))
			if err != nil {
				return fmt.Errorf("unable to convert %s to an int", ctx.ParsedValue)
			}

			// Add amount value to command context
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
		Name:       "replace",
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {

			// Init file info
			info := &FileInfo{}

			// Assign values to info
			if err := info.AssignValues(ctx); err != nil {
				return err
			}

			// Get file contents
			content, err := getFileContents(info)
			if err != nil {
				return err
			}

			// Replace file contents with replacement word
			if err := replaceFileContents(content, info); err != nil {
				return err
			}

			return nil
		},
		Flags: argbin.Flags{
			"-a":       replaceAmountFlag(),
			"--amount": replaceAmountFlag(),
			"-w":       findStringFlag(),
			"--with":   findStringFlag(),
			"-f":       getFileNameFlag(),
			"--file":   getFileNameFlag(),
		},
	}
}
