package replace

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"swiss/utils"

	"github.com/Tyy47/clibox/argbin"
)

// FileInfo stores all the needed information to replace words in a file
type FileInfo struct {
	// Stores the file name
	Name string

	// TargetWord stores the users word they want to replace
	TargetWord string

	// Replacement is the word that replace TargetWord
	Replacement string

	// Seperator adds a space between each entry
	Seperator string

	// Count is the amount that can be manually set to only replace a certain amount of words
	Count int
}

// AssignValues builds all needed FileInfo members and returns an error if any value cannot be retrieved.
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

// getFileContents reads a file based on FileInfo.Name and returns the file contents in bytes. returns an error if unable to read the file.
func getFileContents(info *FileInfo) ([]byte, error) {
	file, err := os.ReadFile(info.Name)
	return file, err
}

// replaceFileContents takes in a files content in bytes along side FileInfo to replaces all or chosen amount of found words in a file.
func replaceFileContents(content []byte, info *FileInfo) error {
	// sliceSplitAndConvert closure to convert byte slice to string slice, trim whitespace and returns the slice.
	sliceSplitAndConvert := func() []string {
		// Split the string converted byte array with FileInfo.Seperator
		stringedSlice := strings.Split(string(content), info.Seperator)

		// Trim all random spaces on words
		for i, single := range stringedSlice {
			stringedSlice[i] = strings.TrimSpace(single)
		}

		// Return stringed slice
		return stringedSlice
	}()

	// Error check to see if it's an empty file
	if len(sliceSplitAndConvert) == 0 {
		return fmt.Errorf("%s is an empty file.", info.Name)
	}

	// Replacement list copied from the search list
	replaceList := sliceSplitAndConvert

	// Counter to count string changes in a file
	var changeCounter int

	// Loop over the stringed search array
	for i, single := range sliceSplitAndConvert {

		// If the change counter meets the Count in fileinfo, then it breaks out of the loop
		if info.Count != 0 {
			if changeCounter == info.Count {
				break
			}
		}

		// grabPunc closure to locate any punctuation in a given word. if found, returns the substring and the cut suffix.
		grabPunc := func() (subString string, suffix string) {
			// All available punctuation
			puncList := []string{
				".", ",", ";", ":", "{", "[", "]", `}`,
				"!", "@", "$", "%", "^", "&", "*",
				"-", "--", "_", "__", "+", "=", "`", "~",
			}

			// Search the string to see if it has punctuation
			for _, suffix := range puncList {
				if subString, ok := strings.CutSuffix(single, suffix); ok {
					return subString, suffix
				}
			}

			return "", ""
		}

		// Shadow single and grab the suffix from the punc closure
		single, suffix := grabPunc()

		// If the target word is found, it is changed and the change counter is incremented
		if single == info.TargetWord {
			if suffix != "" {
				replaceList[i] = info.Replacement + suffix
				changeCounter += 1
				continue
			}

			// Replace word in array
			replaceList[i] = info.Replacement
			changeCounter += 1
		}
	}

	// byteArray closure to convert the string array back into bytes
	byteArray := func(stringArray []string) []byte {
		// Join all array contents to single string
		newString := strings.Join(stringArray, " ") + "\n"

		// Convert and return the bytes from the string
		bytes := []byte(newString)
		return bytes
	}(replaceList)

	// Write the byteArray to the file to replace
	if err := os.WriteFile(info.Name, byteArray, 0o755); err != nil {
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
			"-f":       getFileNameFlag(),
			"--file":   getFileNameFlag(),
			"-h":     utils.HelpFlag(),
			"--help": utils.HelpFlag(),
			"-w":       findStringFlag(),
			"--with":   findStringFlag(),
		},
		HelpMenu: `
╭───────────────────  Swiss  ────────────────────╮
│                                                │
│       The army knife of CLI applications       │
│                                                │
╰────────────────────────────────────────────────╯
Replace module - Replace words in a file via Swiss.

Commands:
	replace <word>: replaces a word with the value given

Flags:
	-a --amount: Specify an amount of times a word should be replaced.
	-f --file: Location to the file.
	-h --help: Opens the help menu.
	-w --with: The word your using as the replacement.`,
	}
}
