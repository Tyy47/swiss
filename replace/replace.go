package replace

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	//	"swiss/utils"

	"github.com/Tyy47/clibox/argbin"
)




func findStringFlag() *argbin.Flag {
	return &argbin.Flag{
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
			ctx.Values["string"] = ctx.ParsedFlagValue
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

func ReplaceCommand() *argbin.Command {
	return &argbin.Command{
		Name: "replace",
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
			return nil
		},
		Flags: argbin.Flags{
			"-a": replaceAmountFlag(),
			"--amount": replaceAmountFlag(),
			"-s": findStringFlag(),
			"--string": findStringFlag(),
		},
	}
}
