package replace

import (
	"fmt"

	"github.com/Tyy47/clibox/argbin"
)






func ReplaceCommand() *argbin.Command {
	return &argbin.Command{
		Name: "replace",
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
			// Gathers value correctly
			fmt.Println(ctx.ParsedValue)		

			return nil
		},

	}
}
