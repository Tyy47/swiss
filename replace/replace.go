package replace

import (
	"github.com/Tyy47/clibox/argbin"
)


func findStringFlag() *argbin.Flag {
	return &argbin.Flag{
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
			// check if the value is a string
			ctx.Values["string"] = ctx.ParsedFlagValue
			return nil
		},
	}
}

func replaceAmountFlag() *argbin.Flag {
	return &argbin.Flag{
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
			// check if it's an int
			ctx.Values["amount"] = ctx.ParsedFlagValue
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

	}
}
