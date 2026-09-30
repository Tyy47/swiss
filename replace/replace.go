package replace

import (
	"github.com/Tyy47/clibox/argbin"
)






func ReplaceCommand() *argbin.Command {
	return &argbin.Command{
		Name: "replace",
		TakesValue: true,
		Execute: func(ctx *argbin.Context) error {
			

			return nil
		},

	}
}
