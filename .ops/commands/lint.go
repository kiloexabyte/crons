package commands

import (
	"context"
	"fmt"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Lint() error {
	ctx := context.Background()
	sh := command.Shell(sys.Machine(), "go", "golangci-lint")

	err := sh.Exec(ctx, "golangci-lint", "run")
	if err != nil {
		return err
	}

	if err := sh.Exec(ctx, "go", "fmt"); err != nil {
		return fmt.Errorf("go fmt: %w", err)
	}

	return nil
}
