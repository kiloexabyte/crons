package commands

import (
	"context"
	"os"

	"lesiw.io/command"
	"lesiw.io/command/sys"

)

func (Ops) Build() error {
	sh := command.Shell(sys.Machine(), "golangci-lint", "go")
	ctx := context.Background()

	// Cross-compile for Raspberry Pi (ARM64)
	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "arm64")

	err := sh.Exec(ctx, "go", "build", "-o", "crons", ".")
	if err != nil {
		return err
	}

	return nil
}
