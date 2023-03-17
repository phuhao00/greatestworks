package ssh

import (
	"context"

	"greatestworks/aop/tool"
	"greatestworks/aop/tool/ssh/impl"
)

var babysitterCmd = tool.Command{
	Name:        "babysitter",
	Description: "The weaver ssh babysitter",
	Help: `Usage:
  weaver ssh babysitter

Flags:
  -h, --help   Print this help message.`,
	Fn: func(ctx context.Context, args []string) error {
		return impl.RunBabysitter(ctx)
	},
}
